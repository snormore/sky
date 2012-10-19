#ifndef _sky_add_action_message_h
#define _sky_add_action_message_h

#include <inttypes.h>
#include <stdbool.h>
#include <netinet/in.h>

#include "bstring.h"
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
// Serialization
//--------------------------------------

int sky_add_action_message_pack(sky_add_action_message *message, FILE *file);

int sky_add_action_message_unpack(sky_add_action_message *message, FILE *file);

//--------------------------------------
// Processing
//--------------------------------------

int sky_add_action_message_process(sky_add_action_message *message,
    sky_table *table, FILE *output);

#endif
