#include <stdlib.h>
#include <stdio.h>
#include <arpa/inet.h>
#include <sys/time.h>
#include <assert.h>

#include "types.h"
#include "next_actions_message.h"
#include "path_iterator.h"
#include "action.h"
#include "minipack.h"
#include "mem.h"
#include "dbg.h"


//==============================================================================
//
// Definitions
//
//==============================================================================

//--------------------------------------
// Structs
//--------------------------------------

// The data object to track event data.
typedef struct {
    sky_timestamp_t timestamp;
    sky_action_id_t action_id;
} sky_next_actions_data;


//--------------------------------------
// String Constants
//--------------------------------------

#define SKY_NEXT_ACTIONS_KEY_COUNT 4

struct tagbstring SKY_NEXT_ACTIONS_STATUS_STR = bsStatic("status");
struct tagbstring SKY_NEXT_ACTIONS_OK_STR     = bsStatic("ok");
struct tagbstring SKY_NEXT_ACTIONS_DATA_STR   = bsStatic("data");
struct tagbstring SKY_NEXT_ACTIONS_COUNT_STR  = bsStatic("count");


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

// Creates a 'Next Action' message object.
//
// Returns a new message.
sky_next_actions_message *sky_next_actions_message_create()
{
    sky_next_actions_message *message = NULL;
    message = calloc(1, sizeof(sky_next_actions_message)); check_mem(message);
    return message;

error:
    sky_next_actions_message_free(message);
    return NULL;
}

// Frees a 'Next Action' message object from memory.
//
// message - The message object to be freed.
//
// Returns nothing.
void sky_next_actions_message_free(sky_next_actions_message *message)
{
    if(message) {
        sky_next_actions_message_free_deps(message);
        free(message);
    }
}

// Frees message object dependencies from memory.
//
// message - The message object.
//
// Returns nothing.
void sky_next_actions_message_free_deps(sky_next_actions_message *message)
{
    if(message) {
        sky_data_descriptor_free(message->data_descriptor);
        message->data_descriptor = NULL;
        free(message->results);
        message->results = NULL;
        free(message->prior_action_ids);
        message->prior_action_ids = NULL;
    }
}


//--------------------------------------
// Message Handler
//--------------------------------------

// Creates a message handler for the 'Next Actions' message.
//
// Returns a message handler.
sky_message_handler *sky_next_actions_message_handler_create()
{
    sky_message_handler *handler = sky_message_handler_create(); check_mem(handler);
    handler->scope = SKY_MESSAGE_HANDLER_SCOPE_TABLE;
    handler->name = bfromcstr("next_actions");
    handler->process = sky_next_actions_message_process;
    return handler;

error:
    sky_message_handler_free(handler);
    return NULL;
}

// Delegates processing of the 'Next Actions' message to a worker.
//
// server  - The server.
// header  - The message header.
// table   - The table the message is working against
// input   - The input file stream.
// output  - The output file stream.
//
// Returns 0 if successful, otherwise returns -1.
int sky_next_actions_message_process(sky_server *server,
                                     sky_message_header *header,
                                     sky_table *table,
                                     FILE *input, FILE *output)
{
    int rc = 0;
    sky_next_actions_message *message = NULL;
    assert(table != NULL);
    assert(header != NULL);
    assert(input != NULL);
    assert(output != NULL);
    
    // Create worker.
    sky_worker *worker = sky_worker_create(); check_mem(worker);
    worker->context = server->context;
    worker->read = sky_next_actions_message_worker_read;
    worker->map = sky_next_actions_message_worker_map;
    worker->map_free = sky_next_actions_message_worker_map_free;
    worker->reduce = sky_next_actions_message_worker_reduce;
    worker->write = sky_next_actions_message_worker_write;
    worker->free = sky_next_actions_message_worker_free;

    // Attach workers.
    worker->input = input;
    worker->output = output;
    
    // Attach servlets.
    rc = sky_server_get_table_servlets(server, table, &worker->servlets, &worker->servlet_count);
    check(rc == 0, "Unable to copy servlets to worker");

    // Create a message object.
    message = sky_next_actions_message_create(); check_mem(message);
    message->action_count = table->action_file->action_count;
    message->results = calloc(message->action_count+1, sizeof(*message->results));
    check_mem(message->results);

    // Initialize data descriptor.
    rc = sky_next_actions_message_init_data_descriptor(message, table->property_file);
    check(rc == 0, "Unable to initialize data descriptor");
    
    // Attach message to worker.
    worker->data = (sky_next_actions_message*)message;
    
    // Start worker.
    rc = sky_worker_start(worker);
    check(rc == 0, "Unable to start worker");
    
    return 0;

error:
    sky_next_actions_message_free(message);
    sky_worker_free(worker);
    return -1;
}

// Initializes a data descriptor on 'Next Actions' message.
//
// message       - The message.
// property_file - The property file to initialize from.
//
// Returns 0 if successful, otherwise returns -1.
int sky_next_actions_message_init_data_descriptor(sky_next_actions_message *message,
                                                  sky_property_file *property_file)
{
    sky_data_descriptor *descriptor = NULL;
    assert(message != NULL);
    assert(property_file != NULL);
    assert(message->data_descriptor == NULL);
    
    // Create data descriptor.
    descriptor = sky_data_descriptor_create(); check_mem(descriptor);
    descriptor->data_sz = (uint32_t)sizeof(sky_next_actions_data);
    descriptor->timestamp_descriptor.offset = offsetof(sky_next_actions_data, timestamp);
    descriptor->action_descriptor.offset = offsetof(sky_next_actions_data, action_id);
    
    // Assign it to message.
    message->data_descriptor = descriptor;

    return 0;

error:
    message->data_descriptor = NULL;
    sky_data_descriptor_free(descriptor);
    return -1;
}


//--------------------------------------
// Serialization
//--------------------------------------

// Calculates the total number of bytes needed to store the message.
//
// message - The message.
//
// Returns the number of bytes required to store the message.
size_t sky_next_actions_message_sizeof(sky_next_actions_message *message)
{
    size_t sz = 0;
    sz += minipack_sizeof_array(message->prior_action_id_count);

    uint32_t i;
    for(i=0; i<message->prior_action_id_count; i++) {
        sz += minipack_sizeof_int(message->prior_action_ids[i]);
    }

    return sz;
}

// Serializes a 'next_actions' message to a file stream.
//
// message - The message.
// file    - The file stream to write to.
//
// Returns 0 if successful, otherwise returns -1.
int sky_next_actions_message_pack(sky_next_actions_message *message, FILE *file)
{
    size_t sz;
    assert(message != NULL);
    assert(file != NULL);

    minipack_fwrite_array(file, message->prior_action_id_count, &sz);
    check(sz > 0, "Unable to pack prior action id array");

    uint32_t i;
    for(i=0; i<message->prior_action_id_count; i++) {
        minipack_fwrite_int(file, message->prior_action_ids[i], &sz);
        check(sz > 0, "Unable to pack prior action id");
    }
    
    return 0;

error:
    return -1;
}

// Deserializes an 'next_actions' message from a file stream.
//
// message - The message.
// file    - The file stream to read from.
//
// Returns 0 if successful, otherwise returns -1.
int sky_next_actions_message_unpack(sky_next_actions_message *message, FILE *file)
{
    size_t sz;
    assert(message != NULL);
    assert(file != NULL);

    message->prior_action_id_count = minipack_fread_array(file, &sz);
    check(sz > 0, "Unable to unpack prior action id array");

    message->prior_action_ids = realloc(message->prior_action_ids, sizeof(*message->prior_action_ids) * message->prior_action_id_count);
    check_mem(message->prior_action_ids);
    
    uint32_t i;
    for(i=0; i<message->prior_action_id_count; i++) {
        message->prior_action_ids[i] = minipack_fread_uint(file, &sz);
        check(sz > 0, "Unable to unpack prior action id");
    }

    return 0;

error:
    return -1;
}


//--------------------------------------
// Worker
//--------------------------------------

// Reads the message from the file stream.
//
// worker - The worker.
// input  - The input stream.
//
// Returns 0 if successful, otherwise returns -1.
int sky_next_actions_message_worker_read(sky_worker *worker, FILE *input)
{
    int rc;
    assert(worker != NULL);
    assert(input != NULL);

    // Parse message.
    sky_next_actions_message *message = (sky_next_actions_message*)worker->data;
    rc = sky_next_actions_message_unpack(message, input);
    check(rc == 0, "Unable to unpack 'next_actions' message");
    check(message->prior_action_id_count > 0, "Prior actions must be specified");

    return 0;

error:
    return -1;
}

// Maps tablet data to a next action summation data structure.
//
// worker - The worker.
// tablet - The tablet to work against.
// ret    - A pointer to where the summation data structure should be returned.
//
// Returns 0 if successful, otherwise returns -1.
int sky_next_actions_message_worker_map(sky_worker *worker, sky_tablet *tablet,
                                        void **ret)
{
    int rc;
    assert(worker != NULL);
    assert(tablet != NULL);
    assert(ret != NULL);

    sky_next_actions_message *message = (sky_next_actions_message*)worker->data;

    // Initialize data object.
    sky_next_actions_data data;
    memset(&data, 0, sizeof(data));
    
    // Create an array to store data.
    uint32_t action_count = message->action_count;
    sky_next_actions_result *results = calloc(action_count+1, sizeof(*results));
    check_mem(results);

    // Initialize the path iterator.
    sky_path_iterator iterator;
    sky_path_iterator_init(&iterator);

    // Attach data and descriptor to cursor.
    iterator.cursor.data_descriptor = message->data_descriptor;
    iterator.cursor.data = (void*)(&data);

    rc = sky_path_iterator_set_tablet(&iterator, tablet);
    check(rc == 0, "Unable to initialize path iterator");

    // Iterate over each path.
    uint64_t event_count = 0;
    while(!iterator.eof) {
        // Loop over each event in the path.
        uint32_t prior_action_index = 0;
        while(!iterator.cursor.eof) {
            // Aggregate if we've reached the match.
            if(prior_action_index == message->prior_action_id_count) {
                if(data.action_id <= action_count) {
                    results[data.action_id].count++;
                }
                prior_action_index = 0;
            }

            // Match against action list.
            if(message->prior_action_ids[prior_action_index] == data.action_id) {
                prior_action_index++;
            }
            else {
                prior_action_index = 0;
            }

            // Find next event.
            rc = sky_cursor_next(&iterator.cursor);
            check(rc == 0, "Unable to find next event");
            
            // Increment event count.
            event_count++;
        }

        // Move to next path.
        rc = sky_path_iterator_next(&iterator);
        check(rc == 0, "Unable to find next path");
    }

    // HACK: Increment the total event count. Note that this is not thread
    // safe however this number is only meant for debugging.
    message->event_count += event_count;

    // Return data.
    *ret = (void*)results;

    sky_path_iterator_uninit(&iterator);
    return 0;

error:
    free(results);
    *ret = NULL;
    sky_path_iterator_uninit(&iterator);
    return -1;
}

// Frees the data structure created and returned in the map() function.
//
// data - A pointer to the data to be freed.
//
// Returns 0 if successful, otherwise returns -1.
int sky_next_actions_message_worker_map_free(void *data)
{
    assert(data != NULL);
    free(data);
    return 0;
}

// Combines the data from a single execution of the map() function into data
// saved against the worker.
//
// worker - The worker.
// data   - Data created and returned in the map() function.
//
// Returns 0 if successful, otherwise returns -1.
int sky_next_actions_message_worker_reduce(sky_worker *worker, void *data)
{
    assert(worker != NULL);
    assert(data != NULL);
    
    // Ease-of-use references.
    sky_next_actions_message *message = (sky_next_actions_message*)worker->data;
    sky_next_actions_result *map_results = (sky_next_actions_result*)data;

    // Merge results.
    uint32_t i;
    for(i=0; i<message->action_count+1; i++) {
        message->results[i].count += map_results[i].count;
    }
    
    return 0;
}

// Writes the results to an output stream.
//
// worker - The worker.
// output - The output stream.
//
// Returns 0 if successful, otherwise returns -1.
int sky_next_actions_message_worker_write(sky_worker *worker, FILE *output)
{
    size_t sz;
    assert(worker != NULL);
    assert(output != NULL);
    
    // Ease-of-use references.
    sky_next_actions_message *message = (sky_next_actions_message*)worker->data;

    // Count the total number of return elements.
    uint32_t i, key_count = 0;
    for(i=0; i<message->action_count+1; i++) {
        if(message->results[i].count > 0) {
            key_count++;
        }
    }
    
    // Return.
    //   {status:"ok", data:{<action_id>:{count:0}, ...}}
    check(minipack_fwrite_map(output, 2, &sz) == 0, "Unable to write root map");
    check(sky_minipack_fwrite_bstring(output, &SKY_NEXT_ACTIONS_STATUS_STR) == 0, "Unable to write status key");
    check(sky_minipack_fwrite_bstring(output, &SKY_NEXT_ACTIONS_OK_STR) == 0, "Unable to write status value");
    check(sky_minipack_fwrite_bstring(output, &SKY_NEXT_ACTIONS_DATA_STR) == 0, "Unable to write data key");
    check(minipack_fwrite_map(output, key_count, &sz) == 0, "Unable to write data key");
    for(i=0; i<message->action_count+1; i++) {
        if(message->results[i].count > 0) {
            check(minipack_fwrite_uint(output, i, &sz) == 0, "Unable to write action id");
            check(minipack_fwrite_map(output, 1, &sz) == 0, "Unable to write result map");
            check(sky_minipack_fwrite_bstring(output, &SKY_NEXT_ACTIONS_COUNT_STR) == 0, "Unable to write result count key");
            check(minipack_fwrite_uint(output, message->results[i].count, &sz) == 0, "Unable to write result count");
        }
    }
    
    // Write total number of events to log.
    printf("[next_actions] events: %" PRIu64 "\n", message->event_count);
    
    return 0;

error:
    return -1;
}

// Frees all data attached to the worker.
//
// worker - The worker.
//
// Returns 0 if successful, otherwise returns -1.
int sky_next_actions_message_worker_free(sky_worker *worker)
{
    assert(worker != NULL);
    
    // Clean up.
    sky_next_actions_message *message = (sky_next_actions_message*)worker->data;
    sky_next_actions_message_free(message);
    worker->data = NULL;
    
    return 0;
}

