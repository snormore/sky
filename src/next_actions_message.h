#ifndef _sky_next_actions_message_h
#define _sky_next_actions_message_h

#include <inttypes.h>
#include <stdbool.h>
#include <netinet/in.h>

#include "bstring.h"
#include "table.h"
#include "tablet.h"
#include "event.h"
#include "message_handler.h"
#include "worker.h"


//==============================================================================
//
// Typedefs
//
//==============================================================================

// The result data to send back to the client.
typedef struct {
    uint32_t count;
} sky_next_actions_result;

// A message for retrieving a count of the next immediate action following a
// series of actions.
typedef struct {
    sky_action_id_t *prior_action_ids;
    uint32_t prior_action_id_count;
    sky_next_actions_result *results;
    uint32_t action_count;
    sky_data_descriptor *data_descriptor;
} sky_next_actions_message;


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

sky_next_actions_message *sky_next_actions_message_create();

void sky_next_actions_message_free(sky_next_actions_message *message);

void sky_next_actions_message_free_deps(sky_next_actions_message *message);

//--------------------------------------
// Message Handler
//--------------------------------------

sky_message_handler *sky_next_actions_message_handler_create();

int sky_next_actions_message_process(sky_server *server, sky_table *table,
    FILE *input, FILE *output);

int sky_next_actions_message_init_data_descriptor(sky_next_actions_message *message,
    sky_property_file *property_file);

//--------------------------------------
// Serialization
//--------------------------------------

int sky_next_actions_message_pack(sky_next_actions_message *message,
    FILE *file);

int sky_next_actions_message_unpack(sky_next_actions_message *message,
    FILE *file);

//--------------------------------------
// Worker
//--------------------------------------

int sky_next_actions_message_worker_read(sky_worker *worker, FILE *input);

int sky_next_actions_message_worker_map(sky_worker *worker, sky_tablet *tablet,
    void **data);

int sky_next_actions_message_worker_map_free(void *data);

int sky_next_actions_message_worker_reduce(sky_worker *worker, void *data);

int sky_next_actions_message_worker_write(sky_worker *worker, FILE *output);

int sky_next_actions_message_worker_free(sky_worker *worker);

#endif
