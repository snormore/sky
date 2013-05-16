package skyd

/*
#cgo LDFLAGS: -lluajit-5.1 -lleveldb
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <inttypes.h>
#include <string.h>
#include <leveldb/c.h>
#include <luajit-2.0/lua.h>
#include <luajit-2.0/lualib.h>
#include <luajit-2.0/lauxlib.h>

int mp_pack(lua_State *L);
int mp_unpack(lua_State *L);

//==============================================================================
//
// Constants
//
//==============================================================================

#define sky_event_flag_t uint8_t
#define EVENT_FLAG       0x92

// The number of microseconds per second.
#define USEC_PER_SEC        1000000

// A bit-mask to extract the microseconds from a Sky timestamp.
#define USEC_MASK           0xFFFFF

// The number of bits that seconds are shifted over in a timestamp.
#define SECONDS_BIT_OFFSET  20

#define SKY_PROPERTY_DESCRIPTOR_PADDING  32


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
    fprintf(stderr, "Cursor pointing at invalid raw event data [" MSG "]: %p->%p\n", cursor->ptr, PTR); \
    cursor->eof = true; \
    return; \
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

struct sky_cursor {
    void *data;
    uint32_t data_sz;
    uint32_t action_data_sz;

    int32_t session_event_index;
    void *startptr;
    void *nextptr;
    void *endptr;
    void *ptr;
    bool eof;
    bool in_session;
    uint32_t last_timestamp;
    uint32_t session_idle_in_sec;

    sky_timestamp_descriptor timestamp_descriptor;
    sky_property_descriptor *property_descriptors;
    sky_property_descriptor *property_zero_descriptor;
    uint32_t property_count;

    void *key_prefix;
    uint32_t key_prefix_sz;
    leveldb_iterator_t* leveldb_iterator;
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

bool sky_cursor_next_object(sky_cursor *cursor);
bool sky_cursor_has_next_object(sky_cursor *cursor);

//--------------------------------------
// Event Iteration
//--------------------------------------

void sky_cursor_set_ptr(sky_cursor *cursor, void *ptr, size_t sz);

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

    // Add one property to account for the zero descriptor.
    min_property_id -= SKY_PROPERTY_DESCRIPTOR_PADDING;
    max_property_id += SKY_PROPERTY_DESCRIPTOR_PADDING;
    int32_t property_count = (max_property_id - min_property_id) + 1;

    // Allocate memory for the descriptors.
    cursor->property_descriptors = calloc(property_count, sizeof(sky_property_descriptor));
    cursor->property_count = property_count;
    cursor->property_zero_descriptor = NULL;
    
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

        if(cursor->data != NULL) free(cursor->data);

        free(cursor);
    }
}


//--------------------------------------
// Data Management
//--------------------------------------

void sky_cursor_set_value(sky_cursor *cursor, void *target,
                          int64_t property_id, void *ptr, size_t *sz)
{
    sky_property_descriptor *property_descriptor = &cursor->property_zero_descriptor[property_id];
    property_descriptor->set_func(target + property_descriptor->offset, ptr, sz);
}


//--------------------------------------
// Descriptor Management
//--------------------------------------

void sky_cursor_set_data_sz(sky_cursor *cursor, uint32_t sz) {
    cursor->data_sz = sz;
    if(cursor->data != NULL) free(cursor->data);
    cursor->data = calloc(1, sz);
}

void sky_cursor_set_timestamp_offset(sky_cursor *cursor, uint32_t offset) {
    cursor->timestamp_descriptor.timestamp_offset = offset;
}

void sky_cursor_set_ts_offset(sky_cursor *cursor, uint32_t offset) {
    cursor->timestamp_descriptor.ts_offset = offset;
}

// Sets the data type and offset for a given property id.
void sky_cursor_set_property(sky_cursor *cursor, int64_t property_id,
                             uint32_t offset, uint32_t sz, const char *data_type)
{
    sky_property_descriptor *property_descriptor = &cursor->property_zero_descriptor[property_id];
    
    // Set the offset and set_func function on the descriptor.
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
    
    // Resize the action data area.
    if(property_id < 0 && offset+sz > cursor->action_data_sz) {
        cursor->action_data_sz = offset+sz;
    }
}


//--------------------------------------
// Object Iteration
//--------------------------------------

// Moves the cursor to point to the next object.
bool sky_cursor_next_object(sky_cursor *cursor)
{
    // Move to next object.
    if(sky_cursor_has_next_object(cursor)) {
        // Retrieve the data for the next object.
        size_t data_length;
        void *data = (void*)leveldb_iter_value(cursor->leveldb_iterator, &data_length);
        
        // Set the pointer on the cursor.
        sky_cursor_set_ptr(cursor, data, data_length);
        
        leveldb_iter_next(cursor->leveldb_iterator);
        return true;
    }
    else {
        return false;
    }
}

bool sky_cursor_has_next_object(sky_cursor *cursor)
{
    if(cursor->leveldb_iterator != NULL && leveldb_iter_valid(cursor->leveldb_iterator)) {
        // If there's no prefix then the key is valid.
        if(cursor->key_prefix_sz == 0) {
            return true;
        }

        // Otherwise check that the key has the given the prefix.
        size_t key_sz;
        void *key = (void*)leveldb_iter_key(cursor->leveldb_iterator, &key_sz);

        return (cursor->key_prefix_sz <= key_sz && memcmp(cursor->key_prefix, key, cursor->key_prefix_sz) == 0);
    }
    
    return false;
}

void sky_cursor_set_key_prefix(sky_cursor *cursor, void *prefix, uint32_t sz)
{
    cursor->key_prefix = prefix;
    cursor->key_prefix_sz = sz;
}


//--------------------------------------
// Event Iteration
//--------------------------------------

void sky_cursor_set_ptr(sky_cursor *cursor, void *ptr, size_t sz)
{
    // Set the start of the path and the length of the data.
    cursor->startptr   = ptr;
    cursor->nextptr    = ptr;
    cursor->endptr     = ptr + sz;
    cursor->ptr        = NULL;
    cursor->in_session = true;
    cursor->last_timestamp      = 0;
    cursor->session_idle_in_sec = 0;
    cursor->session_event_index = -1;
    cursor->eof        = !(ptr != NULL && cursor->startptr < cursor->endptr);
    
    // Clear the data object if set.
    memset(cursor->data, 0, cursor->data_sz);
    
    // The first item is the current state so skip it.
    if(cursor->startptr != NULL && minipack_is_raw(cursor->startptr)) {
        cursor->startptr += minipack_sizeof_elem_and_data(cursor->startptr);
        cursor->nextptr = cursor->startptr;
    }
}

void sky_cursor_next_event(sky_cursor *cursor)
{
    // Ignore any calls when the cursor is out of session or EOF.
    if(cursor->eof || !cursor->in_session) {
        return;
    }

    // Move the pointer to the next position.
    void *prevptr = cursor->ptr;
    cursor->ptr = cursor->nextptr;
    void *ptr = cursor->ptr;

    // If pointer is beyond the last event then set eof.
    if(cursor->ptr >= cursor->endptr) {
        cursor->eof        = true;
        cursor->in_session = false;
        cursor->ptr        = NULL;
        cursor->startptr   = NULL;
        cursor->nextptr    = NULL;
        cursor->endptr     = NULL;
    }
    // Otherwise update the event object with data.
    else {
        sky_event_flag_t flag = *((sky_event_flag_t*)ptr);
        
        // If flag isn't correct then report and exit.
        if(flag != EVENT_FLAG) badcursordata("eflag", ptr);
        ptr += sizeof(sky_event_flag_t);
        
        // Read timestamp.
        size_t sz;
        int64_t ts = minipack_unpack_int(ptr, &sz);
        if(sz == 0) badcursordata("timestamp", ptr);
        uint32_t timestamp = sky_timestamp_to_seconds(ts);
        ptr += sz;

        // Check for session boundry. This only applies if this is not the
        // first event in the session and a session idle time has been set.
        if(cursor->last_timestamp > 0 && cursor->session_idle_in_sec > 0) {
            // If the elapsed time is greater than the idle time then rewind
            // back to the event we started on at the beginning of the function
            // and mark the cursor as being "out of session".
            if(timestamp - cursor->last_timestamp >= cursor->session_idle_in_sec) {
                cursor->ptr = prevptr;
                cursor->in_session = false;
            }
        }
        cursor->last_timestamp = timestamp;

        // Only process the event if we're still in session.
        if(cursor->in_session) {
            cursor->session_event_index++;
            
            // Set timestamp.
            int64_t *data_ts = (int64_t*)(cursor->data + cursor->timestamp_descriptor.ts_offset);
            uint32_t *data_timestamp = (uint32_t*)(cursor->data + cursor->timestamp_descriptor.timestamp_offset);
            *data_ts = ts;
            *data_timestamp = timestamp;
            
            // Clear old action data.
            if(cursor->action_data_sz > 0) {
              memset(cursor->data, 0, cursor->action_data_sz);
            }

            // Read msgpack map!
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
                sky_cursor_set_value(cursor, cursor->data, property_id, ptr, &sz);
                if(sz == 0) {
                  debug("[invalid read, skipping]");
                  sz = minipack_sizeof_elem_and_data(ptr);
                }
                ptr += sz;
            }

            cursor->nextptr = ptr;
        }
    }
}

bool sky_lua_cursor_next_event(sky_cursor *cursor)
{
    sky_cursor_next_event(cursor);
    return (!cursor->eof && cursor->in_session);
}

bool sky_cursor_eof(sky_cursor *cursor)
{
    return cursor->eof;
}

bool sky_cursor_eos(sky_cursor *cursor)
{
    return !cursor->in_session;
}

void sky_cursor_set_session_idle(sky_cursor *cursor, uint32_t seconds)
{
    // Save the idle value.
    cursor->session_idle_in_sec = seconds;

    // If the value is non-zero then start sessionizing the cursor.
    cursor->in_session = (seconds > 0 ? false : !cursor->eof);
}

void sky_cursor_next_session(sky_cursor *cursor)
{
    // Set a flag to allow the cursor to continue iterating unless EOF is set.
    if(!cursor->in_session) {
        cursor->session_event_index = -1;
        cursor->in_session = !cursor->eof;
    }
}

bool sky_lua_cursor_next_session(sky_cursor *cursor)
{
    sky_cursor_next_session(cursor);
    return !cursor->eof;
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
	"github.com/jmhodges/levigo"
	"github.com/ugorji/go-msgpack"
	"regexp"
	"sort"
	"text/template"
	"unsafe"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// An ExecutionEngine is used to iterate over a series of objects.
type ExecutionEngine struct {
	tableName    string
	iterator     *levigo.Iterator
	cursor       *C.sky_cursor
	state        *C.lua_State
	header       string
	source       string
	fullSource   string
	propertyFile *PropertyFile
	propertyRefs []*Property
}

//------------------------------------------------------------------------------
//
// Constructor
//
//------------------------------------------------------------------------------

func NewExecutionEngine(table *Table, source string) (*ExecutionEngine, error) {
	if table == nil {
		return nil, errors.New("skyd.ExecutionEngine: Table required")
	}
	propertyFile := table.propertyFile
	if propertyFile == nil {
		return nil, errors.New("skyd.ExecutionEngine: Property file required")
	}

	// Find a list of all references properties.
	propertyRefs, err := extractPropertyReferences(propertyFile, source)
	if err != nil {
		return nil, err
	}

	// Create the engine.
	e := &ExecutionEngine{
		tableName:    table.Name,
		propertyFile: propertyFile,
		source:       source,
		propertyRefs: propertyRefs,
	}

	// Initialize the engine.
	err = e.init()
	if err != nil {
		fmt.Printf("%s\n\n", e.FullAnnotatedSource())
		e.Destroy()
		return nil, err
	}

	return e, nil
}

//------------------------------------------------------------------------------
//
// Properties
//
//------------------------------------------------------------------------------

// Retrieves the source for the engine.
func (e *ExecutionEngine) Source() string {
	return e.source
}

// Retrieves the generated header for the engine.
func (e *ExecutionEngine) Header() string {
	return e.header
}

// Retrieves the full source sent to the Lua compiler.
func (e *ExecutionEngine) FullSource() string {
	return e.fullSource
}

// Retrieves the full annotated source with line numbers.
func (e *ExecutionEngine) FullAnnotatedSource() string {
	lineNumber := 1
	r, _ := regexp.Compile(`\n`)
	return "00001 " + r.ReplaceAllStringFunc(e.fullSource, func(str string) string {
		lineNumber += 1
		return fmt.Sprintf("%s%05d ", str, lineNumber)
	})
}

// Sets the iterator to use.
func (e *ExecutionEngine) SetIterator(iterator *levigo.Iterator) error {
	// Close the old iterator (if it's not the one being set).
	if e.iterator != nil && e.iterator != iterator {
		e.iterator.Close()
	}

	// Attach the new iterator.
	e.iterator = iterator
	if e.iterator != nil {
		// Determine table prefix.
		prefix, err := TablePrefix(e.tableName)
		if err != nil {
			return err
		}
		prefix = prefix[0 : len(prefix)-1]

		// Set the prefix & seek.
		C.sky_cursor_set_key_prefix(e.cursor, (unsafe.Pointer(&prefix[0])), C.uint32_t(len(prefix)))
		e.iterator.Seek(prefix)

		// Assign the iterator to the cursor.
		e.cursor.leveldb_iterator = e.iterator.Iter
	}

	return nil
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Lifecycle
//--------------------------------------

// Initializes the Lua context and compiles the source code.
func (e *ExecutionEngine) init() error {
	if e.state != nil {
		return nil
	}

	// Initialize the state and open the libraries.
	e.state = C.luaL_newstate()
	if e.state == nil {
		return errors.New("Unable to initialize Lua context.")
	}
	C.luaL_openlibs(e.state)

	// Generate the header file.
	err := e.generateHeader()
	if err != nil {
		e.Destroy()
		return err
	}

	// Compile the script.
	e.fullSource = fmt.Sprintf("%v\n%v", e.header, e.source)
	source := C.CString(e.fullSource)
	defer C.free(unsafe.Pointer(source))
	ret := C.luaL_loadstring(e.state, source)
	if ret != 0 {
		defer e.Destroy()
		errstring := C.GoString(C.lua_tolstring(e.state, -1, nil))
		return fmt.Errorf("skyd.ExecutionEngine: Syntax Error: %v", errstring)
	}

	// Run script once to initialize.
	ret = C.lua_pcall(e.state, 0, 0, 0)
	if ret != 0 {
		defer e.Destroy()
		errstring := C.GoString(C.lua_tolstring(e.state, -1, nil))
		return fmt.Errorf("skyd.ExecutionEngine: Init Error: %v", errstring)
	}

	// Setup cursor.
	err = e.initCursor()
	if err != nil {
		e.Destroy()
		return err
	}

	return nil
}

// Initializes the cursor used by the script.
func (e *ExecutionEngine) initCursor() error {
	// Create the cursor.
	minPropertyId, maxPropertyId := e.propertyFile.NextIdentifiers()
	e.cursor = C.sky_cursor_new((C.int32_t)(minPropertyId), (C.int32_t)(maxPropertyId))

	// Initialize the cursor from within Lua.
	functionName := C.CString("sky_init_cursor")
	defer C.free(unsafe.Pointer(functionName))

	C.lua_getfield(e.state, -10002, functionName)
	C.lua_pushlightuserdata(e.state, unsafe.Pointer(e.cursor))
	//fmt.Printf("%s\n\n", e.FullAnnotatedSource())
	rc := C.lua_pcall(e.state, 1, 0, 0)
	if rc != 0 {
		luaErrString := C.GoString(C.lua_tolstring(e.state, -1, nil))
		return fmt.Errorf("Unable to init cursor: %s", luaErrString)
	}

	return nil
}

// Closes the lua context.
func (e *ExecutionEngine) Destroy() {
	if e.state != nil {
		C.lua_close(e.state)
		e.state = nil
	}
	if e.iterator != nil {
		e.SetIterator(nil)
	}
}

//--------------------------------------
// Execution
//--------------------------------------

// Initializes the data structure used for aggregation.
func (e *ExecutionEngine) Initialize() (interface{}, error) {
	functionName := C.CString("sky_initialize")
	defer C.free(unsafe.Pointer(functionName))

	C.lua_getfield(e.state, -10002, functionName)
	C.lua_pushlightuserdata(e.state, unsafe.Pointer(e.cursor))
	rc := C.lua_pcall(e.state, 1, 1, 0)
	if rc != 0 {
		luaErrString := C.GoString(C.lua_tolstring(e.state, -1, nil))
		fmt.Println(e.FullAnnotatedSource())
		return nil, fmt.Errorf("skyd.ExecutionEngine: Unable to initialize: %s", luaErrString)
	}

	return e.decodeResult()
}

// Executes an aggregation over the iterator.
func (e *ExecutionEngine) Aggregate(data interface{}) (interface{}, error) {
	functionName := C.CString("sky_aggregate")
	defer C.free(unsafe.Pointer(functionName))

	C.lua_getfield(e.state, -10002, functionName)
	C.lua_pushlightuserdata(e.state, unsafe.Pointer(e.cursor))
	err := e.encodeArgument(data)
	if err != nil {
		return nil, err
	}
	rc := C.lua_pcall(e.state, 2, 1, 0)
	if rc != 0 {
		luaErrString := C.GoString(C.lua_tolstring(e.state, -1, nil))
		fmt.Println(e.FullAnnotatedSource())
		return nil, fmt.Errorf("skyd.ExecutionEngine: Unable to aggregate: %s", luaErrString)
	}

	return e.decodeResult()
}

// Executes an merge over the iterator.
func (e *ExecutionEngine) Merge(results interface{}, data interface{}) (interface{}, error) {
	functionName := C.CString("sky_merge")
	defer C.free(unsafe.Pointer(functionName))

	C.lua_getfield(e.state, -10002, functionName)
	err := e.encodeArgument(results)
	if err != nil {
		return results, err
	}
	err = e.encodeArgument(data)
	if err != nil {
		return results, err
	}
	rc := C.lua_pcall(e.state, 2, 1, 0)
	if rc != 0 {
		luaErrString := C.GoString(C.lua_tolstring(e.state, -1, nil))
		fmt.Println(e.FullAnnotatedSource())
		return results, fmt.Errorf("skyd.ExecutionEngine: Unable to merge: %s", luaErrString)
	}

	return e.decodeResult()
}

// Encodes a Go object into Msgpack and adds it to the function arguments.
func (e *ExecutionEngine) encodeArgument(value interface{}) error {
	// Encode Go object into msgpack.
	buffer := new(bytes.Buffer)
	encoder := msgpack.NewEncoder(buffer)
	err := encoder.Encode(value)
	if err != nil {
		return err
	}

	// Push the msgpack data onto the Lua stack.
	data := buffer.String()
	cdata := C.CString(data)
	defer C.free(unsafe.Pointer(cdata))
	C.lua_pushlstring(e.state, cdata, (C.size_t)(len(data)))

	// Convert the argument from msgpack into Lua.
	rc := C.mp_unpack(e.state)
	if rc != 1 {
		return errors.New("skyd.ExecutionEngine: Unable to msgpack encode Lua argument")
	}
	C.lua_remove(e.state, -2)

	return nil
}

// Decodes the result from a function into a Go object.
func (e *ExecutionEngine) decodeResult() (interface{}, error) {
	// Encode Lua object into msgpack.
	rc := C.mp_pack(e.state)
	if rc != 1 {
		return nil, errors.New("skyd.ExecutionEngine: Unable to msgpack decode Lua result")
	}
	sz := C.size_t(0)
	ptr := C.lua_tolstring(e.state, -1, (*C.size_t)(&sz))
	str := C.GoStringN(ptr, (C.int)(sz))
	C.lua_settop(e.state, -(1)-1) // lua_pop()

	// Decode msgpack into a Go object.
	var ret interface{}
	decoder := msgpack.NewDecoder(bytes.NewBufferString(str), nil)
	err := decoder.Decode(&ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

//--------------------------------------
// Codegen
//--------------------------------------

// Generates the header for the script based on a source string.
func (e *ExecutionEngine) generateHeader() error {
	// Parse the header template.
	t := template.New("header.lua")
	t.Funcs(template.FuncMap{"structdef": propertyStructDef, "metatypedef": metatypeFunctionDef, "initdescriptor": initDescriptorDef})
	_, err := t.Parse(LuaHeader)
	if err != nil {
		return err
	}

	// Generate the template from the property references.
	var buffer bytes.Buffer
	err = t.Execute(&buffer, e.propertyRefs)
	if err != nil {
		return err
	}

	// Assign header
	e.header = buffer.String()

	return nil
}

// Extracts the property references from the source string.
func extractPropertyReferences(propertyFile *PropertyFile, source string) ([]*Property, error) {
	// Create a list of properties.
	properties := make([]*Property, 0)
	lookup := make(map[int64]*Property)

	// Find all the event property references in the script.
	r, err := regexp.Compile(`\bevent(?:\.|:)(\w+)`)
	if err != nil {
		return nil, err
	}
	for _, match := range r.FindAllStringSubmatch(source, -1) {
		name := match[1]
		property := propertyFile.GetPropertyByName(name)
		if property == nil {
			return nil, fmt.Errorf("Property not found: '%v'", name)
		}
		if lookup[property.Id] == nil {
			properties = append(properties, property)
			lookup[property.Id] = property
		}
	}
	sort.Sort(PropertyList(properties))

	return properties, nil
}

func propertyStructDef(args ...interface{}) string {
	if property, ok := args[0].(*Property); ok && property.Id != 0 {
		return fmt.Sprintf("%v _%v;", getPropertyCType(property), property.Name)
	}
	return ""
}

func metatypeFunctionDef(args ...interface{}) string {
	if property, ok := args[0].(*Property); ok && property.Id != 0 {
		switch property.DataType {
		case StringDataType:
			return fmt.Sprintf("%v = function(event) return ffi.string(event._%v.data, event._%v.length) end,", property.Name, property.Name, property.Name)
		default:
			return fmt.Sprintf("%v = function(event) return event._%v end,", property.Name, property.Name)
		}
	}
	return ""
}

func initDescriptorDef(args ...interface{}) string {
	if property, ok := args[0].(*Property); ok && property.Id != 0 {
		return fmt.Sprintf("cursor:set_property(%d, ffi.offsetof('sky_lua_event_t', '_%s'), ffi.sizeof('%s'), '%s')", property.Id, property.Name, getPropertyCType(property), property.DataType)
	}
	return ""
}

func getPropertyCType(property *Property) string {
	switch property.DataType {
	case StringDataType:
		return "sky_string_t"
	case FactorDataType, IntegerDataType:
		return "int32_t"
	case FloatDataType:
		return "double"
	case BooleanDataType:
		return "bool"
	default:
		panic(fmt.Sprintf("skyd.ExecutionEngine: Invalid data type: %v", property.DataType))
	}
	return ""
}
