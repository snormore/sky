#ifndef _sky_add_property_message_h
#define _sky_add_property_message_h

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

// A message for adding properties to a table.
typedef struct {
    sky_property* property;
} sky_add_property_message;


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

sky_add_property_message *sky_add_property_message_create();

void sky_add_property_message_free(sky_add_property_message *message);

void sky_add_property_message_free_property(sky_add_property_message *message);

//--------------------------------------
// Serialization
//--------------------------------------

int sky_add_property_message_pack(sky_add_property_message *message, FILE *file);

int sky_add_property_message_unpack(sky_add_property_message *message, FILE *file);

//--------------------------------------
// Processing
//--------------------------------------

int sky_add_property_message_process(sky_add_property_message *message,
    sky_table *table, FILE *output);

#endif
