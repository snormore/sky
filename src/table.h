#ifndef _table_h
#define _table_h

#include <inttypes.h>
#include <stdio.h>
#include <stdbool.h>
#include <sys/stat.h>
#include <sys/mman.h>

typedef struct sky_table sky_table;

#include "bstring.h"
#include "action.h"
#include "event.h"
#include "types.h"
#include "data_file.h"
#include "action_file.h"
#include "property_file.h"


//==============================================================================
//
// Typedefs
//
//==============================================================================

#define SKY_LOCK_NAME ".skylock"

// The table is a reference to the disk location where data is stored. The
// table also maintains a cache of block info and predefined actions and
// properties.
struct sky_table {
    sky_data_file *data_file;
    sky_action_file *action_file;
    sky_property_file *property_file;
    bstring name;
    bstring path;
    bool opened;
    uint32_t default_block_size;
    FILE *lock_file;
};


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

sky_table *sky_table_create();

void sky_table_free(sky_table *table);


//--------------------------------------
// Path Management
//--------------------------------------

int sky_table_set_path(sky_table *table, bstring path);


//--------------------------------------
// State
//--------------------------------------

int sky_table_open(sky_table *table);

int sky_table_close(sky_table *table);


//--------------------------------------
// Event Management
//--------------------------------------

int sky_table_add_event(sky_table *table, sky_event *event);

#endif
