#include <stdlib.h>
#include <stdio.h>
#include <arpa/inet.h>

#include "types.h"
#include "get_property_message.h"
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

// Creates a 'get_property' message object.
//
// Returns a new message.
sky_get_property_message *sky_get_property_message_create()
{
    sky_get_property_message *message = NULL;
    message = calloc(1, sizeof(sky_get_property_message)); check_mem(message);
    return message;

error:
    sky_get_property_message_free(message);
    return NULL;
}

// Frees a 'get_property' message object from memory.
//
// message - The message object to be freed.
//
// Returns nothing.
void sky_get_property_message_free(sky_get_property_message *message)
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
size_t sky_get_property_message_sizeof(sky_get_property_message *message)
{
    size_t sz = 0;
    sz += minipack_sizeof_int(message->property_id);
    return sz;
}

// Serializes a 'get_property' message to a file stream.
//
// message - The message.
// file    - The file stream to write to.
//
// Returns 0 if successful, otherwise returns -1.
int sky_get_property_message_pack(sky_get_property_message *message, FILE *file)
{
    size_t sz;
    check(message != NULL, "Message required");
    check(file != NULL, "File stream required");

    minipack_fwrite_int(file, message->property_id, &sz);
    check(sz > 0, "Unable to pack property id");
    
    return 0;

error:
    return -1;
}

// Deserializes a 'get_property' message from a file stream.
//
// message - The message.
// file    - The file stream to read from.
//
// Returns 0 if successful, otherwise returns -1.
int sky_get_property_message_unpack(sky_get_property_message *message, FILE *file)
{
    size_t sz;
    check(message != NULL, "Message required");
    check(file != NULL, "File stream required");

    message->property_id = (sky_property_id_t)minipack_fread_int(file, &sz);
    check(sz > 0, "Unable to unpack property id");

    return 0;

error:
    return -1;
}


//--------------------------------------
// Processing
//--------------------------------------

// Applies a 'get_property' message to a table.
//
// message - The message.
// table   - The table to apply the message to.
// output  - The output stream to write to.
//
// Returns 0 if successful, otherwise returns -1.
int sky_get_property_message_process(sky_get_property_message *message,
                                     sky_table *table, FILE *output)
{
    int rc;
    size_t sz;
    check(message != NULL, "Message required");
    check(table != NULL, "Table required");
    check(output != NULL, "Output stream required");

    struct tagbstring status_str = bsStatic("status");
    struct tagbstring ok_str = bsStatic("ok");
    struct tagbstring property_str = bsStatic("property");

    // Retrieve property.
    sky_property *property = NULL;
    rc = sky_property_file_find_by_id(table->property_file, message->property_id, &property);
    check(rc == 0, "Unable to add property");
    
    // Return.
    //   {status:"OK", property:{...}}
    minipack_fwrite_map(output, 2, &sz);
    check(sz > 0, "Unable to write output");
    check(sky_minipack_fwrite_bstring(output, &status_str) == 0, "Unable to write status key");
    check(sky_minipack_fwrite_bstring(output, &ok_str) == 0, "Unable to write status value");
    check(sky_minipack_fwrite_bstring(output, &property_str) == 0, "Unable to write property key");
    
    if(property != NULL) {
        check(sky_property_pack(property, output) == 0, "Unable to write property value");
    }
    else {
        minipack_fwrite_nil(output, &sz);
        check(sz > 0, "Unable to write null property value");
    }
    
    return 0;

error:
    return -1;
}