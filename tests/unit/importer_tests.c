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

    struct tagbstring one_str = bsStatic("1");
    sky_tablet_get_path(importer->table->tablets[1], &one_str, &data, &data_length);
    mu_assert_mem(data, "\x01\x00\x00\x80\x02\xf2\xb3\x04\x00\x02\x00\x03\x00\x00\xc0\x03\xf2\xb3\x04\x00\x01\x00\x09\x00\x00\x00\x01\x14\xff\xa3\x66\x6f\x6f\x02\xc3", data_length);
    free(data);

    struct tagbstring two_str = bsStatic("2");
    sky_tablet_get_path(importer->table->tablets[2], &two_str, &data, &data_length);
    mu_assert_mem(data, "\x03\x00\x00\x80\x02\xf2\xb3\x04\x00\x02\x00\x0e\x00\x00\x00\x01\x15\x02\xc2\xfe\xcb\x40\x59\x0c\xcc\xcc\xcc\xcc\xcd", data_length);
    free(data);

    struct tagbstring three_str = bsStatic("3");
    sky_tablet_get_path(importer->table->tablets[3], &three_str, &data, &data_length);
    mu_assert_mem(data, "\x01\x00\x00\x80\x02\xf2\xb3\x04\x00\x01\x00", data_length);
    free(data);

    struct tagbstring four_str = bsStatic("4");
    sky_tablet_get_path(importer->table->tablets[0], &four_str, &data, &data_length);
    mu_assert_mem(data, "\x01\x00\x00\x80\x02\xf2\xb3\x04\x00\x01\x00", data_length);
    free(data);

    struct tagbstring five_str = bsStatic("5");
    sky_tablet_get_path(importer->table->tablets[1], &five_str, &data, &data_length);
    mu_assert_mem(data, "\x01\x00\x00\x80\x02\xf2\xb3\x04\x00\x01\x00", data_length);
    free(data);

    struct tagbstring six_str = bsStatic("6");
    sky_tablet_get_path(importer->table->tablets[2], &six_str, &data, &data_length);
    mu_assert_mem(data, "\x01\x00\x00\x80\x02\xf2\xb3\x04\x00\x01\x00", data_length);
    free(data);

    struct tagbstring seven_str = bsStatic("7");
    sky_tablet_get_path(importer->table->tablets[3], &seven_str, &data, &data_length);
    mu_assert_mem(data, "\x01\x00\x00\x80\x02\xf2\xb3\x04\x00\x01\x00", data_length);
    free(data);

    struct tagbstring eight_str = bsStatic("8");
    sky_tablet_get_path(importer->table->tablets[0], &eight_str, &data, &data_length);
    mu_assert_mem(data, "\x01\x00\x00\x80\x02\xf2\xb3\x04\x00\x01\x00", data_length);
    free(data);

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