#include <inttypes.h>
#include <stdbool.h>
#include <stdlib.h>

#include "worklet.h"
#include "bstring.h"
#include "dbg.h"


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

// Creates a worklet.
//
// worker - The parent worker that the worklet is associated with.
//
// Returns a reference to the worker.
sky_worklet *sky_worklet_create(sky_worker *worker)
{
    sky_worklet *worklet = NULL;
    check(worker != NULL, "Worker required");
    worklet = calloc(1, sizeof(sky_worklet)); check_mem(worklet);
    worklet->worker = worker;
    return worklet;

error:
    sky_worklet_free(worklet);
    return NULL;
}

// Frees a worklet from memory.
//
// worklet - The worklet.
//
// Returns nothing.
void sky_worklet_free(sky_worklet *worklet)
{
    if(worklet) {
        worklet->worker = NULL;
        worklet->data = NULL;
        free(worklet);
    }
}
