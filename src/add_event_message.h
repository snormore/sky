#ifndef _sky_add_event_message_h
#define _sky_add_event_message_h

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

typedef struct sky_add_event_message_data sky_add_event_message_data;

// A message for adding events to the database.
typedef struct sky_add_event_message {
    sky_object_id_t object_id;
    sky_timestamp_t timestamp;
    sky_action_id_t action_id;
    uint32_t data_count;
    sky_add_event_message_data **data;
} sky_add_event_message;

// A key/value used to store event data.
struct sky_add_event_message_data {
    bstring key;
    bstring data_type;
    union {
        bool boolean_value;
        int64_t int_value;
        double float_value;
        bstring string_value;
    };
};


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

sky_add_event_message *sky_add_event_message_create();

sky_add_event_message_data *sky_add_event_message_data_create();

void sky_add_event_message_free(sky_add_event_message *message);

void sky_add_event_message_data_free(sky_add_event_message_data *data);


//--------------------------------------
// Serialization
//--------------------------------------

size_t sky_add_event_message_sizeof(sky_add_event_message *message);

int sky_add_event_message_pack(sky_add_event_message *message, FILE *file);

int sky_add_event_message_unpack(sky_add_event_message *message, FILE *file);

//--------------------------------------
// Processing
//--------------------------------------

int sky_add_event_message_process(sky_add_event_message *message,
    sky_table *table, FILE *output);

#endif
