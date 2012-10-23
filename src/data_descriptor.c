#include <stdlib.h>

#include "dbg.h"
#include "mem.h"
#include "data_descriptor.h"


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
    // Determine the total number of properties to track.
    int32_t property_count = max_property_id - min_property_id;

    // Allocate memory for the descriptor and child descriptors in one block.
    size_t sz = sizeof(sky_data_descriptor) + (sizeof(sky_data_property_descriptor) * property_count);
    sky_data_descriptor *descriptor = calloc(sz, 1);
    check_mem(descriptor);
    descriptor->min_property_id = min_property_id;
    descriptor->max_property_id = max_property_id;
    descriptor->property_descriptors = (sky_data_property_descriptor*)(descriptor + sizeof(sky_data_descriptor));
    descriptor->property_descriptor_count = property_count;
    
    // Initialize all property descriptors to noop.
    uint32_t i;
    for(i=0; i<property_count; i++) {
        descriptor->property_descriptors[i].property_id = min_property_id + (sky_property_id_t)i;
        descriptor->property_descriptors[i].setter = sky_data_descriptor_set_noop;
    }
    
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
// Setters
//--------------------------------------

// Performs no operation on the data. Used for skipping over unused properties
// in the descriptor.
//
// target - Unused.
// value  - The memory location where a MessagePack value is encoded.
// sz     - A pointer to where the number of bytes read should be returned.
//
// Returns the number of bytes read.
size_t sky_data_descriptor_set_noop(void *target, void *value, size_t *sz)
{
    *sz = minipack_sizeof_elem_and_data(value);
}

// Reads a MessagePack integer value from memory and sets the value to the
// given memory location.
//
// target - The place where the integer should be written to.
// value  - The memory location where a MessagePack encoded int is located.
// sz     - A pointer to where the number of bytes read should be returned.
//
// Returns the number of bytes read.
size_t sky_data_descriptor_set_int64(void *target, void *value, size_t *sz)
{
    *((int64_target*)target) = minipack_unpack_int(value, sz);
}

