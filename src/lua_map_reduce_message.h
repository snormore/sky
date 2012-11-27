#ifndef _sky_lua_map_reduce_message_h
#define _sky_lua_map_reduce_message_h

#include <inttypes.h>
#include <stdbool.h>
#include <netinet/in.h>

#include "bstring.h"
#include "message_header.h"
#include "message_handler.h"
#include "sky_lua.h"
#include "table.h"
#include "tablet.h"
#include "event.h"
#include "worker.h"


//==============================================================================
//
// Typedefs
//
//==============================================================================

// A message for executing a distributed map/reduce lua script across a table.
typedef struct {
    bstring results;
    bstring source;
    lua_State *L;
} sky_lua_map_reduce_message;


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

sky_lua_map_reduce_message *sky_lua_map_reduce_message_create();

void sky_lua_map_reduce_message_free(sky_lua_map_reduce_message *message);


//--------------------------------------
// Message Handler
//--------------------------------------

sky_message_handler *sky_lua_map_reduce_message_handler_create();

int sky_lua_map_reduce_message_process(sky_server *server,
    sky_message_header *header, sky_table *table, FILE *input, FILE *output);

//--------------------------------------
// Serialization
//--------------------------------------

int sky_lua_map_reduce_message_pack(sky_lua_map_reduce_message *message,
    FILE *file);

int sky_lua_map_reduce_message_unpack(sky_lua_map_reduce_message *message,
    FILE *file);

//--------------------------------------
// Worker
//--------------------------------------

int sky_lua_map_reduce_message_worker_map(sky_worker *worker,
    sky_tablet *tablet, void **data);

int sky_lua_map_reduce_message_worker_map_free(void *data);

int sky_lua_map_reduce_message_worker_reduce(sky_worker *worker, void *data);

int sky_lua_map_reduce_message_worker_write(sky_worker *worker, FILE *output);

int sky_lua_map_reduce_message_worker_free(sky_worker *worker);

#endif
