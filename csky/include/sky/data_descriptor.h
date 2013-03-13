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

// Defines a collection of descriptors for a struct to serialize data into it.
typedef struct {
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