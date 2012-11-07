#include <stdio.h>
#include <stdlib.h>

#include <multi_message.h>
#include <mem.h>

#include "../minunit.h"


//==============================================================================
//
// Test Cases
//
//==============================================================================

//--------------------------------------
// Serialization
//--------------------------------------

int test_sky_multi_message_pack() {
    cleantmp();
    sky_multi_message *message = sky_multi_message_create();
    message->message_count = 20;
    
    FILE *file = fopen("tmp/message", "w");
    mu_assert_bool(sky_multi_message_pack(message, file) == 0);
    fclose(file);
    mu_assert_file("tmp/message", "tests/fixtures/multi_message/0/message");
    sky_multi_message_free(message);
    return 0;
}

int test_sky_multi_message_unpack() {
    FILE *file = fopen("tests/fixtures/multi_message/0/message", "r");
    sky_multi_message *message = sky_multi_message_create();
    mu_assert_bool(sky_multi_message_unpack(message, file) == 0);
    fclose(file);

    mu_assert_int_equals(message->message_count, 20);
    sky_multi_message_free(message);
    return 0;
}


//==============================================================================
//
// Setup
//
//==============================================================================

int all_tests() {
    mu_run_test(test_sky_multi_message_pack);
    mu_run_test(test_sky_multi_message_unpack);
    return 0;
}

RUN_TESTS()