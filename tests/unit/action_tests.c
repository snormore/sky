#include <stdio.h>
#include <stdlib.h>

#include <action.h>
#include <mem.h>
#include "../minunit.h"


//==============================================================================
//
// Helpers
//
//==============================================================================

#define ASSERT_BAD_DATA_MSG(FILENAME, ERROR_MSG) \
    FILE *file = fopen("tests/fixtures/actions/bad/" FILENAME, "r"); \
    sky_action *action = sky_action_create(); \
    mu_begin_log(); \
    int rc = sky_action_unpack(action, file); \
    mu_end_log(); \
    mu_assert_int_equals(rc, -1); \
    mu_assert_log(ERROR_MSG); \
    sky_action_free(action); \
    fclose(file); \

//==============================================================================
//
// Test Cases
//
//==============================================================================

//--------------------------------------
// Serialization
//--------------------------------------

int test_sky_action_pack() {
    cleantmp();
    sky_action *action = sky_action_create();
    action->id = 10;
    action->name = bfromcstr("foo");
    
    FILE *file = fopen("tmp/data", "w");
    mu_assert_bool(sky_action_pack(action, file) == 0);
    fclose(file);
    mu_assert_file("tmp/data", "tests/fixtures/actions/0/data");
    sky_action_free(action);
    return 0;
}

int test_sky_action_unpack() {
    FILE *file = fopen("tests/fixtures/actions/0/data", "r");
    sky_action *action = sky_action_create();
    mu_assert_bool(sky_action_unpack(action, file) == 0);
    fclose(file);

    mu_assert_int_equals(action->id, 10);
    mu_assert_bstring(action->name, "foo");
    sky_action_free(action);
    return 0;
}

int test_sky_action_unpack_bad_map_data() {
    ASSERT_BAD_DATA_MSG("bad_map", "Unable to read map")
    return 0;
}

int test_sky_action_unpack_bad_key_data() {
    ASSERT_BAD_DATA_MSG("bad_key", "Unable to read map key")
    return 0;
}

int test_sky_action_unpack_bad_id_data() {
    ASSERT_BAD_DATA_MSG("bad_id", "Unable to read action id")
    return 0;
}



//==============================================================================
//
// Setup
//
//==============================================================================

int all_tests() {
    mu_run_test(test_sky_action_pack);
    mu_run_test(test_sky_action_unpack);
    mu_run_test(test_sky_action_unpack_bad_map_data);
    mu_run_test(test_sky_action_unpack_bad_key_data);
    mu_run_test(test_sky_action_unpack_bad_id_data);
    return 0;
}

RUN_TESTS()