#include <stdio.h>
#include <stdlib.h>

#include <next_actions_message.h>
#include <dbg.h>
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
// Worker
//--------------------------------------

int test_sky_next_actions_message_worker_read() {
    sky_next_actions_message *message = sky_next_actions_message_create();
    sky_worker *worker = sky_worker_create();
    worker->data = (void*)message;
    FILE *file = fopen("tests/fixtures/next_actions_message/2/message", "r");

    int rc = sky_next_actions_message_worker_read(worker, file);
    mu_assert_int_equals(rc, 0);
    mu_assert_bool(worker->data != NULL);
    mu_assert_int_equals(message->prior_action_id_count, 2);
    mu_assert_int_equals(message->prior_action_ids[0], 3);
    mu_assert_int_equals(message->prior_action_ids[1], 4);

    fclose(file);
    sky_next_actions_message_free(message);
    sky_worker_free(worker);
    return 0;
}

int test_sky_next_actions_message_worker_map() {
    importtmp("tests/fixtures/next_actions_message/1/import.json");
    sky_table *table = sky_table_create();
    table->path = bfromcstr("tmp");
    sky_table_open(table);

    sky_next_actions_message *message = sky_next_actions_message_create();
    message->action_count = table->action_file->action_count;
    message->prior_action_id_count = 2;
    message->prior_action_ids = calloc(message->prior_action_id_count, sizeof(*message->prior_action_ids));
    message->prior_action_ids[0] = 1;
    message->prior_action_ids[1] = 2;
    sky_next_actions_message_init_data_descriptor(message, table->property_file);
    sky_worker *worker = sky_worker_create();
    worker->data = (void*)message;

    sky_next_actions_result *results;
    int rc = sky_next_actions_message_worker_map(worker, table->tablets[0], (void*)&results);
    mu_assert_int_equals(rc, 0);
    mu_assert_int_equals(message->action_count, 4);
    mu_assert_int_equals(results[0].count, 0);
    mu_assert_int_equals(results[1].count, 0);
    mu_assert_int_equals(results[2].count, 0);
    mu_assert_int_equals(results[3].count, 2);
    mu_assert_int_equals(results[4].count, 1);

    sky_next_actions_message_worker_map_free(results);
    sky_next_actions_message_free(message);
    sky_worker_free(worker);
    sky_table_free(table);
    return 0;
}

int test_sky_next_actions_message_worker_reduce() {
    sky_next_actions_message *message = sky_next_actions_message_create();
    message->action_count = 2;
    message->results = calloc(message->action_count+1, sizeof(*message->results));
    sky_worker *worker = sky_worker_create();
    worker->data = (void*)message;
    sky_next_actions_result *map_results = calloc(message->action_count+1, sizeof(*map_results));
    
    // Reduce #1
    map_results[0].count = 2;
    map_results[1].count = 4;
    map_results[2].count = 1;
    int rc = sky_next_actions_message_worker_reduce(worker, map_results);
    mu_assert_int_equals(rc, 0);
    mu_assert_int_equals(message->results[0].count, 2);
    mu_assert_int_equals(message->results[1].count, 4);
    mu_assert_int_equals(message->results[2].count, 1);

    // Reduce #2
    map_results[0].count = 6;
    map_results[1].count = 0;
    map_results[2].count = 100;
    rc = sky_next_actions_message_worker_reduce(worker, map_results);
    mu_assert_int_equals(rc, 0);
    mu_assert_int_equals(message->results[0].count, 8);
    mu_assert_int_equals(message->results[1].count, 4);
    mu_assert_int_equals(message->results[2].count, 101);

    free(map_results);
    sky_next_actions_message_free(message);
    sky_worker_free(worker);
    return 0;
}

int test_sky_next_actions_message_worker_write() {
    sky_next_actions_message *message = sky_next_actions_message_create();
    message->action_count = 2;
    message->results = calloc(message->action_count+1, sizeof(*message->results));
    message->results[1].count = 2;
    message->results[2].count = 6;
    sky_worker *worker = sky_worker_create();
    worker->data = (void*)message;

    FILE *file = fopen("tmp/message", "w");
    int rc = sky_next_actions_message_worker_write(worker, file);
    mu_assert_int_equals(rc, 0);
    fclose(file);
    mu_assert_file("tmp/message", "tests/fixtures/next_actions_message/1/output");

    sky_next_actions_message_free(message);
    sky_worker_free(worker);
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
    mu_run_test(test_sky_next_actions_message_worker_read);
    mu_run_test(test_sky_next_actions_message_worker_map);
    mu_run_test(test_sky_next_actions_message_worker_reduce);
    mu_run_test(test_sky_next_actions_message_worker_write);
    return 0;
}

RUN_TESTS()