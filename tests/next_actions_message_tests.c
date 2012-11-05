#include <stdio.h>
#include <stdlib.h>

#include <next_actions_message.h>
#include <mem.h>

#include "minunit.h"


//==============================================================================
//
// Test Cases
//
//==============================================================================

//--------------------------------------
// Serialization
//--------------------------------------

int test_sky_next_actions_message_pack() {
    cleantmp();
    sky_next_actions_message *message = sky_next_actions_message_create();
    message->prior_action_id_count = 2;
    message->prior_action_ids = calloc(message->prior_action_id_count, sizeof(*message->prior_action_ids));
    message->prior_action_ids[0] = 1;
    message->prior_action_ids[1] = 2;
    
    FILE *file = fopen("tmp/message", "w");
    mu_assert_bool(sky_next_actions_message_pack(message, file) == 0);
    fclose(file);
    mu_assert_file("tmp/message", "tests/fixtures/next_actions_message/0/message");
    sky_next_actions_message_free(message);
    return 0;
}

int test_sky_next_actions_message_unpack() {
    FILE *file = fopen("tests/fixtures/next_actions_message/0/message", "r");
    sky_next_actions_message *message = sky_next_actions_message_create();
    mu_assert_bool(sky_next_actions_message_unpack(message, file) == 0);
    fclose(file);

    mu_assert_int_equals(message->prior_action_id_count, 2);
    mu_assert_int_equals(message->prior_action_ids[0], 1);
    mu_assert_int_equals(message->prior_action_ids[1], 2);
    sky_next_actions_message_free(message);
    return 0;
}


//--------------------------------------
// Processing
//--------------------------------------

int test_sky_next_actions_message_process() {
    importtmp("tests/fixtures/next_actions_message/1/import.json");
    sky_table *table = sky_table_create();
    table->path = bfromcstr("tmp");
    sky_table_open(table);

    sky_next_actions_message *message = sky_next_actions_message_create();
    message->prior_action_id_count = 2;
    message->prior_action_ids = calloc(message->prior_action_id_count, sizeof(*message->prior_action_ids));
    message->prior_action_ids[0] = 1;
    message->prior_action_ids[1] = 2;
    FILE *output = fopen("tmp/output", "w");
    mu_assert(sky_next_actions_message_process(message, table->tablets[0], output) == 0, "");
    fclose(output);
    mu_assert_file("tmp/output", "tests/fixtures/next_actions_message/1/output");

    sky_next_actions_message_free(message);
    sky_table_free(table);
    return 0;
}


//==============================================================================
//
// Setup
//
//==============================================================================

int all_tests() {
    mu_run_test(test_sky_next_actions_message_pack);
    mu_run_test(test_sky_next_actions_message_unpack);
    mu_run_test(test_sky_next_actions_message_process);
    return 0;
}

RUN_TESTS()