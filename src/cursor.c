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
        if(cursor->paths) free(cursor->paths);
        free(cursor);
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
    check(cursor != NULL, "Cursor required");

    // If data is not null then create an array of one pointer.
    void **ptrs;
    if(ptr != NULL) {
        ptrs = malloc(sizeof(void*)); check_mem(ptrs);
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
    if(ptrs) free(ptrs);
    return -1;
}

// Assigns a list of path pointers to the cursor.
// 
// cursor - The cursor.
// ptrs   - An array to pointers of raw paths.
//
// Returns 0 if successful, otherwise returns -1.
int sky_cursor_set_paths(sky_cursor *cursor, void **ptrs, int count)
{
    int rc;
    check(cursor != NULL, "Cursor required");
    
    // Free old path list.
    if(cursor->paths != NULL) {
        free(cursor->paths);
    }

    // Assign path data list.
    cursor->paths = ptrs;
    cursor->path_count = count;
    cursor->path_index = 0;
    cursor->event_index = 0;
    cursor->eof = (count == 0);
    
    // Position the pointer at the first path if paths are passed.
    if(count > 0) {
        rc = sky_cursor_set_ptr(cursor, cursor->paths[0]);
        check(rc == 0, "Unable to set paths");
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
    check(cursor != NULL, "Cursor required");
    check(ptr != NULL, "Pointer required");
    
    // Store position of first event and store position of end of path.
    cursor->ptr    = ptr + SKY_PATH_HEADER_LENGTH;
    cursor->endptr = ptr + sky_path_sizeof_raw(ptr);
    
    return 0;

error:
    cursor->ptr = NULL;
    cursor->endptr = NULL;
    return -1;
}


//--------------------------------------
// Iteration
//--------------------------------------

int sky_cursor_next(sky_cursor *cursor)
{
    int rc;
    check(cursor != NULL, "Cursor required");
    check(!cursor->eof, "No more events are available");

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

    // Make sure that we are point at an event.
    if(!cursor->eof) {
        sky_event_flag_t flag = *((sky_event_flag_t*)cursor->ptr);
        check(flag & SKY_EVENT_FLAG_ACTION || flag & SKY_EVENT_FLAG_DATA, "Cursor pointing at invalid raw event data: %p", cursor->ptr);
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
    check(cursor != NULL, "Cursor required");
    
    cursor->path_index  = 0;
    cursor->event_index = 0;
    cursor->eof         = true;
    cursor->ptr         = NULL;
    cursor->endptr      = NULL;

    return 0;

error:
    return -1;
}


//--------------------------------------
// Event Management
//--------------------------------------

// Updates a memory location based on the current event and a data descriptor.
//
// cursor    - The cursor.
// action_id - A pointer to where the action id should be returned to.
//
// Returns 0 if successful, otherwise returns -1.
int sky_cursor_set_data(sky_cursor *cursor, sky_data_descriptor *descriptor,
                        void *data)
{
    size_t sz;
    check(cursor != NULL, "Cursor required");
    check(!cursor->eof, "Cursor cannot be EOF");
    check(descriptor != NULL, "Data descriptor required");
    check(data != NULL, "Data pointer required");

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

    // Read data if this event contains data.
    uint32_t property_count = descriptor->property_count;
    if(property_count > 0 && event_flag & SKY_EVENT_FLAG_DATA) {
        uint32_t data_length = *((uint32_t*)ptr);
        void *end_ptr = ptr + data_length;
        
        // Loop over data and assign values to data object.
        while(ptr < end_ptr) {
            // Read property id.
            sky_property_id_t property_id = *((sky_property_id_t*)ptr);
            ptr += sizeof(property_id);

            // Assign value to data object member.
            sky_data_property_descriptor *property_descriptor = &descriptor->property_descriptors[property_id-descriptor->min_property_id];
            void *property_target_ptr = data + property_descriptor->offset;
            property_descriptor->setter(property_target_ptr, ptr, &sz);
            
            // If there is no size then move it forward manually.
            if(sz == 0) {
                sz = minipack_sizeof_elem_and_data(ptr);
            }
            ptr += sz;
        }
    }
    
    return 0;

error:
    return -1;
}

// Retrieves a the action identifier of the current event.
//
// cursor    - The cursor.
// action_id - A pointer to where the action id should be returned to.
//
// Returns 0 if successful, otherwise returns -1.
int sky_cursor_get_action_id(sky_cursor *cursor, sky_action_id_t *action_id)
{
    check(cursor != NULL, "Cursor required");
    check(!cursor->eof, "Cursor cannot be EOF");
    check(action_id != NULL, "Action id return pointer required");

    // Retrieve the action id.
    if(*((sky_event_flag_t*)cursor->ptr) & SKY_EVENT_FLAG_ACTION) {
        *action_id = *((sky_action_id_t*)(cursor->ptr + sizeof(sky_event_flag_t) + sizeof(sky_timestamp_t)));
    }
    else {
        *action_id = 0;
    }
    
    return 0;

error:
    *action_id = 0;
    return -1;
}

// Retrieves the pointer to where the data section of an event starts as
// well of the length of the data.
//
// cursor      - The cursor.
// data_ptr    - A pointer to where the memory location of the data starts.
// data_length - The length of the data section.
//
// Returns 0 if successful, otherwise returns -1.
int sky_cursor_get_data_ptr(sky_cursor *cursor, void **data_ptr,
                            uint32_t *data_length)
{
    check(cursor != NULL, "Cursor required");
    check(!cursor->eof, "Cursor cannot be EOF");
    check(data_ptr != NULL, "Data return pointer required");
    check(data_length != NULL, "Data length return pointer required");

    // Retrieve the data section if this event has data.
    if(*((sky_event_flag_t*)cursor->ptr) & SKY_EVENT_FLAG_DATA) {
        // Move past the initial header (flag, timestamp, action_id).
        void *ptr = cursor->ptr;
        ptr += sizeof(sky_event_flag_t) + sizeof(sky_timestamp_t);
        if(*((sky_event_flag_t*)cursor->ptr) & SKY_EVENT_FLAG_ACTION) {
            ptr += sizeof(sky_action_id_t);
        }

        // At the end of the header is the data length.
        *data_length = *((sky_event_data_length_t*)(ptr));
        
        // After the data length begins the data section.
        *data_ptr = ptr + sizeof(sky_event_data_length_t);
    }
    // If this is an action-only event then return null.
    else {
        *data_ptr = NULL;
        *data_length = 0;
    }
    
    return 0;

error:
    *data_ptr = NULL;
    *data_length = 0;
    return -1;
}

