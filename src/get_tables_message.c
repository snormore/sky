#include <stdlib.h>
#include <stdio.h>
#include <arpa/inet.h>
#include <dirent.h>

#include "types.h"
#include "get_tables_message.h"
#include "table.h"
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

// Creates an 'get_tables' message object.
//
// Returns a new message.
sky_get_tables_message *sky_get_tables_message_create()
{
    sky_get_tables_message *message = NULL;
    message = calloc(1, sizeof(sky_get_tables_message)); check_mem(message);
    return message;

error:
    sky_get_tables_message_free(message);
    return NULL;
}

// Frees a 'get_tables' message object from memory.
//
// message - The message object to be freed.
//
// Returns nothing.
void sky_get_tables_message_free(sky_get_tables_message *message)
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
sky_message_handler *sky_get_tables_message_handler_create()
{
    sky_message_handler *handler = sky_message_handler_create(); check_mem(handler);
    handler->scope = SKY_MESSAGE_HANDLER_SCOPE_SERVER;
    handler->name = bfromcstr("get_tables");
    handler->process = sky_get_tables_message_process;
    return handler;

error:
    sky_message_handler_free(handler);
    return NULL;
}

// Retrieves a table to the server. This function is synchronous and does not use
// a worker.
//
// server - The server.
// header - The message header.
// table  - The table the message is working against
// input  - The input file stream.
// output - The output file stream.
//
// Returns 0 if successful, otherwise returns -1.
int sky_get_tables_message_process(sky_server *server,
                                   sky_message_header *header,
                                   sky_table *_table, FILE *input, FILE *output)
{
    int rc = 0;
    size_t sz;
    uint32_t i;
    sky_get_tables_message *message = NULL;
    bstring path = NULL;
    bstring *table_names = NULL;
    uint32_t table_name_count = 0;
    DIR *dir = NULL;
    check(server != NULL, "Server required");
    check(header != NULL, "Message header required");
    check(input != NULL, "Input stream required");
    check(output != NULL, "Output stream required");
    (void)_table;
    
    struct tagbstring status_str = bsStatic("status");
    struct tagbstring ok_str = bsStatic("ok");
    struct tagbstring tables_str = bsStatic("tables");
    struct tagbstring name_str = bsStatic("name");

    // Parse message.
    message = sky_get_tables_message_create(); check_mem(message);
    rc = sky_get_tables_message_unpack(message, input);
    check(rc == 0, "Unable to parse 'get_tables' message");

    // Retrieve a list of directories from the server data path.
    // Open directory.
    dir = opendir(bdatae(server->path, ""));
    check(dir != NULL, "Unable to open server data directory: %s", bdatae(server->path, ""));
    
    // Copy over contents of directory.
    struct dirent *ent = NULL;
    while((ent = readdir(dir))) {
        if(strcmp(ent->d_name, ".") != 0 && strcmp(ent->d_name, "..") != 0) {
            path = bformat("%s/%s", bdata(server->path), ent->d_name); check_mem(path);
            if(sky_file_is_dir(path)) {
                // Resize table name array.
                table_names = realloc(table_names, sizeof(*table_names) * (table_name_count+1));
                check_mem(table_names);
                table_name_count++;

                // Append table name.
                table_names[table_name_count-1] = bfromcstr(ent->d_name);
                check_mem(table_names[table_name_count-1]);
            }
            bdestroy(path);
            path = NULL;
        }
    }
    closedir(dir);
    dir = NULL;
    
    // Return.
    //   {status:"OK"}
    minipack_fwrite_map(output, 2, &sz);
    check(sz > 0, "Unable to write output");
    check(sky_minipack_fwrite_bstring(output, &status_str) == 0, "Unable to write status key");
    check(sky_minipack_fwrite_bstring(output, &ok_str) == 0, "Unable to write status value");
    check(sky_minipack_fwrite_bstring(output, &tables_str) == 0, "Unable to write tables key");

    // Write all table names.
    minipack_fwrite_array(output, table_name_count, &sz);
    check(sz > 0, "Unable to write tables array");

    for(i=0; i<table_name_count; i++) {
        minipack_fwrite_map(output, 1, &sz);
        check(sz > 0, "Unable to write map");
        check(sky_minipack_fwrite_bstring(output, &name_str) == 0, "Unable to write table name key");
        check(sky_minipack_fwrite_bstring(output, table_names[i]) == 0, "Unable to write table name: %s", bdata(table_names[i]));
    }
    
    for(i=0; i<table_name_count; i++) {
        bdestroy(table_names[i]);
        table_names[i] = NULL;
    }
    free(table_names);
    
    fclose(input);
    fclose(output);
    bdestroy(path);
    sky_get_tables_message_free(message);

    return 0;

error:
    if(dir != NULL) closedir(dir);
    for(i=0; i<table_name_count; i++) {
        bdestroy(table_names[i]);
        table_names[i] = NULL;
    }
    free(table_names);

    if(input) fclose(input);
    if(output) fclose(output);
    bdestroy(path);
    sky_get_tables_message_free(message);
    return -1;
}


//--------------------------------------
// Serialization
//--------------------------------------

// Serializes an 'get_tables' message to a file stream.
//
// message - The message.
// file    - The file stream to write to.
//
// Returns 0 if successful, otherwise returns -1.
int sky_get_tables_message_pack(sky_get_tables_message *message, FILE *file)
{
    check(message != NULL, "Message required");
    check(file != NULL, "File stream required");

    return 0;

error:
    return -1;
}

// Deserializes an 'get_tables' message from a file stream.
//
// message - The message.
// file    - The file stream to read from.
//
// Returns 0 if successful, otherwise returns -1.
int sky_get_tables_message_unpack(sky_get_tables_message *message, FILE *file)
{
    check(message != NULL, "Message required");
    check(file != NULL, "File stream required");

    return 0;

error:
    return -1;
}
