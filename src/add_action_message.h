#ifndef _sky_add_action_message_h
#define _sky_add_action_message_h

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

// A message for adding actions to a table.
typedef struct {
    sky_action* action;
} sky_add_action_message;


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

sky_add_action_message *sky_add_action_message_create();

void sky_add_action_message_free(sky_add_action_message *message);

void sky_add_action_message_free_action(sky_add_action_message *message);

//--------------------------------------
// Message Handler
//--------------------------------------

sky_message_handler *sky_add_action_message_handler_create();

int sky_add_action_message_process(sky_server *server,
    sky_message_header *header, sky_table *table, FILE *input, FILE *output);

//--------------------------------------
// Serialization
//--------------------------------------

int sky_add_action_message_pack(sky_add_action_message *message, FILE *file);

int sky_add_action_message_unpack(sky_add_action_message *message, FILE *file);

#endif
