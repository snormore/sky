#ifndef _sky_string_h
#define _sky_string_h

#include <inttypes.h>
#include <stdbool.h>

#include "bstring.h"

//==============================================================================
//
// Definitions
//
//==============================================================================

// The Sky string stores information about a fixed length string.
typedef struct {
    int32_t length;
    char *data;
} sky_string;


//==============================================================================
//
// Functions
//
//==============================================================================

//======================================
// Lifecycle
//======================================

sky_string sky_string_create(int32_t length, char *data);


//======================================
// Equality
//======================================

bool sky_string_equals(sky_string *a, sky_string *b);

bool sky_string_bequals(sky_string *a, bstring b);

#endif
