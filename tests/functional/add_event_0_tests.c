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
    importtmp_n("tests/functional/fixtures/add_event/0/data.json", 4);
    start_server(1, &thread);
    send_msg("tests/functional/fixtures/add_event/0/input");
    pthread_join(thread, NULL);

    mu_assert_msg("tests/functional/fixtures/add_event/0/output");

    void *data;
    size_t data_length;
    sky_table *table = sky_table_create();
    table->path = bfromcstr("tmp");
    sky_table_open(table);
    sky_tablet_get_path(table->tablets[2], 10, &data, &data_length);
    mu_assert_mem(data, 
        "\x03\xE8\x03\x00\x00\x00\x00\x00\x00\x05\x00\x1F\x00\x00\x00\xFF"
        "\xA6\x7A\x7A\x7A\x7A\x7A\x7A\xFE\x0A\x01\xA3\x78\x79\x7A\x02\xD1"
        "\x00\xC8\x03\xCB\x40\x59\x0C\xCC\xCC\xCC\xCC\xCD\x04\xC3\x01\x40"
        "\x42\x0F\x00\x00\x00\x00\x00\x01\x00\x01\x80\x84\x1E\x00\x00\x00"
        "\x00\x00\x02\x00\x01\xC0\xC6\x2D\x00\x00\x00\x00\x00\x03\x00",
        data_length
    );
    sky_table_free(table);
    free(data);
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
