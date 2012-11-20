#include <ctype.h>
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
    check(rc == 0, "Unable to compile Lua script: %s", lua_tostring(*L, -1));

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
    check(rc == 0, "Unable to execute Lua script: %s", lua_tostring(L, -1));

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


//--------------------------------------
// Property File Integration
//--------------------------------------

// Generates the LuaJIT header given a Lua script and a property file. The
// header file is generated based on the property usage of the 'event'
// variable in the script.
//
// source        - The source code of the Lua script.
// property_file - The property file used to lookup properties.
// ret           - A pointer to where the header contents should be returned.
//
// Returns 0 if successful, otherwise returns -1.
int sky_lua_generate_event_struct_decl(bstring source,
                                       sky_property_file *property_file,
                                       bstring *ret)
{
    int rc;
    bstring identifier = NULL;
    assert(source != NULL);
    assert(property_file != NULL);
    assert(ret != NULL);

    // Initialize returned value.
    *ret = NULL;

    // Setup a lookup of properties.
    bool lookup[SKY_PROPERTY_ID_COUNT+1];
    memset(lookup, 0, sizeof(lookup));

    // Loop over every mention of an "event." property.
    int pos;
    struct tagbstring EVENT_DOT_STR = bsStatic("event.");
    while((pos = binstr(source, pos, &EVENT_DOT_STR)) != BSTR_ERR) {
        // Make sure that this is not part of another identifier.
        bool skip = false;
        if(pos > 0 && (isalnum(bchar(source, pos-1)) || bchar(source, pos-1) == '_')) {
            skip = true;
        }
        
        // Move past the "event." string.
        pos += blength(&EVENT_DOT_STR);

        if(!skip) {
            // Read in identifier.
            int i;
            for(i=pos+1; i<blength(source); i++) {
                char ch = bchar(source, i);
                if(!(isalnum(ch) || ch == '_')) {
                    break;
                }
            }
            identifier = bmidstr(source, pos, i-pos); check_mem(identifier);
            if(blength(identifier)) {
                sky_property *property = NULL;
                rc = sky_property_file_find_by_name(property_file, identifier, &property);
                check(rc == 0, "Unable to find property by name: %s", bdata(identifier));
                check(property != NULL, "Property not found: %s", bdata(identifier));
            
                if(!lookup[property->id-SKY_PROPERTY_ID_MIN]) {
                    if(*ret == NULL) {
                        *ret = bfromcstr("");
                        check_mem(*ret);
                    }
                
                    // Append property definition to *ret.
                    switch(property->data_type) {
                        case SKY_DATA_TYPE_STRING: {
                            rc = bformata(*ret, "char %s[];\n", bdata(property->name));
                            break;
                        }
                        case SKY_DATA_TYPE_INT: {
                            rc = bformata(*ret, "int64_t %s;\n", bdata(property->name));
                            break;
                        }
                        case SKY_DATA_TYPE_DOUBLE: {
                            rc = bformata(*ret, "double %s;\n", bdata(property->name));
                            break;
                        }
                        case SKY_DATA_TYPE_BOOLEAN: {
                            rc = bformata(*ret, "bool %s;\n", bdata(property->name));
                            break;
                        }
                        default:{
                            sentinel("Invalid sky lua type: %d", property->data_type);
                        }
                    }
                    check(rc != BSTR_ERR, "Unable to append event property definition: %d", rc);

                    // Flag the property as already processed.
                    lookup[property->id - SKY_PROPERTY_ID_MIN] = true;
                }
            }

            bdestroy(identifier);
            identifier = NULL;
        }
    }

    // If there are any properties found then wrap the whole thing in a struct.
    if(*ret != NULL) {
        rc = bassignformat(*ret, "typedef struct {\n%s} sky_event_t;", bdata(*ret));
        check(rc != BSTR_ERR, "Unable to append event property definition");
    }

    return 0;

error:
    bdestroy(identifier);
    bdestroy(*ret);
    *ret = NULL;
    return -1;
}

