#ifndef _sky_message_handler_h
#define _sky_message_handler_h

typedef struct sky_message_handler sky_message_handler;

#include "bstring.h"
#include "message_header.h"
#include "server.h"
#include "table.h"
#include "types.h"


//==============================================================================
//
// Typedefs
//
//==============================================================================

// The types of scope that the message can have.
//
// server - Anything that requires serialized access to top-level objects such
//          as tables, actions or properties.
//
// table  - Anything that requires access to the full table. This will
//          distribute the message across all tablets.
//
// object - Anything that requires access to a single object. This will
//          only send the object to the tablet containing the object.
typedef enum {
    SKY_MESSAGE_HANDLER_SCOPE_SERVER = 1,
    SKY_MESSAGE_HANDLER_SCOPE_TABLE  = 2,
    SKY_MESSAGE_HANDLER_SCOPE_OBJECT = 3
} sky_message_handler_scope_e;

// Defines a function that processes an input and output stream on a given
// table.
typedef int (*sky_message_handler_process_func_t)(sky_server *server, sky_message_header *header, sky_table *table, FILE *input, FILE *output);

// A container for a specific type of message.
struct sky_message_handler {
    bstring name;
    sky_message_handler_scope_e scope;
    sky_message_handler_process_func_t process;
};



//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

sky_message_handler *sky_message_handler_create();

void sky_message_handler_free(sky_message_handler *handler);

#endif
