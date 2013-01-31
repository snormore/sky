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

    struct tagbstring ten_str = bsStatic("10");
    sky_tablet_get_path(table->tablets[2], &ten_str, &data, &data_length);
    mu_assert_mem(data, 
        "\x03\xe8\x03\x00\x00\x00\x00\x00\x00\x05\x00\x1f\x00\x00\x00\xff"
        "\xa6\x7a\x7a\x7a\x7a\x7a\x7a\xfe\x0a\x01\xa3\x78\x79\x7a\x02\xd1"
        "\x00\xc8\x03\xcb\x40\x59\x0c\xcc\xcc\xcc\xcc\xcd\x04\xc3\x01\x00"
        "\x00\x10\x00\x00\x00\x00\x00\x01\x00\x01\x00\x00\x20\x00\x00\x00"
        "\x00\x00\x02\x00\x01\x00\x00\x30\x00\x00\x00\x00\x00\x03\x00",
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
