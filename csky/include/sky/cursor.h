#ifndef _sky_cursor_h
#define _sky_cursor_h

#include <inttypes.h>
#include <stdbool.h>
#include "sky/data_descriptor.h"

//==============================================================================
//
// Constants
//
//==============================================================================

#define sky_event_flag_t uint8_t
#define EVENT_FLAG       0x92


//==============================================================================
//
// Typedefs
//
//==============================================================================

typedef struct sky_cursor {
    void *data;
    int32_t session_event_index;
    void *startptr;
    void *nextptr;
    void *endptr;
    void *ptr;
    bool eof;
    bool in_session;
    uint32_t last_timestamp;
    uint32_t session_idle_in_sec;
    sky_data_descriptor *data_descriptor;
} sky_cursor;


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

sky_cursor *sky_cursor_new();

void sky_cursor_free(sky_cursor *cursor);


//--------------------------------------
// Data Management
//--------------------------------------

void sky_cursor_set_ptr(sky_cursor *cursor, void *ptr, size_t sz);

void sky_cursor_next_event(sky_cursor *cursor);

bool sky_lua_cursor_next_event(sky_cursor *cursor);

bool sky_cursor_eof(sky_cursor *cursor);

bool sky_cursor_eos(sky_cursor *cursor);

void sky_cursor_set_session_idle(sky_cursor *cursor, uint32_t seconds);

void sky_cursor_next_session(sky_cursor *cursor);

bool sky_lua_cursor_next_session(sky_cursor *cursor);

void sky_cursor_clear_data(sky_cursor *cursor);

#endif