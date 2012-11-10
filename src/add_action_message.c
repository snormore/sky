#include <stdlib.h>
#include <stdio.h>
#include <arpa/inet.h>

#include "types.h"
#include "add_action_message.h"
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

// Creates an 'add_action' message object.
//
// Returns a new message.
sky_add_action_message *sky_add_action_message_create()
{
    sky_add_action_message *message = NULL;
    message = calloc(1, sizeof(sky_add_action_message)); check_mem(message);
    message->action = sky_action_create(); check_mem(message->action);
    return message;

error:
    sky_add_action_message_free(message);
    return NULL;
}

// Frees a 'add_action' message object from memory.
//
// message - The message object to be freed.
//
// Returns nothing.
void sky_add_action_message_free(sky_add_action_message *message)
{
    if(message) {
        sky_add_action_message_free_action(message);
        free(message);
    }
}

// Frees an action associated with an 'add_action' message.
//
// message - The message object to be freed.
//
// Returns nothing.
void sky_add_action_message_free_action(sky_add_action_message *message)
{
    if(message) {
        // Only free the action if it's not managed by an action file.
        if(message->action && message->action->action_file == NULL) {
            sky_action_free(message->action);
        }
        message->action = NULL;
    }
}


//--------------------------------------
// Message Handler
//--------------------------------------

// Creates a message handler for the 'Add Action' message.
//
// Returns a message handler.
sky_message_handler *sky_add_action_message_handler_create()
{
    sky_message_handler *handler = sky_message_handler_create(); check_mem(handler);
    handler->scope = SKY_MESSAGE_HANDLER_SCOPE_SERVER;
    handler->name = bfromcstr("add_action");
    handler->process = sky_add_action_message_process;
    return handler;

error:
    sky_message_handler_free(handler);
    return NULL;
}

// Delegates processing of the 'Add Event' message to a worker.
//
// server - The server.
// header - The message header.
// table  - The table the message is working against
// input  - The input file stream.
// output - The output file stream.
//
// Returns 0 if successful, otherwise returns -1.
int sky_add_action_message_process(sky_server *server,
                                   sky_message_header *header,
                                   sky_table *table, FILE *input, FILE *output)
{
    int rc = 0;
    size_t sz;
    sky_add_action_message *message = NULL;
    check(server != NULL, "Server required");
    check(header != NULL, "Message header required");
    check(table != NULL, "Table required");
    check(input != NULL, "Input stream required");
    check(output != NULL, "Output stream required");
    
    struct tagbstring status_str = bsStatic("status");
    struct tagbstring ok_str = bsStatic("ok");
    struct tagbstring action_str = bsStatic("action");

    // Parse message.
    message = sky_add_action_message_create(); check_mem(message);
    rc = sky_add_action_message_unpack(message, input);
    check(rc == 0, "Unable to parse 'add_action' message");

    // Add action.
    rc = sky_action_file_add_action(table->action_file, message->action);
    check(rc == 0, "Unable to add action");
    
    // Save actions file.
    rc = sky_action_file_save(table->action_file);
    check(rc == 0, "Unable to save action file");
    
    // Return.
    //   {status:"OK", action:{...}}
    minipack_fwrite_map(output, 2, &sz);
    check(sz > 0, "Unable to write output");
    check(sky_minipack_fwrite_bstring(output, &status_str) == 0, "Unable to write status key");
    check(sky_minipack_fwrite_bstring(output, &ok_str) == 0, "Unable to write status value");
    check(sky_minipack_fwrite_bstring(output, &action_str) == 0, "Unable to write action key");
    check(sky_action_pack(message->action, output) == 0, "Unable to write action value");

    // Clean up.
    fclose(input);
    fclose(output);
    sky_add_action_message_free(message);

    return 0;

error:
    if(input) fclose(input);
    if(output) fclose(output);
    sky_add_action_message_free(message);
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
size_t sky_add_action_message_sizeof(sky_add_action_message *message)
{
    size_t sz = 0;
    sz += sky_action_sizeof(message->action);
    return sz;
}

// Serializes an 'add_action' message to a file stream.
//
// message - The message.
// file    - The file stream to write to.
//
// Returns 0 if successful, otherwise returns -1.
int sky_add_action_message_pack(sky_add_action_message *message, FILE *file)
{
    int rc;
    check(message != NULL, "Message required");
    check(file != NULL, "File stream required");

    rc = sky_action_pack(message->action, file);
    check(rc == 0, "Unable to pack action");
    
    return 0;

error:
    return -1;
}

// Deserializes an 'add_action' message from a file stream.
//
// message - The message.
// file    - The file stream to read from.
//
// Returns 0 if successful, otherwise returns -1.
int sky_add_action_message_unpack(sky_add_action_message *message, FILE *file)
{
    int rc;
    check(message != NULL, "Message required");
    check(file != NULL, "File stream required");

    rc = sky_action_unpack(message->action, file);
    check(rc == 0, "Unable to unpack action");

    return 0;

error:
    return -1;
}
