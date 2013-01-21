#include <stdio.h>
#include <stdlib.h>

#include <lua_aggregate_message.h>
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

int test_sky_lua_aggregate_message_pack() {
    cleantmp();
    sky_lua_aggregate_message *message = sky_lua_aggregate_message_create();
    message->source = bfromcstr("x = 1\ny = 2\nreturn x + y");
    
    FILE *file = fopen("tmp/message", "w");
    mu_assert_bool(sky_lua_aggregate_message_pack(message, file) == 0);
    fclose(file);
    mu_assert_file("tmp/message", "tests/fixtures/lua_aggregate_message/0/message");
    sky_lua_aggregate_message_free(message);
    return 0;
}

int test_sky_lua_aggregate_message_unpack() {
    FILE *file = fopen("tests/fixtures/lua_aggregate_message/0/message", "r");
    sky_lua_aggregate_message *message = sky_lua_aggregate_message_create();
    mu_assert_bool(sky_lua_aggregate_message_unpack(message, file) == 0);
    fclose(file);

    mu_assert_bstring(message->source, "x = 1\ny = 2\nreturn x + y");
    sky_lua_aggregate_message_free(message);
    return 0;
}


//--------------------------------------
// Worker
//--------------------------------------

int sky_lua_test(int x, int y) {
    return x + y;
}

int test_sky_lua_aggregate_message_worker_map() {
    importtmp("tests/fixtures/lua_aggregate_message/0/import.json");
    sky_table *table = sky_table_create();
    table->path = bfromcstr("tmp");
    sky_table_open(table);

    sky_lua_aggregate_message *message = sky_lua_aggregate_message_create();
    message->source = bfromcstr(
        "function aggregate(cursor, data)\n"
        "  event = cursor.event\n"
        "  data.count = data.count or 0\n"
        "  \n"
        "  while cursor:next() do\n"
        "    mystr = event:mystr()\n"
        "    data.count = data.count + 1\n"
        "    data[event.action_id] = (data[event.action_id] or 0) + 1\n"
        "    data[mystr] = (data[mystr] or 0) + 1\n"
        "  end\n"
        "end"
    );
    sky_worker *worker = sky_worker_create();
    worker->data = (void*)message;

    bstring results = NULL;
    int rc = sky_lua_aggregate_message_worker_map(worker, table->tablets[0], (void**)&results);
    mu_assert_int_equals(rc, 0);
    mu_assert_int_equals(blength(results), 31);
    mu_assert_mem(
        bdatae(results, ""), 
        "\x88\x01\x02\x02\x02\x03\x01\x04\x01\xA3" "baz" "\x02\xA5" "count" "\x06\xA3" "foo" "\x03\xA3" "bar" "\x01",
        blength(results)
    );

    bdestroy(results);
    sky_lua_aggregate_message_free(message);
    sky_worker_free(worker);
    sky_table_free(table);
    return 0;
}

int test_sky_lua_aggregate_message_worker_reduce() {
    int rc;
    importtmp("tests/fixtures/lua_aggregate_message/0/import.json");
    sky_table *table = sky_table_create();
    table->path = bfromcstr("tmp");
    sky_table_open(table);

    sky_lua_aggregate_message *message = sky_lua_aggregate_message_create();
    message->results = bfromcstr("\x80");
    message->source = bfromcstr(
        "function reduce(results, data)\n"
        "  for k,v in pairs(data) do\n"
        "    results[k] = (results[k] or 0) + v\n"
        "  end\n"
        "  return results\n"
        "end"
    );
    rc = sky_lua_initscript_with_table(message->source, table, NULL, &message->L);
    mu_assert_int_equals(rc, 0);
    sky_worker *worker = sky_worker_create();
    worker->data = (void*)message;

    struct tagbstring map_data1 = bsStatic("\x81\x02\x08");
    rc = sky_lua_aggregate_message_worker_reduce(worker, &map_data1);
    mu_assert_int_equals(rc, 0);
    mu_assert_int_equals(blength(message->results), 3);
    mu_assert_mem(bdatae(message->results, ""), "\x81\x02\x08", blength(message->results));

    struct tagbstring map_data2 = bsStatic("\x83\x02\x02\x04\x01\xA3" "foo" "\02");
    rc = sky_lua_aggregate_message_worker_reduce(worker, &map_data2);
    mu_assert_int_equals(rc, 0);
    mu_assert_int_equals(blength(message->results), 10);
    mu_assert_mem(bdatae(message->results, ""), "\x83\x02\x0A\x04\x01\xA3" "foo" "\x02", blength(message->results));

    sky_lua_aggregate_message_free(message);
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
    mu_run_test(test_sky_lua_aggregate_message_pack);
    mu_run_test(test_sky_lua_aggregate_message_unpack);
    mu_run_test(test_sky_lua_aggregate_message_worker_map);
    mu_run_test(test_sky_lua_aggregate_message_worker_reduce);
    return 0;
}

RUN_TESTS()