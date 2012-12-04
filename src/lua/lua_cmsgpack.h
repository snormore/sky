#ifndef _lua_cmsgpack_h
#define _lua_cmsgpack_h

#include <lua.h>
#include <lauxlib.h>

LUALIB_API int luaopen_cmsgpack(lua_State *L);

int mp_pack(lua_State *L);

int mp_unpack(lua_State *L);

#endif
