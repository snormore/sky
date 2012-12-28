#include <stdlib.h>
#include <stdio.h>
#include <arpa/inet.h>

#include "types.h"
#include "ping_message.h"
#include "minipack.h"
#include "mem.h"
#include "dbg.h"


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Message Handler
//--------------------------------------

// Creates a message handler for the 'ping' message.
//
// Returns a message handler.
sky_message_handler *sky_ping_message_handler_create()
{
    sky_message_handler *handler = sky_message_handler_create(); check_mem(handler);
    handler->scope = SKY_MESSAGE_HANDLER_SCOPE_SERVER;
    handler->name = bfromcstr("ping");
    handler->process = sky_ping_message_process;
    return handler;

error:
    sky_message_handler_free(handler);
    return NULL;
}

// Checks that the server is up and running. This function is synchronous and
// does not use a worker.
//
// server - The server.
// header - The message header.
// table  - The table the message is working against
// input  - The input file stream.
// output - The output file stream.
//
// Returns 0 if successful, otherwise returns -1.
int sky_ping_message_process(sky_server *server,
                                   sky_message_header *header,
                                   sky_table *table, FILE *input, FILE *output)
{
    size_t sz;
    check(server != NULL, "Server required");
    check(header != NULL, "Message header required");
    check(input != NULL, "Input stream required");
    check(output != NULL, "Output stream required");
    UNUSED(table);
    
    struct tagbstring status_str = bsStatic("status");
    struct tagbstring ok_str = bsStatic("ok");

    // Return:
    //   {status:"OK"}
    minipack_fwrite_map(output, 1, &sz);
    check(sz > 0, "Unable to write output");
    check(sky_minipack_fwrite_bstring(output, &status_str) == 0, "Unable to write status key");
    check(sky_minipack_fwrite_bstring(output, &ok_str) == 0, "Unable to write status value");
    
    // Clean up.
    fclose(input);
    fclose(output);

    return 0;

error:
    if(input) fclose(input);
    if(output) fclose(output);
    return -1;
}
