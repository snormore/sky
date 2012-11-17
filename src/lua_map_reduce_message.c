#include <stdlib.h>
#include <stdio.h>
#include <arpa/inet.h>
#include <sys/time.h>
#include <assert.h>

#include <lua.h>
#include <lualib.h>
#include <lauxlib.h>

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

#define SKY_LUA_MAP_REDUCE_KEY_COUNT 2

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
    worker->read = sky_lua_map_reduce_message_worker_read;
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

    // Create a message object.
    message = sky_lua_map_reduce_message_create(); check_mem(message);
    check_mem(message->results);

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
    assert(message != NULL);
    return 0;
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
    assert(message != NULL);
    assert(file != NULL);
    return 0;
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
    assert(message != NULL);
    assert(file != NULL);

    return 0;
}


//--------------------------------------
// Worker
//--------------------------------------

// Reads the message from the file stream.
//
// worker - The worker.
// input  - The input stream.
//
// Returns 0 if successful, otherwise returns -1.
int sky_lua_map_reduce_message_worker_read(sky_worker *worker, FILE *input)
{
    int rc;
    assert(worker != NULL);
    assert(input != NULL);

    // Parse message.
    sky_lua_map_reduce_message *message = (sky_lua_map_reduce_message*)worker->data;
    rc = sky_lua_map_reduce_message_unpack(message, input);
    check(rc == 0, "Unable to unpack 'lua::map_reduce' message");

    return 0;

error:
    return -1;
}

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
    sky_cursor cursor;
    sky_cursor_init(&cursor);
    assert(worker != NULL);
    assert(tablet != NULL);
    assert(ret != NULL);

    sky_lua_map_reduce_message *message = (sky_lua_map_reduce_message*)worker->data;
    
    // Load Lua.
    L = luaL_newstate(); check_mem(L);
    luaL_openlibs(L);

    // Compile lua script.
    rc = luaL_loadstring(L, bdata(message->source));
    check(rc == 0, "Unable to compile Lua script");

    // Execute the script.
    rc = lua_pcall(L, 0, 0, 0);
    check(rc == 0, "Unable to execute Lua script");

    // Close Lua.
    lua_close(L);
    
    sky_cursor_uninit(&cursor);
    return 0;

error:
    *ret = NULL;
    sky_cursor_uninit(&cursor);
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
    assert(worker != NULL);
    assert(data != NULL);
    
    // Ease-of-use references.
    //sky_lua_map_reduce_message *message = (sky_lua_map_reduce_message*)worker->data;

    // TODO: Merge results.
    
    return 0;
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
    check(minipack_fwrite_nil(output, &sz) == 0, "Unable to write data value");
    
    // Write total number of events to log.
    printf("[lua::map_reduce] events: %lld\n", message->event_count);
    
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

