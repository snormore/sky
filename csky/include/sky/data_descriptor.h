#ifndef _sky_data_descriptor_h
#define _sky_data_descriptor_h

#include <stdlib.h>
#include <inttypes.h>
#include <stdbool.h>
#include <string.h>

//==============================================================================
//
// Typedefs
//
//==============================================================================

typedef void (*sky_data_property_descriptor_set_func)(void *target, void *value, size_t *sz);
typedef void (*sky_data_property_descriptor_clear_func)(void *target);

typedef struct { uint16_t ts_offset; uint16_t timestamp_offset;} sky_data_timestamp_descriptor;

typedef struct {
    int64_t property_id;
    uint16_t offset;
    sky_data_property_descriptor_set_func set_func;
    sky_data_property_descriptor_clear_func clear_func;
} sky_data_property_descriptor;

// Defines a collection of descriptors for a struct to serialize data into it.
typedef struct {
    sky_data_timestamp_descriptor timestamp_descriptor;
    sky_data_property_descriptor *property_descriptors;
    sky_data_property_descriptor *property_zero_descriptor;
    sky_data_property_descriptor **action_property_descriptors;
    uint32_t action_property_descriptor_count;
    uint32_t property_count;
    uint32_t active_property_count;
    uint32_t data_sz;
} sky_data_descriptor;


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

sky_data_descriptor *sky_data_descriptor_new(int64_t min_property_id, int64_t max_property_id);

void sky_data_descriptor_free(sky_data_descriptor *descriptor);


//--------------------------------------
// Value Management
//--------------------------------------

void sky_data_descriptor_set_value(sky_data_descriptor *descriptor,
  void *target, int64_t property_id, void *ptr, size_t *sz);


//--------------------------------------
// Descriptor Management
//--------------------------------------

void sky_data_descriptor_set_data_sz(sky_data_descriptor *descriptor, uint32_t sz);

void sky_data_descriptor_set_timestamp_offset(sky_data_descriptor *descriptor, uint32_t offset);

void sky_data_descriptor_set_ts_offset(sky_data_descriptor *descriptor, uint32_t offset);

void sky_data_descriptor_set_property(sky_data_descriptor *descriptor,
  int64_t property_id, uint32_t offset, const char *data_type);


#endif