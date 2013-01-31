#include <assert.h>

#include "cursor.h"
#include "event.h"
#include "minipack.h"
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
        free(cursor);
    }
}



//--------------------------------------
// Pointer Management
//--------------------------------------

// Initializes the cursor to point at a new path at a given pointer.
//
// cursor - The cursor to update.
// ptr    - The address of the start of a path.
// sz     - The length of the path data, in bytes.
//
// Returns 0 if successful, otherwise returns -1.
int sky_cursor_set_ptr(sky_cursor *cursor, void *ptr, size_t sz)
{
    int rc;
    assert(cursor != NULL);
    assert(ptr != NULL);
    
    // Set the start of the path and the length of the data.
    cursor->startptr   = ptr;
    cursor->endptr     = ptr + sz;
    cursor->ptr        = NULL;
    cursor->in_session = true;
    cursor->last_timestamp      = 0;
    cursor->session_idle_in_sec = 0;
    cursor->session_event_index = -1;
    cursor->eof        = !(ptr != NULL && cursor->startptr < cursor->endptr);
    
    // Clear the data object if set.
    rc = sky_cursor_clear_data(cursor);
    check(rc == 0, "Unable to clear data");

    return 0;

error:
    return -1;
}


//--------------------------------------
// Iteration
//--------------------------------------

// Moves the cursor to the next event in a path.
//
// cursor - The cursor.
//
// Returns 0 if successful, otherwise returns -1.
int sky_cursor_next_event(sky_cursor *cursor)
{
    int rc;
    size_t event_length = 0;
    assert(cursor != NULL);

    // Ignore any calls when the cursor is out of session or EOF.
    if(cursor->eof || !cursor->in_session) {
        return 0;
    }

    // If the cursor hasn't started then initialize it.
    if(cursor->ptr == NULL) {
        cursor->ptr = cursor->startptr;
    }
    // Otherwise move to next event.
    else {
        event_length = sky_event_sizeof_raw(cursor->ptr);
        cursor->ptr += event_length;
    }

    // If pointer is beyond the last event then set eof.
    if(cursor->ptr >= cursor->endptr) {
        cursor->eof        = true;
        cursor->in_session = false;
        cursor->ptr        = NULL;
        cursor->startptr   = NULL;
        cursor->endptr     = NULL;
    }
    // Otherwise update the event object with data.
    else {
        sky_event_flag_t flag = *((sky_event_flag_t*)cursor->ptr);
        check(flag & SKY_EVENT_FLAG_ACTION || flag & SKY_EVENT_FLAG_DATA, "Cursor pointing at invalid raw event data: %p", cursor->ptr);

        // Retrieve current timestamp.
        sky_timestamp_t ts;
        uint32_t timestamp;
        rc = sky_cursor_get_timestamp(cursor, &ts, &timestamp);
        check(rc == 0, "Unable to retrieve current event timestamp");

        // Check for session boundry. This only applies if this is not the
        // first event in the session and a session idle time has been set.
        if(cursor->last_timestamp > 0 && cursor->session_idle_in_sec > 0) {
            // If the elapsed time is greater than the idle time then rewind
            // back to the event we started on at the beginning of the function
            // and mark the cursor as being "out of session".
            if(timestamp - cursor->last_timestamp >= cursor->session_idle_in_sec) {
                cursor->ptr -= event_length;
                cursor->in_session = false;
            }
        }
        cursor->last_timestamp = timestamp;

        // Only process the event if we're still in session.
        if(cursor->in_session) {
            cursor->session_event_index++;
            
            // Update data if it is available.
            if(cursor->data != NULL && cursor->data_descriptor != NULL) {
                rc = sky_cursor_set_data(cursor);
                check(rc == 0, "Unable to set set data on cursor");
            }
        }
    }

    return 0;

error:
    return -1;
}

// Moves the cursor to the next event in a path and returns a flag stating
// if the cursor is still valid (a.k.a. not EOF).
//
// cursor - The cursor.
//
// Returns true if still valid, otherwise returns false.
bool sky_lua_cursor_next_event(sky_cursor *cursor)
{
    assert(cursor != NULL);
    sky_cursor_next_event(cursor);
    return (!cursor->eof && cursor->in_session);
}

// Returns whether the cursor is at the end of the entire path.
//
// cursor - The cursor.
//
// Returns true if at the end, otherwise returns false.
bool sky_cursor_eof(sky_cursor *cursor)
{
    assert(cursor != NULL);
    return cursor->eof;
}

// Returns whether the cursor is at the end of the current session.
//
// cursor - The cursor.
//
// Returns true if at the end, otherwise returns false.
bool sky_cursor_eos(sky_cursor *cursor)
{
    assert(cursor != NULL);
    return !cursor->in_session;
}


//--------------------------------------
// Session Management
//--------------------------------------

// Sets the minimum number of seconds to elapse between two events to make
// a session demarcation.
//
// cursor  - The cursor.
// seconds - The idle time, in seconds.
//
// Returns 0 if successful, otherwise returns -1.
int sky_cursor_set_session_idle(sky_cursor *cursor, uint32_t seconds)
{
    assert(cursor != NULL);
    
    // Save the idle value.
    cursor->session_idle_in_sec = seconds;

    // If the value is non-zero then start sessionizing the cursor.
    cursor->in_session = (seconds > 0 ? false : !cursor->eof);
    
    return 0;
}

// Initializes the next session in the cursor unless the cursor has reached
// EOF.
//
// cursor - The cursor.
//
// Returns 0 if successful, otherwise returns -1.
int sky_cursor_next_session(sky_cursor *cursor)
{
    assert(cursor != NULL);
    
    // Set a flag to allow the cursor to continue iterating unless EOF is set.
    if(!cursor->in_session) {
        cursor->session_event_index = -1;
        cursor->in_session = !cursor->eof;
    }
    
    return 0;
}

// Moves the cursor to the next session in a path and returns a flag stating
// if the cursor is still valid (a.k.a. not EOF).
//
// cursor - The cursor.
//
// Returns true if still valid, otherwise returns false.
bool sky_lua_cursor_next_session(sky_cursor *cursor)
{
    assert(cursor != NULL);
    sky_cursor_next_session(cursor);
    return !cursor->eof;
}


//--------------------------------------
// Data Management
//--------------------------------------

// Retrieves the current unix timestamp & Sky timestamp from the current event.
//
// cursor    - The cursor.
// timestamp - A pointer to where the Unix timestamp should be returned.
// ts        - A pointer to where the Sky timestamp should be returned.
//
// Returns 0 if successful, otherwise returns -1.
int sky_cursor_get_timestamp(sky_cursor *cursor, sky_timestamp_t *ts,
                             uint32_t *timestamp)
{
    assert(cursor != NULL);
    assert(ts != NULL);
    assert(timestamp != NULL);
    
    sky_timestamp_t ts_value = *((sky_timestamp_t*)(cursor->ptr+sizeof(sky_event_flag_t)));
    *ts = ts_value;
    *timestamp = (uint32_t)sky_timestamp_to_seconds(ts_value);
    
    return 0;
}


// Updates a memory location based on the current event and a data descriptor.
//
// cursor    - The cursor.
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
    sky_timestamp_t *ts = (sky_timestamp_t*)(data + descriptor->timestamp_descriptor.ts_offset);
    uint32_t *timestamp = (uint32_t*)(data + descriptor->timestamp_descriptor.timestamp_offset);
    rc = sky_cursor_get_timestamp(cursor, ts, timestamp);
    check(rc == 0, "Unable to retrieve cursor timestamp");
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

// Clears the data object.
//
// cursor    - The cursor.
//
// Returns 0 if successful, otherwise returns -1.
int sky_cursor_clear_data(sky_cursor *cursor)
{
    assert(cursor != NULL);

    if(cursor->data != NULL && cursor->data_descriptor != NULL && cursor->data_descriptor->data_sz > 0) {
        memset(cursor->data, 0, cursor->data_descriptor->data_sz);
    }
    
    return 0;
}
