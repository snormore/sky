#ifndef _sky_get_property_message_h
#define _sky_get_property_message_h

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

// A message for retrieving a property by id from a table.
typedef struct {
    sky_property_id_t property_id;
} sky_get_property_message;


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

sky_get_property_message *sky_get_property_message_create();

void sky_get_property_message_free(sky_get_property_message *message);

//--------------------------------------
// Message Handler
//--------------------------------------

sky_message_handler *sky_get_property_message_handler_create();

int sky_get_property_message_process(sky_server *server,
    sky_message_header *header, sky_table *table, FILE *input, FILE *output);

//--------------------------------------
// Serialization
//--------------------------------------

int sky_get_property_message_pack(sky_get_property_message *message,
    FILE *file);

int sky_get_property_message_unpack(sky_get_property_message *message,
    FILE *file);

#endif
