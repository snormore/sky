#include <stdlib.h>
#include <stdio.h>
#include <arpa/inet.h>

#include "types.h"
#include "get_action_message.h"
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

// Creates a 'get_action' message object.
//
// Returns a new message.
sky_get_action_message *sky_get_action_message_create()
{
    sky_get_action_message *message = NULL;
    message = calloc(1, sizeof(sky_get_action_message)); check_mem(message);
    return message;

error:
    sky_get_action_message_free(message);
    return NULL;
}

// Frees a 'get_action' message object from memory.
//
// message - The message object to be freed.
//
// Returns nothing.
void sky_get_action_message_free(sky_get_action_message *message)
{
    if(message) {
        free(message);
    }
}


//--------------------------------------
// Message Handler
//--------------------------------------

// Creates a message handler for the 'Get Action' message.
//
// Returns a message handler.
sky_message_handler *sky_get_action_message_handler_create()
{
    sky_message_handler *handler = sky_message_handler_create(); check_mem(handler);
    handler->scope = SKY_MESSAGE_HANDLER_SCOPE_SERVER;
    handler->name = bfromcstr("get_action");
    handler->process = sky_get_action_message_process;
    return handler;

error:
    sky_message_handler_free(handler);
    return NULL;
}

// Retrieves an action from a table. This function is synchronous and does
// not use a worker.
//
// server - The server.
// header - The message header.
// table  - The table the message is working against
// input  - The input file stream.
// output - The output file stream.
//
// Returns 0 if successful, otherwise returns -1.
int sky_get_action_message_process(sky_server *server,
                                   sky_message_header *header,
                                   sky_table *table, FILE *input, FILE *output)
{
    int rc = 0;
    size_t sz;
    sky_get_action_message *message = NULL;
    check(server != NULL, "Server required");
    check(header != NULL, "Message header required");
    check(table != NULL, "Table required");
    check(input != NULL, "Input stream required");
    check(output != NULL, "Output stream required");
    
    struct tagbstring status_str = bsStatic("status");
    struct tagbstring ok_str = bsStatic("ok");
    struct tagbstring action_str = bsStatic("action");

    // Parse message.
    message = sky_get_action_message_create(); check_mem(message);
    rc = sky_get_action_message_unpack(message, input);
    check(rc == 0, "Unable to parse 'get_action' message");

    // Retrieve action.
    sky_action *action = NULL;
    rc = sky_action_file_find_action_by_id(table->action_file, message->action_id, &action);
    check(rc == 0, "Unable to get action");
    
    // Return.
    //   {status:"OK", action:{...}}
    minipack_fwrite_map(output, 2, &sz);
    check(sz > 0, "Unable to write output");
    check(sky_minipack_fwrite_bstring(output, &status_str) == 0, "Unable to write status key");
    check(sky_minipack_fwrite_bstring(output, &ok_str) == 0, "Unable to write status value");
    check(sky_minipack_fwrite_bstring(output, &action_str) == 0, "Unable to write action key");
    
    if(action != NULL) {
        check(sky_action_pack(action, output) == 0, "Unable to write action value");
    }
    else {
        minipack_fwrite_nil(output, &sz);
        check(sz > 0, "Unable to write null action value");
    }
    
    // Clean up.
    sky_get_action_message_free(message);
    fclose(input);
    fclose(output);

    return 0;

error:
    sky_get_action_message_free(message);
    if(input) fclose(input);
    if(output) fclose(output);
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
size_t sky_get_action_message_sizeof(sky_get_action_message *message)
{
    size_t sz = 0;
    sz += minipack_sizeof_uint(message->action_id);
    return sz;
}

// Serializes an 'get_action' message to a file stream.
//
// message - The message.
// file    - The file stream to write to.
//
// Returns 0 if successful, otherwise returns -1.
int sky_get_action_message_pack(sky_get_action_message *message, FILE *file)
{
    size_t sz;
    check(message != NULL, "Message required");
    check(file != NULL, "File stream required");

    minipack_fwrite_uint(file, message->action_id, &sz);
    check(sz > 0, "Unable to pack action id");
    
    return 0;

error:
    return -1;
}

// Deserializes an 'get_action' message from a file stream.
//
// message - The message.
// file    - The file stream to read from.
//
// Returns 0 if successful, otherwise returns -1.
int sky_get_action_message_unpack(sky_get_action_message *message, FILE *file)
{
    size_t sz;
    check(message != NULL, "Message required");
    check(file != NULL, "File stream required");

    message->action_id = minipack_fread_uint(file, &sz);
    check(sz > 0, "Unable to unpack action id");

    return 0;

error:
    return -1;
}

