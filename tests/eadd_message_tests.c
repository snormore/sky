#include <stdio.h>
#include <stdlib.h>

#include <add_event_message.h>
#include <mem.h>

#include "minunit.h"


//==============================================================================
//
// Fixtures
//
//==============================================================================

sky_add_event_message *create_message_with_data()
{
    sky_add_event_message *message = sky_add_event_message_create();
    message->object_id = 10;
    message->timestamp = 1000LL;
    message->action_id = 20;
    message->data_count = 4;
    message->data = malloc(sizeof(message->data) * message->data_count);

    message->data[0] = sky_add_event_message_data_create();
    message->data[0]->key = bfromcstr("myString");
    message->data[0]->data_type = SKY_DATA_TYPE_STRING;
    message->data[0]->string_value = bfromcstr("xyz");

    message->data[1] = sky_add_event_message_data_create();
    message->data[1]->key = bfromcstr("myInt");
    message->data[1]->data_type = SKY_DATA_TYPE_INT;
    message->data[1]->int_value = 200;

    message->data[2] = sky_add_event_message_data_create();
    message->data[2]->key = bfromcstr("myFloat");
    message->data[2]->data_type = SKY_DATA_TYPE_DOUBLE;
    message->data[2]->double_value = 100.2;

    message->data[3] = sky_add_event_message_data_create();
    message->data[3]->key = bfromcstr("myBoolean");
    message->data[3]->data_type = SKY_DATA_TYPE_BOOLEAN;
    message->data[3]->boolean_value = true;

    return message;
}


//==============================================================================
//
// Test Cases
//
//==============================================================================

//--------------------------------------
// Serialization
//--------------------------------------

int test_sky_add_event_message_pack() {
    cleantmp();
    sky_add_event_message *message = create_message_with_data();
    FILE *file = fopen("tmp/message", "w");
    mu_assert_bool(sky_add_event_message_pack(message, file) == 0);
    fclose(file);
    mu_assert_file("tmp/message", "tests/fixtures/add_event_message/0/message");
    sky_add_event_message_free(message);
    return 0;
}

int test_sky_add_event_message_unpack() {
    FILE *file = fopen("tests/fixtures/add_event_message/0/message", "r");
    sky_add_event_message *message = sky_add_event_message_create();
    mu_assert_bool(sky_add_event_message_unpack(message, file) == 0);
    fclose(file);

    mu_assert_int_equals(message->object_id, 10);
    mu_assert_int64_equals(message->timestamp, 1000LL);
    mu_assert_int_equals(message->action_id, 20);
    sky_add_event_message_free(message);
    return 0;
}

int test_sky_add_event_message_sizeof() {
    sky_add_event_message *message = create_message_with_data();
    mu_assert_long_equals(sky_add_event_message_sizeof(message), 90L);
    sky_add_event_message_free(message);
    return 0;
}



//--------------------------------------
// Processing
//--------------------------------------

int test_sky_add_event_message_process() {
    loadtmp("tests/fixtures/add_event_message/1/table/pre");
    sky_table *table = sky_table_create();
    table->path = bfromcstr("tmp");
    sky_table_open(table);
    
    sky_add_event_message *message = create_message_with_data();
    FILE *output = fopen("tmp/output", "w");
    mu_assert(sky_add_event_message_process(message, table, output) == 0, "");
    fclose(output);
    mu_assert_file("tmp/0/header", "tests/fixtures/add_event_message/1/table/post/0/header");
    mu_assert_file("tmp/0/data", "tests/fixtures/add_event_message/1/table/post/0/data");
    mu_assert_file("tmp/output", "tests/fixtures/add_event_message/1/output");

    sky_add_event_message_free(message);
    sky_table_free(table);
    return 0;
}


//==============================================================================
//
// Setup
//
//==============================================================================

int all_tests() {
    mu_run_test(test_sky_add_event_message_pack);
    mu_run_test(test_sky_add_event_message_unpack);
    mu_run_test(test_sky_add_event_message_sizeof);
    mu_run_test(test_sky_add_event_message_process);
    return 0;
}

RUN_TESTS()