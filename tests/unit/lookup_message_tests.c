#include <stdio.h>
#include <stdlib.h>

#include <lookup_message.h>
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

int test_sky_lookup_message_pack() {
    cleantmp();
    sky_lookup_message *message = sky_lookup_message_create();
    message->action_name_count = 2;
    message->action_names = calloc(message->action_name_count, sizeof(bstring));
    message->action_names[0] = bfromcstr("foo");
    message->action_names[1] = bfromcstr("bar");
    message->property_name_count = 3;
    message->property_names = calloc(message->property_name_count, sizeof(bstring));
    message->property_names[0] = bfromcstr("xxx");
    message->property_names[1] = bfromcstr("yyy");
    message->property_names[2] = bfromcstr("zzz");
    
    FILE *file = fopen("tmp/message", "w");
    mu_assert_bool(sky_lookup_message_pack(message, file) == 0);
    fclose(file);
    mu_assert_file("tmp/message", "tests/fixtures/lookup_message/0/message");
    sky_lookup_message_free(message);
    return 0;
}

int test_sky_lookup_message_unpack() {
    FILE *file = fopen("tests/fixtures/lookup_message/0/message", "r");
    sky_lookup_message *message = sky_lookup_message_create();
    mu_assert_bool(sky_lookup_message_unpack(message, file) == 0);
    fclose(file);

    mu_assert_int_equals(message->action_name_count, 2);
    mu_assert_bstring(message->action_names[0], "foo");
    mu_assert_bstring(message->action_names[1], "bar");
    mu_assert_int_equals(message->property_name_count, 3);
    mu_assert_bstring(message->property_names[0], "xxx");
    mu_assert_bstring(message->property_names[1], "yyy");
    mu_assert_bstring(message->property_names[2], "zzz");
    sky_lookup_message_free(message);
    return 0;
}


//--------------------------------------
// Processing
//--------------------------------------

int test_sky_lookup_message_process() {
    importtmp("tests/fixtures/lookup_message/1/import.json");
    sky_server *server = sky_server_create(NULL);
    sky_message_header *header = sky_message_header_create();
    sky_table *table = sky_table_create();
    table->path = bfromcstr("tmp");
    sky_table_open(table);
    
    FILE *input = fopen("tests/fixtures/lookup_message/1/input", "r");
    FILE *output = fopen("tmp/output", "w");
    int rc = sky_lookup_message_process(server, header, table, input, output);
    mu_assert_int_equals(rc, 0);
    mu_assert_file("tmp/output", "tests/fixtures/lookup_message/1/output");

    sky_table_free(table);
    sky_message_header_free(header);
    sky_server_free(server);
    return 0;
}


//==============================================================================
//
// Setup
//
//==============================================================================

int all_tests() {
    mu_run_test(test_sky_lookup_message_pack);
    mu_run_test(test_sky_lookup_message_unpack);
    mu_run_test(test_sky_lookup_message_process);
    return 0;
}

RUN_TESTS()