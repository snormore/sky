package engine

/*
#cgo LDFLAGS: -lluajit-5.1
#include <stdlib.h>
#include <luajit-2.0/lua.h>
#include <luajit-2.0/lualib.h>
#include <luajit-2.0/lauxlib.h>

int mp_pack(lua_State *L);
int mp_unpack(lua_State *L);
*/
import "C"

// lua_state creates a new Lua state object based on a source string.
func lua_state(source string) (*C.lua_State, error) {
    // Initialize the Lua state.
    state = C.luaL_newstate()
    if state == nil {
        return errors.New("lua alloc error")
    }
    C.luaL_openlibs(state)

    // Compile the source.
    source := C.CString(e.source)
    defer C.free(unsafe.Pointer(source))

    // Compile the script.
    if rc := C.luaL_loadstring(state, source); rc != 0 {
        defer C.lua_close(state)
        return fmt.Errorf("lua syntax error: %s", C.GoString(C.lua_tolstring(state, -1, nil)))
    }

    // Run script once to initialize.
    if rc := C.lua_pcall(state, 0, 0, 0); rc != 0 {
        defer C.lua_close(state)
        return fmt.Errorf("lua init error: %s", C.GoString(C.lua_tolstring(state, -1, nil)))
    }

    return nil
}

// Encodes a Go object into Msgpack and adds it to the function arguments.
func lua_arg(state *C.lua_State, value interface{}) error {
	// Encode Go object into msgpack.
	var handle codec.MsgpackHandle
	handle.RawToString = true
	buffer := new(bytes.Buffer)
	encoder := codec.NewEncoder(buffer, &handle)
	if err := encoder.Encode(value); err != nil {
		return err
	}

	// Push the msgpack data onto the Lua stack.
	data := buffer.String()
	cdata := C.CString(data)
	defer C.free(unsafe.Pointer(cdata))
	C.lua_pushlstring(state, cdata, (C.size_t)(len(data)))

	// Convert the argument from msgpack into Lua.
	if rc := C.mp_unpack(state); rc != 1 {
		return errors.New("lua encode error")
	}
	C.lua_remove(state, -2)

	return nil
}

// Decodes the result from a function into a Go object.
func lua_ret(state *C.lua_State) (interface{}, error) {
	// Encode Lua object into msgpack.
	if rc := C.mp_pack(state); rc != 1 {
		return nil, errors.New("lua decode error")
	}
	sz := C.size_t(0)
	ptr := C.lua_tolstring(state, -1, (*C.size_t)(&sz))
	str := C.GoStringN(ptr, (C.int)(sz))
	C.lua_settop(state, -(1)-1) // lua_pop()

	// Decode msgpack into a Go object.
	var handle codec.MsgpackHandle
	handle.RawToString = true
	var ret interface{}
	decoder := codec.NewDecoder(bytes.NewBufferString(str), &handle)
	if err := decoder.Decode(&ret); err != nil {
		return nil, err
	}

	return ret, nil
}
