#include <stdio.h>
#include <stdlib.h>

#include <create_table_message.h>
#include <mem.h>
#include <dbg.h>

#include "../minunit.h"


//==============================================================================
//
// Test Cases
//
//==============================================================================

//--------------------------------------
// Serialization
//--------------------------------------

int test_sky_create_table_message_pack() {
    cleantmp();
    sky_create_table_message *message = sky_create_table_message_create();
    message->data_path = bfromcstr("/tmp/data");
    message->table = sky_table_create();
    message->table->path = bfromcstr("/tmp/data/foo");
    message->table->default_tablet_count = 10;
    
    FILE *file = fopen("tmp/message", "w");
    mu_assert_bool(sky_create_table_message_pack(message, file) == 0);
    fclose(file);
    mu_assert_file("tmp/message", "tests/fixtures/create_table_message/0/message");
    sky_table_free(message->table);
    sky_create_table_message_free(message);
    return 0;
}

int test_sky_create_table_message_unpack() {
    FILE *file = fopen("tests/fixtures/create_table_message/0/message", "r");
    sky_create_table_message *message = sky_create_table_message_create();
    message->data_path = bfromcstr("/tmp/data");
    message->table = sky_table_create();
    mu_assert_bool(sky_create_table_message_unpack(message, file) == 0);
    fclose(file);

    mu_assert_bstring(message->table->path, "/tmp/data/foo");
    sky_table_free(message->table);
    sky_create_table_message_free(message);
    return 0;
}


//--------------------------------------
// Processing
//--------------------------------------

int test_sky_create_table_message_process() {
    cleantmp();
    sky_server *server = sky_server_create(NULL);
    server->path = bfromcstr("tmp");
    sky_message_header *header = sky_message_header_create();

    FILE *input = fopen("tests/fixtures/create_table_message/1/input", "r");
    FILE *output = fopen("tmp/output", "w");
    int rc = sky_create_table_message_process(server, header, NULL, input, output);
    mu_assert_int_equals(rc, 0);
    
    struct tagbstring foo_db_path = bsStatic("tmp/foo");
    struct tagbstring foo_db_0_path = bsStatic("tmp/foo/0");
    struct tagbstring foo_db_9_path = bsStatic("tmp/foo/9");
    mu_assert_bool(sky_file_exists(&foo_db_path));
    mu_assert_bool(sky_file_exists(&foo_db_0_path));
    mu_assert_bool(sky_file_exists(&foo_db_9_path));
    mu_assert_file("tmp/output", "tests/fixtures/create_table_message/1/output");

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
    mu_run_test(test_sky_create_table_message_pack);
    mu_run_test(test_sky_create_table_message_unpack);
    mu_run_test(test_sky_create_table_message_process);
    return 0;
}

RUN_TESTS()