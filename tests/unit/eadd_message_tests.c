#include <stdio.h>
#include <stdlib.h>

#include <add_event_message.h>
#include <dbg.h>
#include <mem.h>

#include "../minunit.h"


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
    message->action_name = bfromcstr("foo");
    message->data_count = 4;
    message->data = malloc(sizeof(message->data) * message->data_count);

    message->data[0] = sky_add_event_message_data_create();
    message->data[0]->key = bfromcstr("ostring");
    message->data[0]->data_type = SKY_DATA_TYPE_STRING;
    message->data[0]->string_value = bfromcstr("xyz");

    message->data[1] = sky_add_event_message_data_create();
    message->data[1]->key = bfromcstr("oint");
    message->data[1]->data_type = SKY_DATA_TYPE_INT;
    message->data[1]->int_value = 200;

    message->data[2] = sky_add_event_message_data_create();
    message->data[2]->key = bfromcstr("odouble");
    message->data[2]->data_type = SKY_DATA_TYPE_DOUBLE;
    message->data[2]->double_value = 100.2;

    message->data[3] = sky_add_event_message_data_create();
    message->data[3]->key = bfromcstr("oboolean");
    message->data[3]->data_type = SKY_DATA_TYPE_BOOLEAN;
    message->data[3]->boolean_value = true;


    message->action_data_count = 2;
    message->action_data = malloc(sizeof(message->action_data) * message->action_data_count);

    message->action_data[0] = sky_add_event_message_data_create();
    message->action_data[0]->key = bfromcstr("astring");
    message->action_data[0]->data_type = SKY_DATA_TYPE_STRING;
    message->action_data[0]->string_value = bfromcstr("zzzzzz");

    message->action_data[1] = sky_add_event_message_data_create();
    message->action_data[1]->key = bfromcstr("aint");
    message->action_data[1]->data_type = SKY_DATA_TYPE_INT;
    message->action_data[1]->int_value = 10;

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
    mu_assert_bstring(message->action_name, "foo");
    sky_add_event_message_free(message);
    return 0;
}

int test_sky_add_event_message_sizeof() {
    sky_add_event_message *message = create_message_with_data();
    mu_assert_long_equals(sky_add_event_message_sizeof(message), 115L);
    sky_add_event_message_free(message);
    return 0;
}


//--------------------------------------
// Worker
//--------------------------------------

int test_sky_add_event_message_worker_map() {
    loadtmp("tests/fixtures/add_event_message/1/table/pre");
    sky_table *table = sky_table_create();
    table->path = bfromcstr("tmp");
    table->default_tablet_count = 1;
    sky_table_open(table);

    struct tagbstring XYZ_STR = bsStatic("xyz");
    sky_add_event_message *message = sky_add_event_message_create();
    message->event = sky_event_create(10, 1000L, 20);
    message->event->data_count = 4;
    message->event->data = calloc(message->event->data_count, sizeof(*message->event->data));
    message->event->data[0] = sky_event_data_create_string(1, &XYZ_STR);
    message->event->data[1] = sky_event_data_create_int(2, 200);
    message->event->data[2] = sky_event_data_create_double(3, 100.2);
    message->event->data[3] = sky_event_data_create_boolean(4, true);
    sky_worker *worker = sky_worker_create();
    worker->data = (void*)message;

    void *null = NULL;
    int rc = sky_add_event_message_worker_map(worker, table->tablets[0], &null);
    mu_assert_int_equals(rc, 0);
    mu_assert_file("tmp/0/header", "tests/fixtures/add_event_message/1/table/post/0/header");
    mu_assert_file("tmp/0/data", "tests/fixtures/add_event_message/1/table/post/0/data");

    sky_add_event_message_free(message);
    sky_worker_free(worker);
    sky_table_free(table);
    return 0;
}

int test_sky_add_event_message_worker_write() {
    sky_add_event_message *message = sky_add_event_message_create();
    sky_worker *worker = sky_worker_create();
    worker->data = (void*)message;
    worker->output = fopen("tmp/output", "w");

    int rc = sky_add_event_message_worker_write(worker, worker->output);
    mu_assert_int_equals(rc, 0);
    sky_add_event_message_free(message);
    sky_worker_free(worker);

    mu_assert_file("tmp/output", "tests/fixtures/add_event_message/1/output");

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
    mu_run_test(test_sky_add_event_message_worker_map);
    mu_run_test(test_sky_add_event_message_worker_write);
    return 0;
}

RUN_TESTS()