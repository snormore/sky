#include <stdio.h>
#include <stdlib.h>

#include <sky.h>
#include <dbg.h>
#include <mem.h>

#include "server_helpers.h"


//==============================================================================
//
// Test Cases
//
//==============================================================================

int test() {
    pthread_t thread;
    importtmp_n("tests/functional/fixtures/multi/0/data.json", 4);
    start_server(1, &thread);
    send_msg("tests/functional/fixtures/multi/0/input");
    pthread_join(thread, NULL);
    mu_assert_msg("tests/functional/fixtures/multi/0/output");

    void *data;
    size_t data_length;
    sky_table *table = sky_table_create();
    table->path = bfromcstr("tmp");
    sky_table_open(table);

    sky_tablet_get_path(table->tablets[1], 1, &data, &data_length);
    mu_assert_mem(data, "\x01\x1A\x00\x00\x00\x00\x00\x00\x00\x04\x00", data_length);
    free(data);
    sky_tablet_get_path(table->tablets[2], 2, &data, &data_length);
    mu_assert_mem(data, "\x01\x1A\x00\x00\x00\x00\x00\x00\x00\x03\x00", data_length);
    free(data);
    sky_tablet_get_path(table->tablets[3], 3, &data, &data_length);
    mu_assert_mem(data, "\x01\x1A\x00\x00\x00\x00\x00\x00\x00\x02\x00", data_length);
    free(data);
    sky_tablet_get_path(table->tablets[0], 4, &data, &data_length);
    mu_assert_mem(data, "\x01\x1A\x00\x00\x00\x00\x00\x00\x00\x01\x00", data_length);
    free(data);
    sky_table_free(table);
    return 0;
}

//==============================================================================
//
// Setup
//
//==============================================================================

int all_tests() {
    mu_run_test(test);
    return 0;
}

RUN_TESTS()
