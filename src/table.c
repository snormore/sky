#include <stdlib.h>
#include <fcntl.h>
#include <inttypes.h>
#include <unistd.h>
#include <sys/stat.h>
#include <sys/mman.h>
#include <math.h>
#include <assert.h>

#include "dbg.h"
#include "mem.h"
#include "sky_endian.h"
#include "bstring.h"
#include "file.h"
#include "table.h"

//==============================================================================
//
// Forward Declarations
//
//==============================================================================

//--------------------------------------
// Locking
//--------------------------------------

int sky_table_lock(sky_table *table);

int sky_table_unlock(sky_table *table);

//--------------------------------------
// Tablets
//--------------------------------------

int sky_table_load_tablets(sky_table *table);

int sky_table_unload_tablets(sky_table *table);

//--------------------------------------
// Action file
//--------------------------------------

int sky_table_load_action_file(sky_table *table);

int sky_table_unload_action_file(sky_table *table);

//--------------------------------------
// Property file
//--------------------------------------

int sky_table_load_property_file(sky_table *table);

int sky_table_unload_property_file(sky_table *table);


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------


// Creates a reference to a table.
// 
// Returns a reference to the new table if successful. Otherwise returns
// null.
sky_table *sky_table_create()
{
    sky_table *table = calloc(sizeof(sky_table), 1); check_mem(table);
    table->default_tablet_count = DEFAULT_TABLET_COUNT;
    return table;
    
error:
    sky_table_free(table);
    return NULL;
}

// Removes a table reference from memory.
//
// table - The table to free.
void sky_table_free(sky_table *table)
{
    if(table) {
        bdestroy(table->name);
        table->name = NULL;
        bdestroy(table->path);
        table->path = NULL;
        sky_table_close(table);
        free(table);
    }
}


//--------------------------------------
// Path Management
//--------------------------------------

// Sets the file path of a table.
//
// table - The table.
// path  - The file path to set.
//
// Returns 0 if successful, otherwise returns -1.
int sky_table_set_path(sky_table *table, bstring path)
{
    assert(table != NULL);

    if(table->path) bdestroy(table->path);
    table->path = bstrcpy(path);
    if(path) check_mem(table->path);

    return 0;

error:
    table->path = NULL;
    return -1;
}


//--------------------------------------
// Tablet management
//--------------------------------------

// Retrieves the tablet that an object id would reside in.
//
// table     - The table.
// object_id - The object id to find.
// ret       - A pointer to where the tablet should be returned.
//
// Returns 0 if successful, otherwise returns -1.
int sky_table_get_target_tablet(sky_table *table, sky_object_id_t object_id,
                                sky_tablet **ret)
{
    assert(table != NULL);
    assert(ret != NULL);
    check(table->tablet_count, "Table must have tablets available");
    
    // Calculate the tablet index.
    uint32_t target_index = object_id % table->tablet_count;
    *ret = table->tablets[target_index];
    
    return 0;

error:
    *ret = NULL;
    return -1;
}

// Retrieves the number of existing tablets on a table based on the numeric
// directories in a table's path.
//
// table - The table.
// ret   - A pointer to where the count should be returned.
//
// Returns 0 if successful, otherwise returns -1.
int sky_table_get_tablet_count(sky_table *table, uint32_t *ret)
{
    bstring path = NULL;
    assert(table != NULL);
    
    // Find the last tablet path.
    uint32_t index = 0;
    while(true) {
        path = bformat("%s/%d", bdata(table->path), index); check_mem(path);

        // If we can't find the tablet path then the next 'index' is our
        // count since the index is zero-based.
        if(!sky_file_exists(path)) {
            break;
        }
        
        // Otherwise keep incrementing.
        index++;
        bdestroy(path);
        path = NULL;
    }
    
    // Return the last index we tried.
    *ret = index;
    
    bdestroy(path);
    return 0;

error:
    bdestroy(path);
    *ret = 0;
    return -1;
}

// Initializes and opens the tablets in the table.
//
// table - The table.
//
// Returns 0 if successful, otherwise returns -1.
int sky_table_load_tablets(sky_table *table)
{
    int rc;
    sky_tablet *tablet = NULL;
    assert(table != NULL);
    check(table->path != NULL, "Table path required");
    
    // Unload existing tablets.
    sky_table_unload_tablets(table);
    
    // Retrieve the number of tablets that exist.
    uint32_t tablet_count = 0;
    rc = sky_table_get_tablet_count(table, &tablet_count);
    check(rc == 0, "Unable to determine tablet count");
    
    // If no tablets exist then use the default tablet count and initialize some.
    if(tablet_count == 0) {
        check(table->default_tablet_count > 0, "Default tablet count must be greater than zero");
        tablet_count = table->default_tablet_count;
    }

    // Initialize tablets.
    table->tablets = calloc(tablet_count, sizeof(*table->tablets));
    check_mem(table->tablets);
    table->tablet_count = tablet_count;
    
    uint32_t i;
    for(i=0; i<table->tablet_count; i++) {
        tablet = sky_tablet_create(table); check_mem(tablet);
        tablet->index = i;
        tablet->path = bformat("%s/%d", bdata(table->path), i); check_mem(tablet->path);
        table->tablets[i] = tablet;
        tablet = NULL;
    }

    // Open tablets.
    for(i=0; i<table->tablet_count; i++) {
        rc = sky_tablet_open(table->tablets[i]);
        check(rc == 0, "Unable to open tablet: %s", bdata(table->tablets[i]->path));
    }

    return 0;
error:
    sky_table_unload_tablets(table);
    return -1;
}

// Closes all tablets on the table.
//
// table - The table.
//
// Returns 0 if successful, otherwise returns -1.
int sky_table_unload_tablets(sky_table *table)
{
    assert(table != NULL);

    uint32_t i;
    for(i=0; i<table->tablet_count; i++) {
        sky_tablet *tablet = table->tablets[i];
        sky_tablet_close(tablet);
        sky_tablet_free(tablet);
        table->tablets[i] = NULL;
    }
    free(table->tablets);
    table->tablets = NULL;
    table->tablet_count = 0;

    return 0;
}


//--------------------------------------
// Action file management
//--------------------------------------

// Initializes and opens the action file on the table.
//
// table - The table to initialize the action file for.
//
// Returns 0 if successful, otherwise returns -1.
int sky_table_load_action_file(sky_table *table)
{
    int rc;
    assert(table != NULL);
    check(table->path != NULL, "Table path required");
    
    // Unload any existing action file.
    sky_table_unload_action_file(table);
    
    // Initialize action file.
    table->action_file = sky_action_file_create();
    check_mem(table->action_file);
    table->action_file->path = bformat("%s/actions", bdata(table->path));
    check_mem(table->action_file->path);
    
    // Load data
    rc = sky_action_file_load(table->action_file);
    check(rc == 0, "Unable to load actions");

    return 0;
error:
    sky_table_unload_action_file(table);
    return -1;
}

// Closes the action file on the table.
//
// table - The table to initialize the action file for.
//
// Returns 0 if successful, otherwise returns -1.
int sky_table_unload_action_file(sky_table *table)
{
    assert(table != NULL);

    if(table->action_file) {
        sky_action_file_free(table->action_file);
        table->action_file = NULL;
    }

    return 0;
}


//--------------------------------------
// Property file management
//--------------------------------------

// Initializes and opens the property file on the table.
//
// table - The table to initialize the property file for.
//
// Returns 0 if successful, otherwise returns -1.
int sky_table_load_property_file(sky_table *table)
{
    int rc;
    assert(table != NULL);
    check(table->path != NULL, "Table path required");
    
    // Unload any existing property file.
    sky_table_unload_property_file(table);
    
    // Initialize property file.
    table->property_file = sky_property_file_create();
    check_mem(table->property_file);
    table->property_file->path = bformat("%s/properties", bdata(table->path));
    check_mem(table->property_file->path);
    
    // Load data
    rc = sky_property_file_load(table->property_file);
    check(rc == 0, "Unable to load properties");

    return 0;
error:
    sky_table_unload_property_file(table);
    return -1;
}

// Closes the property file on the table.
//
// table - The table to initialize the property file for.
//
// Returns 0 if successful, otherwise returns -1.
int sky_table_unload_property_file(sky_table *table)
{
    assert(table != NULL);

    if(table->property_file) {
        sky_property_file_free(table->property_file);
        table->property_file = NULL;
    }

    return 0;
}


//--------------------------------------
// State
//--------------------------------------

// Opens the table for reading and writing events.
// 
// table - The table to open.
//
// Returns 0 if successful, otherwise returns -1.
int sky_table_open(sky_table *table)
{
    int rc;
    assert(table != NULL);
    check(table->path != NULL, "Table path is required");
    check(!table->opened, "Table is already open");

    // Create directory if it doesn't exist.
    if(!sky_file_exists(table->path)) {
        rc = mkdir(bdatae(table->path, ""), S_IRWXU);
        check(rc == 0, "Unable to create table directory: %s", bdata(table->path));
    }

    // Obtain a lock.
    rc = sky_table_lock(table);
    check(rc == 0, "Unable to obtain lock");

    // Load data file.
    rc = sky_table_load_tablets(table);
    check(rc == 0, "Unable to load tablets");
    
    // Load action file.
    rc = sky_table_load_action_file(table);
    check(rc == 0, "Unable to load action file");
    
    // Load property file.
    rc = sky_table_load_property_file(table);
    check(rc == 0, "Unable to load property file");
    
    // Flag the table as open.
    table->opened = true;

    return 0;

error:
    sky_table_close(table);
    return -1;
}

// Closes a table.
//
// table - The table to close.
//
// Returns 0 if successful, otherwise returns -1.
int sky_table_close(sky_table *table)
{
    int rc;
    assert(table != NULL);

    // Unload tablets.
    rc = sky_table_unload_tablets(table);
    check(rc == 0, "Unable to unload data file");

    // Unload action data.
    rc = sky_table_unload_action_file(table);
    check(rc == 0, "Unable to unload action file");

    // Unload property data.
    rc = sky_table_unload_property_file(table);
    check(rc == 0, "Unable to unload property file");

    // Update state to closed.
    table->opened = false;

    // Release the lock.
    rc = sky_table_unlock(table);
    check(rc == 0, "Unable to remove lock");

    return 0;
    
error:
    return -1;
}


//--------------------------------------
// Locking
//--------------------------------------

// Creates a lock file for the table.
// 
// table - The table to lock.
//
// Returns 0 if successful, otherwise returns -1.
int sky_table_lock(sky_table *table)
{
    assert(table != NULL);

    // Construct path to lock.
    bstring path = bformat("%s/%s", bdata(table->path), SKY_LOCK_NAME); check_mem(path);

    // Open file with an exclusive lock. If this fails then it means another
    // process already has obtained a lock.
    table->lock_file = fopen(bdata(path), "w");
    check(table->lock_file != NULL, "Table is already locked: %s", bdata(table->path));

    // Clean up.
    bdestroy(path);

    return 0;

error:
    bdestroy(path);
    return -1;
}

// Unlocks a table.
// 
// table - The table to unlock.
//
// Returns 0 if successful, otherwise returns -1.
int sky_table_unlock(sky_table *table)
{
    assert(table != NULL);

    // Generate lock file path.
    bstring path = NULL;

    // Only attempt to remove the lock if one has been obtained.
    if(table->lock_file != NULL) {
        // Close the lock file descriptor.
        fclose(table->lock_file);

        // Remove the lock file. We don't necessarily care if it gets cleaned
        // up so don't check for an error. The write lock is the more
        // important piece.
        path = bformat("%s/%s", bdata(table->path), SKY_LOCK_NAME); check_mem(path);
        unlink(bdatae(path, ""));
    }

    bdestroy(path);
    return 0;

error:
    bdestroy(path);
    return -1;
}


//--------------------------------------
// Event Management
//--------------------------------------

// Adds an event to the table.
//
// table - The table.
// event - The event to add.
//
// Returns 0 if successful, otherwise returns -1.
int sky_table_add_event(sky_table *table, sky_event *event)
{
    int rc;
    assert(table != NULL);
    assert(event != NULL);
    check(table->opened, "Table must be open to add an event");

    // Retrieve the target tablet.
    sky_tablet *tablet = NULL;
    rc = sky_table_get_target_tablet(table, event->object_id, &tablet);
    check(rc == 0, "Unable to determine the target tablet");

    // Delegate to the tablet.
    rc = sky_tablet_add_event(tablet, event);
    check(rc == 0, "Unable to add event to tablet");
    
    return 0;

error:
    return -1;
}

