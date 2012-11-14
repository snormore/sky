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
    mu_assert_file("tmp/0/data", "tests/functional/fixtures/multi/0/0/data");
    mu_assert_file("tmp/0/header", "tests/functional/fixtures/multi/0/0/header");
    mu_assert_file("tmp/1/data", "tests/functional/fixtures/multi/0/1/data");
    mu_assert_file("tmp/1/header", "tests/functional/fixtures/multi/0/1/header");
    mu_assert_file("tmp/2/data", "tests/functional/fixtures/multi/0/2/data");
    mu_assert_file("tmp/2/header", "tests/functional/fixtures/multi/0/2/header");
    mu_assert_file("tmp/3/data", "tests/functional/fixtures/multi/0/3/data");
    mu_assert_file("tmp/3/header", "tests/functional/fixtures/multi/0/3/header");
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