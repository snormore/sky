#include <stdlib.h>

#include "types.h"
#include "sky_string.h"
#include "dbg.h"


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Data Types
//--------------------------------------

// Looks up a data type by name.
// 
// name - The name of the data type.
//
// Returns an enumerator for the data type.
sky_data_type_e sky_data_type_to_enum(bstring name)
{
    if(biseqcstr(name, "String") == 1) {
        return SKY_DATA_TYPE_STRING;
    }
    else if(biseqcstr(name, "Int") == 1) {
        return SKY_DATA_TYPE_INT;
    }
    else if(biseqcstr(name, "Double") == 1) {
        return SKY_DATA_TYPE_DOUBLE;
    }
    else if(biseqcstr(name, "Boolean") == 1) {
        return SKY_DATA_TYPE_BOOLEAN;
    }
    else {
        return SKY_DATA_TYPE_NONE;
    }
}

// Creates a string to represent a data type enum.
// 
// data_type - The enumerated data type value.
//
// Returns a string to represent the data type.
bstring sky_data_type_to_str(sky_data_type_e data_type)
{
    switch(data_type) {
        case SKY_DATA_TYPE_STRING: return bfromcstr("String");
        case SKY_DATA_TYPE_INT: return bfromcstr("Int");
        case SKY_DATA_TYPE_DOUBLE: return bfromcstr("Double");
        case SKY_DATA_TYPE_BOOLEAN: return bfromcstr("Boolean");
        default: return NULL;
    }
}

// Determines the size of the data type, in bytes.
// 
// data_type - The enumerated data type value.
//
// Returns the number of bytes required to store the data type.
size_t sky_data_type_sizeof(sky_data_type_e data_type)
{
    switch(data_type) {
        case SKY_DATA_TYPE_STRING: return sizeof(sky_string);
        case SKY_DATA_TYPE_INT: return sizeof(int64_t);
        case SKY_DATA_TYPE_DOUBLE: return sizeof(double);
        case SKY_DATA_TYPE_BOOLEAN: return sizeof(bool);
        default: return 0;
    }
}

