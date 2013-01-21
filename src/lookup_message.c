#include <stdlib.h>
#include <stdio.h>
#include <arpa/inet.h>

#include "types.h"
#include "lookup_message.h"
#include "action.h"
#include "property.h"
#include "action_file.h"
#include "property_file.h"
#include "minipack.h"
#include "mem.h"
#include "dbg.h"


//==============================================================================
//
// Definitions
//
//==============================================================================

#define SKY_LOOKUP_MESSAGE_TABLE_KEY_COUNT 2

struct tagbstring SKY_LOOKUP_MESSAGE_ACTION_NAMES_STR = bsStatic("actionNames");

struct tagbstring SKY_LOOKUP_MESSAGE_PROPERTY_NAMES_STR = bsStatic("propertyNames");


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

// Creates a 'lookup' message object.
//
// Returns a new message.
sky_lookup_message *sky_lookup_message_create()
{
    sky_lookup_message *message = NULL;
    message = calloc(1, sizeof(sky_lookup_message)); check_mem(message);
    return message;

error:
    sky_lookup_message_free(message);
    return NULL;
}

// Frees a 'lookup' message object from memory.
//
// message - The message object to be freed.
//
// Returns nothing.
void sky_lookup_message_free(sky_lookup_message *message)
{
    uint32_t i;
    if(message) {
        for(i=0; i<message->action_name_count; i++) {
            bdestroy(message->action_names[i]);
            message->action_names[i] = NULL;
        }
        free(message->action_names);

        for(i=0; i<message->property_name_count; i++) {
            bdestroy(message->property_names[i]);
            message->property_names[i] = NULL;
        }
        free(message->property_names);

        free(message);
    }
}


//--------------------------------------
// Message Handler
//--------------------------------------

// Creates a message handler for the 'lookup' message.
//
// Returns a message handler.
sky_message_handler *sky_lookup_message_handler_create()
{
    sky_message_handler *handler = sky_message_handler_create(); check_mem(handler);
    handler->scope = SKY_MESSAGE_HANDLER_SCOPE_TABLE;
    handler->name = bfromcstr("lookup");
    handler->process = sky_lookup_message_process;
    return handler;

error:
    sky_message_handler_free(handler);
    return NULL;
}

// Looks up identifiers for a list of actions and properties by name. This
// function is synchronous and does not use a worker.
//
// server - The server.
// header - The message header.
// table  - The table the message is working against
// input  - The input file stream.
// output - The output file stream.
//
// Returns 0 if successful, otherwise returns -1.
int sky_lookup_message_process(sky_server *server,
                               sky_message_header *header,
                               sky_table *table, FILE *input, FILE *output)
{
    int rc = 0;
    size_t sz;
    uint32_t i;
    sky_lookup_message *message = NULL;
    check(server != NULL, "Server required");
    check(header != NULL, "Message header required");
    check(table != NULL, "Table required");
    check(input != NULL, "Input stream required");
    check(output != NULL, "Output stream required");
    
    struct tagbstring status_str = bsStatic("status");
    struct tagbstring ok_str = bsStatic("ok");
    struct tagbstring action_ids_str = bsStatic("actions");
    struct tagbstring property_ids_str = bsStatic("properties");

    // Parse message.
    message = sky_lookup_message_create(); check_mem(message);
    rc = sky_lookup_message_unpack(message, input);
    check(rc == 0, "Unable to parse 'lookup' message");

    // Return.
    //   {status:"OK", actions:[{...}]}
    minipack_fwrite_map(output, 3, &sz);
    check(sz > 0, "Unable to write output");
    check(sky_minipack_fwrite_bstring(output, &status_str) == 0, "Unable to write status key");
    check(sky_minipack_fwrite_bstring(output, &ok_str) == 0, "Unable to write status value");

    // Loop over action names and serialize them.
    check(sky_minipack_fwrite_bstring(output, &action_ids_str) == 0, "Unable to write action_ids key");
    minipack_fwrite_array(output, message->action_name_count, &sz);
    check(sz > 0, "Unable to write action ids array");
    
    for(i=0; i<message->action_name_count; i++) {
        sky_action *action = NULL;
        rc = sky_action_file_find_by_name(table->action_file, message->action_names[i], &action);
        check(rc == 0, "Unable to search for action by name");

        if(action != NULL) {
            check(sky_action_pack(action, output) == 0, "Unable to write action value");
        }
        else {
            minipack_fwrite_nil(output, &sz);
            check(sz > 0, "Unable to write null action value");
        }
    }

    // Loop over property names and serialize them.
    check(sky_minipack_fwrite_bstring(output, &property_ids_str) == 0, "Unable to write property_ids key");
    minipack_fwrite_array(output, message->property_name_count, &sz);
    check(sz > 0, "Unable to write property ids array");
    
    for(i=0; i<message->property_name_count; i++) {
        sky_property *property = NULL;
        rc = sky_property_file_find_by_name(table->property_file, message->property_names[i], &property);
        check(rc == 0, "Unable to search for property by name");

        if(property != NULL) {
            check(sky_property_pack(property, output) == 0, "Unable to write property value");
        }
        else {
            minipack_fwrite_nil(output, &sz);
            check(sz > 0, "Unable to write null property value");
        }
    }

    // Clean up.
    sky_lookup_message_free(message);
    fclose(input);
    fclose(output);

    return 0;

error:
    if(input) fclose(input);
    if(output) fclose(output);
    sky_lookup_message_free(message);
    return -1;
}


//--------------------------------------
// Serialization
//--------------------------------------

// Serializes a 'lookup' message to a file stream.
//
// message - The message.
// file    - The file stream to write to.
//
// Returns 0 if successful, otherwise returns -1.
int sky_lookup_message_pack(sky_lookup_message *message, FILE *file)
{
    int rc;
    size_t sz;
    uint32_t i;
    check(message != NULL, "Message required");
    check(file != NULL, "File stream required");

    // Map
    minipack_fwrite_map(file, SKY_LOOKUP_MESSAGE_TABLE_KEY_COUNT, &sz);
    check(sz > 0, "Unable to write map");
    
    // Action names.
    check(sky_minipack_fwrite_bstring(file, &SKY_LOOKUP_MESSAGE_ACTION_NAMES_STR) == 0, "Unable to write action names key");
    minipack_fwrite_array(file, message->action_name_count, &sz);
    check(sz > 0, "Unable to write action names array");

    for(i=0; i<message->action_name_count; i++) {
        rc = sky_minipack_fwrite_bstring(file, message->action_names[i]);
        check(rc == 0, "Unable to write action name");
    }

    // Property names.
    check(sky_minipack_fwrite_bstring(file, &SKY_LOOKUP_MESSAGE_PROPERTY_NAMES_STR) == 0, "Unable to write property names key");
    minipack_fwrite_array(file, message->property_name_count, &sz);
    check(sz > 0, "Unable to write property names array");

    for(i=0; i<message->property_name_count; i++) {
        rc = sky_minipack_fwrite_bstring(file, message->property_names[i]);
        check(rc == 0, "Unable to write property name");
    }

    return 0;

error:
    return -1;
}

// Deserializes a 'lookup' message from a file stream.
//
// message - The message.
// file    - The file stream to read from.
//
// Returns 0 if successful, otherwise returns -1.
int sky_lookup_message_unpack(sky_lookup_message *message, FILE *file)
{
    int rc;
    size_t sz;
    bstring key = NULL;
    check(message != NULL, "Message required");
    check(file != NULL, "File stream required");

    // Map
    uint32_t map_length = minipack_fread_map(file, &sz);
    check(sz > 0, "Unable to read map");
    
    // Map items
    uint32_t i, j;
    for(i=0; i<map_length; i++) {
        rc = sky_minipack_fread_bstring(file, &key);
        check(rc == 0, "Unable to read map key");
        
        // Action names.
        if(biseq(key, &SKY_LOOKUP_MESSAGE_ACTION_NAMES_STR)) {
            message->action_name_count = minipack_fread_array(file, &sz);
            check(sz != 0, "Unable to read action names array");

            message->action_names = calloc(message->action_name_count, sizeof(*message->action_names));
            check_mem(message->action_names);
            
            for(j=0; j<message->action_name_count; j++) {
                rc = sky_minipack_fread_bstring(file, &message->action_names[j]);
                check(rc == 0, "Unable to read action name");
            }
        }
        // Property names.
        else if(biseq(key, &SKY_LOOKUP_MESSAGE_PROPERTY_NAMES_STR)) {
            message->property_name_count = minipack_fread_array(file, &sz);
            check(sz != 0, "Unable to read property names array");

            message->property_names = calloc(message->property_name_count, sizeof(*message->property_names));
            check_mem(message->property_names);
            
            for(j=0; j<message->property_name_count; j++) {
                rc = sky_minipack_fread_bstring(file, &message->property_names[j]);
                check(rc == 0, "Unable to read property name");
            }
        }

        bdestroy(key);
    }
    
    return 0;

error:
    return -1;
}


