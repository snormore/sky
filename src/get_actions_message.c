#include <stdlib.h>
#include <stdio.h>
#include <arpa/inet.h>

#include "types.h"
#include "get_actions_message.h"
#include "action.h"
#include "minipack.h"
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

// Creates a 'get_actions' message object.
//
// Returns a new message.
sky_get_actions_message *sky_get_actions_message_create()
{
    sky_get_actions_message *message = NULL;
    message = calloc(1, sizeof(sky_get_actions_message)); check_mem(message);
    return message;

error:
    sky_get_actions_message_free(message);
    return NULL;
}

// Frees a 'get_actions' message object from memory.
//
// message - The message object to be freed.
//
// Returns nothing.
void sky_get_actions_message_free(sky_get_actions_message *message)
{
    if(message) {
        free(message);
    }
}


//--------------------------------------
// Message Handler
//--------------------------------------

// Creates a message handler for the 'Get Actions' message.
//
// Returns a message handler.
sky_message_handler *sky_get_actions_message_handler_create()
{
    sky_message_handler *handler = sky_message_handler_create(); check_mem(handler);
    handler->scope = SKY_MESSAGE_HANDLER_SCOPE_TABLE;
    handler->name = bfromcstr("get_actions");
    handler->process = sky_get_actions_message_process;
    return handler;

error:
    sky_message_handler_free(handler);
    return NULL;
}

// Retrieves a list of actions from the server for a table. This function is
// synchronous and does not use a worker.
//
// server - The server.
// header - The message header.
// table  - The table the message is working against
// input  - The input file stream.
// output - The output file stream.
//
// Returns 0 if successful, otherwise returns -1.
int sky_get_actions_message_process(sky_server *server,
                                    sky_message_header *header,
                                    sky_table *table, FILE *input, FILE *output)
{
    int rc = 0;
    size_t sz;
    sky_get_actions_message *message = NULL;
    check(server != NULL, "Server required");
    check(header != NULL, "Message header required");
    check(table != NULL, "Table required");
    check(input != NULL, "Input stream required");
    check(output != NULL, "Output stream required");
    
    struct tagbstring status_str = bsStatic("status");
    struct tagbstring ok_str = bsStatic("ok");
    struct tagbstring actions_str = bsStatic("actions");

    // Parse message.
    message = sky_get_actions_message_create(); check_mem(message);
    rc = sky_get_actions_message_unpack(message, input);
    check(rc == 0, "Unable to parse 'get_actions' message");

    // Return.
    //   {status:"OK", actions:[{...}]}
    minipack_fwrite_map(output, 2, &sz);
    check(sz > 0, "Unable to write output");
    check(sky_minipack_fwrite_bstring(output, &status_str) == 0, "Unable to write status key");
    check(sky_minipack_fwrite_bstring(output, &ok_str) == 0, "Unable to write status value");
    check(sky_minipack_fwrite_bstring(output, &actions_str) == 0, "Unable to write actions key");

    // Loop over actions and serialize them.
    minipack_fwrite_array(output, table->action_file->action_count, &sz);
    check(sz > 0, "Unable to write actions array");
    
    uint32_t i;
    for(i=0; i<table->action_file->action_count; i++) {
        sky_action *action = table->action_file->actions[i];
        check(sky_action_pack(action, output) == 0, "Unable to write action");
    }

    // Clean up.
    sky_get_actions_message_free(message);
    fclose(input);
    fclose(output);

    return 0;

error:
    if(input) fclose(input);
    if(output) fclose(output);
    sky_get_actions_message_free(message);
    return -1;
}


//--------------------------------------
// Serialization
//--------------------------------------

// Calculates the total number of bytes needed to store the message.
//
// message - The message.
//
// Returns the number of bytes required to store the message.
size_t sky_get_actions_message_sizeof(sky_get_actions_message *message)
{
    size_t sz = 0;
    check(message != NULL, "Message required");
    return sz;

error:
    return 0;
}

// Serializes a 'get_actions' message to a file stream.
//
// message - The message.
// file    - The file stream to write to.
//
// Returns 0 if successful, otherwise returns -1.
int sky_get_actions_message_pack(sky_get_actions_message *message, FILE *file)
{
    check(message != NULL, "Message required");
    check(file != NULL, "File stream required");

    return 0;

error:
    return -1;
}

// Deserializes a 'get_actions' message from a file stream.
//
// message - The message.
// file    - The file stream to read from.
//
// Returns 0 if successful, otherwise returns -1.
int sky_get_actions_message_unpack(sky_get_actions_message *message, FILE *file)
{
    check(message != NULL, "Message required");
    check(file != NULL, "File stream required");

    return 0;

error:
    return -1;
}


