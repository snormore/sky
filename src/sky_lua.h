#ifndef _sky_lua_h
#define _sky_lua_h

#include <lua.h>
#include <lualib.h>
#include <lauxlib.h>
#include "lua/lua_cmsgpack.h"

#include "table.h"
#include "property_file.h"
#include "data_descriptor.h"
#include "bstring.h"

//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Initialization
//--------------------------------------

int sky_lua_initscript(bstring source, lua_State **L);

int sky_lua_initscript_with_table(bstring source, sky_table *table,
    sky_data_descriptor *descriptor, lua_State **L);

//--------------------------------------
// MessagePack
//--------------------------------------

int sky_lua_to_msgpack(lua_State *L, bstring *ret);

//--------------------------------------
// Property File Integration
//--------------------------------------

int sky_lua_generate_header(bstring source, sky_table *table, bstring *ret);

int sky_lua_generate_event_info(bstring source,
    sky_property_file *property_file, bstring *event_decl,
    bstring *init_descriptor_func);

#endif
