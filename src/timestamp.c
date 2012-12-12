#include <stdlib.h>
#include <sys/time.h>

#include "timestamp.h"
#include "dbg.h"

//==============================================================================
//
// Timestamp Parsing
//
//==============================================================================

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
    *ret = value * 1000000;
    
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
