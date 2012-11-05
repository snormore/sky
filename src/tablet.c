#include <stdlib.h>
#include <fcntl.h>
#include <inttypes.h>
#include <unistd.h>
#include <sys/stat.h>
#include <sys/mman.h>
#include <math.h>

#include "tablet.h"
#include "dbg.h"

//==============================================================================
//
// Forward Declarations
//
//==============================================================================

//--------------------------------------
// Data file
//--------------------------------------

int sky_table_load_data_file(sky_table *table);

int sky_table_unload_data_file(sky_table *table);


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------


// Creates a reference to a tablet.
//
// table - The table the tablet belongs to.
//
// Returns a reference to the new tablet if successful. Otherwise returns
// null.
sky_tablet *sky_tablet_create(sky_table *table)
{
    check(table != NULL, "Table required");
    sky_tablet *tablet = calloc(sizeof(sky_tablet), 1); check_mem(tablet);
    tablet->table = table;
    return tablet;
    
error:
    sky_tablet_free(tablet);
    return NULL;
}

// Removes a tablet reference from memory.
//
// tablet - The tablet to free.
//
// Returns nothing.
void sky_tablet_free(sky_tablet *tablet)
{
    if(tablet) {
        tablet->table = NULL;
        bdestroy(tablet->path);
        tablet->path = NULL;
        tablet->index = 0;
        sky_tablet_close(tablet);
        free(tablet);
    }
}


//--------------------------------------
// Path Management
//--------------------------------------

// Sets the file path of a tablet.
//
// tablet - The tablet.
// path   - The file path to set.
//
// Returns 0 if successful, otherwise returns -1.
int sky_tablet_set_path(sky_tablet *tablet, bstring path)
{
    check(tablet != NULL, "Tablet required");

    if(tablet->path) bdestroy(tablet->path);
    tablet->path = bstrcpy(path);
    if(path) check_mem(tablet->path);

    return 0;

error:
    tablet->path = NULL;
    return -1;
}


//--------------------------------------
// Data file management
//--------------------------------------

// Opens the tablet for reading and writing. This should be called directly.
// Only the owning table should open a tablet.
//
// tablet - The tablet.
//
// Returns 0 if successful, otherwise returns -1.
int sky_tablet_open(sky_tablet *tablet)
{
    int rc;
    check(tablet != NULL, "Tablet required");
    check(tablet->path != NULL, "Tablet path required");
    
    // Close the tablet if it's already open.
    sky_tablet_close(tablet);
    
    // Initialize tablet space (0).
    if(!sky_file_exists(tablet->path)) {
        rc = mkdir(bdata(tablet->path), S_IRWXU);
        check(rc == 0, "Unable to create tablet directory: %s", bdata(tablet->path));
    }
    
    // Initialize data file.
    tablet->data_file = sky_data_file_create();
    check_mem(tablet->data_file);
    tablet->data_file->path = bformat("%s/data", bdata(tablet->path));
    check_mem(tablet->data_file->path);
    tablet->data_file->header_path = bformat("%s/header", bdata(tablet->path));
    check_mem(tablet->data_file->header_path);
    
    // Initialize settings on the block.
    if(tablet->table->default_block_size > 0) {
        tablet->data_file->block_size = tablet->table->default_block_size;
    }
    
    // Load data
    rc = sky_data_file_load(tablet->data_file);
    check(rc == 0, "Unable to load data file");

    return 0;
error:
    sky_tablet_close(tablet);
    return -1;
}

// Closes the tablet.
//
// tablet - The tablet.
//
// Returns 0 if successful, otherwise returns -1.
int sky_tablet_close(sky_tablet *tablet)
{
    check(tablet != NULL, "Tablet required");

    if(tablet->data_file) {
        sky_data_file_free(tablet->data_file);
        tablet->data_file = NULL;
    }

    return 0;
error:
    return -1;
}


//--------------------------------------
// Event Management
//--------------------------------------

// Adds an event to the tablet.
//
// tablet - The tablet.
// event  - The event to add.
//
// Returns 0 if successful, otherwise returns -1.
int sky_tablet_add_event(sky_tablet *tablet, sky_event *event)
{
    int rc;
    check(tablet != NULL, "Table required");
    check(event != NULL, "Event required");

    // Make sure that this event is being added to the correct tablet.
    sky_tablet *target_tablet = NULL;
    rc = sky_table_get_target_tablet(tablet->table, event->object_id, &target_tablet);
    check(rc == 0, "Unable to determine target tablet");
    check(tablet == target_tablet, "Event added to invalid tablet; IDX:%d of %d, OID:%d", tablet->index, tablet->table->tablet_count, event->object_id);

    // Delegate to the data file.
    rc = sky_data_file_add_event(tablet->data_file, event);
    check(rc == 0, "Unable to add event to data file");
    
    return 0;

error:
    return -1;
}

