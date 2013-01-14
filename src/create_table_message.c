#include <stdlib.h>
#include <stdio.h>
#include <arpa/inet.h>

#include "types.h"
#include "create_table_message.h"
#include "table.h"
#include "minipack.h"
#include "mem.h"
#include "dbg.h"


//==============================================================================
//
// Definitions
//
//==============================================================================

#define SKY_CREATE_TABLE_MESSAGE_TABLE_KEY_COUNT 2

struct tagbstring SKY_CREATE_TABLE_MESSAGE_NAME_STR = bsStatic("name");
struct tagbstring SKY_CREATE_TABLE_MESSAGE_TABLET_COUNT_STR = bsStatic("tabletCount");


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

// Creates an 'create_table' message object.
//
// Returns a new message.
sky_create_table_message *sky_create_table_message_create()
{
    sky_create_table_message *message = NULL;
    message = calloc(1, sizeof(sky_create_table_message)); check_mem(message);
    return message;

error:
    sky_create_table_message_free(message);
    return NULL;
}

// Frees a 'create_table' message object from memory.
//
// message - The message object to be freed.
//
// Returns nothing.
void sky_create_table_message_free(sky_create_table_message *message)
{
    if(message) {
        free(message);
    }
}



//--------------------------------------
// Message Handler
//--------------------------------------

// Creates a message handler for the message.
//
// Returns a message handler.
sky_message_handler *sky_create_table_message_handler_create()
{
    sky_message_handler *handler = sky_message_handler_create(); check_mem(handler);
    handler->scope = SKY_MESSAGE_HANDLER_SCOPE_SERVER;
    handler->name = bfromcstr("create_table");
    handler->process = sky_create_table_message_process;
    return handler;

error:
    sky_message_handler_free(handler);
    return NULL;
}

// Adds a table to the server. This function is synchronous and does not use
// a worker.
//
// server - The server.
// header - The message header.
// table  - The table the message is working against
// input  - The input file stream.
// output - The output file stream.
//
// Returns 0 if successful, otherwise returns -1.
int sky_create_table_message_process(sky_server *server,
                                   sky_message_header *header,
                                   sky_table *table, FILE *input, FILE *output)
{
    int rc = 0;
    size_t sz;
    sky_create_table_message *message = NULL;
    check(server != NULL, "Server required");
    check(header != NULL, "Message header required");
    check(input != NULL, "Input stream required");
    check(output != NULL, "Output stream required");
    (void)table;
    
    struct tagbstring status_str = bsStatic("status");
    struct tagbstring ok_str = bsStatic("ok");
    struct tagbstring error_str = bsStatic("error");
    struct tagbstring message_str = bsStatic("message");
    struct tagbstring table_exists_str = bsStatic("Table already exists");

    // Parse message.
    message = sky_create_table_message_create(); check_mem(message);
    message->data_path = bstrcpy(server->path); check_mem(message->data_path);
    message->table = sky_table_create();
    rc = sky_create_table_message_unpack(message, input);
    check(rc == 0, "Unable to parse 'create_table' message");

    // If table already exists then return an error.
    if(sky_file_exists(message->table->path)) {
        // Return.
        //   {status:"error", message:""}
        minipack_fwrite_map(output, 2, &sz);
        check(sz > 0, "Unable to write output");
        check(sky_minipack_fwrite_bstring(output, &status_str) == 0, "Unable to write status key");
        check(sky_minipack_fwrite_bstring(output, &error_str) == 0, "Unable to write status value");
        check(sky_minipack_fwrite_bstring(output, &message_str) == 0, "Unable to write message key");
        check(sky_minipack_fwrite_bstring(output, &table_exists_str) == 0, "Unable to write message");
    }
    // Otherwise create the table.
    else {
        // Open the table and close it to create it.
        rc = sky_table_open(message->table);
        check(rc == 0, "Unable to open table for creation: %s", bdata(message->table->path));

        // Close the table.
        rc = sky_table_close(message->table);
        check(rc == 0, "Unable to close table: %s", bdata(message->table->path));

        // Return.
        //   {status:"OK"}
        minipack_fwrite_map(output, 1, &sz);
        check(sz > 0, "Unable to write output");
        check(sky_minipack_fwrite_bstring(output, &status_str) == 0, "Unable to write status key");
        check(sky_minipack_fwrite_bstring(output, &ok_str) == 0, "Unable to write status value");
    }
    
    // Clean up.
    if(!header->multi) {
        fclose(input);
        fclose(output);
    }

    sky_table_free(message->table);
    sky_create_table_message_free(message);
    
    return 0;

error:
    if(!header->multi) {
        if(input) fclose(input);
        if(output) fclose(output);
    }
    sky_table_free(message->table);
    sky_create_table_message_free(message);
    return -1;
}


//--------------------------------------
// Serialization
//--------------------------------------

// Serializes an 'create_table' message to a file stream.
//
// message - The message.
// file    - The file stream to write to.
//
// Returns 0 if successful, otherwise returns -1.
int sky_create_table_message_pack(sky_create_table_message *message, FILE *file)
{
    size_t sz;
    bstring table_name = NULL;
    check(message != NULL, "Message required");
    check(file != NULL, "File stream required");

    struct tagbstring slash_str = bsStatic("/");
    
    // Extract name from the end of the path.
    int index = binchrr(message->table->path, blength(message->table->path), &slash_str);
    table_name = (index == BSTR_ERR ? bstrcpy(message->table->path) : bmidstr(message->table->path, index+1, blength(message->table->path)));

    // Map
    minipack_fwrite_map(file, SKY_CREATE_TABLE_MESSAGE_TABLE_KEY_COUNT, &sz);
    check(sz > 0, "Unable to write map");
    
    // Name
    check(sky_minipack_fwrite_bstring(file, &SKY_CREATE_TABLE_MESSAGE_NAME_STR) == 0, "Unable to write table name key");
    check(sky_minipack_fwrite_bstring(file, table_name) == 0, "Unable to write name value");

    // Tablet count
    check(sky_minipack_fwrite_bstring(file, &SKY_CREATE_TABLE_MESSAGE_TABLET_COUNT_STR) == 0, "Unable to write tablet count key");
    minipack_fwrite_uint(file, message->table->default_tablet_count, &sz);
    check(sz > 0, "Unable to pack table count value");
    
    bdestroy(table_name);
    
    return 0;

error:
    bdestroy(table_name);
    return -1;
}

// Deserializes an 'create_table' message from a file stream.
//
// message - The message.
// file    - The file stream to read from.
//
// Returns 0 if successful, otherwise returns -1.
int sky_create_table_message_unpack(sky_create_table_message *message, FILE *file)
{
    int rc;
    size_t sz;
    bstring key = NULL;
    bstring table_name = NULL;
    check(message != NULL, "Message required");
    check(file != NULL, "File stream required");

    // Map
    uint32_t map_length = minipack_fread_map(file, &sz);
    check(sz > 0, "Unable to read map");
    
    // Map items
    uint32_t i;
    for(i=0; i<map_length; i++) {
        rc = sky_minipack_fread_bstring(file, &key);
        check(rc == 0, "Unable to read map key");
        
        if(biseq(key, &SKY_CREATE_TABLE_MESSAGE_NAME_STR)) {
            rc = sky_minipack_fread_bstring(file, &table_name);
            check(rc == 0, "Unable to read table name");
            
            message->table->path = bformat("%s/%s", bdata(message->data_path), bdata(table_name));
            check_mem(table_name);
            
            bdestroy(table_name);
            table_name = NULL;
        }
        else if(biseq(key, &SKY_CREATE_TABLE_MESSAGE_TABLET_COUNT_STR)) {
            message->table->default_tablet_count = (uint32_t)minipack_fread_uint(file, &sz);
            check(sz > 0, "Unable to read tablet count");
        }

        bdestroy(key);
    }
    
    bdestroy(table_name);
    return 0;

error:
    bdestroy(table_name);
    bdestroy(key);
    return -1;
}
