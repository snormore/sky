#include <stdio.h>
#include <stdlib.h>

#include <dbg.h>
#include <mem.h>
#include <path_iterator.h>

#include "minunit.h"


//==============================================================================
//
// Test Cases
//
//==============================================================================

//--------------------------------------
// Iteration
//--------------------------------------

int test_sky_cursor_next() {
    char data[1024]; memset(data, 0, 1024);
    FILE *file = fopen("tests/fixtures/cursors/0/data", "r");
    fread(data, 1, 1024, file);
    fclose(file);

    sky_cursor *cursor = sky_cursor_create();

    // Event 1
    int rc = sky_cursor_set_path(cursor, data);
    mu_assert_int_equals(rc, 0);
    mu_assert_int_equals(cursor->path_index, 0);
    mu_assert_int_equals(cursor->event_index, 0);
    mu_assert_long_equals(cursor->ptr-((void*)data), 8L);
    mu_assert_bool(!cursor->eof);
    
    // Event 2
    rc = sky_cursor_next(cursor);
    mu_assert_int_equals(rc, 0);
    mu_assert_int_equals(cursor->path_index, 0);
    mu_assert_int_equals(cursor->event_index, 1);
    mu_assert_long_equals(cursor->ptr-((void*)data), 19L);
    mu_assert_bool(!cursor->eof);

    // Event 3
    rc = sky_cursor_next(cursor);
    mu_assert_int_equals(rc, 0);
    mu_assert_int_equals(cursor->path_index, 0);
    mu_assert_int_equals(cursor->event_index, 2);
    mu_assert_long_equals(cursor->ptr-((void*)data), 37L);
    mu_assert_bool(!cursor->eof);
    
    // EOF
    rc = sky_cursor_next(cursor);
    mu_assert_int_equals(rc, 0);
    mu_assert_int_equals(cursor->path_index, 0);
    mu_assert_int_equals(cursor->event_index, 0);
    mu_assert_bool(cursor->eof);

    sky_cursor_free(cursor);
    return 0;
}


//==============================================================================
//
// Setup
//
//==============================================================================

int all_tests() {
    mu_run_test(test_sky_cursor_next);
    return 0;
}

RUN_TESTS()