#ifndef _sky_data_file_h
#define _sky_data_file_h

#include <inttypes.h>
#include <stdbool.h>

#include "bstring.h"
#include "types.h"

//==============================================================================
//
// Definitions
//
//==============================================================================

// Defines a function that assigns the value at a given memory location to
// a property descriptor.
typedef void (*sky_data_property_descriptor_set_func)(void *target, void *value, size_t *sz);

// Defines the offset in memory for where the timestamp should be set.
typedef struct {
    uint16_t offset;
} sky_data_timestamp_descriptor;

// Defines the offset in memory for where the action id should be set.
typedef struct {
    uint16_t offset;
} sky_data_action_descriptor;

// Defines the offset in memory for where property value should be set.
typedef struct {
    sky_property_id_t property_id;
    uint16_t offset;
    sky_data_property_descriptor_set_func setter;
} sky_data_property_descriptor;

// Defines a collection of descriptors for a struct to serialize data into it.
typedef struct {
    sky_data_timestamp_descriptor timestamp_descriptor;
    sky_data_action_descriptor action_descriptor;
    sky_data_property_descriptor *property_descriptors;
    sky_property_id_t min_property_id;
    sky_property_id_t max_property_id;
    uint32_t property_descriptor_count;
} sky_data_descriptor;


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

sky_data_descriptor *sky_data_descriptor_create( sky_property_id_t min_property_id,
    sky_property_id_t max_property_id);

void sky_data_descriptor_free(sky_data_descriptor *descriptor);


#endif
