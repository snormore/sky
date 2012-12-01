#include <stdlib.h>
#include <inttypes.h>

#include "dbg.h"
#include "bstring.h"
#include "action.h"
#include "minipack.h"

//==============================================================================
//
// Definitions
//
//==============================================================================

#define SKY_ACTION_KEY_COUNT 2

struct tagbstring SKY_ACTION_ID_STR = bsStatic("id");
struct tagbstring SKY_ACTION_NAME_STR = bsStatic("name");


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

// Creates a reference to an action.
// 
// Returns a reference to the new action if successful.
sky_action *sky_action_create()
{
    sky_action *action = calloc(sizeof(*action), 1); check_mem(action);
    return action;
    
error:
    sky_action_free(action);
    return NULL;
}

// Removes an action reference from memory.
//
// action - The action to free.
void sky_action_free(sky_action *action)
{
    if(action) {
        action->action_file = NULL;
        action->id = 0;
        if(action->name) bdestroy(action->name);
        action->name = NULL;
        free(action);
    }
}


//--------------------------------------
// Serialization
//--------------------------------------

// Calculates the total number of bytes needed to store the action.
//
// action - The action.
//
// Returns the number of bytes required to store the action.
size_t sky_action_sizeof(sky_action *action)
{
    size_t sz = 0;
    sz += minipack_sizeof_map(SKY_ACTION_KEY_COUNT);
    sz += minipack_sizeof_raw((&SKY_ACTION_ID_STR)->slen) + (&SKY_ACTION_ID_STR)->slen;
    sz += minipack_sizeof_uint(action->id);
    sz += minipack_sizeof_raw((&SKY_ACTION_NAME_STR)->slen) + (&SKY_ACTION_NAME_STR)->slen;
    sz += blength(action->name);
    return sz;
}

// Serializes an action to a file stream.
//
// action - The action.
// file   - The file stream to write to.
//
// Returns 0 if successful, otherwise returns -1.
int sky_action_pack(sky_action *action, FILE *file)
{
    size_t sz;
    check(action != NULL, "Action required");
    check(file != NULL, "File stream required");

    // Map
    minipack_fwrite_map(file, SKY_ACTION_KEY_COUNT, &sz);
    check(sz > 0, "Unable to write map");
    
    // ID
    check(sky_minipack_fwrite_bstring(file, &SKY_ACTION_ID_STR) == 0, "Unable to write id key");
    minipack_fwrite_uint(file, action->id, &sz);
    check(sz > 0, "Unable to write id value");

    // Name
    check(sky_minipack_fwrite_bstring(file, &SKY_ACTION_NAME_STR) == 0, "Unable to write name key");
    check(sky_minipack_fwrite_bstring(file, action->name) == 0, "Unable to write name value");

    return 0;

error:
    return -1;
}

// Deserializes an action from a file stream.
//
// action - The action.
// file   - The file stream to read from.
//
// Returns 0 if successful, otherwise returns -1.
int sky_action_unpack(sky_action *action, FILE *file)
{
    int rc;
    size_t sz;
    bstring key = NULL;
    check(action != NULL, "Message required");
    check(file != NULL, "File stream required");

    // Map
    uint32_t map_length = minipack_fread_map(file, &sz);
    check(sz > 0, "Unable to read map");
    
    // Map items
    uint32_t i;
    for(i=0; i<map_length; i++) {
        rc = sky_minipack_fread_bstring(file, &key);
        check(rc == 0, "Unable to read map key");
        
        if(biseq(key, &SKY_ACTION_ID_STR)) {
            action->id = (sky_action_id_t)minipack_fread_uint(file, &sz);
            check(sz > 0, "Unable to read action id");
        }
        else if(biseq(key, &SKY_ACTION_NAME_STR)) {
            rc = sky_minipack_fread_bstring(file, &action->name);
            check(rc == 0, "Unable to read action id");
        }

        bdestroy(key);
    }
    
    return 0;

error:
    bdestroy(key);
    return -1;
}

