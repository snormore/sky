#ifndef _sky_worklet_h
#define _sky_worklet_h

#include <stdio.h>
#include <inttypes.h>
#include <stdbool.h>

typedef struct sky_worklet sky_worklet;

#include "bstring.h"
#include "worker.h"


//==============================================================================
//
// Typedefs
//
//==============================================================================

struct sky_worklet {
    sky_worker *worker;
    void *data;
};


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

sky_worklet *sky_worklet_create(sky_worker *worker);

void sky_worklet_free(sky_worklet *worklet);

#endif
