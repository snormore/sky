#ifndef _sky_add_table_message_h
#define _sky_add_table_message_h

#include <inttypes.h>
#include <stdbool.h>
#include <netinet/in.h>

#include "bstring.h"
#include "message_handler.h"
#include "table.h"
#include "event.h"


//==============================================================================
//
// Typedefs
//
//==============================================================================

// A message for adding a table to the server.
typedef struct {
    bstring data_path;
    sky_table* table;
} sky_add_table_message;


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

sky_add_table_message *sky_add_table_message_create();

void sky_add_table_message_free(sky_add_table_message *message);

//--------------------------------------
// Message Handler
//--------------------------------------

sky_message_handler *sky_add_table_message_handler_create();

int sky_add_table_message_process(sky_server *server,
    sky_message_header *header, sky_table *table, FILE *input, FILE *output);

//--------------------------------------
// Serialization
//--------------------------------------

int sky_add_table_message_pack(sky_add_table_message *message, FILE *file);

int sky_add_table_message_unpack(sky_add_table_message *message, FILE *file);

#endif
