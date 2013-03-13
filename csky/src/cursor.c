#include <stdio.h>
#include <stdlib.h>
#include "sky/cursor.h"
#include "sky/mem.h"
#include "sky/timestamp.h"
#include "sky/minipack.h"
#include "sky/dbg.h"

//==============================================================================
//
// Macros
//
//==============================================================================

#define badcursordata(MSG, PTR) do {\
    fprintf(stderr, "Cursor pointing at invalid raw event data [" MSG "]: %p->%p\n", cursor->ptr, PTR); \
    memdump(cursor->startptr, (cursor->endptr - cursor->startptr)); \
    cursor->eof = true; \
    return; \
} while(0)


//==============================================================================
//
// Cursor
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

// Creates a reference to a cursor.
sky_cursor *sky_cursor_new()
{
    sky_cursor *cursor = calloc(sizeof(sky_cursor),1 );
    return cursor;
}

// Removes a cursor reference from memory.
void sky_cursor_free(sky_cursor *cursor)
{
    if(cursor) {
        free(cursor);
    }
}


//--------------------------------------
// Data Management
//--------------------------------------

void sky_cursor_set_ptr(sky_cursor *cursor, void *ptr, size_t sz)
{
    // Set the start of the path and the length of the data.
    cursor->startptr   = ptr;
    cursor->nextptr    = ptr;
    cursor->endptr     = ptr + sz;
    cursor->ptr        = NULL;
    cursor->in_session = true;
    cursor->last_timestamp      = 0;
    cursor->session_idle_in_sec = 0;
    cursor->session_event_index = -1;
    cursor->eof        = !(ptr != NULL && cursor->startptr < cursor->endptr);
    
    // Clear the data object if set.
    sky_cursor_clear_data(cursor);
}

void sky_cursor_next_event(sky_cursor *cursor)
{
    // Ignore any calls when the cursor is out of session or EOF.
    if(cursor->eof || !cursor->in_session) {
        return;
    }

    // Move the pointer to the next position.
    void *prevptr = cursor->ptr;
    cursor->ptr = cursor->nextptr;
    void *ptr = cursor->ptr;

    // If pointer is beyond the last event then set eof.
    if(cursor->ptr >= cursor->endptr) {
        cursor->eof        = true;
        cursor->in_session = false;
        cursor->ptr        = NULL;
        cursor->startptr   = NULL;
        cursor->nextptr    = NULL;
        cursor->endptr     = NULL;
    }
    // Otherwise update the event object with data.
    else {
        sky_event_flag_t flag = *((sky_event_flag_t*)ptr);
        
        // If flag isn't correct then report and exit.
        if(flag != EVENT_FLAG) badcursordata("eflag", ptr);
        ptr += sizeof(sky_event_flag_t);
        
        // Read timestamp.
        size_t sz;
        int64_t ts = minipack_unpack_int(ptr, &sz);
        if(sz == 0) badcursordata("timestamp", ptr);
        uint32_t timestamp = sky_timestamp_to_seconds(ts);
        ptr += sz;

        // Check for session boundry. This only applies if this is not the
        // first event in the session and a session idle time has been set.
        if(cursor->last_timestamp > 0 && cursor->session_idle_in_sec > 0) {
            // If the elapsed time is greater than the idle time then rewind
            // back to the event we started on at the beginning of the function
            // and mark the cursor as being "out of session".
            if(timestamp - cursor->last_timestamp >= cursor->session_idle_in_sec) {
                cursor->ptr = prevptr;
                cursor->in_session = false;
            }
        }
        cursor->last_timestamp = timestamp;

        // Only process the event if we're still in session.
        if(cursor->in_session) {
            cursor->session_event_index++;
            
            // Set timestamp.
            int64_t *data_ts = (int64_t*)(cursor->data + cursor->data_descriptor->timestamp_descriptor.ts_offset);
            uint32_t *data_timestamp = (uint32_t*)(cursor->data + cursor->data_descriptor->timestamp_descriptor.timestamp_offset);
            *data_ts = ts;
            *data_timestamp = timestamp;
            
            // Clear old action data.
            uint32_t i;
            for(i=0; i<cursor->data_descriptor->action_property_descriptor_count; i++) {
                sky_data_property_descriptor *property_descriptor = cursor->data_descriptor->action_property_descriptors[i];
                property_descriptor->clear_func(cursor->data + property_descriptor->offset);
            }

            // Read msgpack map!
            uint32_t count = minipack_unpack_map(ptr, &sz);
            if(sz == 0) badcursordata("datamap", ptr);
            ptr += sz;

            // Loop over key/value pairs.
            for(i=0; i<count; i++) {
                // Read property id (key).
                int64_t property_id = minipack_unpack_int(ptr, &sz);
                if(sz == 0) badcursordata("key", ptr);
                ptr += sz;

                // Read property value and set it on the data object.
                sky_data_descriptor_set_value(cursor->data_descriptor, cursor->data, property_id, ptr, &sz);
                if(sz == 0) {
                  sz = minipack_sizeof_elem_and_data(ptr);
                }
                ptr += sz;
            }

            cursor->nextptr = ptr;
        }
    }
}

bool sky_lua_cursor_next_event(sky_cursor *cursor)
{
    sky_cursor_next_event(cursor);
    return (!cursor->eof && cursor->in_session);
}

bool sky_cursor_eof(sky_cursor *cursor)
{
    return cursor->eof;
}

bool sky_cursor_eos(sky_cursor *cursor)
{
    return !cursor->in_session;
}

void sky_cursor_set_session_idle(sky_cursor *cursor, uint32_t seconds)
{
    // Save the idle value.
    cursor->session_idle_in_sec = seconds;

    // If the value is non-zero then start sessionizing the cursor.
    cursor->in_session = (seconds > 0 ? false : !cursor->eof);
}

void sky_cursor_next_session(sky_cursor *cursor)
{
    // Set a flag to allow the cursor to continue iterating unless EOF is set.
    if(!cursor->in_session) {
        cursor->session_event_index = -1;
        cursor->in_session = !cursor->eof;
    }
}

bool sky_lua_cursor_next_session(sky_cursor *cursor)
{
    sky_cursor_next_session(cursor);
    return !cursor->eof;
}

void sky_cursor_clear_data(sky_cursor *cursor)
{
    if(cursor->data != NULL && cursor->data_descriptor != NULL && cursor->data_descriptor->data_sz > 0) {
        memset(cursor->data, 0, cursor->data_descriptor->data_sz);
    }
}

