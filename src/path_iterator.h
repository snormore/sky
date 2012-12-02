#ifndef _path_iterator_h
#define _path_iterator_h

#include <inttypes.h>
#include <stdbool.h>
#include <leveldb/c.h>

#include "bstring.h"
#include "tablet.h"
#include "cursor.h"


//==============================================================================
//
// Typedefs
//
//==============================================================================

typedef struct sky_path_iterator {
    sky_tablet *tablet;
    leveldb_iterator_t* leveldb_iterator;
    bool running;
    bool eof;
    sky_cursor cursor;
} sky_path_iterator;


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

sky_path_iterator *sky_path_iterator_create();

sky_path_iterator *sky_path_iterator_alloc();

void sky_path_iterator_init(sky_path_iterator *iterator);

void sky_path_iterator_uninit(sky_path_iterator *iterator);

void sky_path_iterator_free(sky_path_iterator *iterator);

//--------------------------------------
// Source
//--------------------------------------

int sky_path_iterator_set_tablet(sky_path_iterator *iterator,
    sky_tablet *tablet);

//--------------------------------------
// Iteration
//--------------------------------------

int sky_path_iterator_next(sky_path_iterator *iterator);

bool sky_path_iterator_eof(sky_path_iterator *iterator);

#endif
