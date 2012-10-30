#include <stdio.h>
#include <stdlib.h>

#include <data_descriptor.h>
#include <bstring.h>
#include <mem.h>

#include "minunit.h"


//==============================================================================
//
// Test Cases
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

int test_sky_data_descriptor_create() {
    sky_data_descriptor *descriptor = sky_data_descriptor_create(-10, 21);
    mu_assert_int_equals(descriptor->min_property_id, -10);
    mu_assert_int_equals(descriptor->max_property_id, 21);
    mu_assert_int_equals(descriptor->property_count, 32);
    mu_assert_int_equals(descriptor->property_descriptors[0].property_id, -10);
    mu_assert_int_equals(descriptor->property_descriptors[31].property_id, 21);
    mu_assert_int_equals(descriptor->property_descriptors[10].property_id, 0);
    mu_assert_bool(descriptor->property_zero_descriptor == &descriptor->property_descriptors[10]);
    sky_data_descriptor_free(descriptor);
    return 0;
}


//==============================================================================
//
// Setup
//
//==============================================================================

int all_tests() {
    mu_run_test(test_sky_data_descriptor_create);
    return 0;
}

RUN_TESTS()