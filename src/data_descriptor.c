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

void sky_data_descriptor_set_noop(void *target, void *value, size_t *sz);

void sky_data_descriptor_set_string(void *target, void *value, size_t *sz);

void sky_data_descriptor_set_int64(void *target, void *value, size_t *sz);

void sky_data_descriptor_set_double(void *target, void *value, size_t *sz);

void sky_data_descriptor_set_boolean(void *target, void *value, size_t *sz);


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
    descriptor->property_descriptors = (sky_data_property_descriptor*)(descriptor + sizeof(sky_data_descriptor));
    descriptor->property_count = property_count;
    descriptor->property_zero_descriptor = NULL;
    
    // Initialize all property descriptors to noop.
    uint32_t i;
    for(i=0; i<property_count; i++) {
        sky_property_id_t property_id = min_property_id + (sky_property_id_t)i;
        descriptor->property_descriptors[i].property_id = property_id;
        descriptor->property_descriptors[i].setter = sky_data_descriptor_set_noop;
        
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
        free(descriptor);
    }
}


//--------------------------------------
// Property Management
//--------------------------------------

// Sets the data type and offset for a given property id.
//
// descriptor  - The data descriptor.
// property_id - The property id.
// property_id - The property id.
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
    
    // Set the offset and setter function on the descriptor.
    property_descriptor->offset = offset;
    
    switch(data_type) {
        case SKY_DATA_TYPE_NONE: {
            property_descriptor->setter = sky_data_descriptor_set_noop;
            break;
        }
        case SKY_DATA_TYPE_STRING: {
            property_descriptor->setter = sky_data_descriptor_set_string;
            break;
        }
        case SKY_DATA_TYPE_INT: {
            property_descriptor->setter = sky_data_descriptor_set_int64;
            break;
        }
        case SKY_DATA_TYPE_DOUBLE: {
            property_descriptor->setter = sky_data_descriptor_set_double;
            break;
        }
        case SKY_DATA_TYPE_BOOLEAN: {
            property_descriptor->setter = sky_data_descriptor_set_boolean;
            break;
        }
        default: {
            sentinel("Invalid property data type");
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

