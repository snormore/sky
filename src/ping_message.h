#ifndef _sky_ping_message_h
#define _sky_ping_message_h

#include <inttypes.h>
#include <stdbool.h>
#include <netinet/in.h>

#include "bstring.h"
#include "message_handler.h"
#include "table.h"
#include "event.h"


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Message Handler
//--------------------------------------

sky_message_handler *sky_ping_message_handler_create();

int sky_ping_message_process(sky_server *server,
    sky_message_header *header, sky_table *table, FILE *input, FILE *output);

#endif
