#include <inttypes.h>
#include <stdbool.h>
#include <stdlib.h>
#include <pthread.h>
#include <zmq.h>

#include "servlet.h"
#include "worker.h"
#include "dbg.h"


//==============================================================================
//
// Forward Declarations
//
//==============================================================================

void *sky_servlet_run(void *_servlet);


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

// Creates a servlet.
//
// server - The server the servlet belongs to.
// tablet - The tablet the servlet is attached to.
//
// Returns a reference to the servlet.
sky_servlet *sky_servlet_create(sky_server *server, sky_tablet *tablet)
{
    sky_servlet *servlet = NULL;
    check(server != NULL, "Server required");
    check(tablet != NULL, "Tablet required");

    servlet = calloc(1, sizeof(sky_servlet)); check_mem(servlet);
    servlet->server = server;
    servlet->tablet = tablet;
    
    return servlet;

error:
    sky_servlet_free(servlet);
    return NULL;
}

// Frees a servlet from memory.
//
// servlet - The servlet.
//
// Returns nothing.
void sky_servlet_free(sky_servlet *servlet)
{
    if(servlet) {
        servlet->tablet = NULL;
        servlet->server = NULL;
        if(servlet->uri) bdestroy(servlet->uri);
        servlet->uri = NULL;
        free(servlet);
    }
}


//--------------------------------------
// State
//--------------------------------------

// Starts a servlet. Once a servlet is started, it cannot be stopped until the
// server has stopped.
//
// servlet - The servlet.
//
// Returns 0 if successful, otherwise returns -1.
int sky_servlet_start(sky_servlet *servlet)
{
    int rc;
    check(servlet != NULL, "Servlet required");
    check(servlet->uri != NULL, "Servlet URI required");
    check(servlet->state != SKY_SERVLET_STATE_STOPPED, "Servlet already running");

    // Create the worker thread.
    rc = pthread_create(&servlet->thread, NULL, sky_servlet_run, (void*)servlet);
    check(rc == 0, "Unable to create servlet thread");

    // Update servlet state.
    servlet->state = SKY_SERVLET_STATE_RUNNING;

    return 0;

error:
    return -1;
}


//--------------------------------------
// Processing
//--------------------------------------

// The worker function for the servlet thread.
//
// servlet - The servlet.
//
// Returns NULL.
void *sky_servlet_run(void *_servlet)
{
    int rc;
    sky_servlet *servlet = (sky_servlet *)_servlet;
    check(servlet != NULL, "Servlet required");
    
    // Create a listening queue.
    void *context = servlet->server->context;
    void *socket = zmq_socket(context, ZMQ_PULL); check_mem(socket);
    rc = zmq_connect(socket, bdata(servlet->uri));
    check(rc == 0, "Unable to connect servlet pull socket");
    
    // Read in messages from pull socket.
    while(true) {
        // TODO: Read worklet and process.
    }
    
    return NULL;

error:
    return NULL;
}
