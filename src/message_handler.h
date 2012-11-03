#ifndef _sky_message_handler_h
#define _sky_message_handler_h

#include "bstring.h"


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
} sky_server_message_scope_e;

// A container for a specific type of message.
typedef struct {
    bstring name;
    sky_message_handler_scope_e scope;
    sky_message_handler_unpack_func_t unpack_func;
    sky_message_handler_process_func_t process_func;
} sky_message_handler;



//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

sky_server *sky_message_handler_create();

void sky_message_handler_free(sky_message_handler *handler);

#endif
