-- SKY GENERATED CODE BEGIN --
local ffi = require('ffi')
ffi.cdef([[
typedef struct sky_string_t { int32_t length; char *data; } sky_string_t;
typedef struct sky_data_descriptor_t sky_data_descriptor_t;
typedef struct sky_object_iterator_t sky_object_iterator_t;
typedef struct {
  int64_t ts;
  uint32_t timestamp;
  {{range .}}{{structdef .}}
  {{end}}
} sky_lua_event_t;
typedef struct sky_cursor_t { sky_lua_event_t *event; int32_t session_event_index; } sky_cursor_t;

int sky_data_descriptor_set_data_sz(sky_data_descriptor_t *descriptor, uint32_t sz);
int sky_data_descriptor_set_timestamp_offset(sky_data_descriptor_t *descriptor, uint32_t offset);
int sky_data_descriptor_set_ts_offset(sky_data_descriptor_t *descriptor, uint32_t offset);
int sky_data_descriptor_set_action_id_offset(sky_data_descriptor_t *descriptor, uint32_t offset);
int sky_data_descriptor_set_property(sky_data_descriptor_t *descriptor, int8_t property_id, uint32_t offset, int data_type);

bool sky_object_iterator_eof(sky_object_iterator_t *);
void sky_object_iterator_next(sky_object_iterator_t *);
sky_cursor_t *sky_lua_object_iterator_get_cursor(sky_object_iterator_t *);

bool sky_cursor_eof(sky_cursor_t *);
bool sky_cursor_eos(sky_cursor_t *);
bool sky_lua_cursor_next_event(sky_cursor_t *);
bool sky_lua_cursor_next_session(sky_cursor_t *);
bool sky_cursor_set_session_idle(sky_cursor_t *, uint32_t);
]])
ffi.metatype('sky_data_descriptor_t', {
  __index = {
    set_data_sz = function(descriptor, sz) return ffi.C.sky_data_descriptor_set_data_sz(descriptor, sz) end,
    set_timestamp_offset = function(descriptor, offset) return ffi.C.sky_data_descriptor_set_timestamp_offset(descriptor, offset) end,
    set_ts_offset = function(descriptor, offset) return ffi.C.sky_data_descriptor_set_ts_offset(descriptor, offset) end,
    set_action_id_offset = function(descriptor, offset) return ffi.C.sky_data_descriptor_set_action_id_offset(descriptor, offset) end,
    set_property = function(descriptor, property_id, offset, data_type) return ffi.C.sky_data_descriptor_set_property(descriptor, property_id, offset, data_type) end,
  }
})
ffi.metatype('sky_object_iterator_t', {
  __index = {
    eof = function(iterator) return ffi.C.sky_object_iterator_eof(iterator) end,
    cursor = function(iterator) return ffi.C.sky_lua_object_iterator_get_cursor(iterator) end,
    next = function(iterator) return ffi.C.sky_object_iterator_next(iterator) end,
  }
})
ffi.metatype('sky_cursor_t', {
  __index = {
    eof = function(cursor) return ffi.C.sky_cursor_eof(cursor) end,
    eos = function(cursor) return ffi.C.sky_cursor_eos(cursor) end,
    next = function(cursor) return ffi.C.sky_lua_cursor_next_event(cursor) end,
    next_session = function(cursor) return ffi.C.sky_lua_cursor_next_session(cursor) end,
    set_session_idle = function(cursor, seconds) return ffi.C.sky_cursor_set_session_idle(cursor, seconds) end,
  }
})
ffi.metatype('sky_lua_event_t', {
  __index = {
  {{range .}}{{metatypedef .}}
  {{end}}
  }
})

function sky_init_descriptor(_descriptor)
  descriptor = ffi.cast('sky_data_descriptor_t*', _descriptor)
  {{range .}}{{initdescriptor .}}
  {{end}}
end

function sky_aggregate(_iterator)
  iterator = ffi.cast('sky_object_iterator_t*', _iterator)
  data = {}
  while not iterator:eof() do
    cursor = iterator:cursor()
    aggregate(cursor, data)
    iterator:next()
  end
  return data
end

-- The wrapper for the merge.
function sky_merge(results, data)
  if data ~= nil then
    merge(results, data)
  end"
  return results
end
-- SKY GENERATED CODE END --

