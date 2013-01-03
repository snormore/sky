#include <ctype.h>
#include <assert.h>

#include "sky_lua.h"
#include "path_iterator.h"
#include "cursor.h"
#include "dbg.h"
#include "mem.h"


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
    *L = luaL_newstate(); check_mem(*L);
    luaL_openlibs(*L);
    
    // Load Lua msgpack library.
    rc = luaopen_cmsgpack(*L);
    check(rc == 1, "Unable to load lua-cmsgpack");
    
    //debug("--SOURCE--\n%s", bdata(source));
    
    // Compile lua script.
    rc = luaL_loadstring(*L, bdata(source));
    check(rc == 0, "Unable to compile Lua script: %s", lua_tostring(*L, -1));

    // Call once to make the functions available.
    rc = lua_pcall(*L, 0, 0, 0);
    check(rc == 0, "Unable to initialize lua script");

    return 0;

error:
    if(*L != NULL) lua_close(*L);
    *L = NULL;
    return -1;
}

// Generates a special header to allow LuaJIT and Sky to interact and then
// initializes the script.
//
// source     - The script source code.
// table      - The table used to generate the header.
// descriptor - The descriptor to initialize with the script.
// L          - A reference to where the new Lua state should be returned.
//
// Returns 0 if successful, otherwise returns -1.
int sky_lua_initscript_with_table(bstring source, sky_table *table,
                                  sky_data_descriptor *descriptor,
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

    // Initialize data descriptor.
    if(descriptor != NULL) {
        descriptor->int_type = SKY_DATA_DESCRIPTOR_INT32;
        lua_getglobal(*L, "sky_init_descriptor");
        lua_pushlightuserdata(*L, descriptor);
        rc = lua_pcall(*L, 1, 0, 0);
        check(rc == 0, "Lua error while initializing descriptor: %s", lua_tostring(*L, -1));
    }

    bdestroy(new_source);
    bdestroy(header);
    return 0;

error:
    if(*L != NULL) {
        lua_pop(*L, 1);
        lua_close(*L);
        *L = NULL;
    }
    bdestroy(new_source);
    bdestroy(header);
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
int sky_lua_msgpack_pack(lua_State *L, bstring *ret)
{
    int rc;
    assert(L != NULL);
    assert(ret != NULL);

    // Initialize returned value.
    *ret = NULL;
    
    // Encode result as msgpack.
    rc = mp_pack(L);
    check(rc == 1, "Unable to msgpack encode Lua result");

    size_t sz;
    const char *str = lua_tolstring(L, -1, &sz);
    *ret = blk2bstr(str, sz);
    check_mem(*ret);

    lua_pop(L, 1);

    return 0;

error:
    bdestroy(*ret);
    *ret = NULL;
    return -1;
}

// Converts a MessagePack encoded string to a Lua object on the top of the
// stack.
//
// L    - The lua state.
// data - A pointer to where the msgpack result should be returned.
//
// Returns 0 if successful, otherwise returns -1.
int sky_lua_msgpack_unpack(lua_State *L, bstring data)
{
    int rc;
    assert(L != NULL);
    assert(data != NULL);
    
    // Push msgpack value to the top of the stack.
    lua_pushlstring(L, bdata(data), blength(data));
    
    // Decode into Lua object.
    rc = mp_unpack(L);
    check(rc == 1, "Unable to decode msgpack data into Lua object");

    lua_remove(L, -2);

    return 0;

error:
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
    bstring init_descriptor_func = NULL;
    assert(source != NULL);
    assert(table != NULL);
    assert(ret != NULL);

    // Initialize returned value.
    *ret = NULL;

    // Generate sky_lua_event_t declaration.
    rc = sky_lua_generate_event_info(source, table->property_file, &event_decl, &init_descriptor_func);
    check(rc == 0, "Unable to generate lua event header");
    
    // Generate full header.
    *ret = bformat(
        "-- SKY GENERATED CODE BEGIN --\n"
        "local ffi = require(\"ffi\")\n"
        "ffi.cdef([[\n"
        "typedef struct sky_data_descriptor_t sky_data_descriptor_t;\n"
        "typedef struct sky_path_iterator_t sky_path_iterator_t;\n"
        "typedef struct sky_cursor_t sky_cursor_t;\n"
        "%s\n"
        "\n"
        "int sky_data_descriptor_set_data_sz(sky_data_descriptor_t *descriptor, uint32_t sz);\n"
        "int sky_data_descriptor_set_timestamp_offset(sky_data_descriptor_t *descriptor, uint32_t offset);\n"
        "int sky_data_descriptor_set_ts_offset(sky_data_descriptor_t *descriptor, uint32_t offset);\n"
        "int sky_data_descriptor_set_action_id_offset(sky_data_descriptor_t *descriptor, uint32_t offset);\n"
        "int sky_data_descriptor_set_property(sky_data_descriptor_t *descriptor, int8_t property_id, uint32_t offset, int data_type);\n"
        "\n"
        "bool sky_path_iterator_eof(sky_path_iterator_t *);\n"
        "void sky_path_iterator_next(sky_path_iterator_t *);\n"
        "sky_cursor_t *sky_lua_path_iterator_get_cursor(sky_path_iterator_t *);\n"
        "\n"
        "bool sky_cursor_eof(sky_cursor_t *);\n"
        "bool sky_lua_cursor_next_event(sky_cursor_t *);\n"
        "sky_lua_event_t *sky_lua_cursor_get_event(sky_cursor_t *);\n"
        "]])\n"
        "ffi.metatype(\"sky_data_descriptor_t\", {\n"
        "  __index = {\n"
        "    set_data_sz = function(descriptor, sz) return ffi.C.sky_data_descriptor_set_data_sz(descriptor, sz) end,\n"
        "    set_timestamp_offset = function(descriptor, offset) return ffi.C.sky_data_descriptor_set_timestamp_offset(descriptor, offset) end,\n"
        "    set_ts_offset = function(descriptor, offset) return ffi.C.sky_data_descriptor_set_ts_offset(descriptor, offset) end,\n"
        "    set_action_id_offset = function(descriptor, offset) return ffi.C.sky_data_descriptor_set_action_id_offset(descriptor, offset) end,\n"
        "    set_property = function(descriptor, property_id, offset, data_type) return ffi.C.sky_data_descriptor_set_property(descriptor, property_id, offset, data_type) end,\n"
        "  }\n"
        "})\n"
        "ffi.metatype(\"sky_path_iterator_t\", {\n"
        "  __index = {\n"
        "    eof = function(iterator) return ffi.C.sky_path_iterator_eof(iterator) end,\n"
        "    cursor = function(iterator) return ffi.C.sky_lua_path_iterator_get_cursor(iterator) end,\n"
        "    next = function(iterator) return ffi.C.sky_path_iterator_next(iterator) end,\n"
        "  }\n"
        "})\n"
        "ffi.metatype(\"sky_cursor_t\", {\n"
        "  __index = {\n"
        "    eof = function(cursor) return ffi.C.sky_cursor_eof(cursor) end,\n"
        "    event = function(cursor) return ffi.C.sky_lua_cursor_get_event(cursor) end,\n"
        "    next = function(cursor) return ffi.C.sky_lua_cursor_next_event(cursor) end,\n"
        "  }\n"
        "})\n"
        "\n"
        "%s\n"
        "function sky_map_all(_iterator)\n"
        "  iterator = ffi.cast(\"sky_path_iterator_t*\", _iterator)\n"
        "  data = {}\n"
        "  while not iterator:eof() do\n"
        "    cursor = iterator:cursor()\n"
        "    map(cursor, data)\n"
        "    iterator:next()\n"
        "  end\n"
        "  return data\n"
        "end\n"
        "\n"
        "-- Automatically deep merges two tables. Numeric values will be summed.\n"
        "-- Strings and booleans are first come first serve.\n"
        "function sky_deep_merge(a, b)\n"
        "  if a == nil or b == nil then return end\n"
        "  \n"
        "  for k,v in pairs(b) do\n"
        "    vtype = type(v)\n"
        "    if vtype == 'table' then\n"
        "      a[k] = {}\n"
        "      sky_deep_merge(a[k], v)\n"
        "    elseif vtype == 'number' then\n"
        "      if a[k] == nil then\n"
        "        a[k] = v\n"
        "      else\n"
        "        a[k] = a[k] + v\n"
        "      end\n"
        "    elseif (vtype == 'string' or vtype == 'boolean') and a[k] == nil then\n"
        "      a[k] = v\n"
        "    end\n"
        "  end\n"
        "end\n"
        "\n"
        "-- The default reduce implementation simply deep merges the results.\n"
        "function reduce(results, data)\n"
        "  sky_deep_merge(results, data)\n"
        "  return results\n"
        "end\n"
        "-- SKY GENERATED CODE END --\n"
        ,
        bdata(event_decl),
        bdata(init_descriptor_func)
    );
    check_mem(*ret);

    bdestroy(event_decl);
    bdestroy(init_descriptor_func);
    return 0;

error:
    bdestroy(event_decl);
    bdestroy(init_descriptor_func);
    bdestroy(*ret);
    *ret = NULL;
    return -1;
}


// Generates the LuaJIT header given a Lua script and a property file. The
// header file is generated based on the property usage of the 'event'
// variable in the script.
//
// source          - The source code of the Lua script.
// property_file   - The property file used to lookup properties.
// event_decl      - A pointer to where the struct def should be returned.
// init_descriptor_func - A pointer to where the descriptor init function should be returned.
//
// Returns 0 if successful, otherwise returns -1.
int sky_lua_generate_event_info(bstring source,
                                sky_property_file *property_file,
                                bstring *event_decl,
                                bstring *init_descriptor_func)
{
    int rc;
    bstring identifier = NULL;
    assert(source != NULL);
    assert(property_file != NULL);
    assert(event_decl != NULL);
    assert(init_descriptor_func != NULL);

    // Initialize returned value.
    *event_decl = bfromcstr(
        "  int64_t ts;\n"
        "  uint32_t timestamp;\n"
        "  uint16_t action_id;\n"
    );
    check_mem(*event_decl);
    *init_descriptor_func = bfromcstr(
        "  descriptor:set_data_sz(ffi.sizeof(\"sky_lua_event_t\"));\n"
        "  descriptor:set_ts_offset(ffi.offsetof(\"sky_lua_event_t\", \"ts\"));\n"
        "  descriptor:set_timestamp_offset(ffi.offsetof(\"sky_lua_event_t\", \"timestamp\"));\n"
        "  descriptor:set_action_id_offset(ffi.offsetof(\"sky_lua_event_t\", \"action_id\"));\n"
    );
    check_mem(*init_descriptor_func);

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
        pos += (&EVENT_DOT_STR)->slen;

        // Read in identifier.
        int i;
        for(i=pos+1; i<blength(source); i++) {
            char ch = bchar(source, i);
            if(!(isalnum(ch) || ch == '_')) {
                break;
            }
        }
        identifier = bmidstr(source, pos, i-pos); check_mem(identifier);

        // Skip if the identifier refers to the action id or timestamp.
        if(biseqcstr(identifier, "action_id") == 1 || biseqcstr(identifier, "timestamp") == 1 || biseqcstr(identifier, "ts") == 1) {
            skip = true;
        }

        if(!skip) {
            if(blength(identifier)) {
                sky_property *property = NULL;
                rc = sky_property_file_find_by_name(property_file, identifier, &property);
                check(rc == 0, "Unable to find property by name: %s", bdata(identifier));
                check(property != NULL, "Property not found: %s", bdata(identifier));
            
                if(!lookup[property->id-SKY_PROPERTY_ID_MIN]) {
                    // Append property definition to event decl and function.
                    switch(property->data_type) {
                        case SKY_DATA_TYPE_STRING: {
                            bformata(*event_decl, "  char %s[];\n", bdata(property->name));
                            break;
                        }
                        case SKY_DATA_TYPE_INT: {
                            bformata(*event_decl, "  int32_t %s;\n", bdata(property->name));
                            break;
                        }
                        case SKY_DATA_TYPE_DOUBLE: {
                            bformata(*event_decl, "  double %s;\n", bdata(property->name));
                            break;
                        }
                        case SKY_DATA_TYPE_BOOLEAN: {
                            bformata(*event_decl, "  bool %s;\n", bdata(property->name));
                            break;
                        }
                        default:{
                            sentinel("Invalid sky lua type: %d", property->data_type);
                        }
                    }
                    check_mem(*event_decl);

                    bformata(*init_descriptor_func, "  descriptor:set_property(%d, ffi.offsetof(\"sky_lua_event_t\", \"%s\"), %d);\n", property->id, bdata(property->name), property->data_type);
                    check_mem(*init_descriptor_func);

                    // Flag the property as already processed.
                    lookup[property->id - SKY_PROPERTY_ID_MIN] = true;
                }
            }
        }

        bdestroy(identifier);
        identifier = NULL;
    }

    // Wrap properties in a struct.
    bassignformat(*event_decl, "typedef struct {\n%s} sky_lua_event_t;", bdata(*event_decl));
    check_mem(*event_decl);

    // Wrap info function.
    bassignformat(*init_descriptor_func,
        "function sky_init_descriptor(_descriptor)\n"
        "  descriptor = ffi.cast(\"sky_data_descriptor_t*\", _descriptor)\n"
        "%s"
        "end\n",
        bdata(*init_descriptor_func)
    );
    check_mem(*init_descriptor_func);

    return 0;

error:
    bdestroy(identifier);
    bdestroy(*event_decl);
    *event_decl = NULL;
    bdestroy(*init_descriptor_func);
    *init_descriptor_func = NULL;
    return -1;
}


//--------------------------------------
// Path Iterator/Cursor Management
//--------------------------------------

// Moves to the next path in the iterator and returns a cursor to traverse it.
//
// iterator - The path iterator.
//
// Returns a new cursor if there are more paths available. Otherwise returns
// NULL.
sky_cursor *sky_lua_path_iterator_get_cursor(sky_path_iterator *iterator)
{
    assert(iterator != NULL);
    return &iterator->cursor;
}

// Moves the cursor to the next event in a path.
//
// cursor - The cursor.
//
// Returns a pointer to the event data if there are remaining events.
// Otherwise returns NULL.
void *sky_lua_cursor_get_event(sky_cursor *cursor)
{
    assert(cursor != NULL);
    return cursor->data;
}
