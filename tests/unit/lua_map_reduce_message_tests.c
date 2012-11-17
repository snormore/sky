#include <stdio.h>
#include <stdlib.h>

#include <lua_map_reduce_message.h>
#include <dbg.h>
#include <mem.h>

#include "../minunit.h"


//==============================================================================
//
// Test Cases
//
//==============================================================================

//--------------------------------------
// Worker
//--------------------------------------

int test_sky_lua_map_reduce_message_worker_map() {
    importtmp("tests/fixtures/lua_map_reduce_message/0/import.json");
    sky_table *table = sky_table_create();
    table->path = bfromcstr("tmp");
    sky_table_open(table);

    sky_lua_map_reduce_message *message = sky_lua_map_reduce_message_create();
    message->source = bfromcstr(
        "x = 12\n"
        "y = 100\n"
        "io.write(\"TESTING!\" .. (x*y) .. \"\\n\")"
    );
    sky_worker *worker = sky_worker_create();
    worker->data = (void*)message;

    void *results = NULL;
    int rc = sky_lua_map_reduce_message_worker_map(worker, table->tablets[0], &results);
    mu_assert_int_equals(rc, 0);

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
    mu_run_test(test_sky_lua_map_reduce_message_worker_map);
    return 0;
}

RUN_TESTS()