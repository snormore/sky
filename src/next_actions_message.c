#include <stdlib.h>
#include <stdio.h>
#include <arpa/inet.h>
#include <sys/time.h>

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

// The data object to track event data.
typedef struct {
    sky_timestamp_t timestamp;
    sky_action_id_t action_id;
} sky_next_actions_data;

// The result data to send back to the client.
typedef struct {
    uint32_t count;
} sky_next_actions_result;


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
        free(message->prior_action_ids);
        message->prior_action_ids = NULL;
    }
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
    check(message != NULL, "Message required");
    check(file != NULL, "File stream required");

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
    check(message != NULL, "Message required");
    check(file != NULL, "File stream required");

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
// Processing
//--------------------------------------

// Queries a table to determine the number of occurrences of the action
// immediately following a series of actions.
//
// message - The message.
// table   - The table to apply the message to.
// output  - The output stream to write to.
//
// Returns 0 if successful, otherwise returns -1.
int sky_next_actions_message_process(sky_next_actions_message *message,
                                     sky_table *table, FILE *output)
{
    int rc;
    uint32_t i;
    size_t sz;
    check(message != NULL, "Message required");
    check(message->prior_action_id_count > 0, "Prior actions must be specified");
    check(table != NULL, "Table required");
    check(output != NULL, "Output stream required");

    struct tagbstring status_str = bsStatic("status");
    struct tagbstring ok_str = bsStatic("ok");
    struct tagbstring data_str = bsStatic("data");
    struct tagbstring count_str = bsStatic("count");

    // Initialize data object.
    sky_next_actions_data data;
    memset(&data, 0, sizeof(data));
    
    // Initialize data descriptor.
    sky_data_descriptor *descriptor = NULL;
    rc = sky_property_file_create_data_descriptor(table->property_file, &descriptor);
    check(rc == 0, "Unable to create data descriptor");
    descriptor->timestamp_descriptor.offset = offsetof(sky_next_actions_data, timestamp);
    descriptor->action_descriptor.offset = offsetof(sky_next_actions_data, action_id);
    
    // Create an array to store data.
    uint32_t action_count = table->action_file->action_count;
    sky_next_actions_result *results = calloc(action_count+1, sizeof(*results));
    check_mem(results);

    // Initialize the path iterator.
    sky_path_iterator iterator;
    sky_path_iterator_init(&iterator);
    rc = sky_path_iterator_set_data_file(&iterator, table->data_file);
    check(rc == 0, "Unable to initialze path iterator");

    // Start benchmark.
    struct timeval tv;
    gettimeofday(&tv, NULL);
    int64_t t0 = (tv.tv_sec*1000) + (tv.tv_usec/1000);

    // Iterate over each path.
    uint64_t event_count = 0;
    while(!iterator.eof) {
        // Retrieve the path pointer.
        void *path_ptr = NULL;
        rc = sky_path_iterator_get_ptr(&iterator, &path_ptr);
        check(rc == 0, "Unable to retrieve the path iterator pointer");
    
        // Initialize the cursor.
        sky_cursor cursor;
        sky_cursor_init(&cursor);
        rc = sky_cursor_set_path(&cursor, path_ptr);
        check(rc == 0, "Unable to set cursor path");

        // Loop over each event in the path.
        uint32_t prior_action_index = 0;
        while(!cursor.eof) {
            // Retrieve action.
            rc = sky_cursor_set_data(&cursor, descriptor, (void*)(&data));
            check(rc == 0, "Unable to retrieve first action");

            // Aggregate if we've reached the match.
            if(prior_action_index == message->prior_action_id_count) {
                results[data.action_id].count++;
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
            rc = sky_cursor_next(&cursor);
            check(rc == 0, "Unable to find next event");
            
            // Increment event count.
            event_count++;
        }

        // Move to next path.
        rc = sky_path_iterator_next(&iterator);
        check(rc == 0, "Unable to find next path");
    }
    
    // End benchmark.
    gettimeofday(&tv, NULL);
    int64_t t1 = (tv.tv_sec*1000) + (tv.tv_usec/1000);
    debug("'Next Action' queried %lld events in: %.3f seconds\n", event_count, ((float)(t1-t0))/1000);
    
    // Count the total number of return elements.
    uint32_t key_count = 0;
    for(i=0; i<action_count; i++) {
        if(results[i].count > 0) {
            key_count++;
        }
    }
    
    // Return.
    //   {status:"ok", data:{<action_id>:{count:0}, ...}}
    check(minipack_fwrite_map(output, 2, &sz) == 0, "Unable to write root map");
    check(sky_minipack_fwrite_bstring(output, &status_str) == 0, "Unable to write status key");
    check(sky_minipack_fwrite_bstring(output, &ok_str) == 0, "Unable to write status value");
    check(sky_minipack_fwrite_bstring(output, &data_str) == 0, "Unable to write data key");
    check(minipack_fwrite_map(output, key_count, &sz) == 0, "Unable to write data key");
    for(i=0; i<action_count+1; i++) {
        if(results[i].count > 0) {
            check(minipack_fwrite_uint(output, i, &sz) == 0, "Unable to write action id");
            check(minipack_fwrite_map(output, 1, &sz) == 0, "Unable to write result map");
            check(sky_minipack_fwrite_bstring(output, &count_str) == 0, "Unable to write result count key");
            check(minipack_fwrite_uint(output, results[i].count, &sz) == 0, "Unable to write result count");
        }
    }
    
    free(results);
    return 0;

error:
    free(results);
    return -1;
}