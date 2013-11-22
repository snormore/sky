package query

/*
#cgo LDFLAGS: -lluajit-5.1 -llmdb
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <inttypes.h>
#include <string.h>
#include <lmdb.h>
#include <luajit-2.0/lua.h>
#include <luajit-2.0/lualib.h>
#include <luajit-2.0/lauxlib.h>

//==============================================================================
//
// Constants
//
//==============================================================================

// The number of microseconds per second.
#define USEC_PER_SEC        1000000

// A bit-mask to extract the microseconds from a Sky timestamp.
#define USEC_MASK           0xFFFFF

// The number of bits that seconds are shifted over in a timestamp.
#define SECONDS_BIT_OFFSET  20


//==============================================================================
//
// Macros
//
//==============================================================================

#define memdump(PTR, LENGTH) do {\
    char *address = (char*)PTR;\
    int length = LENGTH;\
    int i = 0;\
    char *line = (char*)address;\
    unsigned char ch;\
    fprintf(stderr, "%"PRIX64" | ", (int64_t)address);\
    while (length-- > 0) {\
        fprintf(stderr, "%02X ", (unsigned char)*address++);\
        if (!(++i % 16) || (length == 0 && i % 16)) {\
            if (length == 0) { while (i++ % 16) { fprintf(stderr, "__ "); } }\
            fprintf(stderr, "| ");\
            while (line < address) {\
                ch = *line++;\
                fprintf(stderr, "%c", (ch < 33 || ch == 255) ? 0x2E : ch);\
            }\
            if (length > 0) { fprintf(stderr, "\n%09X | ", (int)address); }\
        }\
    }\
    fprintf(stderr, "\n\n");\
} while(0)

// Removed from macro below.
// memdump(cursor->startptr, (cursor->endptr - cursor->startptr));

#define badcursordata(MSG, PTR) do {\
    fprintf(stderr, "Cursor pointing at invalid raw event data [" MSG "]: %p\n", PTR); \
    cursor->next_event->timestamp = 0; \
    return false; \
} while(0)

#define debug(M, ...) fprintf(stderr, "DEBUG %s:%d: " M "\n", __FILE__, __LINE__, ##__VA_ARGS__)


//==============================================================================
//
// Typedefs
//
//==============================================================================

typedef struct {
  int32_t length;
  char *data;
} sky_string;

typedef struct sky_cursor sky_cursor;

typedef int (*sky_cursor_next_object_func)(void *cursor);

typedef void (*sky_property_descriptor_set_func)(void *target, void *value, size_t *sz);
typedef void (*sky_property_descriptor_clear_func)(void *target);

typedef struct { uint16_t ts_offset; uint16_t timestamp_offset;} sky_timestamp_descriptor;

typedef struct {
    int64_t property_id;
    uint16_t offset;
    sky_property_descriptor_set_func set_func;
    sky_property_descriptor_clear_func clear_func;
} sky_property_descriptor;

typedef struct {
  bool eos;
  bool eof;
  uint32_t timestamp;
  int64_t ts;
} sky_event;

struct sky_cursor {
    sky_event *event;
    sky_event *next_event;

    uint32_t max_timestamp;
    uint32_t session_idle_in_sec;

    bool eos_wait;
    uint32_t event_sz;
    uint32_t action_event_sz;
    uint32_t variable_event_sz;

    sky_timestamp_descriptor timestamp_descriptor;
    sky_property_descriptor *property_descriptors;
    sky_property_descriptor *property_zero_descriptor;
    uint32_t property_count;

    int32_t min_property_id;
    int32_t max_property_id;

    void *key_prefix;
    uint32_t key_prefix_sz;
    MDB_cursor* lmdb_cursor;
};

//==============================================================================
//
// Forward Declarations
//
//==============================================================================

//--------------------------------------
// Setters
//--------------------------------------

void sky_set_noop(void *target, void *value, size_t *sz);

void sky_set_string(void *target, void *value, size_t *sz);

void sky_set_int(void *target, void *value, size_t *sz);

void sky_set_double(void *target, void *value, size_t *sz);

void sky_set_boolean(void *target, void *value, size_t *sz);

//--------------------------------------
// Clear Functions
//--------------------------------------

void sky_clear_string(void *target);

void sky_clear_int(void *target);

void sky_clear_double(void *target);

void sky_clear_boolean(void *target);

//--------------------------------------
// Object Iteration
//--------------------------------------

bool sky_cursor_first_object(sky_cursor *cursor);

bool sky_cursor_next_object(sky_cursor *cursor);

bool sky_cursor_read(sky_cursor *cursor, sky_event *event, void *ptr);

bool sky_cursor_next_event(sky_cursor *cursor);

bool sky_cursor_eof(sky_cursor *cursor);

bool sky_cursor_eos(sky_cursor *cursor);

void sky_cursor_update_eos(sky_cursor *cursor);

//--------------------------------------
// Timestamps
//--------------------------------------

int64_t sky_timestamp_shift(int64_t value);

int64_t sky_timestamp_unshift(int64_t value);

int64_t sky_timestamp_to_seconds(int64_t value);

//--------------------------------------
// Minipack
//--------------------------------------

size_t minipack_sizeof_elem_and_data(void *ptr);

bool minipack_is_raw(void *ptr);

int64_t minipack_unpack_int(void *ptr, size_t *sz);

double minipack_unpack_double(void *ptr, size_t *sz);

bool minipack_unpack_bool(void *ptr, size_t *sz);

uint32_t minipack_unpack_raw(void *ptr, size_t *sz);

uint32_t minipack_unpack_map(void *ptr, size_t *sz);

void minipack_unpack_nil(void *ptr, size_t *sz);


//==============================================================================
//
// Byte Order
//
//==============================================================================

#include <sys/types.h>

#ifndef BYTE_ORDER
#if defined(linux) || defined(__linux__)
# include <endian.h>
#else
# include <machine/endian.h>
#endif
#endif

#if !defined(BYTE_ORDER) && !defined(__BYTE_ORDER)
#error "Undefined byte order"
#endif

uint64_t bswap64(uint64_t value);

#if (BYTE_ORDER == LITTLE_ENDIAN) || (__BYTE_ORDER == __LITTLE_ENDIAN)
#define htonll(x) bswap64(x)
#define ntohll(x) bswap64(x)
#else
#define htonll(x) x
#define ntohll(x) x
#endif


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

// Creates a reference to a cursor.
sky_cursor *sky_cursor_new(int32_t min_property_id,
                           int32_t max_property_id)
{
    sky_cursor *cursor = calloc(1, sizeof(sky_cursor));
    if(cursor == NULL) debug("[malloc] Unable to allocate cursor.");

    // Add one property to account for the zero descriptor.
    int32_t property_count = (max_property_id - min_property_id) + 1;

    // Allocate memory for the descriptors.
    cursor->property_descriptors = calloc(property_count, sizeof(sky_property_descriptor));
    if(cursor->property_descriptors == NULL) debug("[malloc] Unable to allocate property descriptors.");
    cursor->property_count = property_count;
    cursor->property_zero_descriptor = NULL;

    cursor->min_property_id = min_property_id;
    cursor->max_property_id = max_property_id;

    // Initialize all property descriptors to noop.
    int32_t i;
    for(i=0; i<property_count; i++) {
        int64_t property_id = min_property_id + (int64_t)i;
        cursor->property_descriptors[i].property_id = property_id;
        cursor->property_descriptors[i].set_func = sky_set_noop;

        // Save a pointer to the descriptor that points to property zero.
        if(property_id == 0) {
            cursor->property_zero_descriptor = &cursor->property_descriptors[i];
        }
    }

    return cursor;
}

// Removes a cursor reference from memory.
void sky_cursor_free(sky_cursor *cursor)
{
    if(cursor) {
        if(cursor->property_descriptors != NULL) free(cursor->property_descriptors);
        cursor->property_zero_descriptor = NULL;
        cursor->property_count = 0;

        if(cursor->event != NULL) free(cursor->event);
        cursor->event = NULL;
        if(cursor->next_event != NULL) free(cursor->next_event);
        cursor->next_event = NULL;
        if(cursor->key_prefix != NULL) free(cursor->key_prefix);
        cursor->key_prefix = NULL;

        free(cursor);
    }
}


//--------------------------------------
// Data Management
//--------------------------------------

void sky_cursor_set_value(sky_cursor *cursor, void *target,
                          int64_t property_id, void *ptr, size_t *sz)
{
    if(property_id >= cursor->min_property_id && property_id <= cursor->max_property_id) {
        sky_property_descriptor *property_descriptor = &cursor->property_zero_descriptor[property_id];
        property_descriptor->set_func(target + property_descriptor->offset, ptr, sz);
    } else {
        sky_set_noop(NULL, ptr, sz);
    }
}


//--------------------------------------
// Descriptor Management
//--------------------------------------

void sky_cursor_set_event_sz(sky_cursor *cursor, uint32_t sz) {
    cursor->event_sz = sz;

    if(cursor->event != NULL) free(cursor->event);
    cursor->event = calloc(1, sz);
    if(cursor->event == NULL) debug("[malloc] Unable to allocate cursor event.");

    if(cursor->next_event != NULL) free(cursor->next_event);
    cursor->next_event = calloc(1, sz);
    if(cursor->next_event == NULL) debug("[malloc] Unable to allocate cursor next event.");
}

// Sets the data type and offset for a given property id.
void sky_cursor_set_property(sky_cursor *cursor, int64_t property_id,
                             uint32_t offset, uint32_t sz, const char *data_type)
{
    sky_property_descriptor *property_descriptor = &cursor->property_zero_descriptor[property_id];

    // Set the offset and set_func function on the descriptor.
    if(property_id != 0) {
        property_descriptor->offset = offset;
        if(strlen(data_type) == 0) {
            property_descriptor->set_func = sky_set_noop;
            property_descriptor->clear_func = NULL;
        }
        else if(strcmp(data_type, "string") == 0) {
            property_descriptor->set_func = sky_set_string;
            property_descriptor->clear_func = sky_clear_string;
        }
        else if(strcmp(data_type, "factor") == 0 || strcmp(data_type, "integer") == 0) {
            property_descriptor->set_func = sky_set_int;
            property_descriptor->clear_func = sky_clear_int;
        }
        else if(strcmp(data_type, "float") == 0) {
            property_descriptor->set_func = sky_set_double;
            property_descriptor->clear_func = sky_clear_double;
        }
        else if(strcmp(data_type, "boolean") == 0) {
            property_descriptor->set_func = sky_set_boolean;
            property_descriptor->clear_func = sky_clear_boolean;
        }
        else {
            property_descriptor->set_func = sky_set_boolean;
            property_descriptor->clear_func = sky_clear_boolean;
        }
    }

    // Resize the action data area. This area occurs after the
    // fixed fields in the struct.
    int32_t new_action_event_sz = (offset + sz) - sizeof(sky_event);
    if(property_id < 0 && new_action_event_sz > cursor->action_event_sz) {
        cursor->action_event_sz = (uint32_t)new_action_event_sz;
    }

    // Resize the variable data area. This area occurs after the
    // action data fields in the struct.
    int32_t new_variable_event_sz = (offset + sz) - sizeof(sky_event);
    if(property_id == 0 && new_variable_event_sz > 0 && new_variable_event_sz > cursor->variable_event_sz) {
        cursor->variable_event_sz = (uint32_t)new_variable_event_sz;
    }
}


//--------------------------------------
// Object Iteration
//--------------------------------------

// Sets up object after cursor has already been positioned.
bool sky_cursor_iter_object(sky_cursor *cursor, MDB_val *key, MDB_val *data)
{
    if(cursor->key_prefix != NULL && (key->mv_size < cursor->key_prefix_sz || memcmp(cursor->key_prefix, key->mv_data, cursor->key_prefix_sz) != 0)) {
        return false;
    }
    // fprintf(stderr, "\nOBJ (%.*s) [%d]\n", (int)key->mv_size, (char*)key->mv_data, (int)key->mv_size);

    // Clear the data object if set.
    cursor->session_idle_in_sec = 0;
    cursor->eos_wait = false;
    cursor->next_event->eof = false;
    memset(cursor->event, 0, cursor->event_sz);

    // Read first event into "next" event.
    if(!sky_cursor_read(cursor, cursor->next_event, data->mv_data)) {
        return false;
    }

    // Move "next" event to current event and put the next event in buffer.
    return sky_cursor_next_event(cursor);
}

// Moves the cursor to point to the first object. If a prefix is set then
// move to the first object that with the given prefix.
bool sky_cursor_first_object(sky_cursor *cursor)
{
    int rc;
    MDB_val key, data;

    if(cursor->key_prefix == NULL) {
        if((rc = mdb_cursor_get(cursor->lmdb_cursor, &key, &data, MDB_FIRST)) != 0) {
            if(rc != MDB_NOTFOUND) fprintf(stderr, "MDB_FIRST error: %d\n", rc);
            return false;
        }

    } else {
        key.mv_data = cursor->key_prefix;
        key.mv_size = cursor->key_prefix_sz;
        if((rc = mdb_cursor_get(cursor->lmdb_cursor, &key, &data, MDB_SET_RANGE)) != 0) {
            if(rc != MDB_NOTFOUND) fprintf(stderr, "MDB_SET_RANGE error: %d\n", rc);
            return false;
        }
    }

//    fprintf(stderr, "DATA: sz=%d\n", (int)data.mv_size);

    return sky_cursor_iter_object(cursor, &key, &data);
}

// Moves the cursor to point to the next object.
bool sky_cursor_next_object(sky_cursor *cursor)
{
    // Move to next object.
    MDB_val key, data;
    int rc = mdb_cursor_get(cursor->lmdb_cursor, &key, &data, MDB_NEXT_NODUP);
    if(rc != 0) {
        return false;
    }

    return sky_cursor_iter_object(cursor, &key, &data);
}

// Moves the cursor to point to the next event.
// Returns true if the cursor moved forward, otherwise false.
bool sky_cursor_next_event(sky_cursor *cursor)
{
    // Don't allow cursor to move if we're EOF or marked as EOS wait.
    if(cursor->event->eof || (cursor->event->eos && cursor->eos_wait)) {
        return false;
    }
    cursor->eos_wait = true;

    // Copy variable state from current event to next event.
    if(cursor->variable_event_sz > 0) {
        uint32_t variable_event_offset = sizeof(sky_event) + cursor->action_event_sz;
        memcpy(((void*)cursor->next_event) + variable_event_offset, ((void*)cursor->event) + variable_event_offset, cursor->variable_event_sz - cursor->action_event_sz);
    }

    // Copy the next event to the current event.
    memcpy(cursor->event, cursor->next_event, cursor->event_sz);

    // Read the next event.
    if(!cursor->next_event->eof) {
        MDB_val key, data;
        int rc = mdb_cursor_get(cursor->lmdb_cursor, &key, &data, MDB_NEXT_DUP);
        if(rc != 0) {
            // Clear next event if there isn't one.
            memset(cursor->next_event, 0, cursor->event_sz);
            cursor->next_event->eof = true;

            if(rc != MDB_NOTFOUND) {
                printf("lmdb cursor error: %d\n", rc);
            }
        } else {
            cursor->next_event->eof = false;
            if(!sky_cursor_read(cursor, cursor->next_event, data.mv_data)) {
                return true;
            }
        }
    }

    // Update eos.
    sky_cursor_update_eos(cursor);

    return true;
}

// Reads the data at a given pointer into a data object.
bool sky_cursor_read(sky_cursor *cursor, sky_event *event, void *ptr)
{
    // Set timestamp.
    event->ts = htonll(*((int64_t*)ptr));
    event->timestamp = sky_timestamp_to_seconds(event->ts);
    ptr += 8;

    // Clear old action data.
    if(cursor->action_event_sz > 0) {
        memset(&event[1], 0, cursor->action_event_sz);
    }

    // Read msgpack map!
    size_t sz;
    uint32_t count = minipack_unpack_map(ptr, &sz);
    if(sz == 0) {
      minipack_unpack_nil(ptr, &sz);
      if(sz == 0) {
        badcursordata("datamap", ptr);
      }
    }
    ptr += sz;

    // Loop over key/value pairs.
    uint32_t i;
    for(i=0; i<count; i++) {
        // Read property id (key).
        int64_t property_id = minipack_unpack_int(ptr, &sz);
        if(sz == 0) badcursordata("key", ptr);
        ptr += sz;

        // Read property value and set it on the data object.
        sky_cursor_set_value(cursor, event, property_id, ptr, &sz);
        if(sz == 0) {
          debug("[invalid read, skipping]");
          sz = minipack_sizeof_elem_and_data(ptr);
        }
        ptr += sz;
    }
}

bool sky_lua_cursor_next_event(sky_cursor *cursor)
{
    return sky_cursor_next_event(cursor);
}

bool sky_cursor_eof(sky_cursor *cursor)
{
    return cursor->event->eof;
}

// End-of-session (EOS) is defined by idle time between the current event and the next event.
bool sky_cursor_eos(sky_cursor *cursor)
{
    return cursor->event->eos;
}

// Updates the end-of-session flag on the current event.
void sky_cursor_update_eos(sky_cursor *cursor)
{
    if(cursor->next_event->eof) {
        cursor->event->eos = true;
    } else if(cursor->session_idle_in_sec == 0) {
        cursor->event->eos = false;
    } else {
        cursor->event->eos = (cursor->next_event->timestamp - cursor->event->timestamp >= cursor->session_idle_in_sec);
    }
}

void sky_cursor_set_session_idle(sky_cursor *cursor, uint32_t seconds)
{
    cursor->session_idle_in_sec = seconds;
    sky_cursor_update_eos(cursor);
}

void sky_cursor_next_session(sky_cursor *cursor)
{
    cursor->eos_wait = false;
}

bool sky_lua_cursor_next_session(sky_cursor *cursor)
{
    sky_cursor_next_session(cursor);
    return !cursor->next_event->eof;
}



//--------------------------------------
// Setters
//--------------------------------------

void sky_set_noop(void *target, void *value, size_t *sz)
{
    ((void)(target));
    *sz = minipack_sizeof_elem_and_data(value);
}

void sky_set_string(void *target, void *value, size_t *sz)
{
    size_t _sz;
    sky_string *string = (sky_string*)target;
    string->length = minipack_unpack_raw(value, &_sz);
    string->data = (_sz > 0 ? value + _sz : NULL);
    *sz = _sz + string->length;
}

void sky_set_int(void *target, void *value, size_t *sz)
{
    *((int32_t*)target) = (int32_t)minipack_unpack_int(value, sz);
    if(*sz == 0) {
      minipack_unpack_nil(value, sz);
      if(*sz != 0) {
        *((int32_t*)target) = 0;
      }
    }
}

void sky_set_double(void *target, void *value, size_t *sz)
{
    *((double*)target) = minipack_unpack_double(value, sz);
    if(*sz == 0) {
      minipack_unpack_nil(value, sz);
      if(*sz != 0) {
        *((double*)target) = 0;
      }
    }
}

void sky_set_boolean(void *target, void *value, size_t *sz)
{
    *((bool*)target) = minipack_unpack_bool(value, sz);
    if(*sz == 0) {
      minipack_unpack_nil(value, sz);
      if(*sz != 0) {
        *((bool*)target) = false;
      }
    }
}


//--------------------------------------
// Clear Functions
//--------------------------------------

void sky_clear_string(void *target)
{
    sky_string *string = (sky_string*)target;
    string->length = 0;
    string->data = NULL;
}

void sky_clear_int(void *target)
{
    *((int32_t*)target) = 0;
}

void sky_clear_double(void *target)
{
    *((double*)target) = 0;
}

void sky_clear_boolean(void *target)
{
    *((bool*)target) = false;
}

//--------------------------------------
// Timestamps
//--------------------------------------

// Converts a timestamp from the number of microseconds since the epoch to
// a bit-shifted Sky timestamp.
//
// value - Microseconds since the unix epoch.
//
// Returns a bit-shifted Sky timestamp.
int64_t sky_timestamp_shift(int64_t value)
{
    int64_t usec = value % USEC_PER_SEC;
    int64_t sec  = (value / USEC_PER_SEC);

    return (sec << SECONDS_BIT_OFFSET) + usec;
}

// Converts a bit-shifted Sky timestamp to the number of microseconds since
// the Unix epoch.
//
// value - Sky timestamp.
//
// Returns the number of microseconds since the Unix epoch.
int64_t sky_timestamp_unshift(int64_t value)
{
    int64_t usec = value & USEC_MASK;
    int64_t sec  = value >> SECONDS_BIT_OFFSET;

    return (sec * USEC_PER_SEC) + usec;
}

// Converts a bit-shifted Sky timestamp to seconds since the epoch.
//
// value - Sky timestamp.
//
// Returns the number of seconds since the Unix epoch.
int64_t sky_timestamp_to_seconds(int64_t value)
{
    return (value >> SECONDS_BIT_OFFSET);
}

*/
import "C"

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/skydb/sky/core"
	"github.com/szferi/gomdb"
	"github.com/ugorji/go/codec"
	"regexp"
	"sync"
	"unsafe"
)

var (
	cstr_sky_initialize = C.CString("sky_initialize")
	cstr_sky_aggregate = C.CString("sky_aggregate")
)

// Engine is an execution engine used to iterate over a series of objects.
type Engine struct {
    sync.Mutex
	cursor  *mdb.Cursor
    source  string
    prefix  string
}

// New creates a new engine.
func New(source string, prefix string) *Engine {
	e := &Engine{
        source: source,
        prefix: prefix,
    }
}

// Source returns the Lua source for the engine.
func (e *Engine) Source() string {
    return e.source
}

// Prefix returns the object prefix used to limit the engine execution.
func (e *Engine) Prefix() string {
    return e.prefix
}

// Initializes the data structure used for aggregation.
func (e *Engine) Initialize() (interface{}, error) {
	state, err := e.lua_state()
	if err != nil {
		return err
	}
	defer C.lua_close(state)

	C.lua_getfield(state, -10002, cstr_sky_initialize)
	C.lua_pushlightuserdata(state, unsafe.Pointer(e.cursor))
	if rc := C.lua_pcall(state, 1, 1, 0); rc != 0 {
		return fmt.Errorf("engine init error: %s", C.GoString(C.lua_tolstring(state, -1, nil)))
	}

	return lua_ret(state)
}

// Executes an aggregation over the entire database.
func (e *Engine) Aggregate(data interface{}) (interface{}, error) {
	state, err := e.lua_state()
	if err != nil {
		return err
	}
	defer C.lua_close(state)

	C.lua_getfield(state, -10002, cstr_sky_aggregate)
	C.lua_pushlightuserdata(state, unsafe.Pointer(e.cursor))
	if err := lua_arg(state, data); err != nil {
		return nil, err
	}
	if rc := C.lua_pcall(state, 2, 1, 0); rc != 0 {
		return fmt.Errorf("engine aggregate error: %s", C.GoString(C.lua_tolstring(state, -1, nil)))
	}

	return lua_ret(state)
}

// Executes an merge over the aggregated data.
func (e *Engine) Merge(results interface{}, data interface{}) (interface{}, error) {
	state, err := e.lua_state()
	if err != nil {
		return err
	}
	defer C.lua_close(state)

	C.lua_getfield(state, -10002, cstr_sky_merge)
	if err := lua_arg(state, results); err != nil {
		return results, err
	}
	if err := lua_arg(state, data); err != nil {
		return results, err
	}
	if rc := C.lua_pcall(state, 2, 1, 0); rc != 0 {
		return fmt.Errorf("engine merge error: %s", C.GoString(C.lua_tolstring(state, -1, nil)))
	}

	return lua_ret(state)
}
