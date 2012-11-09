#include <inttypes.h>
#include <stdbool.h>
#include <stdlib.h>
#include <zmq.h>

#include "sky_zmq.h"
#include "dbg.h"
#include "mem.h"


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Send / Receive
//--------------------------------------

// Sends a pointer to a socket.
//
// socket - The socket.
// ptr    - The pointer to send.
//
// Returns 0 if successful, otherwise returns -1.
int sky_zmq_send_ptr(void *socket, void *ptr)
{
    int rc;
    
    // Initialize message.
    zmq_msg_t message;
    rc = zmq_msg_init_size(&message, sizeof(ptr));
    check(rc == 0, "Unable to initialize zmq message");
    memcpy(zmq_msg_data(&message), ptr, sizeof(ptr));

    // Send message.
    rc = zmq_msg_send(&message, socket, 0);
    check(rc != -1, "Unable to send zmq message");
    
    // Close message.
    zmq_msg_close(&message);

    return 0;
    
error:
    zmq_msg_close(&message);
    return -1;
}

// Receives a pointer from a socket.
//
// socket - The socket.
// ptr    - A pointer to where the pointer should be returned.
//
// Returns 0 if successful, otherwise returns -1.
int sky_zmq_recv_ptr(void *socket, void **ptr)
{
    int rc;

    // Initialize return value.
    *ptr = NULL;

    // Initialize the message.
    zmq_msg_t message;
    rc = zmq_msg_init(&message);
    check(rc == 0, "Unable to initialize zmq message");
    
    // Receive the message.
    int size = zmq_msg_recv(&message, socket, 0);
    check(size == sizeof(*ptr), "Unable to receive zmq pointer message: recv %d, exp %ld", size, sizeof(*ptr));
    memcpy(ptr, zmq_msg_data(&message), size);
    
    // Close message.
    zmq_msg_close(&message);

    return 0;
    
error:
    *ptr = NULL;
    zmq_msg_close(&message);
    return -1;
}


//--------------------------------------
// Send / Receive
//--------------------------------------

// Sets a socket to not linger.
//
// socket - The socket.
//
// Returns 0 if successful, otherwise returns -1.
int sky_zmq_no_linger(void *socket)
{
    int zero = 0;
    int rc = zmq_setsockopt(socket, ZMQ_LINGER, &zero, sizeof(zero));
    check(rc == 0, "Unable to set no linger on socket");

    return 0;
    
error:
    return -1;
}
