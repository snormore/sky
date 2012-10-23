#ifndef _types_h
#define _types_h

#include <inttypes.h>

#include "bstring.h"

//==============================================================================
//
// Overview
//
//==============================================================================

// This header file includes definitions of the most basic types used in the Sky
// database.


//==============================================================================
//
// Typedefs
//
//==============================================================================

//--------------------------------------
// Timestamp
//--------------------------------------

// Stores the number of microseconds since the epoch (Jan 1, 1970 UTC).
#define sky_timestamp_t int64_t

#define SKY_TIMESTAMP_MIN (INT64_MIN + 1)

#define SKY_TIMESTAMP_MAX INT64_MAX


//--------------------------------------
// Object ID
//--------------------------------------

// Stores an object identifier.
#define sky_object_id_t uint32_t

#define SKY_OBJECT_ID_MIN 1

#define SKY_OBJECT_ID_MAX UINT32_MAX


//--------------------------------------
// Actions
//--------------------------------------

// Stores an action identifier.
#define sky_action_id_t uint16_t


//--------------------------------------
// Properties
//--------------------------------------

// Stores a property identifier.
#define sky_property_id_t int8_t

// The lowest possible property identifier.
#define SKY_PROPERTY_ID_MIN INT8_MIN

// The highest possible property identifier.
#define SKY_PROPERTY_ID_MAX INT8_MAX

// The total possible properties.
#define SKY_PROPERTY_ID_COUNT (SKY_PROPERTY_ID_MAX - SKY_PROPERTY_ID_MIN)


#endif
