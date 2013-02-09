#include <stdlib.h>
#include <fcntl.h>
#include <inttypes.h>
#include <unistd.h>
#include <sys/stat.h>
#include <sys/mman.h>
#include <math.h>
#include <assert.h>

#include "tablet.h"
#include "cursor.h"
#include "sky_string.h"
#include "timestamp.h"
#include "mem.h"
#include "dbg.h"

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
    sky_tablet *tablet = NULL;
    assert(table != NULL);
    tablet = calloc(sizeof(sky_tablet), 1); check_mem(tablet);
    tablet->table = table;
    tablet->readoptions = leveldb_readoptions_create();
    tablet->writeoptions = leveldb_writeoptions_create();
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
        tablet->index = 0;

        bdestroy(tablet->path);
        tablet->path = NULL;

        if(tablet->readoptions) leveldb_readoptions_destroy(tablet->readoptions);
        tablet->readoptions = NULL;

        if(tablet->writeoptions) leveldb_writeoptions_destroy(tablet->writeoptions);
        tablet->writeoptions = NULL;

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
    assert(tablet != NULL);

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
    char* errptr = NULL;
    leveldb_options_t* options = NULL;

    assert(tablet != NULL);
    check(tablet->path != NULL, "Tablet path required");
    
    // Close the tablet if it's already open.
    sky_tablet_close(tablet);
    
    // Initialize data file.
    options = leveldb_options_create();
    leveldb_options_set_create_if_missing(options, true);
    tablet->leveldb_db = leveldb_open(options, bdata(tablet->path), &errptr);
    check(errptr == NULL, "LevelDB Error: %s", errptr);
    check(tablet->leveldb_db != NULL, "Unable to create LevelDB data file");
    leveldb_options_destroy(options);

    return 0;
error:
    if(errptr) leveldb_free(errptr);
    if(options) leveldb_options_destroy(options);
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
    assert(tablet != NULL);

    if(tablet->leveldb_db) {
        leveldb_close(tablet->leveldb_db);
        tablet->leveldb_db = NULL;
    }

    return 0;
}


//--------------------------------------
// Path Management
//--------------------------------------

// Retrieves a path for an object in the tablet.
//
// tablet      - The tablet.
// object_id   - The object identifier for the path.
// data        - A pointer to where the path data should be returned.
// data_length - A pointer to where the length of the path data should be returned.
// 
//
// Returns 0 if successful, otherwise returns -1.
int sky_tablet_get_path(sky_tablet *tablet, bstring object_id,
                        void **data, size_t *data_length)
{
    char *errptr = NULL;
    assert(tablet != NULL);

    // Retrieve the existing value.
    *data = (void*)leveldb_get(tablet->leveldb_db, tablet->readoptions, (const char*)bdata(object_id), blength(object_id), data_length, &errptr);
    check(errptr == NULL, "LevelDB get path error: %s", errptr);
    return 0;

error:
    *data = NULL;
    *data_length = 0;
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
    char *errptr = NULL;
    void *new_data = NULL;
    void *data = NULL;
    sky_data_object *data_object = NULL;
    sky_data_descriptor *descriptor = NULL;
    sky_cursor cursor; memset(&cursor, 0, sizeof(cursor));
    assert(tablet != NULL);
    assert(event != NULL);

    // Make sure that this event is being added to the correct tablet.
    sky_tablet *target_tablet = NULL;
    rc = sky_table_get_target_tablet(tablet->table, event->object_id, &target_tablet);
    check(rc == 0, "Unable to determine target tablet");
    check(tablet == target_tablet, "Event added to invalid tablet; IDX:%d of %d", tablet->index, tablet->table->tablet_count);

    // Retrieve the existing value.
    size_t data_length;
    data = (void*)leveldb_get(tablet->leveldb_db, tablet->readoptions, (const char*)bdata(event->object_id), blength(event->object_id), &data_length, &errptr);
    check(errptr == NULL, "LevelDB get error: %s", errptr);
    
    // Find the insertion point on the path.
    size_t insert_offset = 0;
    size_t event_length;
    
    // If the object doesn't exist yet then just set the single event. Easy peasy.
    if(data == NULL) {
        event_length = sky_event_sizeof(event);
        new_data = calloc(1, event_length); check_mem(new_data);
        insert_offset = 0;
    }
    // If the object does exist, we need to find where to insert the event data.
    // Also, we need to strip off any state which is redundant at the point of
    // insertion.
    else {
        void *insert_ptr = NULL;
        sky_timestamp_t event_ts = sky_timestamp_shift(event->timestamp);
        
        // Initialize data descriptor.
        descriptor = sky_data_descriptor_create(); check_mem(descriptor);
        rc = sky_data_descriptor_init_with_event(descriptor, event);
        check(rc == 0, "Unable to initialize data descriptor for event insert");
    
        // Initialize data object.
        data_object = calloc(1, descriptor->data_sz); check_mem(data);

        // Attach data & descriptor to the cursor.
        cursor.data_descriptor = descriptor;
        cursor.data = (void*)data_object;

        // Initialize the cursor.
        rc = sky_cursor_set_ptr(&cursor, data, data_length);
        check(rc == 0, "Unable to set pointer on cursor");
        check(sky_cursor_next_event(&cursor) == 0, "Unable to move to next event");
        
        // Loop over cursor until we reach the event insertion point.
        while(!cursor.eof) {
            // Retrieve event insertion pointer once the timestamp is reached.
            if(data_object->ts >= event_ts) {
                insert_ptr = cursor.ptr;
                break;
            }
            
            // Move to next event.
            check(sky_cursor_next_event(&cursor) == 0, "Unable to move to next event");
        }

        // If no insertion point was found then append the event to the
        // end of the path.
        if(insert_ptr == NULL) {
            insert_ptr = data + data_length;
        }
        insert_offset = insert_ptr - data;

        // Clear off any object data on the event that matches
        // what is the current state of the event in the database.
        uint32_t i;
        for(i=0; i<event->data_count; i++) {
            // Ignore any action properties.
            if(event->data[i]->key > 0) {
                sky_data_property_descriptor *property_descriptor = &descriptor->property_zero_descriptor[event->data[i]->key];

                // If the values match then splice this from the array.
                // Compare strings.
                void *a = &event->data[i]->value;
                void *b = ((void*)data_object)+property_descriptor->offset;
                size_t n = sky_data_type_sizeof(event->data[i]->data_type);
                
                bool is_equal = false;
                if(event->data[i]->data_type == SKY_DATA_TYPE_STRING) {
                    is_equal = sky_string_bequals((sky_string*)b, event->data[i]->string_value);
                }
                // Compare other types.
                else if(memcmp(a, b, n) == 0) {
                    is_equal = true;
                }

                // If the data is equal then remove it.
                if(is_equal) {
                    sky_event_data_free(event->data[i]);
                    if(i < event->data_count - 1) {
                        memmove(&event->data[i], &event->data[i+1], (event->data_count-i-1) * sizeof(*event->data));
                    }
                    i--;
                    event->data_count--;
                }
            }
        }

        // Determine the serialized size of the event. If the event is
        // completely redundant (e.g. it is a data-only event and the event
        // matches the current object state) then don't allocate space for a
        // new path value.
        event_length = sky_event_sizeof(event);
        if(event_length > 0) {
            // Allocate space for the existing data plus the new data.
            new_data = calloc(1, data_length + event_length); check_mem(new_data);

            // Copy in data before event.
            if(insert_offset > 0) {
                memmove(new_data, data, insert_ptr-data);
            }

            // Copy in data after event.
            if(insert_offset < data_length) {
                event_length = sky_event_sizeof(event);
                memmove(new_data+insert_offset+event_length, data+insert_offset, data_length-insert_offset);
            }
        }
    }
    
    // If no space was allocated then it means the event is redundant and
    // should be ignored.
    if(new_data != NULL) {
        // If the object doesn't exist then just set the event as the data.
        size_t event_sz;
        rc = sky_event_pack(event, new_data + insert_offset, &event_sz);
        check(rc == 0, "Unable to pack event");
        check(event_sz == event_length, "Expected event size (%ld) does not match actual event size (%ld)", event_length, event_sz);

        leveldb_put(tablet->leveldb_db, tablet->writeoptions, (const char*)bdata(event->object_id), blength(event->object_id), new_data, data_length + event_length, &errptr);
        check(errptr == NULL, "LevelDB put error: %s", errptr);
    }
    
    free(data_object);
    sky_data_descriptor_free(descriptor);
    free(data);
    free(new_data);
    
    return 0;

error:
    if(errptr) leveldb_free(errptr);
    sky_data_descriptor_free(descriptor);
    if(data) free(data);
    if(new_data) free(new_data);
    if(data_object) free(data_object);
    return -1;
}
