#include <stdlib.h>
#include <sys/time.h>

#include "timestamp.h"
#include "dbg.h"

//==============================================================================
//
// Constants
//
//==============================================================================

// The number of microseconds per second.
#define USEC_PER_SEC        1000000

// A bit-mask to extract the microseconds from a Sky timestamp.
#define USEC_MASK           0xFFFFF

// The number of bits that seconds are shifted over in a timestamp.
#define SECONDS_BIT_OFFSET  20


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Parsing
//--------------------------------------

// Parses a timestamp from a C string. The return value is the number of
// microseconds before or after the epoch (Jan 1, 1970).
// 
// NOTE: Parsing seems to only work back to around the first decade of the
//       1900's. Need to investigate further why this is.
// 
// str - The string containing an ISO 8601 formatted date.
int sky_timestamp_parse(bstring str, sky_timestamp_t *ret)
{
    char *tz = NULL;

    // Validate string.
    if(str == NULL) {
        return -1;
    }
    
    // Parse date.
    struct tm tp; memset(&tp, 0, sizeof(tp));
    char *ch;
    ch = strptime(bdata(str), "%Y-%m-%dT%H:%M:%SZ", &tp);
    check(ch != NULL, "Unable to parse timestamp");
    
    // Convert to microseconds since epoch in UTC.
    char buffer[100];
    tz = getenv("TZ");
    setenv("TZ","",1);
    tzset();
    strftime(buffer, 100, "%s", &tp);
    if(tz) {
        setenv("TZ",tz,1);
    }
    else {
        unsetenv("TZ");
    }

    // Convert to an integer.
    sky_timestamp_t value = atoll(buffer);
    value = value * 1000000;
    
    // Convert to a Sky timestamp and return.
    *ret = value;
    
    return 0;

error:
    return -1;
}

// Returns the number of microseconds since the epoch.
// 
// ret - The reference to the variable that will be assigned the timestamp.
int sky_timestamp_now(sky_timestamp_t *ret)
{
    struct timeval tv;
    check(gettimeofday(&tv, NULL) == 0, "Cannot obtain current time");
    *ret = (tv.tv_sec*1000000) + (tv.tv_usec);
    return 0;

error:
    return -1;
}


//--------------------------------------
// Shifting
//--------------------------------------

// Converts a timestamp from the number of microseconds since the epoch to
// a bit-shifted Sky timestamp.
//
// value - Microseconds since the unix epoch.
//
// Returns a bit-shifted Sky timestamp.
int64_t sky_timestamp_shift(int64_t value)
{
    int64_t usec = value % USEC_PER_SEC;
    int64_t sec  = (value / USEC_PER_SEC);
    
    return (sec << SECONDS_BIT_OFFSET) + usec;
}

// Converts a bit-shifted Sky timestamp to the number of microseconds since
// the Unix epoch.
//
// value - Sky timestamp.
//
// Returns the number of microseconds since the Unix epoch.
int64_t sky_timestamp_unshift(int64_t value)
{
    int64_t usec = value & USEC_MASK;
    int64_t sec  = value >> SECONDS_BIT_OFFSET;
    
    return (sec * USEC_PER_SEC) + usec;
}

// Converts a bit-shifted Sky timestamp to seconds since the epoch.
//
// value - Sky timestamp.
//
// Returns the number of seconds since the Unix epoch.
int64_t sky_timestamp_to_seconds(int64_t value)
{
    return (value >> SECONDS_BIT_OFFSET);
}

