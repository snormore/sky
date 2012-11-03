#ifndef _server_h
#define _server_h

#include <stdio.h>
#include <inttypes.h>
#include <stdbool.h>
#include <netinet/in.h>

#include "bstring.h"
#include "table.h"
#include "event.h"
#include "message_handler.h"


//==============================================================================
//
// Overview
//
//==============================================================================

// The server acts as the interface to external applications. It communicates
// over TCP sockets using a specific Sky protocol. See the message.h file for
// more detail on the protocol.


//==============================================================================
//
// Definitions
//
//==============================================================================

#define SKY_DEFAULT_PORT 8585

#define SKY_LISTEN_BACKLOG 511


//==============================================================================
//
// Typedefs
//
//==============================================================================

// The various states that the server can be in.
typedef enum sky_server_state_e {
    SKY_SERVER_STATE_STOPPED,
    SKY_SERVER_STATE_RUNNING,
} sky_server_state_e;


typedef struct {
    sky_server_state_e state;
    bstring path;
    int port;
    struct sockaddr_in* sockaddr;
    int socket;
    sky_table **tables;
    uint32_t table_count;
    sky_message_handler **message_handlers;
    uint32_t message_handler_count;
} sky_server;



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
// Query Messages
//--------------------------------------

int sky_server_process_next_actions_message(sky_server *server,
    sky_table *table, FILE *input, FILE *output);

//--------------------------------------
// Action Messages
//--------------------------------------

int sky_server_process_add_action_message(sky_server *server,
    sky_table *table, FILE *input, FILE *output);

int sky_server_process_get_action_message(sky_server *server,
    sky_table *table, FILE *input, FILE *output);

int sky_server_process_get_actions_message(sky_server *server,
    sky_table *table, FILE *input, FILE *output);

//--------------------------------------
// Property Messages
//--------------------------------------

int sky_server_process_add_property_message(sky_server *server,
    sky_table *table, FILE *input, FILE *output);

int sky_server_process_get_property_message(sky_server *server,
    sky_table *table, FILE *input, FILE *output);

int sky_server_process_get_properties_message(sky_server *server,
    sky_table *table, FILE *input, FILE *output);

//--------------------------------------
// Multi Message
//--------------------------------------

int sky_server_process_multi_message(sky_server *server, FILE *input,
    FILE *output);

#endif
