#include <stdlib.h>
#include <stdio.h>
#include <arpa/inet.h>

#include "types.h"
#include "multi_message.h"
#include "minipack.h"
#include "endian.h"
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

// Creates an MULTI message object.
//
// Returns a new MULTI message.
sky_multi_message *sky_multi_message_create()
{
    sky_multi_message *message = NULL;
    message = calloc(1, sizeof(sky_multi_message)); check_mem(message);
    return message;

error:
    sky_multi_message_free(message);
    return NULL;
}

// Frees a MULTI message object from memory.
//
// message - The message object to be freed.
//
// Returns nothing.
void sky_multi_message_free(sky_multi_message *message)
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
size_t sky_multi_message_sizeof(sky_multi_message *message)
{
    size_t sz = 0;
    sz += minipack_sizeof_uint(message->message_count);
    return sz;
}

// Serializes a MULTI message to a file stream.
//
// message - The message.
// file    - The file stream to write to.
//
// Returns 0 if successful, otherwise returns -1.
int sky_multi_message_pack(sky_multi_message *message, FILE *file)
{
    size_t sz;
    check(message != NULL, "Message required");
    check(file != NULL, "File stream required");

    minipack_fwrite_uint(file, message->message_count, &sz);
    check(sz > 0, "Unable to write message count");
    
    return 0;

error:
    return -1;
}

// Deserializes a MULTI message from a file stream.
//
// message - The message.
// file    - The file stream to read from.
//
// Returns 0 if successful, otherwise returns -1.
int sky_multi_message_unpack(sky_multi_message *message, FILE *file)
{
    size_t sz;
    bstring key = NULL;
    check(message != NULL, "Message required");
    check(file != NULL, "File stream required");

    message->message_count = minipack_fread_uint(file, &sz);
    check(sz > 0, "Unable to read message count");

    return 0;

error:
    bdestroy(key);
    return -1;
}
