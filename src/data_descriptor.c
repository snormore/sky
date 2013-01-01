#include <stdlib.h>
#include <assert.h>

#include "data_descriptor.h"
#include "sky_string.h"
#include "minipack.h"
#include "dbg.h"
#include "mem.h"


//==============================================================================
//
// Forward Declarations
//
//==============================================================================

//--------------------------------------
// Setters
//--------------------------------------

void sky_data_descriptor_set_noop(void *target, void *value, size_t *sz);

void sky_data_descriptor_set_string(void *target, void *value, size_t *sz);

void sky_data_descriptor_set_int32(void *target, void *value, size_t *sz);

void sky_data_descriptor_set_int64(void *target, void *value, size_t *sz);

void sky_data_descriptor_set_double(void *target, void *value, size_t *sz);

void sky_data_descriptor_set_boolean(void *target, void *value, size_t *sz);


//--------------------------------------
// Clear Functions
//--------------------------------------

void sky_data_descriptor_clear_string(void *target);

void sky_data_descriptor_clear_int32(void *target);

void sky_data_descriptor_clear_int64(void *target);

void sky_data_descriptor_clear_double(void *target);

void sky_data_descriptor_clear_boolean(void *target);


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

// Creates a data descriptor with a given number of properties.
// 
// Returns a reference to the new descriptor.
sky_data_descriptor *sky_data_descriptor_create()
{
    sky_data_descriptor *descriptor = NULL;

    // Determine the total number of properties to track.
    uint32_t property_count = SKY_PROPERTY_ID_COUNT + 1;

    // Allocate memory for the descriptor and child descriptors in one block.
    size_t sz = sizeof(sky_data_descriptor) + (sizeof(sky_data_property_descriptor) * property_count);
    descriptor = calloc(sz, 1);
    check_mem(descriptor);
    descriptor->int_type = SKY_DATA_DESCRIPTOR_INT64;
    descriptor->property_descriptors = (sky_data_property_descriptor*)(((void*)descriptor) + sizeof(sky_data_descriptor));
    descriptor->property_count = property_count;
    descriptor->property_zero_descriptor = NULL;
    
    // Initialize all property descriptors to noop.
    uint32_t i;
    for(i=0; i<property_count; i++) {
        sky_property_id_t property_id = SKY_PROPERTY_ID_MIN + (sky_property_id_t)i;
        descriptor->property_descriptors[i].property_id = property_id;
        descriptor->property_descriptors[i].set_func = sky_data_descriptor_set_noop;
        
        // Save a pointer to the descriptor that points to property zero.
        if(property_id == 0) {
            descriptor->property_zero_descriptor = &descriptor->property_descriptors[i];
        }
    }

    // Make sure the descriptor got set.
    check(descriptor->property_zero_descriptor != NULL, "Property zero descriptor not initialized");
    
    return descriptor;
    
error:
    sky_data_descriptor_free(descriptor);
    return NULL;
}

// Removes a data descriptor reference from memory.
//
// descriptor - The data descriptor to free.
void sky_data_descriptor_free(sky_data_descriptor *descriptor)
{
    if(descriptor) {
        if(descriptor->action_property_descriptors != NULL) free(descriptor->action_property_descriptors);
        descriptor->action_property_descriptors = NULL;
        descriptor->action_property_descriptor_count = 0;
        
        free(descriptor);
    }
}

// Initializes a data descriptor to track the data properties that are set on
// an event. This assumes that the timestamp and action are being tracked
// using the sky_data_object type at the head of the data structure.
//
// This function will also return the size of the data structure needed to
// track using this descriptor.
//
// event      - The event.
// descriptor - The data descriptor.
//
// Returns 0 if successful, otherwise -1.
int sky_data_descriptor_init_with_event(sky_data_descriptor *descriptor,
                                        sky_event *event)
{
    int rc;
    assert(descriptor != NULL);
    assert(event != NULL);
    
    // Initialize the offset to start right after timestamp and action.
    descriptor->data_sz = (uint32_t)sizeof(sky_data_object);
    
    // Set the standard timestamp and action offsets.
    descriptor->timestamp_descriptor.timestamp_offset = offsetof(sky_data_object, timestamp);
    descriptor->timestamp_descriptor.ts_offset = offsetof(sky_data_object, ts);
    descriptor->action_descriptor.offset = offsetof(sky_data_object, action_id);
    
    // Loop over event data properties and set them on the descriptor.
    uint32_t i;
    for(i=0; i<event->data_count; i++) {
        // Set the property on the descriptor.
        sky_event_data *data = event->data[i];
        rc = sky_data_descriptor_set_property(descriptor, data->key, descriptor->data_sz, data->data_type);
        check(rc == 0, "Unable to set property on data descriptor");

        // Increment data type offset.
        size_t _sz = sky_data_type_sizeof(data->data_type);
        if(_sz < 8) _sz = 8;
        descriptor->data_sz += _sz;
    }
    
    return 0;

error:
    descriptor->data_sz = 0;
    return -1;
}


//--------------------------------------
// Value Management
//--------------------------------------

// Assigns a value to a struct through the data descriptor.
//
// descriptor  - The data descriptor.
// target      - The target struct.
// property_id - The property id to set.
// ptr         - A pointer to the memory location to read the value from.
// sz          - A pointer to where the number of bytes read are returned to.
//
// Returns 0 if successful, otherwise returns -1.
int sky_data_descriptor_set_value(sky_data_descriptor *descriptor,
                                  void *target, sky_property_id_t property_id,
                                  void *ptr, size_t *sz)
{
    assert(descriptor != NULL);
    assert(target != NULL);
    assert(ptr != NULL);
    
    // Find the property descriptor and call the set_func.
    sky_data_property_descriptor *property_descriptor = &descriptor->property_zero_descriptor[property_id];
    property_descriptor->set_func(target + property_descriptor->offset, ptr, sz);
    
    return 0;
}

// Clears all the action values that are being managed by the descriptor.
//
// descriptor - The data descriptor.
// target     - The target data.
//
// Returns 0 if successful, otherwise return -1.
int sky_data_descriptor_clear_action_data(sky_data_descriptor *descriptor,
                                          void *target)
{
    assert(descriptor != NULL);
    assert(target != NULL);
    
    uint32_t i;
    for(i=0; i<descriptor->action_property_descriptor_count; i++) {
        sky_data_property_descriptor *property_descriptor = descriptor->action_property_descriptors[i];
        property_descriptor->clear_func(target + property_descriptor->offset);
    }

    return 0;
}


//--------------------------------------
// Descriptor Management
//--------------------------------------

// Sets the total size of the data object, in bytes.
//
// descriptor  - The data descriptor.
// sz          - The size of the data struct, in bytes.
//
// Returns 0 if successful, otherwise returns -1.
int sky_data_descriptor_set_data_sz(sky_data_descriptor *descriptor,
                                    uint32_t sz)
{
    assert(descriptor != NULL);
    descriptor->data_sz = sz;
    return 0;
}

// Sets the offset of the unix timestamp property on the data object.
//
// descriptor  - The data descriptor.
// offset      - The offset in the struct where the timestamp should be set.
//
// Returns 0 if successful, otherwise returns -1.
int sky_data_descriptor_set_timestamp_offset(sky_data_descriptor *descriptor,
                                             uint32_t offset)
{
    assert(descriptor != NULL);
    descriptor->timestamp_descriptor.timestamp_offset = offset;
    return 0;
}

// Sets the offset of the internal Sky timestamp property on the data object.
//
// descriptor  - The data descriptor.
// offset      - The offset in the struct where the Sky timestamp should be set.
//
// Returns 0 if successful, otherwise returns -1.
int sky_data_descriptor_set_ts_offset(sky_data_descriptor *descriptor,
                                      uint32_t offset)
{
    assert(descriptor != NULL);
    descriptor->timestamp_descriptor.ts_offset = offset;
    return 0;
}

// Sets the offset of the action_id property on the data object.
//
// descriptor  - The data descriptor.
// offset      - The offset in the struct where the action_id should be set.
//
// Returns 0 if successful, otherwise returns -1.
int sky_data_descriptor_set_action_id_offset(sky_data_descriptor *descriptor,
                                             uint32_t offset)
{
    assert(descriptor != NULL);
    descriptor->action_descriptor.offset = offset;
    return 0;
}


// Sets the data type and offset for a given property id.
//
// descriptor  - The data descriptor.
// property_id - The property id.
// offset      - The offset in the struct where the data should be set.
// data_type   - The data type to set.
//
// Returns 0 if successful, otherwise returns -1.
int sky_data_descriptor_set_property(sky_data_descriptor *descriptor,
                                     sky_property_id_t property_id,
                                     uint32_t offset,
                                     sky_data_type_e data_type)
{
    assert(descriptor != NULL);

    // Retrieve property descriptor.
    sky_data_property_descriptor *property_descriptor = &descriptor->property_zero_descriptor[property_id];
    check(property_descriptor->property_id == property_id, "Property descriptor id (%d) does not match (%d)", property_descriptor->property_id, property_id);
    
    // Increment active property count if a new property is being set.
    if(property_descriptor->set_func == sky_data_descriptor_set_noop) {
        descriptor->active_property_count++;
    }
    
    
    // Set the offset and set_func function on the descriptor.
    property_descriptor->offset = offset;
    switch(data_type) {
        case SKY_DATA_TYPE_NONE: {
            property_descriptor->set_func = sky_data_descriptor_set_noop;
            property_descriptor->clear_func = NULL;
            break;
        }
        case SKY_DATA_TYPE_STRING: {
            property_descriptor->set_func = sky_data_descriptor_set_string;
            property_descriptor->clear_func = sky_data_descriptor_clear_string;
            break;
        }
        case SKY_DATA_TYPE_INT: {
            if(descriptor->int_type == SKY_DATA_DESCRIPTOR_INT32) {
                property_descriptor->set_func = sky_data_descriptor_set_int32;
                property_descriptor->clear_func = sky_data_descriptor_clear_int32;
            }
            else {
                property_descriptor->set_func = sky_data_descriptor_set_int64;
                property_descriptor->clear_func = sky_data_descriptor_clear_int64;
            }
            break;
        }
        case SKY_DATA_TYPE_DOUBLE: {
            property_descriptor->set_func = sky_data_descriptor_set_double;
            property_descriptor->clear_func = sky_data_descriptor_clear_double;
            break;
        }
        case SKY_DATA_TYPE_BOOLEAN: {
            property_descriptor->set_func = sky_data_descriptor_set_boolean;
            property_descriptor->clear_func = sky_data_descriptor_clear_boolean;
            break;
        }
        default: {
            sentinel("Invalid property data type");
        }
    }
    
    // If descriptor is an action descriptor and it doesn't exist in our list then add it.
    if(property_id < 0) {
        uint32_t i;
        bool found = false;
        for(i=0; i<descriptor->action_property_descriptor_count; i++) {
            if(descriptor->action_property_descriptors[i]->property_id == property_id) {
                found = true;
                break;
            }
        }
        
        // If it's not in the list then add it.
        if(!found) {
            descriptor->action_property_descriptor_count++;
            descriptor->action_property_descriptors = realloc(descriptor->action_property_descriptors, sizeof(*descriptor->action_property_descriptors) * descriptor->action_property_descriptor_count);
            check_mem(descriptor->action_property_descriptors);
            descriptor->action_property_descriptors[descriptor->action_property_descriptor_count-1] = property_descriptor;
        }
    }
    
    return 0;

error:
    return -1;
}


//--------------------------------------
// Setters
//--------------------------------------

// Performs no operation on the data. Used for skipping over unused properties
// in the descriptor.
//
// target - Unused.
// value  - The memory location where a MessagePack value is encoded.
// sz     - A pointer to where the number of bytes read should be returned.
//
// Returns nothing.
void sky_data_descriptor_set_noop(void *target, void *value, size_t *sz)
{
    check(target != NULL, "Target pointer cannot be null");
    *sz = minipack_sizeof_elem_and_data(value);
    
error:
    return;
}

// Reads a MessagePack string value from memory and sets the value to the
// given memory location.
//
// target - The place where the sky string should be written to.
// value  - The memory location where a MessagePack encoded raw is located.
// sz     - A pointer to where the number of bytes read should be returned.
//
// Returns nothing.
void sky_data_descriptor_set_string(void *target, void *value, size_t *sz)
{
    size_t _sz;
    sky_string *string = (sky_string*)target;
    string->length = minipack_unpack_raw(value, &_sz);
    string->data = (_sz > 0 ? value + _sz : NULL);
    *sz = _sz + string->length;
}

// Reads a MessagePack integer value from memory and sets the value to the
// given memory location.
//
// target - The place where the integer should be written to.
// value  - The memory location where a MessagePack encoded int is located.
// sz     - A pointer to where the number of bytes read should be returned.
//
// Returns nothing.
void sky_data_descriptor_set_int32(void *target, void *value, size_t *sz)
{
    *((int32_t*)target) = (int32_t)minipack_unpack_int(value, sz);
}

// Reads a MessagePack integer value from memory and sets the value to the
// given memory location.
//
// target - The place where the integer should be written to.
// value  - The memory location where a MessagePack encoded int is located.
// sz     - A pointer to where the number of bytes read should be returned.
//
// Returns nothing.
void sky_data_descriptor_set_int64(void *target, void *value, size_t *sz)
{
    *((int64_t*)target) = minipack_unpack_int(value, sz);
}

// Reads a MessagePack double value from memory and sets the value to the
// given memory location.
//
// target - The place where the integer should be written to.
// value  - The memory location where a MessagePack encoded double is located.
// sz     - A pointer to where the number of bytes read should be returned.
//
// Returns nothing.
void sky_data_descriptor_set_double(void *target, void *value, size_t *sz)
{
    *((double*)target) = minipack_unpack_double(value, sz);
}

// Reads a MessagePack boolean value from memory and sets the value to the
// given memory location.
//
// target - The place where the integer should be written to.
// value  - The memory location where a MessagePack encoded boolean is located.
// sz     - A pointer to where the number of bytes read should be returned.
//
// Returns nothing.
void sky_data_descriptor_set_boolean(void *target, void *value, size_t *sz)
{
    *((bool*)target) = minipack_unpack_bool(value, sz);
}


//--------------------------------------
// Clear Functions
//--------------------------------------

// Clears a string.
//
// target - The location of the string data to reset.
//
// Returns nothing.
void sky_data_descriptor_clear_string(void *target)
{
    sky_string *string = (sky_string*)target;
    string->length = 0;
    string->data = NULL;
}

// Clears an integer.
//
// target - The location of the int data to reset.
//
// Returns nothing.
void sky_data_descriptor_clear_int32(void *target)
{
    *((int32_t*)target) = 0;
}

// Clears an integer.
//
// target - The location of the int data to reset.
//
// Returns nothing.
void sky_data_descriptor_clear_int64(void *target)
{
    *((int64_t*)target) = 0LL;
}

// Clears a double.
//
// target - The location of the double data to reset.
//
// Returns nothing.
void sky_data_descriptor_clear_double(void *target)
{
    *((double*)target) = 0;
}

// Clears a boolean.
//
// target - The location of the boolean data to reset.
//
// Returns nothing.
void sky_data_descriptor_clear_boolean(void *target)
{
    *((bool*)target) = false;
}

