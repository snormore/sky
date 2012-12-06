#include <stdlib.h>

#include "sky_string.h"
#include "dbg.h"

//==============================================================================
//
// Functions
//
//==============================================================================

//======================================
// Lifecycle
//======================================

// Creates a fixed-length string.
//
// length - The number of characters in the string.
// data   - A pointer to the character data.
//
// Returns a fixed-length string
sky_string sky_string_create(int32_t length, char *data)
{
    sky_string string;
    string.length = length;
    string.data = data;
    return string;
}


//======================================
// Equality
//======================================

// Evaluates whether the string length and string contents are equal between
// two strings.
//
// a - The first string.
// b - The second string.
//
// Returns true if the strings are equal. Otherwise returns false.
bool sky_string_equals(sky_string *a, sky_string *b)
{
    // Check length first. We could have two strings pointing at the same data
    // but one string only points to a subset of the data.
    if(a->length == b->length) {
        // If both strings point at the same data it must be the same string.
        if(a->data == b->data) {
            return true;
        }
        // If these are different data pointers then check the contents of
        // the data pointers.
        else {
            return (memcmp(a->data, b->data, a->length)) == 0;
        }
    }
    
    return false;
}

// Evaluates whether a sky string is equal to a bstring.
//
// a - The sky string.
// b - The bstring.
//
// Returns true if the strings are equal. Otherwise returns false.
bool sky_string_bequals(sky_string *a, bstring b)
{
    // Check length first. We could have two strings pointing at the same data
    // but one string only points to a subset of the data.
    int32_t length;
    memmove(&length, &a->length, sizeof(length));
    if(length == blength(b)) {
        return (memcmp(a->data, bdatae(b, ""), length)) == 0;
    }
    
    return false;
}
