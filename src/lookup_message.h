#ifndef _sky_lookup_message_h
#define _sky_lookup_message_h

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

// A message for looking up actions and properties by name.
typedef struct {
    bstring *action_names;
    uint32_t action_name_count;
    bstring *property_names;
    uint32_t property_name_count;
} sky_lookup_message;


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

sky_lookup_message *sky_lookup_message_create();

void sky_lookup_message_free(sky_lookup_message *message);

//--------------------------------------
// Message Handler
//--------------------------------------

sky_message_handler *sky_lookup_message_handler_create();

int sky_lookup_message_process(sky_server *server,
    sky_message_header *header, sky_table *table, FILE *input, FILE *output);

//--------------------------------------
// Serialization
//--------------------------------------

int sky_lookup_message_pack(sky_lookup_message *message, FILE *file);

int sky_lookup_message_unpack(sky_lookup_message *message, FILE *file);

#endif
