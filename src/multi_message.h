#ifndef _sky_multi_message_h
#define _sky_multi_message_h

#include <inttypes.h>
#include <stdbool.h>
#include <netinet/in.h>

#include "bstring.h"
#include "message_header.h"
#include "message_handler.h"
#include "server.h"
#include "table.h"


//==============================================================================
//
// Typedefs
//
//==============================================================================

typedef struct sky_multi_message_data sky_multi_message_data;

// A message for processing multiple messages on the server at once.
typedef struct sky_multi_message {
    uint32_t message_count;
} sky_multi_message;


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

sky_multi_message *sky_multi_message_create();

void sky_multi_message_free(sky_multi_message *message);

//--------------------------------------
// Message Handler
//--------------------------------------

sky_message_handler *sky_multi_message_handler_create();

int sky_multi_message_process(sky_server *server,
    sky_message_header *header, sky_table *table, FILE *input, FILE *output);

//--------------------------------------
// Serialization
//--------------------------------------

size_t sky_multi_message_sizeof(sky_multi_message *message);

int sky_multi_message_pack(sky_multi_message *message, FILE *file);

int sky_multi_message_unpack(sky_multi_message *message, FILE *file);

#endif
