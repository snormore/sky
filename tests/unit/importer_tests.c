#include <stdio.h>
#include <stdlib.h>

#include <importer.h>
#include <mem.h>
#include "../minunit.h"


//==============================================================================
//
// Test Cases
//
//==============================================================================

//--------------------------------------
// Import
//--------------------------------------

int test_sky_importer_import() {
    cleantmp();
    sky_importer *importer = sky_importer_create();
    importer->path = bfromcstr("tmp");
    importer->tablet_count = 4;
    
    FILE *file = fopen("tests/fixtures/importer/0/data.json", "r");
    int rc = sky_importer_import(importer, file);
    mu_assert_int_equals(rc, 0);
    fclose(file);

    // Validate.
    void *data;
    size_t data_length;
    mu_assert_bool(importer->table != NULL);
    mu_assert_file("tmp/actions", "tests/fixtures/importer/0/table/actions");
    mu_assert_file("tmp/properties", "tests/fixtures/importer/0/table/properties");

    sky_tablet_get_path(importer->table->tablets[1], 1, &data, &data_length);
    mu_assert_mem(data, "\x01\x00\x5A\x6A\xF8\x2B\x7C\x04\x00\x02\x00\x03\x00\x87\x9B\xF9\x2B\x7C\x04\x00\x01\x00\x09\x00\x00\x00\x01\x14\xFF\xA3\x66\x6F\x6F\x02\xC3", data_length);
    sky_tablet_get_path(importer->table->tablets[2], 2, &data, &data_length);
    mu_assert_mem(data, "\x03\x00\x5A\x6A\xF8\x2B\x7C\x04\x00\x02\x00\x0E\x00\x00\x00\x01\x15\x02\xC2\xFE\xCB\x40\x59\x0C\xCC\xCC\xCC\xCC\xCD", data_length);
    sky_tablet_get_path(importer->table->tablets[3], 3, &data, &data_length);
    mu_assert_mem(data, "\x01\x00\x5A\x6A\xF8\x2B\x7C\x04\x00\x01\x00", data_length);
    sky_tablet_get_path(importer->table->tablets[0], 4, &data, &data_length);
    mu_assert_mem(data, "\x01\x00\x5A\x6A\xF8\x2B\x7C\x04\x00\x01\x00", data_length);
    sky_tablet_get_path(importer->table->tablets[1], 5, &data, &data_length);
    mu_assert_mem(data, "\x01\x00\x5A\x6A\xF8\x2B\x7C\x04\x00\x01\x00", data_length);
    sky_tablet_get_path(importer->table->tablets[2], 6, &data, &data_length);
    mu_assert_mem(data, "\x01\x00\x5A\x6A\xF8\x2B\x7C\x04\x00\x01\x00", data_length);
    sky_tablet_get_path(importer->table->tablets[3], 7, &data, &data_length);
    mu_assert_mem(data, "\x01\x00\x5A\x6A\xF8\x2B\x7C\x04\x00\x01\x00", data_length);
    sky_tablet_get_path(importer->table->tablets[0], 8, &data, &data_length);
    mu_assert_mem(data, "\x01\x00\x5A\x6A\xF8\x2B\x7C\x04\x00\x01\x00", data_length);

    sky_importer_free(importer);
    return 0;
}


//==============================================================================
//
// Setup
//
//==============================================================================

int all_tests() {
    mu_run_test(test_sky_importer_import);
    return 0;
}

RUN_TESTS()