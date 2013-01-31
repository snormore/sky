#include "bstring.h"

//==============================================================================
//
// Constants
//
//==============================================================================

const uint32_t FNV_PRIME = 16777619;
const uint32_t FNV_OFFSET_BASIS = 2166136261;


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Object ID Hash
//--------------------------------------

// Calculates a numeric hash for a bstring. Thanks for the implementation
// from Zed Shaw (http://c.learncodethehardway.org/book/ex38.html).
//
// str - The string to calculate the hash code for.
//
// Returns a numeric hash code.
uint32_t sky_bstring_fnv1a(bstring str)
{
    uint32_t hash = FNV_OFFSET_BASIS;
    int i=0;
    for(i=0; i<blength(str); i++) {
        hash ^= bchare(str, i, 0);
        hash *= FNV_PRIME;
    }

    return hash;
}
