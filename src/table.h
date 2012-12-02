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
#include "tablet.h"
#include "action_file.h"
#include "property_file.h"


//==============================================================================
//
// Typedefs
//
//==============================================================================

#define SKY_LOCK_NAME ".skylock"

#define DEFAULT_TABLET_COUNT 4

// The table is a collection of tablets. Data is shared into tablets but
// action and property data is managed by the table as a whole.
struct sky_table {
    sky_action_file *action_file;
    sky_property_file *property_file;
    sky_tablet **tablets;
    uint32_t tablet_count;
    bstring name;
    bstring path;
    bool opened;
    uint32_t default_tablet_count;
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
// Tablet Management
//--------------------------------------

int sky_table_get_tablet_count(sky_table *table, uint32_t *ret);

int sky_table_get_target_tablet(sky_table *table, sky_object_id_t object_id,
    sky_tablet **ret);

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
