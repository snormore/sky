#include <stdlib.h>

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

void sky_data_descriptor_set_int(void *target, void *value, size_t *sz);

void sky_data_descriptor_set_double(void *target, void *value, size_t *sz);

void sky_data_descriptor_set_boolean(void *target, void *value, size_t *sz);


//--------------------------------------
// Clear Functions
//--------------------------------------

void sky_data_descriptor_clear_string(void *target);

void sky_data_descriptor_clear_int(void *target);

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
// min_property_id - The lowest property id.
// max_property_id - The highest property id.
//
// Returns a reference to the new descriptor.
sky_data_descriptor *sky_data_descriptor_create(sky_property_id_t min_property_id,
                                                sky_property_id_t max_property_id)
{
    check(max_property_id >= min_property_id, "Max property id must be greater than or equal to min property id");

    // Normalize min/max so that range so that either min or max is on zero
    // if it doesn't span across zero.
    if(min_property_id == 1) {
        min_property_id = 0;
    }
    else if(max_property_id == -1) {
        max_property_id = 0;
    }
    check(min_property_id <= 0 && max_property_id >= 0, "Property range must touch or cross zero boundry");

    // Determine the total number of properties to track.
    uint32_t property_count = max_property_id - min_property_id + 1;

    // Allocate memory for the descriptor and child descriptors in one block.
    size_t sz = sizeof(sky_data_descriptor) + (sizeof(sky_data_property_descriptor) * property_count);
    sky_data_descriptor *descriptor = calloc(sz, 1);
    check_mem(descriptor);
    descriptor->min_property_id = min_property_id;
    descriptor->max_property_id = max_property_id;
    descriptor->property_descriptors = (sky_data_property_descriptor*)(((void*)descriptor) + sizeof(sky_data_descriptor));
    descriptor->property_count = property_count;
    descriptor->property_zero_descriptor = NULL;
    
    // Initialize all property descriptors to noop.
    uint32_t i;
    for(i=0; i<property_count; i++) {
        sky_property_id_t property_id = min_property_id + (sky_property_id_t)i;
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
    check(descriptor != NULL, "Descriptor required");
    check(target != NULL, "Target struct required");
    check(property_id >= descriptor->min_property_id && property_id <= descriptor->max_property_id, "Property ID (%d) out of range (%d,%d)", property_id, descriptor->min_property_id, descriptor->max_property_id);
    check(ptr != NULL, "Value pointer required");
    
    // Find the property descriptor and call the set_func.
    sky_data_property_descriptor *property_descriptor = &descriptor->property_zero_descriptor[property_id];
    property_descriptor->set_func(target + property_descriptor->offset, ptr, sz);
    
    return 0;

error:
    return -1;
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
    check(descriptor != NULL, "Data descriptor required");
    check(target != NULL, "Target pointer required");
    
    uint32_t i;
    for(i=0; i<descriptor->action_property_descriptor_count; i++) {
        sky_data_property_descriptor *property_descriptor = descriptor->action_property_descriptors[i];
        property_descriptor->clear_func(target + property_descriptor->offset);
    }

    return 0;

error:
    return -1;
}


//--------------------------------------
// Property Management
//--------------------------------------

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
                                     uint16_t offset,
                                     sky_data_type_e data_type)
{
    check(descriptor != NULL, "Descriptor required");
    check(property_id >= descriptor->min_property_id && property_id <= descriptor->max_property_id, "Property ID out of range for descriptor");

    // Retrieve property descriptor.
    sky_data_property_descriptor *property_descriptor = &descriptor->property_zero_descriptor[property_id];
    check(property_descriptor->property_id == property_id, "Property descriptor id (%d) does not match (%d)", property_descriptor->property_id, property_id);
    
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
            property_descriptor->set_func = sky_data_descriptor_set_int;
            property_descriptor->clear_func = sky_data_descriptor_clear_int;
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
void sky_data_descriptor_set_int(void *target, void *value, size_t *sz)
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
void sky_data_descriptor_clear_int(void *target)
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

