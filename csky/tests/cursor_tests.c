#include <stdio.h>
#include <stdlib.h>
#include <stddef.h>
#include <math.h>

#include <sky/cursor.h>
#include <sky/data_descriptor.h>
#include <sky/sky_string.h>
#include <sky/timestamp.h>
#include <sky/mem.h>

#include "minunit.h"

//==============================================================================
//
// Fixtures
//
//==============================================================================

int DATA0_LENGTH = 128;

char *DATA0 = 
  // 1970-01-01T00:00:00Z, {1:"john doe", 2:1000, 3:100.2, 4:true}
  "\x92" "\xD3\x00\x00\x00\x00\x00\x00\x00\x00" "\x84" "\x01\xA8""john doe" "\x02\xD1\x03\xE8" "\x03\xCB\x40\x59\x0C\xCC\xCC\xCC\xCC\xCD" "\x04\xC3"
  // 1970-01-01T00:00:01Z, {-1:"A1", -2:"super", -3:21, -4:100, -5:true}
  "\x92" "\xD3\x00\x00\x00\x00\x00\x10\x00\x00" "\x85" "\xFF\xA2""A1" "\xFE\xA5""super" "\xFD\x15" "\xFC\xCB\x40\x59\x00\x00\x00\x00\x00\x00" "\xFB\xC3"
  // 1970-01-01T00:00:02Z, {-1:"A2"}
  "\x92" "\xD3\x00\x00\x00\x00\x00\x20\x00\x00" "\x81" "\xFF\xA2""A2"
  // 1970-01-01T00:00:03Z, {1:"frank sinatra", 2:20, 3:-100, 4:false}
  "\x92" "\xD3\x00\x00\x00\x00\x00\x30\x00\x00" "\x84" "\x01\xAD""frank sinatra" "\x02\x14" "\x03\xCB\xC0\x59\x00\x00\x00\x00\x00\x00" "\x04\xC2"
;


//==============================================================================
//
// Declarations
//
//==============================================================================

#define ASSERT_OBJ_STATE(OBJ, TS, ACTION, OSTRING, OINT, ODOUBLE, OBOOLEAN, ASTRING, AINT, ADOUBLE, ABOOLEAN) do {\
    mu_assert_int64_equals(OBJ.ts, TS); \
    mu_assert_int_equals(OBJ.action.length, (int)strlen(ACTION)); \
    mu_assert_bool(memcmp(OBJ.action.data, ACTION, strlen(ACTION)) == 0); \
    mu_assert_int_equals(OBJ.object_string.length, (int)strlen(OSTRING)); \
    mu_assert_bool(memcmp(OBJ.object_string.data, OSTRING, strlen(OSTRING)) == 0); \
    mu_assert_int64_equals(OBJ.object_int, OINT); \
    mu_assert_bool(fabs(OBJ.object_double-ODOUBLE) < 0.1); \
    mu_assert_bool(OBJ.object_boolean == OBOOLEAN); \
    mu_assert_int_equals(OBJ.action_string.length, (int)strlen(ASTRING)); \
    mu_assert_bool(memcmp(OBJ.action_string.data, ASTRING, strlen(ASTRING)) == 0); \
    mu_assert_int64_equals(OBJ.action_int, AINT); \
    mu_assert_bool(fabs(OBJ.action_double-ADOUBLE) < 0.1); \
    mu_assert_bool(OBJ.action_boolean == ABOOLEAN); \
} while(0)

#define ASSERT_OBJ_STATE2(OBJ, TIMESTAMP, ACTION, OINT, AINT) do {\
    mu_assert_int64_equals(OBJ.timestamp, TIMESTAMP); \
    mu_assert_int_equals(OBJ.action.length, (int)strlen(ACTION)); \
    mu_assert_bool(memcmp(OBJ.action.data, ACTION, strlen(ACTION)) == 0); \
    mu_assert_int64_equals(OBJ.object_int, OINT); \
    mu_assert_int64_equals(OBJ.action_int, AINT); \
} while(0)

typedef struct {
    uint32_t timestamp;
    int64_t ts;
    sky_string action;
    sky_string action_string;
    int64_t    action_int;
    double     action_double;
    bool       action_boolean;
    sky_string object_string;
    int64_t    object_int;
    double     object_double;
    bool       object_boolean;
} test_t;


//==============================================================================
//
// Test Cases
//
//==============================================================================

//--------------------------------------
// Set Data
//--------------------------------------

int test_sky_cursor_set_data() {
    // Setup data object & data descriptor.
    test_t obj; memset(&obj, 0, sizeof(obj));
    sky_data_descriptor *descriptor = sky_data_descriptor_new(-4, 4);
    descriptor->timestamp_descriptor.timestamp_offset = offsetof(test_t, timestamp);
    descriptor->timestamp_descriptor.ts_offset = offsetof(test_t, ts);
    sky_data_descriptor_set_property(descriptor, -5, offsetof(test_t, action_boolean), "boolean");
    sky_data_descriptor_set_property(descriptor, -4, offsetof(test_t, action_double), "float");
    sky_data_descriptor_set_property(descriptor, -3, offsetof(test_t, action_int), "integer");
    sky_data_descriptor_set_property(descriptor, -2, offsetof(test_t, action_string), "string");
    sky_data_descriptor_set_property(descriptor, -1, offsetof(test_t, action), "string");
    sky_data_descriptor_set_property(descriptor, 1, offsetof(test_t, object_string), "string");
    sky_data_descriptor_set_property(descriptor, 2, offsetof(test_t, object_int), "integer");
    sky_data_descriptor_set_property(descriptor, 3, offsetof(test_t, object_double), "float");
    sky_data_descriptor_set_property(descriptor, 4, offsetof(test_t, object_boolean), "boolean");

    sky_cursor *cursor = sky_cursor_new();
    cursor->data_descriptor = descriptor;
    cursor->data = &obj;

    sky_cursor_set_ptr(cursor, DATA0, DATA0_LENGTH);
    ASSERT_OBJ_STATE(obj, 0LL, "", "", 0LL, 0, false, "", 0LL, 0, false);

    // Event 1 (State-Only)
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    ASSERT_OBJ_STATE(obj, 0LL, "", "john doe", 1000LL, 100.2, true, "", 0LL, 0, false);
    
    // Event 2 (Action + Action Data)
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    ASSERT_OBJ_STATE(obj, sky_timestamp_shift(1000000LL), "A1", "john doe", 1000LL, 100.2, true, "super", 21LL, 100, true);
    
    // Event 3 (Action-Only)
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    ASSERT_OBJ_STATE(obj, sky_timestamp_shift(2000000LL), "A2", "john doe", 1000LL, 100.2, true, "", 0LL, 0, false);

    // Event 4 (Data-Only)
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    ASSERT_OBJ_STATE(obj, sky_timestamp_shift(3000000LL), "", "frank sinatra", 20LL, -100, false, "", 0LL, 0, false);

    // EOF
    mu_assert_bool(!sky_lua_cursor_next_event(cursor));

    sky_cursor_free(cursor);
    sky_data_descriptor_free(descriptor);
    return 0;
}


//--------------------------------------
// Sessionize
//--------------------------------------

int test_sky_cursor_sessionize() {
    // Setup data object & data descriptor.
    test_t obj; memset(&obj, 0, sizeof(obj));
    sky_data_descriptor *descriptor = sky_data_descriptor_new(-1, 1);
    descriptor->timestamp_descriptor.timestamp_offset = offsetof(test_t, timestamp);
    descriptor->timestamp_descriptor.ts_offset = offsetof(test_t, ts);
    sky_data_descriptor_set_property(descriptor, -2, offsetof(test_t, action_int), "integer");
    sky_data_descriptor_set_property(descriptor, -1, offsetof(test_t, action), "string");
    sky_data_descriptor_set_property(descriptor, 1, offsetof(test_t, object_int), "integer");

    sky_cursor *cursor = sky_cursor_new();
    cursor->data_descriptor = descriptor;
    cursor->data = &obj;

    // Initialize data and set a 10 second idle time.
    sky_cursor_set_ptr(cursor, DATA0, DATA0_LENGTH);
    sky_cursor_set_session_idle(cursor, 10);
    mu_assert_int_equals(cursor->session_event_index, -1);
    ASSERT_OBJ_STATE2(obj, 0, "", 0LL, 0LL);
    
    // Pre-session
    mu_assert_bool(sky_lua_cursor_next_event(cursor) == false);
    mu_assert_int_equals(cursor->session_event_index, -1);
    ASSERT_OBJ_STATE2(obj, 0, "", 0LL, 0LL);

    // Session 1
    mu_assert_bool(sky_lua_cursor_next_session(cursor));
    mu_assert_int_equals(cursor->session_event_index, -1);
    ASSERT_OBJ_STATE2(obj, 0, "", 0LL, 0LL);

    // Session 1, Event 1
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    mu_assert_int_equals(cursor->session_event_index, 0);
    ASSERT_OBJ_STATE2(obj, 0, "A1", 1000LL, 0LL);
    
    // Session 1, Event 2
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    mu_assert_int_equals(cursor->session_event_index, 1);
    ASSERT_OBJ_STATE2(obj, 1, "A2", 1000LL, 100LL);
    
    // Session 1, Event 3
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    mu_assert_int_equals(cursor->session_event_index, 2);
    ASSERT_OBJ_STATE2(obj, 10, "A3", 1000LL, 200LL);
    
    // Prevent next session!
    mu_assert_bool(sky_lua_cursor_next_event(cursor) == false);
    mu_assert_int_equals(cursor->session_event_index, 2);
    ASSERT_OBJ_STATE2(obj, 10, "A3", 1000LL, 200LL);
    

    // Session 2 (Single Event)
    mu_assert_bool(sky_lua_cursor_next_session(cursor));
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    mu_assert_int_equals(cursor->session_event_index, 0);
    ASSERT_OBJ_STATE2(obj, 20, "A1", 1000LL, 300LL);
    mu_assert_bool(sky_lua_cursor_next_event(cursor) == false);


    // Session 3 (with same data)
    mu_assert_bool(sky_lua_cursor_next_session(cursor));
    mu_assert_int_equals(cursor->session_event_index, -1);
    ASSERT_OBJ_STATE2(obj, 20, "A1", 1000LL, 300LL);

    // Session 3, Event 1
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    mu_assert_int_equals(cursor->session_event_index, 0);
    ASSERT_OBJ_STATE2(obj, 60, "A1", 2000LL, 0LL);

    // Session 3, Event 2
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    mu_assert_int_equals(cursor->session_event_index, 1);
    ASSERT_OBJ_STATE2(obj, 63, "A2", 2000LL, 400LL);

    // Prevent next session!
    mu_assert_bool(sky_lua_cursor_next_event(cursor) == false);
    mu_assert_bool(sky_lua_cursor_next_session(cursor) == false);

    // EOF!
    mu_assert_bool(cursor->eof == true);
    mu_assert_bool(cursor->in_session == false);

    // Reuse cursor.
    sky_cursor_set_ptr(cursor, DATA0, DATA0_LENGTH);
    mu_assert_int_equals(cursor->session_event_index, -1);
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    mu_assert_int_equals(cursor->session_event_index, 0);
    ASSERT_OBJ_STATE2(obj, 0, "A1", 1000LL, 0LL);
    
    sky_cursor_free(cursor);
    sky_data_descriptor_free(descriptor);
    return 0;
}



//==============================================================================
//
// Setup
//
//==============================================================================

int all_tests() {
    mu_run_test(test_sky_cursor_set_data);
    mu_run_test(test_sky_cursor_sessionize);
    return 0;
}

RUN_TESTS()
