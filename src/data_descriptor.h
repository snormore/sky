#ifndef _sky_data_file_h
#define _sky_data_file_h

#include <inttypes.h>
#include <stdbool.h>

#include "bstring.h"
#include "types.h"
#include "event.h"

//==============================================================================
//
// Definitions
//
//==============================================================================

// Defines a function that assigns the value at a given memory location to
// a property descriptor.
typedef void (*sky_data_property_descriptor_set_func)(void *target, void *value, size_t *sz);

// Defines a function that clears the value at a given memory location to
// a property descriptor.
typedef void (*sky_data_property_descriptor_clear_func)(void *target);

// Defines a structure for the beginning of a data descriptor data object that
// includes the timestamp and action id.
typedef struct {
    sky_timestamp_t timestamp;
    sky_action_id_t action_id;
} sky_data_object;

// The type of integer to use for storage.
typedef enum {
    SKY_DATA_DESCRIPTOR_INT64 = 0,
    SKY_DATA_DESCRIPTOR_INT32 = 1,
} sky_data_descriptor_int_type_e;

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
    sky_data_property_descriptor_set_func set_func;
    sky_data_property_descriptor_clear_func clear_func;
} sky_data_property_descriptor;

// Defines a collection of descriptors for a struct to serialize data into it.
typedef struct {
    sky_data_timestamp_descriptor timestamp_descriptor;
    sky_data_action_descriptor action_descriptor;
    sky_data_property_descriptor *property_descriptors;
    sky_data_property_descriptor *property_zero_descriptor;
    sky_data_property_descriptor **action_property_descriptors;
    uint32_t action_property_descriptor_count;
    uint32_t property_count;
    uint32_t active_property_count;
    uint32_t data_sz;
    sky_data_descriptor_int_type_e int_type;
} sky_data_descriptor;


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

sky_data_descriptor *sky_data_descriptor_create();

void sky_data_descriptor_free(sky_data_descriptor *descriptor);

int sky_data_descriptor_init_with_event(sky_data_descriptor *descriptor,
    sky_event *event);

//--------------------------------------
// Value Management
//--------------------------------------

int sky_data_descriptor_set_value(sky_data_descriptor *descriptor,
    void *target, sky_property_id_t property_id, void *ptr, size_t *sz);

int sky_data_descriptor_clear_action_data(sky_data_descriptor *descriptor,
    void *target);

//--------------------------------------
// Descriptor Management
//--------------------------------------

int sky_data_descriptor_set_data_sz(sky_data_descriptor *descriptor,
    uint32_t sz);

int sky_data_descriptor_set_timestamp_offset(sky_data_descriptor *descriptor,
    uint32_t offset);

int sky_data_descriptor_set_action_id_offset(sky_data_descriptor *descriptor,
    uint32_t offset);

int sky_data_descriptor_set_property(sky_data_descriptor *descriptor,
    sky_property_id_t property_id, uint32_t offset, sky_data_type_e data_type);

#endif
