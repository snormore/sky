#ifndef _sky_servlet_h
#define _sky_servlet_h

#include <stdio.h>
#include <inttypes.h>
#include <stdbool.h>

typedef struct sky_servlet sky_servlet;

#include "bstring.h"
#include "server.h"
#include "tablet.h"


//==============================================================================
//
// Typedefs
//
//==============================================================================

// The various states that the servlet can be in.
typedef enum {
    SKY_SERVLET_STATE_STOPPED = 0,
    SKY_SERVLET_STATE_RUNNING = 1,
} sky_servlet_state_e;

struct sky_servlet {
    uint32_t id;
    sky_servlet_state_e state;
    sky_server *server;
    sky_tablet *tablet;
    bstring name;
    bstring uri;
    void *pull_socket;
    pthread_t thread;
};


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

sky_servlet *sky_servlet_create(sky_server *server, sky_tablet *tablet);

void sky_servlet_free(sky_servlet *servlet);

//--------------------------------------
// State
//--------------------------------------

int sky_servlet_start(sky_servlet *servlet);

#endif
