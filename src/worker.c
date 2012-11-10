#include <inttypes.h>
#include <stdbool.h>
#include <stdlib.h>
#include <sys/time.h>
#include <pthread.h>
#include <zmq.h>

#include "worker.h"
#include "worklet.h"
#include "bstring.h"
#include "sky_zmq.h"
#include "dbg.h"
#include "mem.h"


//==============================================================================
//
// Forward Declarations
//
//==============================================================================

void *sky_worker_run(void *_worker);


//==============================================================================
//
// Global Variables
//
//==============================================================================

// A counter to track the next available worker id. Currently this is not
// thread-safe so workers should only be created in the main thread.
uint64_t next_worker_id = 0;


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

// Creates a worker.
//
// id - The worker identifier.
//
// Returns a reference to the worker.
sky_worker *sky_worker_create()
{
    sky_worker *worker = calloc(1, sizeof(sky_worker)); check_mem(worker);
    worker->id = next_worker_id++;
    return worker;

error:
    sky_worker_free(worker);
    return NULL;
}

// Closes and frees the push sockets on the worker.
//
// worker - The worker.
//
// Returns nothing.
void sky_worker_free_push_sockets(sky_worker *worker)
{
    int rc;
    if(worker) {
        uint32_t i;
        for(i=0; i<worker->push_socket_count; i++) {
            rc = zmq_close(worker->push_sockets[i]);
            check(rc == 0, "Unable to close worker push socket: %d", i);
            worker->push_sockets[i] = NULL;
        }
        free(worker->push_sockets);
        worker->push_sockets = NULL;
        worker->push_socket_count = 0;
    }

error:
    return;
}

// Closes and frees the pull socket on the worker.
//
// worker - The worker.
//
// Returns nothing.
void sky_worker_free_pull_socket(sky_worker *worker)
{
    int rc;
    if(worker) {
        if(worker->pull_socket) {
            rc = zmq_close(worker->pull_socket);
            check(rc == 0, "Unable to close worker pull socket");
            worker->pull_socket = NULL;
        }
        bdestroy(worker->pull_socket_uri);
        worker->pull_socket_uri = NULL;
    }

error:
    return;
}

// Frees the list of servlets the worker is attached to.
//
// worker - The worker.
//
// Returns nothing.
void sky_worker_free_servlets(sky_worker *worker)
{
    if(worker) {
        uint32_t i;
        for(i=0; i<worker->servlet_count; i++) {
            worker->servlets[i] = NULL;
        }
        free(worker->servlets);
        worker->servlets = NULL;
        worker->servlet_count = 0;
    }
}

// Frees a worker from memory.
//
// worker - The worker.
//
// Returns nothing.
void sky_worker_free(sky_worker *worker)
{
    if(worker) {
        sky_worker_free_servlets(worker);
        sky_worker_free_push_sockets(worker);
        sky_worker_free_pull_socket(worker);
        worker->context = NULL;
        if(worker->input) fclose(worker->input);
        worker->input = NULL;
        if(worker->output) fclose(worker->output);
        worker->output = NULL;
        free(worker);
    }
}


//--------------------------------------
// Servlet Management
//--------------------------------------

// Assigns a list of servlets to a worker.
//
// worker   - The worker.
// servlets - An array of servlets.
// count    - The number of servlets.
//
// Returns 0 if successful, otherwise returns -1.
int sky_worker_set_servlets(sky_worker *worker, sky_servlet **servlets,
                            uint32_t count)
{
    check(worker != NULL, "Worker required");
    check(worker->state == SKY_WORKER_STATE_STOPPED, "Worker must be stopped to assign servlets");
    check(servlets != NULL || count == 0, "Servlets required");
    
    // Free existing servlets.
    sky_worker_free_servlets(worker);
    
    // Copy servlets.
    if(count > 0) {
        worker->servlets = calloc(count, sizeof(*worker->servlets));
        check_mem(worker->servlets);

        uint32_t i;
        for(i=0; i<count; i++) {
            worker->servlets[i] = servlets[i];
        }
        worker->servlet_count = count;
    }

    return 0;

error:
    return -1;
}


//--------------------------------------
// State
//--------------------------------------

// Starts a worker and opens push/pull sockets to communicate with the
// servlets.
//
// worker - The worker.
//
// Returns 0 if successful, otherwise returns -1.
int sky_worker_start(sky_worker *worker)
{
    int rc;
    check(worker != NULL, "Worker required");
    check(worker->state == SKY_WORKER_STATE_STOPPED, "Cannot start a running worker");
    check(worker->servlet_count > 0, "Worker must be associated with servlets before starting");
    
    // Allocate space for sockets.
    worker->push_socket_count = worker->servlet_count;
    worker->push_sockets = calloc(worker->push_socket_count, sizeof(*worker->push_sockets));
    check_mem(worker->push_sockets);
    
    // Generate a push socket for each servlet.
    uint32_t i;
    for(i=0; i<worker->push_socket_count; i++) {
        // Create push socket.
        worker->push_sockets[i] = zmq_socket(worker->context, ZMQ_PUSH);
        check(worker->push_sockets[i] != NULL, "Unable to create worker push socket");
        
        // Connect to servlet.
        rc = zmq_connect(worker->push_sockets[i], bdata(worker->servlets[i]->uri));
        check(rc == 0, "Unable to connect worker to servlet");
    }

    // Create pull socket.
    worker->pull_socket = zmq_socket(worker->context, ZMQ_PULL);
    check(worker->pull_socket != NULL, "Unable to create worker pull socket");

    // Bind pull socket.
    worker->pull_socket_uri = bformat("inproc://worker.%lld.pull", worker->id);
    check_mem(worker->pull_socket_uri);
    rc = zmq_bind(worker->pull_socket, bdata(worker->pull_socket_uri));
    check(rc == 0, "Unable to bind worker pull socket");

    // Create the worker thread.
    rc = pthread_create(&worker->thread, NULL, sky_worker_run, worker);
    check(rc == 0, "Unable to create worker thread");

    return 0;

error:
    sky_worker_free_push_sockets(worker);
    sky_worker_free_pull_socket(worker);
    return -1;
}


//--------------------------------------
// Processing
//--------------------------------------

// The worker thread function.
//
// worker - The worker.
//
// Returns 0 if successful.
void *sky_worker_run(void *_worker)
{
    int rc;
    uint32_t i;
    sky_worker *worker = (sky_worker*)_worker;
    sky_worklet *worklet = NULL;
    check(worker != NULL, "Worker required");
    
    // Start benchmark.
    struct timeval tv;
    gettimeofday(&tv, NULL);
    int64_t t0 = (tv.tv_sec*1000) + (tv.tv_usec/1000);

    // Read data from stream.
    if(worker->read != NULL) {
        rc = worker->read(worker, worker->input);
        check(rc == 0, "Worker unable to read from stream");
    }

    // Push a message to each servlet.
    for(i=0; i<worker->push_socket_count; i++) {
        worklet = sky_worklet_create(worker); check_mem(worklet);
        rc = sky_zmq_send_ptr(worker->push_sockets[i], &worklet);
        check(rc == 0, "Worker unable to send worklet");
    }
    
    // Read in one pull message for every push message sent.
    for(i=0; i<worker->push_socket_count; i++) {
        // Receive worker back from servlet.
        rc = sky_zmq_recv_ptr(worker->pull_socket, (void**)&worklet);
        check(rc == 0 && worklet != NULL, "Worker unable to receive worklet");
        
        // Reduce worklet.
        if(worker->reduce != NULL) {
            rc = worker->reduce(worker, worklet->data);
            check(rc == 0, "Worker unable to reduce");
        }
        
        // Free worklet.
        worker->map_free(worklet->data);
        worklet->data = NULL;
        sky_worklet_free(worklet);
        worklet = NULL;
    }
    
    // Output data to stream.
    if(worker->write != NULL) {
        rc = worker->write(worker, worker->output);
        check(rc == 0, "Worker unable to reduce");
    }
    
    // Close streams.
    if(worker->input) fclose(worker->input);
    if(worker->output) fclose(worker->output);

    // End benchmark.
    gettimeofday(&tv, NULL);
    int64_t t1 = (tv.tv_sec*1000) + (tv.tv_usec/1000);
    debug("Worker ran in: %.3f seconds\n", ((float)(t1-t0))/1000);
    
    // Clean up worker.
    worker->free(worker);
    sky_worker_free(worker);
    
    return NULL;

error:
    sky_worker_free(worker);
    return NULL;
}
