#include <stdlib.h>
#include <stdio.h>
#include <arpa/inet.h>
#include <assert.h>

#include "types.h"
#include "add_event_message.h"
#include "worker.h"
#include "timestamp.h"
#include "minipack.h"
#include "mem.h"
#include "dbg.h"


//==============================================================================
//
// Definitions
//
//==============================================================================

#define SKY_ADD_EVENT_KEY_COUNT 4

struct tagbstring SKY_ADD_EVENT_KEY_OBJECT_ID = bsStatic("objectId");

struct tagbstring SKY_ADD_EVENT_KEY_TIMESTAMP = bsStatic("timestamp");

struct tagbstring SKY_ADD_EVENT_KEY_ACTION = bsStatic("action");

struct tagbstring SKY_ADD_EVENT_KEY_NAME = bsStatic("name");

struct tagbstring SKY_ADD_EVENT_KEY_DATA = bsStatic("data");


//==============================================================================
//
// Forward Declarations
//
//==============================================================================

int sky_add_event_message_copy_data(sky_property_file *property_file,
    bool is_action, sky_add_event_message_data **msg_data, uint32_t msg_data_count,
    sky_event_data ***event_data, uint32_t *event_data_count);

size_t sky_add_event_message_sizeof_action(sky_add_event_message *message);

size_t sky_add_event_message_sizeof_data(sky_add_event_message *message);

int sky_add_event_message_pack_action(sky_add_event_message *message, FILE *file);

int sky_add_event_message_pack_data(sky_add_event_message *message, FILE *file);

int sky_add_event_message_unpack_action(sky_add_event_message *message, FILE *file);

int sky_add_event_message_unpack_data(sky_add_event_message *message, FILE *file);


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

// Creates an 'add_event' message object.
//
// Returns a new message.
sky_add_event_message *sky_add_event_message_create()
{
    sky_add_event_message *message = NULL;
    message = calloc(1, sizeof(sky_add_event_message)); check_mem(message);
    return message;

error:
    sky_add_event_message_free(message);
    return NULL;
}

// Creates an 'add_event' message data object.
//
// Returns a new message data.
sky_add_event_message_data *sky_add_event_message_data_create()
{
    sky_add_event_message_data *data = NULL;
    data = calloc(1, sizeof(sky_add_event_message_data)); check_mem(data);
    return data;

error:
    sky_add_event_message_data_free(data);
    return NULL;
}

// Frees an 'add_event' message object from memory.
//
// message - The message object to be freed.
//
// Returns nothing.
void sky_add_event_message_free(sky_add_event_message *message)
{
    if(message) {
        uint32_t i;
        for(i=0; i<message->data_count; i++) {
            sky_add_event_message_data_free(message->data[i]);
            message->data[i] = NULL;
        }
        free(message->data);
        message->data = NULL;
        for(i=0; i<message->action_data_count; i++) {
            sky_add_event_message_data_free(message->action_data[i]);
            message->action_data[i] = NULL;
        }
        free(message->action_data);
        message->action_data = NULL;

        bdestroy(message->action_name);
        message->action_name = NULL;

        sky_event_free(message->event);
        message->event = NULL;
        
        free(message);
    }
}

// Frees an 'add_event' message data object from memory.
//
// data - The message data object to be freed.
//
// Returns nothing.
void sky_add_event_message_data_free(sky_add_event_message_data *data)
{
    if(data) {
        bdestroy(data->key);
        data->key = NULL;
        if(data->data_type == SKY_DATA_TYPE_STRING) {
            bdestroy(data->string_value);
            data->string_value = NULL;
        }
        data->data_type = SKY_DATA_TYPE_NONE;
        free(data);
    }
}


//--------------------------------------
// Message Handler
//--------------------------------------

// Creates a message handler for the 'Add Event' message.
//
// Returns a message handler.
sky_message_handler *sky_add_event_message_handler_create()
{
    sky_message_handler *handler = sky_message_handler_create(); check_mem(handler);
    handler->scope = SKY_MESSAGE_HANDLER_SCOPE_OBJECT;
    handler->name = bfromcstr("add_event");
    handler->process = sky_add_event_message_process;
    return handler;

error:
    sky_message_handler_free(handler);
    return NULL;
}

// Delegates processing of the 'Add Event' message to a worker.
//
// server - The server.
// header - The message header.
// table  - The table the message is working against
// input  - The input file stream.
// output - The output file stream.
//
// Returns 0 if successful, otherwise returns -1.
int sky_add_event_message_process(sky_server *server,
                                  sky_message_header *header,
                                  sky_table *table, 
                                  FILE *input, FILE *output)
{
    int rc = 0;
    sky_add_event_message *message = NULL;
    sky_worker *worker = NULL;
    sky_action *action = NULL;
    assert(header != NULL);
    assert(table != NULL);
    assert(input != NULL);
    assert(output != NULL);
    
    // Create worker.
    worker = sky_worker_create(); check_mem(worker);
    worker->context = server->context;
    worker->map = sky_add_event_message_worker_map;
    worker->write = sky_add_event_message_worker_write;
    worker->free = sky_add_event_message_worker_free;
    worker->multi = header->multi;
    worker->input = input;
    worker->output = output;
    
    // Parse message.
    message = sky_add_event_message_create(); check_mem(message);
    rc = sky_add_event_message_unpack(message, input);
    check(rc == 0, "Unable to unpack 'add_event' message");
    check(message->object_id > 0, "Object ID must be greater than zero");

    // TODO: Reject action if there is no 'name' key.

    // Find or create action.
    rc = sky_action_file_find_by_name(table->action_file, message->action_name, &action);
    check(rc == 0, "Unable to search for action by name");
    
    // If action doesn't exist then create it.
    if(action == NULL) {
        action = sky_action_create(); check_mem(action);
        action->name = bstrcpy(message->action_name);
        if(message->action_name) check_mem(action->name);

        // Add action.
        rc = sky_action_file_add_action(table->action_file, action);
        check(rc == 0, "Unable to add action");

        // Save action file.
        rc = sky_action_file_save(table->action_file);
        check(rc == 0, "Unable to save action file");
    }

    // Create event object and attach it to the message.
    message->event = sky_event_create(message->object_id, message->timestamp, action->id);
    check_mem(message->event);
    
    // Copy action data from message.
    rc = sky_add_event_message_copy_data(table->property_file, true, message->action_data, message->action_data_count, &message->event->data, &message->event->data_count);
    check(rc == 0, "Unable to copy event action data");
    
    // Copy data from message.
    rc = sky_add_event_message_copy_data(table->property_file, false, message->data, message->data_count, &message->event->data, &message->event->data_count);
    check(rc == 0, "Unable to copy event data");
    
    // Attach the message to the worker.
    worker->data = (void*)message;
    
    // Attach servlets.
    worker->servlets = calloc(1, sizeof(*worker->servlets)); check_mem(worker->servlets);
    worker->servlet_count = 1;
    sky_tablet *tablet = NULL;
    rc = sky_table_get_target_tablet(table, message->object_id, &tablet);
    check(rc == 0 && tablet != NULL, "Unable to find target tablet: %d", message->object_id);
    rc = sky_server_get_tablet_servlet(server, tablet, &worker->servlets[0]);
    check(rc == 0, "Unable to copy servlet to worker");

    // Start worker.
    rc = sky_worker_start(worker);
    check(rc == 0, "Unable to start worker");

    return 0;

error:
    sky_add_event_message_free(message);
    sky_worker_free(worker);
    if(action && action->action_file == NULL) sky_action_free(action);
    return -1;
}

// Copies a set of data from the message to the event.
//
// property_file    - The property file.
// is_action        - A flag stating if action data is being copied.
// msg_data         - An array of message data.
// msg_data_count   - The number of message data items.
// event_data       - A pointer to where the event data should be returned.
// event_data_count - A pointer to where the event data count should be
//                    returned.
//
// Returns 0 if successful, otherwise returns -1.
int sky_add_event_message_copy_data(sky_property_file *property_file,
                                    bool is_action,
                                    sky_add_event_message_data **msg_data,
                                    uint32_t msg_data_count,
                                    sky_event_data ***event_data,
                                    uint32_t *event_data_count)
{
    int rc = 0;
    sky_property *property = NULL;
    assert(property_file != NULL);
    assert(msg_data != NULL || msg_data_count == 0);
    assert(event_data != NULL);
    assert(event_data_count != NULL);

    // Allocate event data space.
    size_t sz = (*event_data_count + msg_data_count) * sizeof(**event_data);
    if(sz > 0) {
        *event_data = realloc(*event_data, sz);
        check_mem(*event_data);
    }
    else {
        if(*event_data) free(*event_data);
        *event_data = NULL;
    }

    // Copy message data to event data.
    uint32_t i;
    for(i=0; i<msg_data_count; i++) {
        sky_event_data *event_item = NULL;
        sky_add_event_message_data *msg_item = msg_data[i];
        
        // Look up property id by name.
        rc = sky_property_file_find_by_name(property_file, msg_item->key, &property);
        check(rc == 0, "Unable to find property '%s'", bdata(msg_item->key));
        
        // Create property if it doesn't exist.
        if(property == NULL) {
            property = sky_property_create(); check_mem(property);
            property->data_type = msg_item->data_type;
            property->type = (is_action ? SKY_PROPERTY_TYPE_ACTION : SKY_PROPERTY_TYPE_OBJECT);
            property->name = bstrcpy(msg_item->key); check_mem(property->name);

            // Add property.
            rc = sky_property_file_add_property(property_file, property);
            check(rc == 0, "Unable to add property");
            
            // Save property file.
            rc = sky_property_file_save(property_file);
            check(rc == 0, "Unable to save property file");
        }
        
        // Create event data based on data type.
        switch(msg_item->data_type) {
            case SKY_DATA_TYPE_STRING:
                event_item = sky_event_data_create_string(property->id, msg_item->string_value);
                break;
            case SKY_DATA_TYPE_INT:
                event_item = sky_event_data_create_int(property->id, msg_item->int_value);
                break;
            case SKY_DATA_TYPE_DOUBLE:
                event_item = sky_event_data_create_double(property->id, msg_item->double_value);
                break;
            case SKY_DATA_TYPE_BOOLEAN:
                event_item = sky_event_data_create_boolean(property->id, msg_item->boolean_value);
                break;
            default:
                sentinel("Invalid data type in 'add_event' message");
        }
        
        // Append event item.
        (*event_data)[*event_data_count] = event_item;
        (*event_data_count)++;
    }

    return 0;

error:
    if(property->property_file == NULL) sky_property_free(property);
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
size_t sky_add_event_message_sizeof(sky_add_event_message *message)
{
    size_t sz = 0;
    sz += minipack_sizeof_map(SKY_ADD_EVENT_KEY_COUNT);
    sz += minipack_sizeof_raw((&SKY_ADD_EVENT_KEY_OBJECT_ID)->slen) + (&SKY_ADD_EVENT_KEY_OBJECT_ID)->slen;
    sz += minipack_sizeof_uint(message->object_id);
    sz += minipack_sizeof_raw((&SKY_ADD_EVENT_KEY_TIMESTAMP)->slen) + (&SKY_ADD_EVENT_KEY_TIMESTAMP)->slen;
    sz += minipack_sizeof_int(message->timestamp);
    sz += minipack_sizeof_raw((&SKY_ADD_EVENT_KEY_ACTION)->slen) + (&SKY_ADD_EVENT_KEY_ACTION)->slen;
    sz += sky_add_event_message_sizeof_action(message);
    sz += minipack_sizeof_raw((&SKY_ADD_EVENT_KEY_DATA)->slen) + (&SKY_ADD_EVENT_KEY_DATA)->slen;
    sz += sky_add_event_message_sizeof_data(message);
    return sz;
}

// Calculates the total number of bytes needed to store the action property of
// the message.
//
// message - The message.
//
// Returns the number of bytes required to store the data property of the
// message.
size_t sky_add_event_message_sizeof_action(sky_add_event_message *message)
{
    uint32_t i;
    size_t sz = 0;
    sz += minipack_sizeof_map(message->action_data_count + 1);

    // Action name.
    sz += minipack_sizeof_raw((&SKY_ADD_EVENT_KEY_NAME)->slen) + (&SKY_ADD_EVENT_KEY_NAME)->slen;
    sz += minipack_sizeof_raw(blength(message->action_name)) + blength(message->action_name);

    // Action data.
    for(i=0; i<message->action_data_count; i++) {
        sky_add_event_message_data *data = message->action_data[i];
        sz += minipack_sizeof_raw(blength(data->key)) + blength(data->key);
        
        switch(data->data_type) {
            case SKY_DATA_TYPE_STRING:
                sz += minipack_sizeof_raw(blength(data->string_value)) + blength(data->string_value);
                break;
            case SKY_DATA_TYPE_INT:
                sz += minipack_sizeof_int(data->int_value);
                break;
            case SKY_DATA_TYPE_DOUBLE:
                sz += minipack_sizeof_double(data->double_value);
                break;
            case SKY_DATA_TYPE_BOOLEAN:
                sz += minipack_sizeof_bool(data->boolean_value);
                break;
            default:
                return 0;
        }
    }
    return sz;
}

// Calculates the total number of bytes needed to store the data property of
// the message.
//
// message - The message.
//
// Returns the number of bytes required to store the data property of the
// message.
size_t sky_add_event_message_sizeof_data(sky_add_event_message *message)
{
    uint32_t i;
    size_t sz = 0;
    sz += minipack_sizeof_map(message->data_count);
    for(i=0; i<message->data_count; i++) {
        sky_add_event_message_data *data = message->data[i];
        sz += minipack_sizeof_raw(blength(data->key)) + blength(data->key);
        
        switch(data->data_type) {
            case SKY_DATA_TYPE_STRING:
                sz += minipack_sizeof_raw(blength(data->string_value)) + blength(data->string_value);
                break;
            case SKY_DATA_TYPE_INT:
            sz += minipack_sizeof_int(data->int_value);
                break;
            case SKY_DATA_TYPE_DOUBLE:
                sz += minipack_sizeof_double(data->double_value);
                break;
            case SKY_DATA_TYPE_BOOLEAN:
                sz += minipack_sizeof_bool(data->boolean_value);
                break;
            default:
                return 0;
        }
    }
    return sz;
}

// Serializes an 'add_event' message to a file stream.
//
// message - The message.
// file    - The file stream to write to.
//
// Returns 0 if successful, otherwise returns -1.
int sky_add_event_message_pack(sky_add_event_message *message, FILE *file)
{
    int rc;
    size_t sz;
    assert(message != NULL);
    assert(file != NULL);

    // Map
    minipack_fwrite_map(file, SKY_ADD_EVENT_KEY_COUNT, &sz);
    check(sz > 0, "Unable to write map");
    
    // Object ID
    check(sky_minipack_fwrite_bstring(file, &SKY_ADD_EVENT_KEY_OBJECT_ID) == 0, "Unable to pack object id key");
    minipack_fwrite_int(file, message->object_id, &sz);
    check(sz != 0, "Unable to pack object id");

    // Timestamp
    check(sky_minipack_fwrite_bstring(file, &SKY_ADD_EVENT_KEY_TIMESTAMP) == 0, "Unable to pack timestamp key");
    minipack_fwrite_int(file, message->timestamp, &sz);
    check(sz != 0, "Unable to pack timestamp");

    // Action
    check(sky_minipack_fwrite_bstring(file, &SKY_ADD_EVENT_KEY_ACTION) == 0, "Unable to pack action key");
    rc = sky_add_event_message_pack_action(message, file);
    check(rc == 0, "Unable to pack 'add_event' action");
    
    // Data
    check(sky_minipack_fwrite_bstring(file, &SKY_ADD_EVENT_KEY_DATA) == 0, "Unable to pack data key");
    rc = sky_add_event_message_pack_data(message, file);
    check(rc == 0, "Unable to pack 'add_event' data");
    
    return 0;

error:
    return -1;
}

// Serializes the action map of an 'add_event' message.
//
// message - The message.
// file    - The file stream to read from.
//
// Returns 0 if successful, otherwise returns -1.
int sky_add_event_message_pack_action(sky_add_event_message *message, FILE *file)
{
    int rc;
    size_t sz;
    assert(message != NULL);
    assert(file != NULL);

    // Map
    minipack_fwrite_map(file, message->action_data_count + 1, &sz);
    check(sz > 0, "Unable to write map");
    
    // Action name
    check(sky_minipack_fwrite_bstring(file, &SKY_ADD_EVENT_KEY_NAME) == 0, "Unable to pack action name key");
    check(sky_minipack_fwrite_bstring(file, message->action_name) == 0, "Unable to pack action name");

    // Map items
    uint32_t i;
    for(i=0; i<message->action_data_count; i++) {
        sky_add_event_message_data *data = message->action_data[i];
        
        // Write key.
        rc = sky_minipack_fwrite_bstring(file, data->key);
        check(rc == 0, "Unable to pack data key");
        
        // Write in the appropriate data type.
        switch(data->data_type) {
            case SKY_DATA_TYPE_STRING:
                rc = sky_minipack_fwrite_bstring(file, data->string_value);
                check(rc == 0, "Unable to pack string value");
                break;
            case SKY_DATA_TYPE_INT:
                minipack_fwrite_int(file, data->int_value, &sz);
                check(sz > 0, "Unable to pack int value");
                break;
            case SKY_DATA_TYPE_DOUBLE:
                minipack_fwrite_double(file, data->double_value, &sz);
                check(sz > 0, "Unable to pack float value");
                break;
            case SKY_DATA_TYPE_BOOLEAN:
                minipack_fwrite_bool(file, data->boolean_value, &sz);
                check(sz > 0, "Unable to pack boolean value");
                break;
            default:
                sentinel("Unsupported data type in 'add_event' action data message struct");
        }
    }

    return 0;

error:
    return -1;
}

// Serializes the data map of an 'add_event' message.
//
// message - The message.
// file    - The file stream to read from.
//
// Returns 0 if successful, otherwise returns -1.
int sky_add_event_message_pack_data(sky_add_event_message *message, FILE *file)
{
    int rc;
    size_t sz;
    assert(message != NULL);
    assert(file != NULL);

    // Map
    minipack_fwrite_map(file, message->data_count, &sz);
    check(sz > 0, "Unable to write map");
    
    // Map items
    uint32_t i;
    for(i=0; i<message->data_count; i++) {
        sky_add_event_message_data *data = message->data[i];
        
        // Write key.
        rc = sky_minipack_fwrite_bstring(file, data->key);
        check(rc == 0, "Unable to pack data key");
        
        // Write in the appropriate data type.
        switch(data->data_type) {
            case SKY_DATA_TYPE_STRING:
                rc = sky_minipack_fwrite_bstring(file, data->string_value);
                check(rc == 0, "Unable to pack string value");
                break;
            case SKY_DATA_TYPE_INT:
                minipack_fwrite_int(file, data->int_value, &sz);
                check(sz > 0, "Unable to pack int value");
                break;
            case SKY_DATA_TYPE_DOUBLE:
                minipack_fwrite_double(file, data->double_value, &sz);
                check(sz > 0, "Unable to pack float value");
                break;
            case SKY_DATA_TYPE_BOOLEAN:
                minipack_fwrite_bool(file, data->boolean_value, &sz);
                check(sz > 0, "Unable to pack boolean value");
                break;
            default:
                sentinel("Unsupported data type in 'add_event' data message struct");
        }
    }

    return 0;

error:
    return -1;
}

// Deserializes an 'add_event' message from a file stream.
//
// message - The message.
// file    - The file stream to read from.
//
// Returns 0 if successful, otherwise returns -1.
int sky_add_event_message_unpack(sky_add_event_message *message, FILE *file)
{
    int rc;
    size_t sz;
    bstring key = NULL;
    assert(message != NULL);
    assert(file != NULL);

    // Map
    uint32_t map_length = minipack_fread_map(file, &sz);
    check(sz > 0, "Unable to read map");
    
    // Map items
    uint32_t i;
    for(i=0; i<map_length; i++) {
        rc = sky_minipack_fread_bstring(file, &key);
        check(rc == 0, "Unable to read map key");
        
        if(biseq(key, &SKY_ADD_EVENT_KEY_OBJECT_ID) == 1) {
            message->object_id = (sky_object_id_t)minipack_fread_uint(file, &sz);
            check(sz != 0, "Unable to unpack object id");
        }
        else if(biseq(key, &SKY_ADD_EVENT_KEY_TIMESTAMP) == 1) {
            message->timestamp = (sky_timestamp_t)minipack_fread_int(file, &sz);
            check(sz != 0, "Unable to unpack timestamp");
        }
        else if(biseq(key, &SKY_ADD_EVENT_KEY_ACTION) == 1) {
            rc = sky_add_event_message_unpack_action(message, file);
            check(rc == 0, "Unable to unpack 'add_event' action value");
        }
        else if(biseq(key, &SKY_ADD_EVENT_KEY_DATA) == 1) {
            rc = sky_add_event_message_unpack_data(message, file);
            check(rc == 0, "Unable to unpack 'add_event' data value");
        }
        
        bdestroy(key);
    }

    return 0;

error:
    bdestroy(key);
    return -1;
}

// Deserializes the action map of an 'add_event' message.
//
// message - The message.
// file    - The file stream to read from.
//
// Returns 0 if successful, otherwise returns -1.
int sky_add_event_message_unpack_action(sky_add_event_message *message, FILE *file)
{
    int rc;
    size_t sz;
    sky_add_event_message_data *data = NULL;
    assert(message != NULL);
    assert(file != NULL);

    // Map
    uint32_t map_length = minipack_fread_map(file, &sz);
    check(sz > 0, "Unable to read map");
    
    // Allocate action data array.
    message->action_data_count = map_length;
    message->action_data = calloc(1, sizeof(*message->action_data) * message->action_data_count); check_mem(message->action_data);
    
    // Map items
    uint32_t i, index = 0;
    for(i=0; i<map_length; i++) {
        bstring key = NULL;
        rc = sky_minipack_fread_bstring(file, &key);
        check(rc == 0, "Unable to read data key");

        // If this is the action name then save it to the message.
        if(biseq(key, &SKY_ADD_EVENT_KEY_NAME) == 1) {
            bdestroy(key); key = NULL;

            rc = sky_minipack_fread_bstring(file, &message->action_name);
            check(rc == 0, "Unable to unpack action name");

            // Decrement the total action data count if one of the data items
            // is the name.
            message->action_data_count--;
        }
        // Otherwise create an action data item.
        else {
            data = sky_add_event_message_data_create(); check_mem(data);
            message->action_data[index++] = data;
            data->key = key;
        
            // Read the first byte of the message to determine the type.
            uint8_t buffer[1];
            check(fread(buffer, sizeof(*buffer), 1, file) == 1, "Unable to read data type");
            ungetc(buffer[0], file);
        
            // Read in the appropriate data type.
            if(minipack_is_raw((void*)buffer)) {
                data->data_type = SKY_DATA_TYPE_STRING;
                rc = sky_minipack_fread_bstring(file, &data->string_value);
                check(rc == 0, "Unable to unpack string value");
            }
            else if(minipack_is_bool((void*)buffer)) {
                data->data_type = SKY_DATA_TYPE_BOOLEAN;
                data->boolean_value = minipack_fread_bool(file, &sz);
                check(sz != 0, "Unable to unpack boolean value");
            }
            else if(minipack_is_double((void*)buffer)) {
                data->data_type = SKY_DATA_TYPE_DOUBLE;
                data->double_value = minipack_fread_double(file, &sz);
                check(sz != 0, "Unable to unpack float value");
            }
            else {
                data->data_type = SKY_DATA_TYPE_INT;
                data->int_value = minipack_fread_int(file, &sz);
                check(sz != 0, "Unable to unpack int value");
            }
        }
    }

    return 0;

error:
    return -1;
}

// Deserializes the data map of an 'add_event' message.
//
// message - The message.
// file    - The file stream to read from.
//
// Returns 0 if successful, otherwise returns -1.
int sky_add_event_message_unpack_data(sky_add_event_message *message, FILE *file)
{
    int rc;
    size_t sz;
    assert(message != NULL);
    assert(file != NULL);

    // Map
    uint32_t map_length = minipack_fread_map(file, &sz);
    check(sz > 0, "Unable to read map");
    
    // Allocate data array.
    message->data_count = map_length;
    message->data = calloc(1, sizeof(*message->data) * message->data_count); check_mem(message->data);
    
    // Map items
    uint32_t i;
    for(i=0; i<map_length; i++) {
        sky_add_event_message_data *data = sky_add_event_message_data_create(); check_mem(data);
        message->data[i] = data;
        
        rc = sky_minipack_fread_bstring(file, &data->key);
        check(rc == 0, "Unable to read data key");
        
        // Read the first byte of the message to determine the type.
        uint8_t buffer[1];
        check(fread(buffer, sizeof(*buffer), 1, file) == 1, "Unable to read data type");
        ungetc(buffer[0], file);
        
        // Read in the appropriate data type.
        if(minipack_is_raw((void*)buffer)) {
            data->data_type = SKY_DATA_TYPE_STRING;
            rc = sky_minipack_fread_bstring(file, &data->string_value);
            check(rc == 0, "Unable to unpack string value");
        }
        else if(minipack_is_bool((void*)buffer)) {
            data->data_type = SKY_DATA_TYPE_BOOLEAN;
            data->boolean_value = minipack_fread_bool(file, &sz);
            check(sz != 0, "Unable to unpack boolean value");
        }
        else if(minipack_is_double((void*)buffer)) {
            data->data_type = SKY_DATA_TYPE_DOUBLE;
            data->double_value = minipack_fread_double(file, &sz);
            check(sz != 0, "Unable to unpack float value");
        }
        else {
            data->data_type = SKY_DATA_TYPE_INT;
            data->int_value = minipack_fread_int(file, &sz);
            check(sz != 0, "Unable to unpack int value");
        }
    }

    return 0;

error:
    return -1;
}


//--------------------------------------
// Worker
//--------------------------------------

// Adds the event to a given tablet.
//
// worker - The worker.
// tablet - The tablet to add the event to.
// ret    - Unused.
//
// Returns 0 if successful, otherwise returns -1.
int sky_add_event_message_worker_map(sky_worker *worker, sky_tablet *tablet,
                                     void **ret)
{
    int rc;
    assert(worker != NULL);
    assert(tablet != NULL);
    assert(ret != NULL);

    // Add event to tablet.
    sky_add_event_message *message = (sky_add_event_message*)worker->data;
    rc = sky_tablet_add_event(tablet, message->event);
    check(rc == 0, "Unable to add event to table");

    *ret = NULL;
    return 0;

error:
    if(ret) *ret = NULL;
    return -1;
}

// Writes the results to an output stream.
//
// worker - The worker.
// output - The output stream.
//
// Returns 0 if successful, otherwise returns -1.
int sky_add_event_message_worker_write(sky_worker *worker, FILE *output)
{
    size_t sz;
    assert(worker != NULL);
    assert(output != NULL);
    
    struct tagbstring status_str = bsStatic("status");
    struct tagbstring ok_str = bsStatic("ok");

    // Return {status:"OK"}
    check(minipack_fwrite_map(output, 1, &sz) == 0, "Unable to write output");
    check(sky_minipack_fwrite_bstring(output, &status_str) == 0, "Unable to write output");
    check(sky_minipack_fwrite_bstring(output, &ok_str) == 0, "Unable to write output");
    
    return 0;

error:
    return -1;
}

// Frees all data attached to the worker.
//
// worker - The worker.
//
// Returns 0 if successful, otherwise returns -1.
int sky_add_event_message_worker_free(sky_worker *worker)
{
    assert(worker != NULL);
    
    // Clean up.
    sky_add_event_message *message = (sky_add_event_message*)worker->data;
    sky_add_event_message_free(message);
    worker->data = NULL;
    
    return 0;
}
