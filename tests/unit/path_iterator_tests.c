#include <stdio.h>
#include <stdlib.h>
#include <math.h>

#include <path_iterator.h>
#include <dbg.h>
#include <mem.h>

#include "../minunit.h"

//==============================================================================
//
// Test Cases
//
//==============================================================================

//--------------------------------------
// Iteration
//--------------------------------------

int test_sky_path_iterator_next() {
    int rc;
    importtmp("tests/fixtures/path_iterator/0/data.json");
    sky_table *table = sky_table_create();
    table->path = bfromcstr("tmp");
    sky_table_open(table);

    sky_path_iterator *iterator = sky_path_iterator_create();
    sky_path_iterator_set_tablet(iterator, table->tablets[0]);
    mu_assert_mem(iterator->cursor.startptr, "\x01\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00", iterator->cursor.endptr-iterator->cursor.startptr);

    rc = sky_path_iterator_next(iterator);
    mu_assert_int_equals(rc, 0);
    mu_assert_mem(iterator->cursor.startptr, "\x01\x40\x42\x0F\x00\x00\x00\x00\x00\x02\x00", iterator->cursor.endptr-iterator->cursor.startptr);
    
    rc = sky_path_iterator_next(iterator);
    mu_assert_int_equals(rc, 0);
    mu_assert_mem(iterator->cursor.startptr, "\x01\x80\x84\x1E\x00\x00\x00\x00\x00\x01\x00", iterator->cursor.endptr-iterator->cursor.startptr);
    
    rc = sky_path_iterator_next(iterator);
    mu_assert_int_equals(rc, 0);
    mu_assert_mem(iterator->cursor.startptr, "\x01\x80\x84\x1E\x00\x00\x00\x00\x00\x01\x00", iterator->cursor.endptr-iterator->cursor.startptr);
    
    sky_path_iterator_free(iterator);
    sky_table_free(table);
    return 0;
}


//==============================================================================
//
// Setup
//
//==============================================================================

int all_tests() {
    mu_run_test(test_sky_path_iterator_next);
    return 0;
}

RUN_TESTS()