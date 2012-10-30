#include <stdio.h>
#include <stdlib.h>

#include <sky_string.h>

#include "minunit.h"


//==============================================================================
//
// Test Cases
//
//==============================================================================

int test_sky_fixed_array_equals_same_var() {
    sky_string a = sky_string_create(4, "foo");
    mu_assert_bool(sky_string_equals(a, a) == true);
    return 0;
}

int test_sky_fixed_array_equals_same_data() {
    sky_string a = sky_string_create(4, "foo");
    sky_string b = sky_string_create(4, "foo");
    mu_assert_bool(sky_string_equals(a, b) == true);
    return 0;
}

int test_sky_fixed_array_equals_different_lengths() {
    sky_string a = sky_string_create(4, "foo");
    sky_string b = sky_string_create(5, "foo");
    mu_assert_bool(sky_string_equals(a, b) == false);
    return 0;
}

int test_sky_fixed_array_equals_different_data() {
    sky_string a = sky_string_create(4, "foo");
    sky_string b = sky_string_create(4, "bar");
    mu_assert_bool(sky_string_equals(a, b) == false);
    return 0;
}


//==============================================================================
//
// Setup
//
//==============================================================================

int all_tests() {
    mu_run_test(test_sky_fixed_array_equals_same_var);
    mu_run_test(test_sky_fixed_array_equals_same_data);
    mu_run_test(test_sky_fixed_array_equals_different_lengths);
    mu_run_test(test_sky_fixed_array_equals_different_data);
    return 0;
}

RUN_TESTS()