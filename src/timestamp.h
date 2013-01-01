#ifndef _timestamp_h
#define _timestamp_h

#include <inttypes.h>

#include "bstring.h"
#include "types.h"


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Parsing
//--------------------------------------

int sky_timestamp_parse(bstring str, sky_timestamp_t *ret);

int sky_timestamp_now(sky_timestamp_t *ret);

//--------------------------------------
// Shifting
//--------------------------------------

int64_t sky_timestamp_shift(int64_t value);

int64_t sky_timestamp_unshift(int64_t value);

int64_t sky_timestamp_to_seconds(int64_t value);

#endif

