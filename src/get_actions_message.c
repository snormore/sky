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


//--------------------------------------
// Processing
//--------------------------------------

// Applies a 'get_actions' message to a table.
//
// message - The message.
// table   - The table to apply the message to.
// output  - The output stream to write to.
//
// Returns 0 if successful, otherwise returns -1.
int sky_get_actions_message_process(sky_get_actions_message *message,
                                    sky_table *table, FILE *output)
{
    size_t sz;
    check(message != NULL, "Message required");
    check(table != NULL, "Table required");
    check(output != NULL, "Output stream required");

    struct tagbstring status_str = bsStatic("status");
    struct tagbstring ok_str = bsStatic("ok");
    struct tagbstring actions_str = bsStatic("actions");

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
    
    return 0;

error:
    return -1;
}