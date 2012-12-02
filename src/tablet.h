#ifndef _sky_table_h
#define _sky_table_h

#include <inttypes.h>
#include <stdio.h>
#include <stdbool.h>
#include <leveldb/c.h>

typedef struct sky_tablet sky_tablet;

#include "bstring.h"
#include "table.h"
#include "event.h"


//==============================================================================
//
// Typedefs
//
//==============================================================================

// The tablet is a reference to the disk location where data is stored.
struct sky_tablet {
    sky_table *table;
    leveldb_t *leveldb_db;
    uint32_t index;
    bstring path;
    leveldb_readoptions_t* readoptions;
    leveldb_writeoptions_t* writeoptions;
};


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

sky_tablet *sky_tablet_create(sky_table *table);

void sky_tablet_free(sky_tablet *tablet);


//--------------------------------------
// Path Management
//--------------------------------------

int sky_tablet_set_path(sky_tablet *tablet, bstring path);


//--------------------------------------
// State
//--------------------------------------

int sky_tablet_open(sky_tablet *tablet);

int sky_tablet_close(sky_tablet *tablet);


//--------------------------------------
// Event Management
//--------------------------------------

int sky_tablet_get_path(sky_tablet *tablet, sky_object_id_t object_id,
    void **data, size_t *data_length);

//--------------------------------------
// Event Management
//--------------------------------------

int sky_tablet_add_event(sky_tablet *tablet, sky_event *event);

#endif
