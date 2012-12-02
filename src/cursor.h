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
    uint32_t event_index;
    void *ptr;
    void *endptr;
    bool eof;
    void *data;
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

int sky_cursor_next(sky_cursor *cursor);

bool sky_cursor_eof(sky_cursor *cursor);

//--------------------------------------
// Event Management
//--------------------------------------

int sky_cursor_set_data(sky_cursor *cursor);

int sky_cursor_clear_data(sky_cursor *cursor);

#endif
