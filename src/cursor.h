#ifndef _cursor_h
#define _cursor_h

#include <inttypes.h>
#include <stdbool.h>
#include <stdlib.h>

#include "bstring.h"
#include "data_descriptor.h"
#include "types.h"


//==============================================================================
//
// Typedefs
//
//==============================================================================

typedef struct sky_cursor {
    void *startptr;
    void *endptr;
    void *ptr;
    bool eof;
    void *data;
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

sky_cursor *sky_cursor_create();

sky_cursor *sky_cursor_alloc();

void sky_cursor_init(sky_cursor *cursor);

void sky_cursor_free(sky_cursor *cursor);

//--------------------------------------
// Path Management
//--------------------------------------

int sky_cursor_set_ptr(sky_cursor *cursor, void *ptr, size_t sz);

//--------------------------------------
// Iteration
//--------------------------------------

int sky_cursor_next_event(sky_cursor *cursor);

bool sky_lua_cursor_next_event(sky_cursor *cursor);

bool sky_cursor_eof(sky_cursor *cursor);

//--------------------------------------
// Session Management
//--------------------------------------

int sky_cursor_set_session_idle(sky_cursor *cursor, uint32_t seconds);

int sky_cursor_next_session(sky_cursor *cursor);

bool sky_lua_cursor_next_session(sky_cursor *cursor);

//--------------------------------------
// Data Management
//--------------------------------------

int sky_cursor_get_timestamp(sky_cursor *cursor, sky_timestamp_t *ts,
    uint32_t *timestamp);

int sky_cursor_set_data(sky_cursor *cursor);

int sky_cursor_clear_data(sky_cursor *cursor);

#endif
