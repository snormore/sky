#include <inttypes.h>
#include <stdbool.h>
#include <stdlib.h>
#include <pthread.h>
#include <zmq.h>

#include "servlet.h"
#include "worker.h"
#include "worklet.h"
#include "sky_zmq.h"
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
    check(tablet->table != NULL, "Tablet must be attached to a table");
    check(blength(tablet->table->name), "Servlet must be attached to a named table");

    servlet = calloc(1, sizeof(sky_servlet)); check_mem(servlet);
    servlet->id = server->next_servlet_id++;
    servlet->name = bformat("%s.%d.%d", bdata(tablet->table->name), server->id, tablet->index);
    check_mem(servlet->name);
    servlet->uri = bformat("inproc://servlet.%s", bdata(servlet->name));
    check_mem(servlet->uri);
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
        if(servlet->name) bdestroy(servlet->name);
        servlet->name = NULL;
        if(servlet->uri) bdestroy(servlet->uri);
        servlet->uri = NULL;
        if(servlet->pull_socket) zmq_close(servlet->pull_socket);
        servlet->pull_socket = NULL;
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
    check(servlet->state == SKY_SERVLET_STATE_STOPPED, "Servlet already running");

    debug("servlet.start.1: %s", bdata(servlet->uri));

    // Create a listening queue.
    void *context = servlet->server->context;
    servlet->pull_socket = zmq_socket(context, ZMQ_PULL);
    check_mem(servlet->pull_socket);
    rc = zmq_bind(servlet->pull_socket, bdata(servlet->uri));
    check(rc == 0, "Unable to connect servlet pull socket");
    
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
    sky_worklet *worklet = NULL;
    sky_servlet *servlet = (sky_servlet *)_servlet;
    check(servlet != NULL, "Servlet required");
    
    // Read in messages from pull socket.
    void *context = servlet->server->context;
    while(true) {
        // Read in worklet.
        rc = sky_zmq_recv_ptr(servlet->pull_socket, (void**)(&worklet));
        check(rc == 0, "Unable to receive worklet message");
        
        debug("servlet.recv.1: %p", worklet);

        // If worklet is NULL then stop the servlet.
        if(worklet == NULL) {
            break;
        }
        
        // Process worklet.
        sky_worker *worker = worklet->worker;
        worker->map(worker, servlet->tablet, &worklet->data);
        debug("servlet.postmap.1: %p | %p", worklet, worklet->data);
        
        // Connect back to worker.
        void *push_socket = zmq_socket(context, ZMQ_PUSH); check_mem(push_socket);
        rc = zmq_connect(push_socket, bdata(worker->pull_socket_uri));
        check(rc == 0, "Unable to connect servlet push socket");

        // Send back worklet.
        rc = sky_zmq_send_ptr(push_socket, (void*)(&worklet));
        check(rc == 0, "Unable to send worklet message");
    }
    
    debug("servlet.end: %s", bdata(servlet->name));
    sky_servlet_free(servlet);
    return NULL;

error:
    return NULL;
}
