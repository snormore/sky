#include <assert.h>

#include "cursor.h"
#include "path.h"
#include "event.h"
#include "minipack.h"
#include "mem.h"
#include "dbg.h"


//==============================================================================
//
// Forward Declarations
//
//==============================================================================

int sky_cursor_set_ptr(sky_cursor *cursor, void *ptr);
int sky_cursor_set_eof(sky_cursor *cursor);


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

// Creates a reference to a cursor.
// 
// Returns a reference to the new cursor if successful.
sky_cursor *sky_cursor_create()
{
    sky_cursor *cursor = calloc(sizeof(sky_cursor),1 ); check_mem(cursor);
    return cursor;
    
error:
    sky_cursor_free(cursor);
    return NULL;
}

// Allocates memory for the cursor.
// 
// Returns a reference to the new cursor if successful.
sky_cursor *sky_cursor_alloc()
{
    return malloc(sizeof(sky_cursor));
}

// Initializes a cursor.
void sky_cursor_init(sky_cursor *cursor)
{
    memset(cursor, 0, sizeof(sky_cursor));
}

// Removes a cursor reference from memory.
//
// cursor - The cursor to free.
void sky_cursor_free(sky_cursor *cursor)
{
    if(cursor) {
        sky_cursor_uninit(cursor);
        free(cursor);
    }
}

// Deallocates child members.
//
// cursor - The cursor.
void sky_cursor_uninit(sky_cursor *cursor)
{
    if(cursor) {
        if(cursor->paths) free(cursor->paths);
        cursor->paths = NULL;
        cursor->path_bcount = 0;
    }
}


//--------------------------------------
// Path Management
//--------------------------------------

// Assigns a single path to the cursor.
// 
// cursor - The cursor.
// ptr    - A pointer to the raw data where the path starts.
//
// Returns 0 if successful, otherwise returns -1.
int sky_cursor_set_path(sky_cursor *cursor, void *ptr)
{
    int rc;
    assert(cursor != NULL);

    // If data is not null then create an array of one pointer.
    if(ptr != NULL) {
        void *ptrs[1];
        ptrs[0] = ptr;
        rc = sky_cursor_set_paths(cursor, ptrs, 1);
        check(rc == 0, "Unable to set path data to cursor");
    }
    // Otherwise clear out pointer paths.
    else {
        rc = sky_cursor_set_paths(cursor, NULL, 0);
        check(rc == 0, "Unable to remove path data");
    }

    return 0;

error:
    return -1;
}

// Assigns a list of path pointers to the cursor.
// 
// cursor - The cursor.
// ptrs   - An array to pointers of raw paths.
// count  - The number of paths in the array.
//
// Returns 0 if successful, otherwise returns -1.
int sky_cursor_set_paths(sky_cursor *cursor, void **ptrs, uint32_t count)
{
    int rc;
    assert(cursor != NULL);
    
    // Reallocate if the paths buffer is too small.
    if(count > cursor->path_bcount) {
        cursor->path_bcount = count;
        cursor->paths = realloc(cursor->paths, sizeof(*cursor->paths) * cursor->path_bcount);
        check_mem(cursor->paths);
    }
    
    // Clear cursor data.
    if(cursor->data != NULL && cursor->data_sz > 0) {
        memset(cursor->data, 0, cursor->data_sz);
    }
    
    // Assign path data list.
    uint32_t i;
    for(i=0; i<count; i++) {
        cursor->paths[i] = ptrs[i];
    }
    cursor->path_count = count;
    cursor->path_index = 0;
    cursor->event_index = 0;
    cursor->eof = (count == 0);
    
    // Position the pointer at the first path if paths are passed.
    if(count > 0) {
        rc = sky_cursor_set_ptr(cursor, cursor->paths[0]);
        check(rc == 0, "Unable to set paths");

        // Update the data on the cursor's data record.
        if(cursor->data_descriptor != NULL && cursor->data != NULL) {
            sky_cursor_set_data(cursor);
        }
    }
    // Otherwise clear out the pointer.
    else {
        cursor->ptr = NULL;
        cursor->endptr = NULL;
    }
    
    return 0;

error:
    return -1;
}


//--------------------------------------
// Pointer Management
//--------------------------------------

// Initializes the cursor to point at a new path at a given pointer.
//
// cursor - The cursor to update.
// ptr    - The address of the start of a path.
//
// Returns 0 if successful, otherwise returns -1.
int sky_cursor_set_ptr(sky_cursor *cursor, void *ptr)
{
    assert(cursor != NULL);
    assert(ptr != NULL);
    
    // Store position of first event and store position of end of path.
    cursor->ptr    = ptr + SKY_PATH_HEADER_LENGTH;
    cursor->endptr = ptr + sky_path_sizeof_raw(ptr);
    
    return 0;
}


//--------------------------------------
// Iteration
//--------------------------------------

int sky_cursor_next(sky_cursor *cursor)
{
    int rc;
    assert(cursor != NULL);
    assert(!cursor->eof);

    // Move to next event.
    size_t event_length = sky_event_sizeof_raw(cursor->ptr);
    cursor->ptr += event_length;
    cursor->event_index++;

    // If pointer is beyond the last event then move to next path.
    if(cursor->ptr >= cursor->endptr) {
        cursor->path_index++;

        // Move to the next path if more paths are remaining.
        if(cursor->path_index < cursor->path_count) {
            rc = sky_cursor_set_ptr(cursor, cursor->paths[cursor->path_index]);
            check(rc == 0, "Unable to set pointer to path");
        }
        // Otherwise set EOF.
        else {
            rc = sky_cursor_set_eof(cursor);
            check(rc == 0, "Unable to set EOF on cursor");
        }
    }

    // Validate and update data.
    if(!cursor->eof) {
        sky_event_flag_t flag = *((sky_event_flag_t*)cursor->ptr);
        check(flag & SKY_EVENT_FLAG_ACTION || flag & SKY_EVENT_FLAG_DATA, "Cursor pointing at invalid raw event data: %p", cursor->ptr);

        if(cursor->data_descriptor != NULL && cursor->data != NULL) {
            sky_cursor_set_data(cursor);
        }
    }

    return 0;

error:
    return -1;
}

// Flags a cursor to say that it is at the end of all its paths.
//
// cursor - The cursor to set EOF on.
//
// Returns 0 if successful, otherwise returns -1.
int sky_cursor_set_eof(sky_cursor *cursor)
{
    assert(cursor != NULL);
    
    cursor->path_index  = 0;
    cursor->event_index = 0;
    cursor->eof         = true;
    cursor->ptr         = NULL;
    cursor->endptr      = NULL;

    return 0;
}


//--------------------------------------
// Event Management
//--------------------------------------

// Updates a memory location based on the current event and a data descriptor.
//
// cursor - The cursor.
//
// Returns 0 if successful, otherwise returns -1.
int sky_cursor_set_data(sky_cursor *cursor)
{
    size_t sz;
    int rc;
    assert(cursor != NULL);
    assert(!cursor->eof);
    assert(cursor->data_descriptor != NULL);
    assert(cursor->data != NULL);

    // Localize variables.
    sky_data_descriptor *descriptor = cursor->data_descriptor;
    void *data = cursor->data;

    // Retrieve the flag off the event.
    void *ptr = cursor->ptr;
    sky_event_flag_t event_flag = *((sky_event_flag_t*)ptr);
    ptr += sizeof(sky_event_flag_t);
    
    // Assign timestamp.
    sky_timestamp_t *timestamp = (sky_timestamp_t*)(data + descriptor->timestamp_descriptor.offset);
    *timestamp = *((sky_timestamp_t*)ptr);
    ptr += sizeof(sky_timestamp_t);

    // Read action if this event contains an action.
    sky_action_id_t *action_id = (sky_action_id_t*)(data + descriptor->action_descriptor.offset);
    if(event_flag & SKY_EVENT_FLAG_ACTION) {
        *action_id = *((sky_action_id_t*)ptr);
        ptr += sizeof(sky_action_id_t);
    }
    else {
        *action_id = 0;
    }

    // Process data descriptor if there are properties being tracked.
    if(descriptor->active_property_count > 0) {
        // Clear old action data.
        rc = sky_data_descriptor_clear_action_data(descriptor, data);
        check(rc == 0, "Unable to clear action data via descriptor");

        // Read data if this event contains data.
        if(event_flag & SKY_EVENT_FLAG_DATA) {
            uint32_t data_length = *((uint32_t*)ptr);
            ptr += sizeof(uint32_t);
            void *end_ptr = ptr + data_length;

            // Loop over data and assign values to data object.
            while(ptr < end_ptr) {
                // Read property id.
                sky_property_id_t property_id = *((sky_property_id_t*)ptr);
                ptr += sizeof(property_id);

                // Assign value to data object member.
                rc = sky_data_descriptor_set_value(descriptor, data, property_id, ptr, &sz);
                check(rc == 0, "Unable to set value via data descriptor");
            
                // If there is no size then move it forward manually.
                if(sz == 0) {
                    sz = minipack_sizeof_elem_and_data(ptr);
                }
                ptr += sz;
            }
        }
    }
    
    return 0;

error:
    return -1;
}
