#ifndef _sky_add_event_message_h
#define _sky_add_event_message_h

#include <inttypes.h>
#include <stdbool.h>
#include <netinet/in.h>

#include "bstring.h"
#include "message_header.h"
#include "message_handler.h"
#include "server.h"
#include "table.h"
#include "tablet.h"
#include "event.h"
#include "worker.h"


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
    bstring action_name;
    uint32_t data_count;
    sky_add_event_message_data **data;
    uint32_t action_data_count;
    sky_add_event_message_data **action_data;
    sky_event *event;
} sky_add_event_message;

// A key/value used to store event data.
struct sky_add_event_message_data {
    bstring key;
    sky_data_type_e data_type;
    union {
        bool boolean_value;
        int64_t int_value;
        double double_value;
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
// Message Handler
//--------------------------------------

sky_message_handler *sky_add_event_message_handler_create();

int sky_add_event_message_process(sky_server *server,
    sky_message_header *header, sky_table *table, FILE *input, FILE *output);

//--------------------------------------
// Serialization
//--------------------------------------

size_t sky_add_event_message_sizeof(sky_add_event_message *message);

int sky_add_event_message_pack(sky_add_event_message *message, FILE *file);

int sky_add_event_message_unpack(sky_add_event_message *message, FILE *file);

//--------------------------------------
// Worker
//--------------------------------------

int sky_add_event_message_worker_map(sky_worker *worker, sky_tablet *tablet,
    void **data);

int sky_add_event_message_worker_write(sky_worker *worker, FILE *output);

int sky_add_event_message_worker_free(sky_worker *worker);


#endif
