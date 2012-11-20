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
        "char _test123[];\n"
        "int64_t hello;\n"
        "} sky_event_t;"
    );
    bdestroy(decl);
    sky_table_free(table);
    return 0;
}


//==============================================================================
//
// Setup
//
//==============================================================================

int all_tests() {
    mu_run_test(test_sky_lua_generate_event_struct_decl);
    return 0;
}

RUN_TESTS()