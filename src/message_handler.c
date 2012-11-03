#include <inttypes.h>
#include <stdbool.h>
#include <stdlib.h>

#include "message_handler.h"
#include "dbg.h"


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

// Creates a message handler.
//
// Returns a message handler.
sky_message_handler *sky_message_handler_create()
{
    sky_message_handler *handler = NULL;
    handler = calloc(1, sizeof(sky_message_handler)); check_mem(handler);
    return handler;

error:
    sky_message_handler_free(handler);
    return NULL;
}

// Frees a message handler.
//
// handler - The message handler.
//
// Returns nothing.
void sky_message_handler_free(sky_message_handler *handler)
{
    if(handler) {
        if(handler->name) bdestroy(handler->name);
        handler->name = NULL;
        free(handler);
    }
}

