#include <stdio.h>
#include <stdlib.h>

#include <sky_lua.h>
#include <mem.h>

#include "../minunit.h"


//==============================================================================
//
// Test Cases
//
//==============================================================================

//--------------------------------------
// Initialization
//--------------------------------------

typedef struct { int64_t x; int64_t y; } sky_lua_event_0_t;

int test_sky_lua_initscript_with_table() {
    int rc;
    importtmp("tests/fixtures/sky_lua/0/data.json");
    sky_table *table = sky_table_create();
    table->path = bfromcstr("tmp");
    sky_table_open(table);

    sky_lua_event_0_t lua_event = {20, 30};

    struct tagbstring source = bsStatic(
        "function map(_event)\n"
        "  event = ffi.cast(\"sky_lua_event_t*\", _event)\n"
        "  return tonumber(event.x) + tonumber(event.y)\n"
        "end\n"
    );
    lua_State *L = NULL;
    rc = sky_lua_initscript_with_table(&source, table, &L);
    mu_assert_int_equals(rc, 0);

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
// Property File Integration
//--------------------------------------

int test_sky_lua_generate_event_struct_decl() {
    importtmp("tests/fixtures/sky_lua/0/data.json");
    sky_table *table = sky_table_create();
    table->path = bfromcstr("tmp");
    sky_table_open(table);

    struct tagbstring source = bsStatic(
        "function map(event)\n"
        "  return event._test123 + event.hello\n"
        "end\n"
        "q = super_event.not_found\n"
        "x = event._test123\n"
        "y = event."
    );
    bstring decl = NULL;
    int rc = sky_lua_generate_event_struct_decl(&source, table->property_file, &decl);
    mu_assert_int_equals(rc, 0);
    mu_assert_bstring(decl, 
        "typedef struct {\n"
        "  char _test123[];\n"
        "  int64_t hello;\n"
        "} sky_lua_event_t;"
    );
    bdestroy(decl);
    sky_table_free(table);
    return 0;
}

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
    bstring header = NULL;
    int rc = sky_lua_generate_header(&source, table, &header);
    mu_assert_int_equals(rc, 0);
    mu_assert_bstring(header, 
        "-- SKY GENERATED CODE BEGIN --\n"
        "local ffi = require(\"ffi\")\n"
        "ffi.cdef([[\n"
        "typedef struct sky_lua_path_iterator_t sky_lua_path_iterator_t;\n"
        "typedef struct sky_lua_cursor_t sky_lua_cursor_t;\n"
        "typedef struct {\n"
        "  int64_t x;\n"
        "  int64_t y;\n"
        "} sky_lua_event_t;\n"
        "\n"
        "sky_lua_cursor_t *sky_lua_path_iterator_next(sky_lua_path_iterator_t *);\n"
        "sky_lua_event_t *sky_lua_cursor_next(sky_lua_cursor_t *);\n"
        "]])\n"
        "-- SKY GENERATED CODE END --\n\n"
    );
    bdestroy(header);
    sky_table_free(table);
    return 0;
}


//==============================================================================
//
// Setup
//
//==============================================================================

int all_tests() {
    mu_run_test(test_sky_lua_initscript_with_table);
    mu_run_test(test_sky_lua_generate_event_struct_decl);
    mu_run_test(test_sky_lua_generate_header);
    return 0;
}

RUN_TESTS()