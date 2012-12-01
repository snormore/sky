#include <stdlib.h>
#include <stdio.h>
#include <arpa/inet.h>
#include <assert.h>

#include "types.h"
#include "multi_message.h"
#include "minipack.h"
#include "mem.h"
#include "dbg.h"


//==============================================================================
//
// Definitions
//
//==============================================================================

#define SKY_MULTI_KEY_COUNT 1

struct tagbstring SKY_MULTI_KEY_COUNT_STR = bsStatic("count");


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

// Creates an 'multi' message object.
//
// Returns a new message.
sky_multi_message *sky_multi_message_create()
{
    sky_multi_message *message = calloc(1, sizeof(sky_multi_message)); check_mem(message);
    return message;

error:
    sky_multi_message_free(message);
    return NULL;
}

// Frees an 'multi' message object from memory.
//
// message - The message object to be freed.
//
// Returns nothing.
void sky_multi_message_free(sky_multi_message *message)
{
    if(message) {
        message->message_count = 0;
        free(message);
    }
}


//--------------------------------------
// Message Handler
//--------------------------------------

// Creates a message handler for the 'Multi' message.
//
// Returns a message handler.
sky_message_handler *sky_multi_message_handler_create()
{
    sky_message_handler *handler = sky_message_handler_create(); check_mem(handler);
    handler->scope = SKY_MESSAGE_HANDLER_SCOPE_SERVER;
    handler->name = bfromcstr("multi");
    handler->process = sky_multi_message_process;
    return handler;

error:
    sky_message_handler_free(handler);
    return NULL;
}

// Delegates processing of the 'Multi' message to a worker.
//
// server - The server.
// header - The message header.
// table  - The table the message is working against
// input  - The input file stream.
// output - The output file stream.
//
// Returns 0 if successful, otherwise returns -1.
int sky_multi_message_process(sky_server *server,
                              sky_message_header *header,
                              sky_table *table, 
                              FILE *input, FILE *output)
{
    int rc = 0;
    size_t sz;
    sky_multi_message *message = NULL;
    assert(header != NULL);
    assert(table == NULL);
    assert(input != NULL);
    assert(output != NULL);
    
    // Parse message.
    message = sky_multi_message_create(); check_mem(message);
    rc = sky_multi_message_unpack(message, input);
    check(rc == 0, "Unable to unpack 'multi' message");

    // Send back an array header with the number of responses so the client
    // will know how many to expect.
    minipack_fwrite_array(output, message->message_count, &sz);
    check(sz > 0, "Unable to write multi message array");

    // Loop over child messages and process.
    uint32_t i;
    for(i=0; i<message->message_count; i++) {
        rc = sky_server_process_message(server, true, input, output);
        check(rc == 0, "Unable to process child message");
    }

    // Close streams.
    fclose(input);
    fclose(output);
    
    return 0;

error:
    sky_multi_message_free(message);
    fclose(input);
    fclose(output);
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
size_t sky_multi_message_sizeof(sky_multi_message *message)
{
    size_t sz = 0;
    sz += minipack_sizeof_map(SKY_MULTI_KEY_COUNT);
    sz += minipack_sizeof_raw((&SKY_MULTI_KEY_COUNT_STR)->slen) + (&SKY_MULTI_KEY_COUNT_STR)->slen;
    sz += minipack_sizeof_uint(message->message_count);
    return sz;
}

// Serializes an 'multi' message to a file stream.
//
// message - The message.
// file    - The file stream to write to.
//
// Returns 0 if successful, otherwise returns -1.
int sky_multi_message_pack(sky_multi_message *message, FILE *file)
{
    size_t sz;
    assert(message != NULL);
    assert(file != NULL);

    // Map
    minipack_fwrite_map(file, SKY_MULTI_KEY_COUNT, &sz);
    check(sz > 0, "Unable to write map");
    
    // Object ID
    check(sky_minipack_fwrite_bstring(file, &SKY_MULTI_KEY_COUNT_STR) == 0, "Unable to pack count key");
    minipack_fwrite_int(file, message->message_count, &sz);
    check(sz != 0, "Unable to pack count");

    return 0;

error:
    return -1;
}

// Deserializes an 'multi' message from a file stream.
//
// message - The message.
// file    - The file stream to read from.
//
// Returns 0 if successful, otherwise returns -1.
int sky_multi_message_unpack(sky_multi_message *message, FILE *file)
{
    int rc;
    size_t sz;
    bstring key = NULL;
    assert(message != NULL);
    assert(file != NULL);

    // Map
    uint32_t map_length = minipack_fread_map(file, &sz);
    check(sz > 0, "Unable to read map");
    
    // Map items
    uint32_t i;
    for(i=0; i<map_length; i++) {
        rc = sky_minipack_fread_bstring(file, &key);
        check(rc == 0, "Unable to read map key");
        
        if(biseq(key, &SKY_MULTI_KEY_COUNT_STR) == 1) {
            message->message_count = (uint32_t)minipack_fread_uint(file, &sz);
            check(sz != 0, "Unable to unpack count");
        }
        
        bdestroy(key);
    }

    return 0;

error:
    bdestroy(key);
    return -1;
}

