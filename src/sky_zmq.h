#ifndef _sky_zmq_h
#define _sky_zmq_h

#include <stdio.h>
#include <inttypes.h>
#include <stdbool.h>

#include "bstring.h"


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Send / Receive
//--------------------------------------

int sky_zmq_send_ptr(void *socket, void *ptr);

int sky_zmq_recv_ptr(void *socket, void **ptr);

//--------------------------------------
// Utility
//--------------------------------------

int sky_zmq_no_linger(void *socket);

#endif
