#ifndef _sky_get_actions_message_h
#define _sky_get_actions_message_h

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

// A message for retrieving an action by id from a table.
typedef struct {
    int64_t dummy;
} sky_get_actions_message;


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

sky_get_actions_message *sky_get_actions_message_create();

void sky_get_actions_message_free(sky_get_actions_message *message);

//--------------------------------------
// Serialization
//--------------------------------------

int sky_get_actions_message_pack(sky_get_actions_message *message, FILE *file);

int sky_get_actions_message_unpack(sky_get_actions_message *message, FILE *file);

//--------------------------------------
// Processing
//--------------------------------------

int sky_get_actions_message_process(sky_get_actions_message *message,
    sky_table *table, FILE *output);

#endif
