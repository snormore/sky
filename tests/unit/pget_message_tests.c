#include <stdio.h>
#include <stdlib.h>

#include <get_property_message.h>
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

int test_sky_get_property_message_pack() {
    cleantmp();
    sky_get_property_message *message = sky_get_property_message_create();
    message->property_id = 20;
    
    FILE *file = fopen("tmp/message", "w");
    mu_assert_bool(sky_get_property_message_pack(message, file) == 0);
    fclose(file);
    mu_assert_file("tmp/message", "tests/fixtures/get_property_message/0/message");
    sky_get_property_message_free(message);
    return 0;
}

int test_sky_get_property_message_unpack() {
    FILE *file = fopen("tests/fixtures/get_property_message/0/message", "r");
    sky_get_property_message *message = sky_get_property_message_create();
    mu_assert_bool(sky_get_property_message_unpack(message, file) == 0);
    fclose(file);

    mu_assert_int_equals(message->property_id, 20);
    sky_get_property_message_free(message);
    return 0;
}


//--------------------------------------
// Processing
//--------------------------------------

int test_sky_get_property_message_process() {
    loadtmp("tests/fixtures/get_property_message/1/table");
    sky_server *server = sky_server_create(NULL);
    sky_message_header *header = sky_message_header_create();
    sky_table *table = sky_table_create();
    table->path = bfromcstr("tmp");
    sky_table_open(table);
    
    FILE *input = fopen("tests/fixtures/get_property_message/1/input", "r");
    FILE *output = fopen("tmp/output", "w");
    int rc = sky_get_property_message_process(server, header, table, input, output);
    mu_assert_int_equals(rc, 0);
    mu_assert_file("tmp/output", "tests/fixtures/get_property_message/1/output");

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
    mu_run_test(test_sky_get_property_message_pack);
    mu_run_test(test_sky_get_property_message_unpack);
    mu_run_test(test_sky_get_property_message_process);
    return 0;
}

RUN_TESTS()