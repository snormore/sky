#include <stdio.h>
#include <stdlib.h>

#include <lua_map_reduce_message.h>
#include <sky_string.h>
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

int test_sky_lua_map_reduce_message_pack() {
    cleantmp();
    sky_lua_map_reduce_message *message = sky_lua_map_reduce_message_create();
    message->source = bfromcstr("x = 1\ny = 2\nreturn x + y");
    
    FILE *file = fopen("tmp/message", "w");
    mu_assert_bool(sky_lua_map_reduce_message_pack(message, file) == 0);
    fclose(file);
    mu_assert_file("tmp/message", "tests/fixtures/lua_map_reduce_message/0/message");
    sky_lua_map_reduce_message_free(message);
    return 0;
}

int test_sky_lua_map_reduce_message_unpack() {
    FILE *file = fopen("tests/fixtures/lua_map_reduce_message/0/message", "r");
    sky_lua_map_reduce_message *message = sky_lua_map_reduce_message_create();
    mu_assert_bool(sky_lua_map_reduce_message_unpack(message, file) == 0);
    fclose(file);

    mu_assert_bstring(message->source, "x = 1\ny = 2\nreturn x + y");
    sky_lua_map_reduce_message_free(message);
    return 0;
}


//--------------------------------------
// Worker
//--------------------------------------

int test_sky_lua_map_reduce_message_worker_read() {
    sky_lua_map_reduce_message *message = sky_lua_map_reduce_message_create();
    sky_worker *worker = sky_worker_create();
    worker->data = (void*)message;
    FILE *file = fopen("tests/fixtures/lua_map_reduce_message/1/message", "r");

    int rc = sky_lua_map_reduce_message_worker_read(worker, file);
    mu_assert_int_equals(rc, 0);
    mu_assert_bool(worker->data != NULL);
    mu_assert_bstring(message->source, "foo");

    fclose(file);
    sky_lua_map_reduce_message_free(message);
    sky_worker_free(worker);
    return 0;
}

int sky_lua_test(int x, int y) {
    return x + y;
}

int test_sky_lua_map_reduce_message_worker_map() {
    importtmp("tests/fixtures/lua_map_reduce_message/0/import.json");
    sky_table *table = sky_table_create();
    table->path = bfromcstr("tmp");
    sky_table_open(table);

    sky_lua_map_reduce_message *message = sky_lua_map_reduce_message_create();
    message->source = bfromcstr(
        "function map(event)\n"
        "  return 10\n"
        "end"
    );
    sky_worker *worker = sky_worker_create();
    worker->data = (void*)message;

    bstring results = NULL;
    int rc = sky_lua_map_reduce_message_worker_map(worker, table->tablets[0], (void**)&results);
    mu_assert_int_equals(rc, 0);
    mu_assert_int_equals(blength(results), 1);
    mu_assert_mem(bdata(results), "\x0a", 1);

    bdestroy(results);
    sky_lua_map_reduce_message_free(message);
    sky_worker_free(worker);
    sky_table_free(table);
    return 0;
}


//==============================================================================
//
// Setup
//
//==============================================================================

int all_tests() {
    mu_run_test(test_sky_lua_map_reduce_message_pack);
    mu_run_test(test_sky_lua_map_reduce_message_unpack);
    mu_run_test(test_sky_lua_map_reduce_message_worker_read);
    //mu_run_test(test_sky_lua_map_reduce_message_worker_map);
    return 0;
}

RUN_TESTS()