#include <stdlib.h>
#include <stdio.h>
#include <arpa/inet.h>

#include "types.h"
#include "add_property_message.h"
#include "property.h"
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

// Creates an 'add_property' message object.
//
// Returns a new message.
sky_add_property_message *sky_add_property_message_create()
{
    sky_add_property_message *message = NULL;
    message = calloc(1, sizeof(sky_add_property_message)); check_mem(message);
    message->property = sky_property_create(); check_mem(message->property);
    return message;

error:
    sky_add_property_message_free(message);
    return NULL;
}

// Frees an 'add_property' message object from memory.
//
// message - The message object to be freed.
//
// Returns nothing.
void sky_add_property_message_free(sky_add_property_message *message)
{
    if(message) {
        sky_add_property_message_free_property(message);
        free(message);
    }
}

// Frees an property associated with an 'add_property' message.
//
// message - The message object to be freed.
//
// Returns nothing.
void sky_add_property_message_free_property(sky_add_property_message *message)
{
    if(message) {
        // Only free the property if it's not managed by an property file.
        if(message->property && message->property->property_file == NULL) {
            sky_property_free(message->property);
        }
        message->property = NULL;
    }
}


//--------------------------------------
// Message Handler
//--------------------------------------

// Creates a message handler for the 'Add Property' message.
//
// Returns a message handler.
sky_message_handler *sky_add_property_message_handler_create()
{
    sky_message_handler *handler = sky_message_handler_create(); check_mem(handler);
    handler->scope = SKY_MESSAGE_HANDLER_SCOPE_TABLE;
    handler->name = bfromcstr("add_property");
    handler->process = sky_add_property_message_process;
    return handler;

error:
    sky_message_handler_free(handler);
    return NULL;
}

// Adds an property to a table. This function is synchronous and does not use
// a worker.
//
// server - The server.
// header - The message header.
// table  - The table the message is working against
// input  - The input file stream.
// output - The output file stream.
//
// Returns 0 if successful, otherwise returns -1.
int sky_add_property_message_process(sky_server *server,
                                     sky_message_header *header,
                                     sky_table *table, FILE *input, FILE *output)
{
    int rc = 0;
    size_t sz;
    sky_add_property_message *message = NULL;
    check(server != NULL, "Server required");
    check(header != NULL, "Message header required");
    check(table != NULL, "Table required");
    check(input != NULL, "Input stream required");
    check(output != NULL, "Output stream required");
    
    struct tagbstring status_str = bsStatic("status");
    struct tagbstring ok_str = bsStatic("ok");
    struct tagbstring property_str = bsStatic("property");

    // Parse message.
    message = sky_add_property_message_create(); check_mem(message);
    rc = sky_add_property_message_unpack(message, input);
    check(rc == 0, "Unable to parse 'add_property' message");

    // Add property.
    rc = sky_property_file_add_property(table->property_file, message->property);
    check(rc == 0, "Unable to add property");
    
    // Save property file.
    rc = sky_property_file_save(table->property_file);
    check(rc == 0, "Unable to save property file");
    
    // Return.
    //   {status:"OK", property:{...}}
    minipack_fwrite_map(output, 2, &sz);
    check(sz > 0, "Unable to write output");
    check(sky_minipack_fwrite_bstring(output, &status_str) == 0, "Unable to write status key");
    check(sky_minipack_fwrite_bstring(output, &ok_str) == 0, "Unable to write status value");
    check(sky_minipack_fwrite_bstring(output, &property_str) == 0, "Unable to write property key");
    check(sky_property_pack(message->property, output) == 0, "Unable to write property value");

    // Clean up.
    sky_add_property_message_free(message);
    fclose(input);
    fclose(output);

    return 0;

error:
    sky_add_property_message_free(message);
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
size_t sky_add_property_message_sizeof(sky_add_property_message *message)
{
    size_t sz = 0;
    sz += sky_property_sizeof(message->property);
    return sz;
}

// Serializes an 'add_property' message to a file stream.
//
// message - The message.
// file    - The file stream to write to.
//
// Returns 0 if successful, otherwise returns -1.
int sky_add_property_message_pack(sky_add_property_message *message, FILE *file)
{
    int rc;
    check(message != NULL, "Message required");
    check(file != NULL, "File stream required");

    rc = sky_property_pack(message->property, file);
    check(rc == 0, "Unable to pack property");
    
    return 0;

error:
    return -1;
}

// Deserializes an 'add_property' message from a file stream.
//
// message - The message.
// file    - The file stream to read from.
//
// Returns 0 if successful, otherwise returns -1.
int sky_add_property_message_unpack(sky_add_property_message *message, FILE *file)
{
    int rc;
    check(message != NULL, "Message required");
    check(file != NULL, "File stream required");

    rc = sky_property_unpack(message->property, file);
    check(rc == 0, "Unable to unpack property");

    return 0;

error:
    return -1;
}


