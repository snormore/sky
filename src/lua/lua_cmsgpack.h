#ifndef _lua_cmsgpack_h
#define _lua_cmsgpack_h

#include "luajit-2.0/lua.h"
#include "luajit-2.0/lauxlib.h"

LUALIB_API int luaopen_cmsgpack(lua_State *L);

int mp_pack(lua_State *L);

int mp_unpack(lua_State *L);

#endif