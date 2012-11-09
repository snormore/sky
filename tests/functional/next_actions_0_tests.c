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
    importtmp_n("tests/functional/fixtures/next_actions/0/data.json", 4);
    start_server(1, &thread);
    send_msg("tests/functional/fixtures/next_actions/0/input");
    pthread_join(thread, NULL);
    mu_assert_msg("tests/functional/fixtures/next_actions/0/output");
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