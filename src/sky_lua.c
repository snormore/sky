#include <assert.h>

#include "sky_lua.h"
#include "dbg.h"


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Initialization
//--------------------------------------

// Initializes a state and loads the source of the Lua script.
//
// source - The script source code.
// L      - A reference to where the new Lua state should be returned.
//
// Returns 0 if successful, otherwise returns -1.
int sky_lua_initscript(bstring source, lua_State **L)
{
    int rc;
    assert(source != NULL);
    assert(L != NULL);
    
    // Load Lua with standard library.
    *L = luaL_newstate(); check_mem(L);
    luaL_openlibs(*L);
    
    // Load Lua msgpack library.
    rc = luaopen_cmsgpack(*L);
    check(rc == 1, "Unable to load lua-cmsgpack");
    
    // Compile lua script.
    rc = luaL_loadstring(*L, bdata(source));
    check(rc == 0, "Unable to compile Lua script");

    return 0;

error:
    lua_close(*L);
    *L = NULL;
    return -1;
}

//--------------------------------------
// Execution
//--------------------------------------

// Executes a function and returns a MessagePack encoded response.
//
// L     - The lua state.
// nargs - The number of arguments.
// ret   - A pointer to where the msgpack result should be returned.
//
// Returns 0 if successful, otherwise returns -1.
int sky_lua_pcall_msgpack(lua_State *L, int nargs, bstring *ret)
{
    int rc;
    assert(L != NULL);
    assert(ret != NULL);

    // Initialize returned value.
    *ret = NULL;
    
    // Execute function.
    rc = lua_pcall(L, nargs, 1, 0);
    check(rc == 0, "Unable to execute Lua script");

    // Encode result as msgpack.
    rc = mp_pack(L);
    check(rc == 1, "Unable to msgpack encode Lua result");
    *ret = (void*)bfromcstr(lua_tostring(L, -1));
    check_mem(*ret);

    return 0;

error:
    bdestroy(*ret);
    *ret = NULL;
    return -1;
}

