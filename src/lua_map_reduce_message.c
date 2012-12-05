#include <stdlib.h>
#include <stdio.h>
#include <arpa/inet.h>
#include <sys/time.h>
#include <assert.h>

#include "types.h"
#include "lua_map_reduce_message.h"
#include "path_iterator.h"
#include "action.h"
#include "minipack.h"
#include "mem.h"
#include "dbg.h"


//==============================================================================
//
// Definitions
//
//==============================================================================

//--------------------------------------
// String Constants
//--------------------------------------

#define SKY_LUA_MAP_REDUCE_KEY_COUNT 1

struct tagbstring SKY_LUA_MAP_REDUCE_KEY_SOURCE     = bsStatic("source");

struct tagbstring SKY_LUA_MAP_REDUCE_STATUS_STR = bsStatic("status");
struct tagbstring SKY_LUA_MAP_REDUCE_OK_STR     = bsStatic("ok");
struct tagbstring SKY_LUA_MAP_REDUCE_DATA_STR   = bsStatic("data");


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

// Creates a 'lua::map_reduce' message object.
//
// Returns a new message.
sky_lua_map_reduce_message *sky_lua_map_reduce_message_create()
{
    sky_lua_map_reduce_message *message = NULL;
    message = calloc(1, sizeof(sky_lua_map_reduce_message)); check_mem(message);
    return message;

error:
    sky_lua_map_reduce_message_free(message);
    return NULL;
}

// Frees a 'lua::map_reduce' message object from memory.
//
// message - The message object to be freed.
//
// Returns nothing.
void sky_lua_map_reduce_message_free(sky_lua_map_reduce_message *message)
{
    if(message) {
        if(message->L) lua_close(message->L);
        message->L = NULL;
        
        bdestroy(message->source);
        message->source = NULL;
        
        bdestroy(message->results);
        message->results = NULL;

        free(message);
    }
}


//--------------------------------------
// Message Handler
//--------------------------------------

// Creates a message handler for the 'lua::map_reduce' message.
//
// Returns a message handler.
sky_message_handler *sky_lua_map_reduce_message_handler_create()
{
    sky_message_handler *handler = sky_message_handler_create(); check_mem(handler);
    handler->scope = SKY_MESSAGE_HANDLER_SCOPE_TABLE;
    handler->name = bfromcstr("lua::map_reduce");
    handler->process = sky_lua_map_reduce_message_process;
    return handler;

error:
    sky_message_handler_free(handler);
    return NULL;
}

// Delegates processing of the 'Next Actions' message to a worker.
//
// server  - The server.
// header  - The message header.
// table   - The table the message is working against
// input   - The input file stream.
// output  - The output file stream.
//
// Returns 0 if successful, otherwise returns -1.
int sky_lua_map_reduce_message_process(sky_server *server,
                                       sky_message_header *header,
                                       sky_table *table,
                                       FILE *input, FILE *output)
{
    int rc = 0;
    sky_lua_map_reduce_message *message = NULL;
    assert(table != NULL);
    assert(header != NULL);
    assert(input != NULL);
    assert(output != NULL);
    
    // Create worker.
    sky_worker *worker = sky_worker_create(); check_mem(worker);
    worker->context = server->context;
    worker->map = sky_lua_map_reduce_message_worker_map;
    worker->map_free = sky_lua_map_reduce_message_worker_map_free;
    worker->reduce = sky_lua_map_reduce_message_worker_reduce;
    worker->write = sky_lua_map_reduce_message_worker_write;
    worker->free = sky_lua_map_reduce_message_worker_free;
    worker->input = input;
    worker->output = output;

    // Attach servlets.
    rc = sky_server_get_table_servlets(server, table, &worker->servlets, &worker->servlet_count);
    check(rc == 0, "Unable to copy servlets to worker");

    // Create and parse message object.
    message = sky_lua_map_reduce_message_create(); check_mem(message);
    message->results = bfromcstr("\x80"); check_mem(message->results);
    rc = sky_lua_map_reduce_message_unpack(message, input);
    check(rc == 0, "Unable to unpack 'lua::map_reduce' message");
    check(message->source != NULL, "Lua source required");

    // Compile Lua script.
    rc = sky_lua_initscript_with_table(message->source, table, NULL, &message->L);
    check(rc == 0, "Unable to initialize script");

    // Attach message to worker.
    worker->data = (sky_lua_map_reduce_message*)message;
    
    // Start worker.
    rc = sky_worker_start(worker);
    check(rc == 0, "Unable to start worker");
    
    return 0;

error:
    sky_lua_map_reduce_message_free(message);
    sky_worker_free(worker);
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
size_t sky_lua_map_reduce_message_sizeof(sky_lua_map_reduce_message *message)
{
    size_t sz = 0;
    sz += minipack_sizeof_map(SKY_LUA_MAP_REDUCE_KEY_COUNT);
    sz += minipack_sizeof_raw((&SKY_LUA_MAP_REDUCE_KEY_SOURCE)->slen) + (&SKY_LUA_MAP_REDUCE_KEY_SOURCE)->slen;
    sz += minipack_sizeof_raw(blength(message->source)) + blength(message->source);
    return sz;
}

// Serializes a 'lua::map_reduce' message to a file stream.
//
// message - The message.
// file    - The file stream to write to.
//
// Returns 0 if successful, otherwise returns -1.
int sky_lua_map_reduce_message_pack(sky_lua_map_reduce_message *message,
                                    FILE *file)
{
    size_t sz;
    assert(message != NULL);
    assert(file != NULL);

    // Map
    minipack_fwrite_map(file, SKY_LUA_MAP_REDUCE_KEY_COUNT, &sz);
    check(sz > 0, "Unable to write map");
    
    // Source
    check(sky_minipack_fwrite_bstring(file, &SKY_LUA_MAP_REDUCE_KEY_SOURCE) == 0, "Unable to pack source key");
    check(sky_minipack_fwrite_bstring(file, message->source) == 0, "Unable to pack source");

    return 0;

error:
    return -1;
}

// Deserializes an 'lua::map_reduce' message from a file stream.
//
// message - The message.
// file    - The file stream to read from.
//
// Returns 0 if successful, otherwise returns -1.
int sky_lua_map_reduce_message_unpack(sky_lua_map_reduce_message *message,
                                      FILE *file)
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
        
        if(biseq(key, &SKY_LUA_MAP_REDUCE_KEY_SOURCE) == 1) {
            rc = sky_minipack_fread_bstring(file, &message->source);
            check(rc == 0, "Unable to read source");
        }
        
        bdestroy(key);
        key = NULL;
    }

    return 0;

error:
    bdestroy(key);
    return -1;
}


//--------------------------------------
// Worker
//--------------------------------------

// Maps tablet data to a next action summation data structure.
//
// worker - The worker.
// tablet - The tablet to work against.
// ret    - A pointer to where the summation data structure should be returned.
//
// Returns 0 if successful, otherwise returns -1.
int sky_lua_map_reduce_message_worker_map(sky_worker *worker, sky_tablet *tablet,
                                          void **ret)
{
    int rc;
    lua_State *L = NULL;
    bstring msgpack_ret = NULL;
    sky_data_descriptor *descriptor = NULL;
    assert(worker != NULL);
    assert(tablet != NULL);
    assert(ret != NULL);

    sky_lua_map_reduce_message *message = (sky_lua_map_reduce_message*)worker->data;
    
    // Initialize the path iterator.
    sky_path_iterator iterator;
    sky_path_iterator_init(&iterator);

    // Compile Lua script.
    descriptor = sky_data_descriptor_create(); check_mem(descriptor);
    rc = sky_lua_initscript_with_table(message->source, tablet->table, descriptor, &L);
    check(rc == 0, "Unable to initialize script");
    
    iterator.cursor.data_descriptor = descriptor;
    iterator.cursor.data = calloc(1, descriptor->data_sz); check_mem(iterator.cursor.data);

    // Assign the data file to iterate over.
    rc = sky_path_iterator_set_tablet(&iterator, tablet);
    check(rc == 0, "Unable to initialize path iterator");

    // Execute function.
    lua_getglobal(L, "sky_map_all");
    lua_pushlightuserdata(L, &iterator);
    rc = lua_pcall(L, 1, 1, 0);
    check(rc == 0, "Unable to execute Lua script: %s", lua_tostring(L, -1));

    // Execute the script and return a msgpack variable.
    rc = sky_lua_msgpack_pack(L, &msgpack_ret);
    check(rc == 0, "Unable to execute Lua script");

    // Close Lua.
    lua_close(L);
    
    // Return msgpack encoded response.
    *ret = (void*)msgpack_ret;
    
    free(iterator.cursor.data);
    sky_data_descriptor_free(descriptor);
    sky_path_iterator_uninit(&iterator);
    return 0;

error:
    *ret = NULL;
    bdestroy(msgpack_ret);
    if(L) lua_close(L);
    free(iterator.cursor.data);
    sky_data_descriptor_free(descriptor);
    sky_path_iterator_uninit(&iterator);
    return -1;
}

// Frees the data structure created and returned in the map() function.
//
// data - A pointer to the data to be freed.
//
// Returns 0 if successful, otherwise returns -1.
int sky_lua_map_reduce_message_worker_map_free(void *data)
{
    assert(data != NULL);
    free(data);
    return 0;
}

// Combines the data from a single execution of the map() function into data
// saved against the worker.
//
// worker - The worker.
// data   - Data created and returned in the map() function.
//
// Returns 0 if successful, otherwise returns -1.
int sky_lua_map_reduce_message_worker_reduce(sky_worker *worker, void *data)
{
    int rc;
    assert(worker != NULL);
    assert(data != NULL);
    
    sky_lua_map_reduce_message *message = (sky_lua_map_reduce_message*)worker->data;

    // Retrieve ref to 'reduce()' function.
    lua_getglobal(message->L, "reduce");

    // Push 'results' table to the function.
    rc = sky_lua_msgpack_unpack(message->L, message->results);
    check(rc == 0, "Unable to push results table to Lua");
    bdestroy(message->results);
    message->results = NULL;

    // Push 'map_data' table to the function.
    rc = sky_lua_msgpack_unpack(message->L, (bstring)data);
    check(rc == 0, "Unable to push map data table to Lua");
    
    // Execute 'reduce(results, data)'.
    rc = lua_pcall(message->L, 2, 1, 0);
    check(rc == 0, "Unable to execute Lua script: %s", lua_tostring(message->L, -1));

    // Execute the script and return a msgpack variable.
    rc = sky_lua_msgpack_pack(message->L, &message->results);
    check(rc == 0, "Unable to unpack results table from Lua script");

    return 0;

error:
    return -1;
}

// Writes the results to an output stream.
//
// worker - The worker.
// output - The output stream.
//
// Returns 0 if successful, otherwise returns -1.
int sky_lua_map_reduce_message_worker_write(sky_worker *worker, FILE *output)
{
    size_t sz;
    assert(worker != NULL);
    assert(output != NULL);
    
    // Ease-of-use references.
    sky_lua_map_reduce_message *message = (sky_lua_map_reduce_message*)worker->data;

    // Return.
    //   {status:"ok", data:{<action_id>:{count:0}, ...}}
    check(minipack_fwrite_map(output, 2, &sz) == 0, "Unable to write root map");
    check(sky_minipack_fwrite_bstring(output, &SKY_LUA_MAP_REDUCE_STATUS_STR) == 0, "Unable to write status key");
    check(sky_minipack_fwrite_bstring(output, &SKY_LUA_MAP_REDUCE_OK_STR) == 0, "Unable to write status value");
    check(sky_minipack_fwrite_bstring(output, &SKY_LUA_MAP_REDUCE_DATA_STR) == 0, "Unable to write data key");
    check(fwrite(bdatae(message->results, ""), blength(message->results), 1, output) == 1, "Unable to write data value");
    
    return 0;

error:
    return -1;
}

// Frees all data attached to the worker.
//
// worker - The worker.
//
// Returns 0 if successful, otherwise returns -1.
int sky_lua_map_reduce_message_worker_free(sky_worker *worker)
{
    assert(worker != NULL);
    
    // Clean up.
    sky_lua_map_reduce_message *message = (sky_lua_map_reduce_message*)worker->data;
    sky_lua_map_reduce_message_free(message);
    worker->data = NULL;
    
    return 0;
}

