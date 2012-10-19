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


//--------------------------------------
// Processing
//--------------------------------------

// Applies an 'add_action' message to a table.
//
// message - The message.
// table   - The table to apply the message to.
// output  - The output stream to write to.
//
// Returns 0 if successful, otherwise returns -1.
int sky_add_action_message_process(sky_add_action_message *message,
                                   sky_table *table, FILE *output)
{
    int rc;
    size_t sz;
    check(message != NULL, "Message required");
    check(table != NULL, "Table required");
    check(output != NULL, "Output stream required");

    struct tagbstring status_str = bsStatic("status");
    struct tagbstring ok_str = bsStatic("ok");
    struct tagbstring action_str = bsStatic("action");

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
    
    return 0;

error:
    return -1;
}