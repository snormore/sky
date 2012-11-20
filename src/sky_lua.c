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

    // Call once to make the functions available.
    lua_call(*L, 0, 0);

    return 0;

error:
    lua_close(*L);
    *L = NULL;
    return -1;
}

// Generates a special header to allow LuaJIT and Sky to interact and then
// initializes the script.
//
// source - The script source code.
// table  - The table used to generate the header.
// L      - A reference to where the new Lua state should be returned.
//
// Returns 0 if successful, otherwise returns -1.
int sky_lua_initscript_with_table(bstring source, sky_table *table,
                                  lua_State **L)
{
    int rc;
    bstring header = NULL;
    bstring new_source = NULL;
    assert(source != NULL);
    assert(table != NULL);
    assert(L != NULL);
    
    // Generate header.
    rc = sky_lua_generate_header(source, table, &header);
    check(rc == 0, "Unable to generate header");
    new_source = bformat("%s%s", bdata(header), bdata(source)); check_mem(new_source);
    
    // Initialize script.
    rc = sky_lua_initscript(new_source, L);
    check(rc == 0, "Unable to initialize Lua script");

    return 0;

error:
    *L = NULL;
    bdestroy(new_source);
    return -1;
}

//--------------------------------------
// MessagePack
//--------------------------------------

// Converts the top of the stack to a MessagePack encoded bstring.
//
// L     - The lua state.
// ret   - A pointer to where the msgpack result should be returned.
//
// Returns 0 if successful, otherwise returns -1.
int sky_lua_to_msgpack(lua_State *L, bstring *ret)
{
    int rc;
    assert(L != NULL);
    assert(ret != NULL);

    // Initialize returned value.
    *ret = NULL;
    
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
// source - The source code of the Lua script.
// table  - The table used for generation.
// ret    - A pointer to where the header contents should be returned.
//
// Returns 0 if successful, otherwise returns -1.
int sky_lua_generate_header(bstring source, sky_table *table, bstring *ret)
{
    int rc;
    bstring event_decl = NULL;
    assert(source != NULL);
    assert(table != NULL);
    assert(ret != NULL);

    // Initialize returned value.
    *ret = NULL;

    // Generate sky_lua_event_t declaration.
    rc = sky_lua_generate_event_struct_decl(source, table->property_file, &event_decl);
    check(rc == 0, "Unable to generate lua event declaration");
    
    // Generate full header.
    *ret = bformat(
        "-- SKY GENERATED CODE BEGIN --\n"
        "local ffi = require(\"ffi\")\n"
        "ffi.cdef([[\n"
        "typedef struct sky_lua_path_iterator_t sky_lua_path_iterator_t;\n"
        "typedef struct sky_lua_cursor_t sky_lua_cursor_t;\n"
        "%s\n"
        "\n"
        "sky_lua_cursor_t *sky_lua_path_iterator_next(sky_lua_path_iterator_t *);\n"
        "sky_lua_event_t *sky_lua_cursor_next(sky_lua_cursor_t *);\n"
        "]])\n"
        "-- SKY GENERATED CODE END --\n\n",
        bdata(event_decl)
    );
    check_mem(*ret);

    bdestroy(event_decl);
    return 0;

error:
    bdestroy(event_decl);
    bdestroy(*ret);
    *ret = NULL;
    return -1;
}


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
    int pos = 0;
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
                            rc = bformata(*ret, "  char %s[];\n", bdata(property->name));
                            break;
                        }
                        case SKY_DATA_TYPE_INT: {
                            rc = bformata(*ret, "  int64_t %s;\n", bdata(property->name));
                            break;
                        }
                        case SKY_DATA_TYPE_DOUBLE: {
                            rc = bformata(*ret, "  double %s;\n", bdata(property->name));
                            break;
                        }
                        case SKY_DATA_TYPE_BOOLEAN: {
                            rc = bformata(*ret, "  bool %s;\n", bdata(property->name));
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

    // Wrap properties in a struct.
    if(*ret != NULL) {
        bassignformat(*ret, "typedef struct {\n%s} sky_lua_event_t;", bdata(*ret));
        check_mem(*ret);
    }
    else {
        *ret = bfromcstr("typedef struct sky_lua_event_t sky_lua_event_t;"); check_mem(*ret);
    }

    return 0;

error:
    bdestroy(identifier);
    bdestroy(*ret);
    *ret = NULL;
    return -1;
}

