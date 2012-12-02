#include <stdio.h>
#include <stdlib.h>
#include <sys/time.h>

#include <sky_lua.h>
#include <path_iterator.h>
#include <cursor.h>
#include <types.h>
#include <mem.h>
#include <dbg.h>

#include "../minunit.h"


//==============================================================================
//
// Test Cases
//
//==============================================================================

//--------------------------------------
// Initialization
//--------------------------------------

typedef struct { sky_timestamp_t timestamp; sky_action_id_t action_id; int32_t x; int32_t y; } sky_lua_event_0_t;

int test_sky_lua_initscript_with_table() {
    int rc;
    importtmp("tests/fixtures/sky_lua/0/data.json");
    sky_table *table = sky_table_create();
    table->path = bfromcstr("tmp");
    sky_table_open(table);
    sky_data_descriptor *descriptor = sky_data_descriptor_create();

    sky_lua_event_0_t lua_event = {0, 1, 20, 30};

    struct tagbstring source = bsStatic(
        "function map(_event)\n"
        "  event = ffi.cast(\"sky_lua_event_t*\", _event)\n"
        "  return event.x + event.y\n"
        "end\n"
    );
    lua_State *L = NULL;
    rc = sky_lua_initscript_with_table(&source, table, descriptor, &L);
    mu_assert_int_equals(rc, 0);
    mu_assert_int_equals(descriptor->timestamp_descriptor.offset, 0);
    mu_assert_int_equals(descriptor->action_descriptor.offset, 8);
    mu_assert_int_equals(descriptor->property_zero_descriptor[1].offset, 0);
    mu_assert_int_equals(descriptor->property_zero_descriptor[2].offset, (int)offsetof(sky_lua_event_0_t, x));
    mu_assert_int_equals(descriptor->property_zero_descriptor[3].offset, (int)offsetof(sky_lua_event_0_t, y));
    mu_assert_int_equals(descriptor->property_zero_descriptor[4].offset, 0);

    // Call map() function with the event pointer.
    lua_getglobal(L, "map");
    lua_pushlightuserdata(L, &lua_event);
    lua_call(L, 1, 1);
    mu_assert_int_equals(rc, 0);
    mu_assert_long_equals(lua_tointeger(L, -1), 50L);
    
    sky_table_free(table);
    return 0;
}


//--------------------------------------
// Table Integration
//--------------------------------------

int test_sky_lua_generate_header() {
    importtmp("tests/fixtures/sky_lua/0/data.json");
    sky_table *table = sky_table_create();
    table->path = bfromcstr("tmp");
    sky_table_open(table);

    struct tagbstring source = bsStatic(
        "function map(event)\n"
        "  return event.x + event.y\n"
        "end\n"
    );
    bstring event_decl = NULL;
    bstring init_descriptor_func = NULL;
    int rc = sky_lua_generate_event_info(&source, table->property_file, &event_decl, &init_descriptor_func);
    mu_assert_int_equals(rc, 0);
    mu_assert_bstring(event_decl,
        "typedef struct {\n"
        "  int64_t timestamp;\n"
        "  uint16_t action_id;\n"
        "  int32_t x;\n"
        "  int32_t y;\n"
        "} sky_lua_event_t;"
    );
    mu_assert_bstring(init_descriptor_func,
        "function sky_init_descriptor(_descriptor)\n"
        "  descriptor = ffi.cast(\"sky_data_descriptor_t*\", _descriptor)\n"
        "  descriptor:set_data_sz(ffi.sizeof(\"sky_lua_event_t\"));\n"
        "  descriptor:set_timestamp_offset(ffi.offsetof(\"sky_lua_event_t\", \"timestamp\"));\n"
        "  descriptor:set_action_id_offset(ffi.offsetof(\"sky_lua_event_t\", \"action_id\"));\n"
        "  descriptor:set_property(2, ffi.offsetof(\"sky_lua_event_t\", \"x\"), 2);\n"
        "  descriptor:set_property(3, ffi.offsetof(\"sky_lua_event_t\", \"y\"), 2);\n"
        "end\n"
    );
    bdestroy(event_decl);
    bdestroy(init_descriptor_func);
    sky_table_free(table);
    return 0;
}


//--------------------------------------
// Map All
//--------------------------------------

int test_sky_map_all() {
    importtmp("tests/fixtures/sky_lua/1/data.json");
    sky_table *table = sky_table_create();
    table->path = bfromcstr("tmp");
    sky_table_open(table);
    sky_data_descriptor *descriptor = sky_data_descriptor_create();

    // Initialize the path iterator.
    sky_path_iterator iterator;
    sky_path_iterator_init(&iterator);
    iterator.cursor.data_descriptor = descriptor;
    sky_path_iterator_set_tablet(&iterator, table->tablets[0]);

    struct tagbstring source = bsStatic(
        "function map(cursor, data)\n"
        "  data.path_count = (data.path_count or 0) + 1\n"
        "  while not cursor:eof() do\n"
        "    event = cursor:event()\n"
        "    data.event_count = (data.event_count or 0) + 1\n"
        "    data.z = (data.z or 0) + event.x + event.y\n"
        "    cursor:next()\n"
        "  end\n"
        "end\n"
    );
    lua_State *L = NULL;
    int rc = sky_lua_initscript_with_table(&source, table, descriptor, &L);
    mu_assert_int_equals(rc, 0);

    // Allocate data.
    mu_assert_int_equals(descriptor->data_sz, 24);
    iterator.cursor.data = calloc(1, descriptor->data_sz);

    // Start benchmark.
    struct timeval tv;
    gettimeofday(&tv, NULL);
    int64_t t0 = (tv.tv_sec*1000) + (tv.tv_usec/1000);

    // Call sky_map_all() function.
    uint32_t i;
    for(i=0; i<1; i++) {
        sky_path_iterator_set_tablet(&iterator, table->tablets[0]);
        lua_getglobal(L, "sky_map_all");
        lua_pushlightuserdata(L, &iterator);
        lua_call(L, 1, 0);
    }

    // End benchmark.
    gettimeofday(&tv, NULL);
    int64_t t1 = (tv.tv_sec*1000) + (tv.tv_usec/1000);
    printf("[lua] t=%.3fs\n", ((float)(t1-t0))/1000);

    sky_table_free(table);
    free(iterator.cursor.data);
    return 0;
}


//==============================================================================
//
// Setup
//
//==============================================================================

int all_tests() {
    mu_run_test(test_sky_lua_initscript_with_table);
    mu_run_test(test_sky_lua_generate_header);
    mu_run_test(test_sky_map_all);
    return 0;
}

RUN_TESTS()