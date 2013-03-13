#ifndef _sky_object_iterator_h
#define _sky_object_iterator_h

#include <stdbool.h>
#include <leveldb/c.h>
#include "sky/cursor.h"

typedef struct sky_object_iterator {
    leveldb_iterator_t* leveldb_iterator;
    bool running;
    bool eof;
    sky_cursor cursor;
} sky_object_iterator;

#endif