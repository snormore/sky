#include <stdio.h>
#include <timestamp.h>
#include <bstring.h>

#include "../minunit.h"

//==============================================================================
//
// Helpers
//
//==============================================================================

#define mu_timestamp_assert(STR, VALUE) do {\
    sky_timestamp_t timestamp = 0; \
    struct tagbstring str = bsStatic(STR); \
    int rc = sky_timestamp_parse(&str, &timestamp); \
    mu_assert_int_equals(rc, 0); \
    mu_assert_int64_equals(timestamp, VALUE); \
} while(0)


//==============================================================================
//
// Test Cases
//
//==============================================================================

int test_sky_timestamp_parse()
{
    // Parse ISO 8601 date on the epoch.
    mu_timestamp_assert("1970-01-01T00:00:00Z", 0LL);

    // Parse ISO 8601 date.
    mu_timestamp_assert("2010-01-02T10:30:20Z", 1262428220000000LL);

    // Parse ISO 8601 date past 2036.
    mu_timestamp_assert("2080-01-01T00:00:00Z", 3471292800000000LL);

    // Parse ISO 8601 date just before the epoch.
    mu_timestamp_assert("1969-12-31T23:59:59Z", -1000000LL);

    // Parse ISO 8601 date a year before the epoch.
    mu_timestamp_assert("1969-01-01T00:00:00Z", -31536000000000LL);

    // Parse ISO 8601 date a decade before the epoch.
    mu_timestamp_assert("1960-01-01T00:00:00Z", -315619200000000LL);

    // Parse ISO 8601 date a long time before the epoch.
    mu_timestamp_assert("1910-01-01T00:00:00Z", -1893456000000000LL);

    return 0;
}

int test_sky_timestamp_shift()
{
    // One second after the epoch.
    mu_assert_int64_equals(sky_timestamp_shift(1000000), 0x100000);

    // One second before the epoch.
    mu_assert_int64_equals(sky_timestamp_shift(-1000000), -0x100000);

    // One and a half seconds after the epoch.
    mu_assert_int64_equals(sky_timestamp_shift(1500000), 0x17a120);

    return 0;
}

int test_sky_timestamp_unshift()
{
    // One second after the epoch.
    mu_assert_int64_equals(sky_timestamp_unshift(0x100000), 1000000);

    // One second before the epoch.
    mu_assert_int64_equals(sky_timestamp_unshift(-0x100000), -1000000);

    // One and a half seconds after the epoch.
    mu_assert_int64_equals(sky_timestamp_unshift(0x17a120), 1500000);

    return 0;
}

//==============================================================================
//
// Setup
//
//==============================================================================

int all_tests() {
    mu_run_test(test_sky_timestamp_parse);
    mu_run_test(test_sky_timestamp_shift);
    mu_run_test(test_sky_timestamp_unshift);
    return 0;
}

RUN_TESTS()