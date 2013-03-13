#include "sky/data_descriptor.h"
#include "sky/minipack.h"
#include "sky/sky_string.h"

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
sky_data_descriptor *sky_data_descriptor_new(int64_t min_property_id, int64_t max_property_id)
{
    sky_data_descriptor *descriptor = NULL;

    // Add one property to account for the zero descriptor.
    int64_t property_count = (max_property_id - min_property_id) + 1;

    // Allocate memory for the descriptor and child descriptors in one block.
    size_t sz = sizeof(sky_data_descriptor) + (sizeof(sky_data_property_descriptor) * property_count);
    descriptor = calloc(sz, 1);
    descriptor->property_descriptors = (sky_data_property_descriptor*)(((void*)descriptor) + sizeof(sky_data_descriptor));
    descriptor->property_count = property_count;
    descriptor->property_zero_descriptor = NULL;
    
    // Initialize all property descriptors to noop.
    uint32_t i;
    for(i=0; i<property_count; i++) {
        int64_t property_id = min_property_id + (int64_t)i;
        descriptor->property_descriptors[i].property_id = property_id;
        descriptor->property_descriptors[i].set_func = sky_data_descriptor_set_noop;
        
        // Save a pointer to the descriptor that points to property zero.
        if(property_id == 0) {
            descriptor->property_zero_descriptor = &descriptor->property_descriptors[i];
        }
    }

    return descriptor;
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
void sky_data_descriptor_set_value(sky_data_descriptor *descriptor,
                                   void *target, int64_t property_id,
                                   void *ptr, size_t *sz)
{
    // Find the property descriptor and call the set_func.
    sky_data_property_descriptor *property_descriptor = &descriptor->property_zero_descriptor[property_id];
    property_descriptor->set_func(target + property_descriptor->offset, ptr, sz);
}


//--------------------------------------
// Descriptor Management
//--------------------------------------

void sky_data_descriptor_set_data_sz(sky_data_descriptor *descriptor, uint32_t sz) {
    descriptor->data_sz = sz;
}

void sky_data_descriptor_set_timestamp_offset(sky_data_descriptor *descriptor, uint32_t offset) {
    descriptor->timestamp_descriptor.timestamp_offset = offset;
}

void sky_data_descriptor_set_ts_offset(sky_data_descriptor *descriptor, uint32_t offset) {
    descriptor->timestamp_descriptor.ts_offset = offset;
}

// Sets the data type and offset for a given property id.
//
// descriptor  - The data descriptor.
// property_id - The property id.
// offset      - The offset in the struct where the data should be set.
// data_type   - The data type to set.
//
// Returns 0 if successful, otherwise returns -1.
void sky_data_descriptor_set_property(sky_data_descriptor *descriptor,
                                      int64_t property_id,
                                      uint32_t offset,
                                      const char *data_type)
{
    // Retrieve property descriptor.
    sky_data_property_descriptor *property_descriptor = &descriptor->property_zero_descriptor[property_id];
    
    // Increment active property count if a new property is being set.
    if(property_descriptor->set_func == sky_data_descriptor_set_noop) {
        descriptor->active_property_count++;
    }
    
    // Set the offset and set_func function on the descriptor.
    property_descriptor->offset = offset;
    if(strlen(data_type) == 0) {
        property_descriptor->set_func = sky_data_descriptor_set_noop;
        property_descriptor->clear_func = NULL;
    }
    else if(strcmp(data_type, "string") == 0) {
        property_descriptor->set_func = sky_data_descriptor_set_string;
        property_descriptor->clear_func = sky_data_descriptor_clear_string;
    }
    else if(strcmp(data_type, "integer") == 0) {
        property_descriptor->set_func = sky_data_descriptor_set_int;
        property_descriptor->clear_func = sky_data_descriptor_clear_int;
    }
    else if(strcmp(data_type, "float") == 0) {
        property_descriptor->set_func = sky_data_descriptor_set_double;
        property_descriptor->clear_func = sky_data_descriptor_clear_double;
    }
    else if(strcmp(data_type, "boolean") == 0) {
        property_descriptor->set_func = sky_data_descriptor_set_boolean;
        property_descriptor->clear_func = sky_data_descriptor_clear_boolean;
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
            descriptor->action_property_descriptors[descriptor->action_property_descriptor_count-1] = property_descriptor;
        }
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
// Returns nothing.
void sky_data_descriptor_set_noop(void *target, void *value, size_t *sz)
{
    ((void)(target));
    *sz = minipack_sizeof_elem_and_data(value);
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
    *((int32_t*)target) = (int32_t)minipack_unpack_int(value, sz);
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
    *((int32_t*)target) = 0;
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

