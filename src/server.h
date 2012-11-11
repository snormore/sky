#ifndef _sky_server_h
#define _sky_server_h

#include <stdio.h>
#include <inttypes.h>
#include <stdbool.h>
#include <netinet/in.h>

typedef struct sky_server sky_server;

#include "bstring.h"
#include "servlet.h"
#include "table.h"
#include "event.h"
#include "message_handler.h"


//==============================================================================
//
// Definitions
//
//==============================================================================

#define SKY_DEFAULT_PORT 8585

#define SKY_LISTEN_BACKLOG 511

#define SKY_SERVER_SHUTDOWN_URI "inproc://server.shutdown"


//==============================================================================
//
// Typedefs
//
//==============================================================================

// The various states that the server can be in.
typedef enum {
    SKY_SERVER_STATE_STOPPED = 0,
    SKY_SERVER_STATE_RUNNING = 1,
} sky_server_state_e;

struct sky_server {
    sky_server_state_e state;
    uint32_t id;
    bstring path;
    int port;
    struct sockaddr_in* sockaddr;
    int socket;
    sky_servlet **servlets;
    uint32_t servlet_count;
    uint32_t next_servlet_id;
    sky_table **tables;
    uint32_t table_count;
    sky_message_handler **message_handlers;
    uint32_t message_handler_count;
    void *context;
};



//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

sky_server *sky_server_create(bstring path);

void sky_server_free(sky_server *server);

//--------------------------------------
// State
//--------------------------------------

int sky_server_start(sky_server *server);

int sky_server_stop(sky_server *server);

int sky_server_accept(sky_server *server);

//--------------------------------------
// Servlet Management
//--------------------------------------

int sky_server_get_tablet_servlet(sky_server *server, sky_tablet *tablet,
    sky_servlet **servlet);

int sky_server_get_table_servlets(sky_server *server, sky_table *table,
    sky_servlet ***servlets, uint32_t *servlet_count);

//--------------------------------------
// Message Processing
//--------------------------------------

int sky_server_process_message(sky_server *server, FILE *input, FILE *output);

//--------------------------------------
// Message Handlers
//--------------------------------------

int sky_server_get_message_handler(sky_server *server, bstring name,
    sky_message_handler **ret);

int sky_server_add_message_handler(sky_server *server,
    sky_message_handler *handler);

int sky_server_remove_message_handler(sky_server *server,
    sky_message_handler *handler);

int sky_server_add_default_message_handlers(sky_server *server);

//--------------------------------------
// Multi Message
//--------------------------------------

int sky_server_process_multi_message(sky_server *server, FILE *input,
    FILE *output);

#endif
