package skyd

/*
#cgo CFLAGS:-Wno-pointer-to-int-cast
#include <stdio.h>
#include <stdlib.h>
#include <stddef.h>
#include <string.h>
#include <inttypes.h>
#include <stdbool.h>

typedef void (*sky_data_property_descriptor_set_func)(void *target, void *value, size_t *sz);
typedef void (*sky_data_property_descriptor_clear_func)(void *target);
typedef struct { int32_t length; char *data; } sky_string;
typedef struct { uint16_t ts_offset; uint16_t timestamp_offset;} sky_data_timestamp_descriptor;
typedef struct {
    int64_t property_id;
    uint16_t offset;
    sky_data_property_descriptor_set_func set_func;
    sky_data_property_descriptor_clear_func clear_func;
} sky_data_property_descriptor;

// Defines a collection of descriptors for a struct to serialize data into it.
typedef struct {
    sky_data_timestamp_descriptor timestamp_descriptor;
    sky_data_property_descriptor *property_descriptors;
    sky_data_property_descriptor *property_zero_descriptor;
    sky_data_property_descriptor **action_property_descriptors;
    uint32_t action_property_descriptor_count;
    uint32_t property_count;
    uint32_t active_property_count;
    uint32_t data_sz;
} sky_data_descriptor;

typedef struct sky_cursor {
    void *data;
    int32_t session_event_index;
    void *startptr;
    void *nextptr;
    void *endptr;
    void *ptr;
    bool eof;
    bool in_session;
    uint32_t last_timestamp;
    uint32_t session_idle_in_sec;
    sky_data_descriptor *data_descriptor;
} sky_cursor;


#define sky_event_flag_t uint8_t
#define EVENT_FLAG       0x92


//==============================================================================
//
// Macros
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

uint64_t bswap64(uint64_t value)
{
    return (
        ((value & 0x00000000000000FF) << 56) |
        ((value & 0x000000000000FF00) << 40) |
        ((value & 0x0000000000FF0000) << 24) |
        ((value & 0x00000000FF000000) << 8) |
        ((value & 0x000000FF00000000) >> 8) |
        ((value & 0x0000FF0000000000) >> 24) |
        ((value & 0x00FF000000000000) >> 40) |
        ((value & 0xFF00000000000000) >> 56)
    );
}

#if (BYTE_ORDER == LITTLE_ENDIAN) || (__BYTE_ORDER == __LITTLE_ENDIAN)
#define htonll(x) bswap64(x)
#define ntohll(x) bswap64(x)
#else
#define htonll(x) x
#define ntohll(x) x
#endif


#define memdump(PTR, LENGTH) do {\
    char *address = (char*)PTR;\
    int length = LENGTH;\
    int i = 0;\
    char *line = (char*)address;\
    unsigned char ch;\
    fprintf(stderr, "%09X | ", (int)address);\
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

#define badcursordata(MSG) do {\
    fprintf(stderr, "Cursor pointing at invalid raw event data [" MSG "]: %p", cursor->ptr); \
    memdump(cursor->startptr, (cursor->endptr - cursor->startptr)); \
    cursor->eof = true; \
    return; \
} while(0)


//==============================================================================
//
// Constants
//
//==============================================================================

//--------------------------------------
// Timestamp
//--------------------------------------

// The number of bits that seconds are shifted over in a timestamp.
#define SECONDS_BIT_OFFSET  20

//--------------------------------------
// Fixnum
//--------------------------------------

#define POS_FIXNUM_MIN          0
#define POS_FIXNUM_MAX          127
#define POS_FIXNUM_TYPE         0x00
#define POS_FIXNUM_TYPE_MASK    0x80
#define POS_FIXNUM_VALUE_MASK   0x7F
#define POS_FIXNUM_SIZE         1

#define NEG_FIXNUM_MIN          -32
#define NEG_FIXNUM_MAX          -1
#define NEG_FIXNUM_TYPE         0xE0
#define NEG_FIXNUM_TYPE_MASK    0xE0
#define NEG_FIXNUM_VALUE_MASK   0x1F
#define NEG_FIXNUM_SIZE         1


//--------------------------------------
// Unsigned integers
//--------------------------------------

#define UINT8_TYPE              0xCC
#define UINT8_SIZE              2

#define UINT16_TYPE             0xCD
#define UINT16_SIZE             3

#define UINT32_TYPE             0xCE
#define UINT32_SIZE             5

#define UINT64_TYPE             0xCF
#define UINT64_SIZE             9


//--------------------------------------
// Signed integers
//--------------------------------------

#define INT8_TYPE               0xD0
#define INT8_SIZE               2

#define INT16_TYPE              0xD1
#define INT16_SIZE              3

#define INT32_TYPE              0xD2
#define INT32_SIZE              5

#define INT64_TYPE              0xD3
#define INT64_SIZE              9


//--------------------------------------
// Nil
//--------------------------------------

#define NIL_TYPE                0xC0
#define NIL_SIZE                1


//--------------------------------------
// Boolean
//--------------------------------------

#define TRUE_TYPE                0xC3
#define FALSE_TYPE               0xC2
#define BOOL_SIZE                1


//--------------------------------------
// Floating point
//--------------------------------------

#define FLOAT_TYPE              0xCA
#define FLOAT_SIZE              5

#define DOUBLE_TYPE             0xCB
#define DOUBLE_SIZE             9


//--------------------------------------
// Raw bytes
//--------------------------------------

#define FIXRAW_TYPE             0xA0
#define FIXRAW_TYPE_MASK        0xE0
#define FIXRAW_VALUE_MASK       0x1F
#define FIXRAW_SIZE             1
#define FIXRAW_MAXSIZE          31

#define RAW16_TYPE              0xDA
#define RAW16_SIZE              3
#define RAW16_MAXSIZE           65535

#define RAW32_TYPE              0xDB
#define RAW32_SIZE              5
#define RAW32_MAXSIZE           4294967295


//--------------------------------------
// Array
//--------------------------------------

#define FIXARRAY_TYPE           0x90
#define FIXARRAY_TYPE_MASK      0xF0
#define FIXARRAY_VALUE_MASK     0x0F
#define FIXARRAY_SIZE           1
#define FIXARRAY_MAXSIZE        15

#define ARRAY16_TYPE            0xDC
#define ARRAY16_SIZE            3
#define ARRAY16_MAXSIZE         65535

#define ARRAY32_TYPE            0xDD
#define ARRAY32_SIZE            5
#define ARRAY32_MAXSIZE         4294967295


//--------------------------------------
// Map
//--------------------------------------

#define FIXMAP_TYPE             0x80
#define FIXMAP_TYPE_MASK        0xF0
#define FIXMAP_VALUE_MASK       0x0F
#define FIXMAP_SIZE             1
#define FIXMAP_MAXSIZE          15

#define MAP16_TYPE              0xDE
#define MAP16_SIZE              3
#define MAP16_MAXSIZE           65535

#define MAP32_TYPE              0xDF
#define MAP32_SIZE              5
#define MAP32_MAXSIZE           4294967295


//==============================================================================
//
// Forward Declarations
//
//==============================================================================

//--------------------------------------
// Setters
//--------------------------------------

void sky_data_descriptor_set_noop(void *target, void *value, size_t *sz);

void sky_data_descriptor_set_string(void *target, void *value, size_t *sz);

void sky_data_descriptor_set_int(void *target, void *value, size_t *sz);

void sky_data_descriptor_set_double(void *target, void *value, size_t *sz);

void sky_data_descriptor_set_boolean(void *target, void *value, size_t *sz);


//--------------------------------------
// Clear Functions
//--------------------------------------

void sky_data_descriptor_clear_string(void *target);

void sky_data_descriptor_clear_int(void *target);

void sky_data_descriptor_clear_double(void *target);

void sky_data_descriptor_clear_boolean(void *target);


//==============================================================================
//
// Minipack
//
//==============================================================================

bool minipack_is_pos_fixnum(void *ptr) {
    return (*((uint8_t*)ptr) & POS_FIXNUM_TYPE_MASK) == POS_FIXNUM_TYPE;
}

uint8_t minipack_unpack_pos_fixnum(void *ptr, size_t *sz)
{
    *sz = POS_FIXNUM_SIZE;
    uint8_t value = *((uint8_t*)ptr);
    return value & POS_FIXNUM_VALUE_MASK;
}

void minipack_pack_pos_fixnum(void *ptr, uint8_t value, size_t *sz)
{
    *sz = POS_FIXNUM_SIZE;
    *((uint8_t*)ptr) = value & POS_FIXNUM_VALUE_MASK;
}

bool minipack_is_neg_fixnum(void *ptr) {
    return (*((uint8_t*)ptr) & NEG_FIXNUM_TYPE_MASK) == NEG_FIXNUM_TYPE;
}

int8_t minipack_unpack_neg_fixnum(void *ptr, size_t *sz)
{
    *sz = NEG_FIXNUM_SIZE;
    int8_t value = *((int8_t*)ptr) & NEG_FIXNUM_VALUE_MASK;
    return (32-value) * -1;
}

void minipack_pack_neg_fixnum(void *ptr, int8_t value, size_t *sz)
{
    *sz = NEG_FIXNUM_SIZE;
    *((int8_t*)ptr) = (32 + value) | NEG_FIXNUM_TYPE;
}

bool minipack_is_uint8(void *ptr) {
    return (*((uint8_t*)ptr) == UINT8_TYPE);
}

uint8_t minipack_unpack_uint8(void *ptr, size_t *sz)
{
    *sz = UINT8_SIZE;
    return *((uint8_t*)(ptr+1));
}

void minipack_pack_uint8(void *ptr, uint8_t value, size_t *sz)
{
    *sz = UINT8_SIZE;
    *((uint8_t*)ptr)     = UINT8_TYPE;
    *((uint8_t*)(ptr+1)) = value;
}

bool minipack_is_uint16(void *ptr) {
    return (*((uint8_t*)ptr) == UINT16_TYPE);
}

uint16_t minipack_unpack_uint16(void *ptr, size_t *sz)
{
    *sz = UINT16_SIZE;
    uint16_t value = *((uint16_t*)(ptr+1));
    return ntohs(value);
}

void minipack_pack_uint16(void *ptr, uint16_t value, size_t *sz)
{
    *sz = UINT16_SIZE;
    *((uint8_t*)ptr)      = UINT16_TYPE;
    *((uint16_t*)(ptr+1)) = htons(value);
}

bool minipack_is_uint32(void *ptr) {
    return (*((uint8_t*)ptr) == UINT32_TYPE);
}

uint32_t minipack_unpack_uint32(void *ptr, size_t *sz)
{
    *sz = UINT32_SIZE;
    uint32_t value = *((uint32_t*)(ptr+1));
    return ntohl(value);
}

void minipack_pack_uint32(void *ptr, uint32_t value, size_t *sz)
{
    *sz = UINT32_SIZE;
    *((uint8_t*)ptr)      = UINT32_TYPE;
    *((uint32_t*)(ptr+1)) = htonl(value);
}

bool minipack_is_uint64(void *ptr) {
    return (*((uint8_t*)ptr) == UINT64_TYPE);
}

uint64_t minipack_unpack_uint64(void *ptr, size_t *sz)
{
    *sz = UINT64_SIZE;
    uint64_t value = *((uint64_t*)(ptr+1));
    return ntohll(value);
}

void minipack_pack_uint64(void *ptr, uint64_t value, size_t *sz)
{
    *sz = UINT64_SIZE;
    *((uint8_t*)ptr)      = UINT64_TYPE;
    *((uint64_t*)(ptr+1)) = htonll(value);
}

size_t minipack_sizeof_uint_elem(void *ptr)
{
    if(minipack_is_pos_fixnum(ptr)) {
        return POS_FIXNUM_SIZE;
    }
    else if(minipack_is_uint8(ptr)) {
        return UINT8_SIZE;
    }
    else if(minipack_is_uint16(ptr)) {
        return UINT16_SIZE;
    }
    else if(minipack_is_uint32(ptr)) {
        return UINT32_SIZE;
    }
    else if(minipack_is_uint64(ptr)) {
        return UINT64_SIZE;
    }
    else {
        return 0;
    }
}

bool minipack_is_int8(void *ptr) {
    return (*((uint8_t*)ptr) == INT8_TYPE);
}

int8_t minipack_unpack_int8(void *ptr, size_t *sz)
{
    *sz = INT8_SIZE;
    return *((int8_t*)(ptr+1));
}

void minipack_pack_int8(void *ptr, int8_t value, size_t *sz)
{
    *sz = INT8_SIZE;
    *((uint8_t*)ptr)    = INT8_TYPE;
    *((int8_t*)(ptr+1)) = value;
}

bool minipack_is_int16(void *ptr) {
    return (*((uint8_t*)ptr) == INT16_TYPE);
}

int16_t minipack_unpack_int16(void *ptr, size_t *sz)
{
    *sz = INT16_SIZE;
    return ntohs(*((int16_t*)(ptr+1)));
}

void minipack_pack_int16(void *ptr, int16_t value, size_t *sz)
{
    *sz = INT16_SIZE;
    *((uint8_t*)ptr)     = INT16_TYPE;
    *((int16_t*)(ptr+1)) = htons(value);
}

bool minipack_is_int32(void *ptr) {
    return (*((uint8_t*)ptr) == INT32_TYPE);
}

int32_t minipack_unpack_int32(void *ptr, size_t *sz)
{
    *sz = INT32_SIZE;
    return ntohl(*((int32_t*)(ptr+1)));
}

void minipack_pack_int32(void *ptr, int32_t value, size_t *sz)
{
    *sz = INT32_SIZE;
    *((uint8_t*)ptr)     = INT32_TYPE;
    *((int32_t*)(ptr+1)) = htonl(value);
}

bool minipack_is_int64(void *ptr) {
    return (*((uint8_t*)ptr) == INT64_TYPE);
}

int64_t minipack_unpack_int64(void *ptr, size_t *sz)
{
    *sz = INT64_SIZE;
    return ntohll(*((int64_t*)(ptr+1)));
}

void minipack_pack_int64(void *ptr, int64_t value, size_t *sz)
{
    *sz = INT64_SIZE;
    *((uint8_t*)ptr)     = INT64_TYPE;
    *((int64_t*)(ptr+1)) = htonll(value);
}


int64_t minipack_unpack_int(void *ptr, size_t *sz)
{
    if(minipack_is_pos_fixnum(ptr)) {
        return (int64_t)minipack_unpack_pos_fixnum(ptr, sz);
    }
    if(minipack_is_neg_fixnum(ptr)) {
        return (int64_t)minipack_unpack_neg_fixnum(ptr, sz);
    }
    else if(minipack_is_int8(ptr)) {
        return (int64_t)minipack_unpack_int8(ptr, sz);
    }
    else if(minipack_is_int16(ptr)) {
        return (int64_t)minipack_unpack_int16(ptr, sz);
    }
    else if(minipack_is_int32(ptr)) {
        return (int64_t)minipack_unpack_int32(ptr, sz);
    }
    else if(minipack_is_int64(ptr)) {
        return minipack_unpack_int64(ptr, sz);
    }
    // Fallback to unsigned ints.
    else if(minipack_is_uint8(ptr)) {
        return (int64_t)minipack_unpack_uint8(ptr, sz);
    }
    else if(minipack_is_uint16(ptr)) {
        return (int64_t)minipack_unpack_uint16(ptr, sz);
    }
    else if(minipack_is_uint32(ptr)) {
        return (int64_t)minipack_unpack_uint32(ptr, sz);
    }
    else if(minipack_is_uint64(ptr)) {
        return minipack_unpack_uint64(ptr, sz);
    }
    else {
        *sz = 0;
        return 0;
    }
}

void minipack_pack_int(void *ptr, int64_t value, size_t *sz)
{
    if(value >= POS_FIXNUM_MIN && value <= POS_FIXNUM_MAX) {
        minipack_pack_pos_fixnum(ptr, (int8_t)value, sz);
    }
    else if(value >= NEG_FIXNUM_MIN && value <= NEG_FIXNUM_MAX) {
        minipack_pack_neg_fixnum(ptr, (int8_t)value, sz);
    }
    else if(value >= INT8_MIN && value <= INT8_MAX) {
        minipack_pack_int8(ptr, (int8_t)value, sz);
    }
    else if(value >= INT16_MIN && value <= INT16_MAX) {
        minipack_pack_int16(ptr, (int16_t)value, sz);
    }
    else if(value >= INT32_MIN && value <= INT32_MAX) {
        minipack_pack_int32(ptr, (int32_t)value, sz);
    }
    else if(value >= INT64_MIN && value <= INT64_MAX) {
        minipack_pack_int64(ptr, value, sz);
    }
    else {
        *sz = 0;
    }
}

size_t minipack_sizeof_int_elem(void *ptr)
{
    if(minipack_is_pos_fixnum(ptr)) {
        return POS_FIXNUM_SIZE;
    }
    else if(minipack_is_neg_fixnum(ptr)) {
        return NEG_FIXNUM_SIZE;
    }
    else if(minipack_is_int8(ptr)) {
        return INT8_SIZE;
    }
    else if(minipack_is_int16(ptr)) {
        return INT16_SIZE;
    }
    else if(minipack_is_int32(ptr)) {
        return INT32_SIZE;
    }
    else if(minipack_is_int64(ptr)) {
        return INT64_SIZE;
    }
    // Fallback to unsigned ints.
    else if(minipack_is_uint8(ptr)) {
        return UINT8_SIZE;
    }
    else if(minipack_is_uint16(ptr)) {
        return UINT16_SIZE;
    }
    else if(minipack_is_uint32(ptr)) {
        return UINT32_SIZE;
    }
    else if(minipack_is_uint64(ptr)) {
        return UINT64_SIZE;
    }
    else {
        return 0;
    }
}

size_t minipack_sizeof_nil() {
    return NIL_SIZE;
}

bool minipack_is_nil(void *ptr) {
    return (*((uint8_t*)ptr) == NIL_TYPE);
}

void minipack_unpack_nil(void *ptr, size_t *sz)
{
    if(minipack_is_nil(ptr)) {
        *sz = NIL_SIZE;
    }
    else {
        *sz = 0;
    }
}

void minipack_pack_nil(void *ptr, size_t *sz)
{
    *sz = NIL_SIZE;
    *((uint8_t*)ptr) = NIL_TYPE;
}

bool minipack_is_true(void *ptr) {
    return (*((uint8_t*)ptr) == TRUE_TYPE);
}

bool minipack_is_false(void *ptr) {
    return (*((uint8_t*)ptr) == FALSE_TYPE);
}

size_t minipack_sizeof_bool() {
    return BOOL_SIZE;
}

bool minipack_is_bool(void *ptr) {
    return minipack_is_true(ptr) || minipack_is_false(ptr);
}

bool minipack_unpack_bool(void *ptr, size_t *sz)
{
    if(minipack_is_true(ptr)) {
        *sz = BOOL_SIZE;
        return true;
    }
    else if(minipack_is_false(ptr)) {
        *sz = BOOL_SIZE;
        return false;
    }
    else {
        *sz = 0;
        return false;
    }
}

void minipack_pack_bool(void *ptr, bool value, size_t *sz)
{
    *sz = BOOL_SIZE;
    
    if(value) {
        *((uint8_t*)ptr) = TRUE_TYPE;
    }
    else {
        *((uint8_t*)ptr) = FALSE_TYPE;
    }
}

size_t minipack_sizeof_float() {
    return FLOAT_SIZE;
}

bool minipack_is_float(void *ptr) {
    return (*((uint8_t*)ptr) == FLOAT_TYPE);
}

float minipack_unpack_float(void *ptr, size_t *sz)
{
    *sz = FLOAT_SIZE;
    
    // Cast bytes to int32 to use ntohl.
    uint32_t value = *((uint32_t*)(ptr+1));
    value = ntohl(value);
    float *float_value = (float*)&value;
    return *float_value;
}

void minipack_pack_float(void *ptr, float value, size_t *sz)
{
    *sz = FLOAT_SIZE;
    
    uint32_t *bytes_ptr = (uint32_t*)&value;
    uint32_t bytes = *bytes_ptr;
    bytes = htonl(bytes);
    *((uint8_t*)ptr)   = FLOAT_TYPE;
    float *float_ptr = (float*)&bytes;
    *((float*)(ptr+1)) = *float_ptr;
}

size_t minipack_sizeof_double() {
    return DOUBLE_SIZE;
}

bool minipack_is_double(void *ptr) {
    return (*((uint8_t*)ptr) == DOUBLE_TYPE);
}

double minipack_unpack_double(void *ptr, size_t *sz)
{
    *sz = DOUBLE_SIZE;
    // Cast bytes to int64 to use ntohll.
    uint64_t value = *((uint64_t*)(ptr+1));
    value = ntohll(value);
    double *double_ptr = (double*)&value;
    return *double_ptr;
}

void minipack_pack_double(void *ptr, double value, size_t *sz)
{
    *sz = DOUBLE_SIZE;
    uint64_t *bytes_ptr = (uint64_t*)&value;
    uint64_t bytes = htonll(*bytes_ptr);
    *((uint8_t*)ptr)    = DOUBLE_TYPE;
    double *double_ptr = (double*)&bytes;
    *((double*)(ptr+1)) = *double_ptr;
}

bool minipack_is_fixraw(void *ptr) {
    return (*((uint8_t*)ptr) & FIXRAW_TYPE_MASK) == FIXRAW_TYPE;
}

uint8_t minipack_unpack_fixraw(void *ptr, size_t *sz)
{
    *sz = FIXRAW_SIZE;
    return *((uint8_t*)ptr) & FIXRAW_VALUE_MASK;
}

void minipack_pack_fixraw(void *ptr, uint8_t length, size_t *sz)
{
    *sz = FIXRAW_SIZE;
    *((uint8_t*)ptr) = (length & FIXRAW_VALUE_MASK) | FIXRAW_TYPE;
}

bool minipack_is_raw16(void *ptr) {
    return (*((uint8_t*)ptr) == RAW16_TYPE);
}

uint16_t minipack_unpack_raw16(void *ptr, size_t *sz)
{
    *sz = RAW16_SIZE;
    return ntohs(*((uint16_t*)(ptr+1)));
}

void minipack_pack_raw16(void *ptr, uint16_t length, size_t *sz)
{
    *sz = RAW16_SIZE;
    *((uint8_t*)ptr)      = RAW16_TYPE;
    *((uint16_t*)(ptr+1)) = htons(length);
}

bool minipack_is_raw32(void *ptr) {
    return (*((uint8_t*)ptr) == RAW32_TYPE);
}

uint32_t minipack_unpack_raw32(void *ptr, size_t *sz)
{
    *sz = RAW32_SIZE;
    return ntohl(*((uint32_t*)(ptr+1)));
}

void minipack_pack_raw32(void *ptr, uint32_t length, size_t *sz)
{
    *sz = RAW32_SIZE;
    *((uint8_t*)ptr)      = RAW32_TYPE;
    *((uint32_t*)(ptr+1)) = htonl(length);
}

bool minipack_is_raw(void *ptr)
{
    return minipack_is_fixraw(ptr) || minipack_is_raw16(ptr) || minipack_is_raw32(ptr);
}

size_t minipack_sizeof_raw(uint32_t length)
{
    if(length <= FIXRAW_MAXSIZE) {
        return FIXRAW_SIZE;
    }
    else if(length <= RAW16_MAXSIZE) {
        return RAW16_SIZE;
    }

    return RAW32_SIZE;
}

size_t minipack_sizeof_raw_elem(void *ptr)
{
    if(minipack_is_fixraw(ptr)) {
        return FIXRAW_SIZE;
    }
    else if(minipack_is_raw16(ptr)) {
        return RAW16_SIZE;
    }
    else if(minipack_is_raw32(ptr)) {
        return RAW32_SIZE;
    }
    else {
        return 0;
    }
}

uint32_t minipack_unpack_raw(void *ptr, size_t *sz)
{
    if(minipack_is_fixraw(ptr)) {
        return (uint32_t)minipack_unpack_fixraw(ptr, sz);
    }
    else if(minipack_is_raw16(ptr)) {
        return (uint32_t)minipack_unpack_raw16(ptr, sz);
    }
    else if(minipack_is_raw32(ptr)) {
        return minipack_unpack_raw32(ptr, sz);
    }
    else {
        *sz = 0;
        return 0;
    }
}

void minipack_pack_raw(void *ptr, uint32_t length, size_t *sz)
{
    if(length <= FIXRAW_MAXSIZE) {
        minipack_pack_fixraw(ptr, (uint8_t)length, sz);
        return;
    }
    else if(length <= RAW16_MAXSIZE) {
        minipack_pack_raw16(ptr, (uint16_t)length, sz);
        return;
    }
    minipack_pack_raw32(ptr, length, sz);
}

bool minipack_is_fixarray(void *ptr) {
    return (*((uint8_t*)ptr) & FIXARRAY_TYPE_MASK) == FIXARRAY_TYPE;
}

uint8_t minipack_unpack_fixarray(void *ptr, size_t *sz)
{
    *sz = FIXARRAY_SIZE;
    return *((uint8_t*)ptr) & FIXARRAY_VALUE_MASK;
}

void minipack_pack_fixarray(void *ptr, uint8_t count, size_t *sz)
{
    *sz = FIXARRAY_SIZE;
    *((uint8_t*)ptr) = (count & FIXARRAY_VALUE_MASK) | FIXARRAY_TYPE;
}

bool minipack_is_array16(void *ptr) {
    return (*((uint8_t*)ptr) == ARRAY16_TYPE);
}

uint16_t minipack_unpack_array16(void *ptr, size_t *sz)
{
    *sz = ARRAY16_SIZE;
    return ntohs(*((uint16_t*)(ptr+1)));
}

void minipack_pack_array16(void *ptr, uint16_t count, size_t *sz)
{
    *sz = ARRAY16_SIZE;
    
    // Write header.
    *((uint8_t*)ptr)      = ARRAY16_TYPE;
    *((uint16_t*)(ptr+1)) = htons(count);
}

bool minipack_is_array32(void *ptr) {
    return (*((uint8_t*)ptr) == ARRAY32_TYPE);
}

uint32_t minipack_unpack_array32(void *ptr, size_t *sz)
{
    *sz = ARRAY32_SIZE;
    return ntohl(*((uint32_t*)(ptr+1)));
}

void minipack_pack_array32(void *ptr, uint32_t count, size_t *sz) {
    *sz = ARRAY32_SIZE;
    
    // Write header.
    *((uint8_t*)ptr)      = ARRAY32_TYPE;
    *((uint32_t*)(ptr+1)) = htonl(count);
}

bool minipack_is_array(void *ptr) {
    return minipack_is_fixarray(ptr) || minipack_is_array16(ptr) || minipack_is_array32(ptr);
}

size_t minipack_sizeof_array(uint32_t count)
{
    if(count <= FIXARRAY_MAXSIZE) {
        return FIXARRAY_SIZE;
    }
    else if(count <= ARRAY16_MAXSIZE) {
        return ARRAY16_SIZE;
    }
    return ARRAY32_SIZE;
}

size_t minipack_sizeof_array_elem(void *ptr)
{
    if(minipack_is_fixarray(ptr)) {
        return FIXARRAY_SIZE;
    }
    else if(minipack_is_array16(ptr)) {
        return ARRAY16_SIZE;
    }
    else if(minipack_is_array32(ptr)) {
        return ARRAY32_SIZE;
    }
    else {
        return 0;
    }
}

uint32_t minipack_unpack_array(void *ptr, size_t *sz)
{
    if(minipack_is_fixarray(ptr)) {
        return (uint32_t)minipack_unpack_fixarray(ptr, sz);
    }
    else if(minipack_is_array16(ptr)) {
        return (uint32_t)minipack_unpack_array16(ptr, sz);
    }
    else if(minipack_is_array32(ptr)) {
        return minipack_unpack_array32(ptr, sz);
    }
    else {
        *sz = 0;
        return 0;
    }
}

void minipack_pack_array(void *ptr, uint32_t count, size_t *sz)
{
    if(count <= FIXARRAY_MAXSIZE) {
        minipack_pack_fixarray(ptr, (uint8_t)count, sz);
        return;
    }
    else if(count <= ARRAY16_MAXSIZE) {
        minipack_pack_array16(ptr, (uint16_t)count, sz);
        return;
    }
    minipack_pack_array32(ptr, count, sz);
}

bool minipack_is_fixmap(void *ptr) {
    return (*((uint8_t*)ptr) & FIXMAP_TYPE_MASK) == FIXMAP_TYPE;
}

uint8_t minipack_unpack_fixmap(void *ptr, size_t *sz)
{
    *sz = FIXMAP_SIZE;
    return *((uint8_t*)ptr) & FIXMAP_VALUE_MASK;
}

void minipack_pack_fixmap(void *ptr, uint8_t count, size_t *sz)
{
    *sz = FIXMAP_SIZE;
    *((uint8_t*)ptr) = (count & FIXMAP_VALUE_MASK) | FIXMAP_TYPE;
}

bool minipack_is_map16(void *ptr) {
    return (*((uint8_t*)ptr) == MAP16_TYPE);
}

uint16_t minipack_unpack_map16(void *ptr, size_t *sz)
{
    *sz = MAP16_SIZE;
    return ntohs(*((uint16_t*)(ptr+1)));
}

void minipack_pack_map16(void *ptr, uint16_t count, size_t *sz)
{
    *sz = MAP16_SIZE;
    
    // Write header.
    *((uint8_t*)ptr)      = MAP16_TYPE;
    *((uint16_t*)(ptr+1)) = htons(count);
}

bool minipack_is_map32(void *ptr) {
    return (*((uint8_t*)ptr) == MAP32_TYPE);
}

uint32_t minipack_unpack_map32(void *ptr, size_t *sz)
{
    *sz = MAP32_SIZE;
    return ntohl(*((uint32_t*)(ptr+1)));
}

void minipack_pack_map32(void *ptr, uint32_t count, size_t *sz)
{
    *sz = MAP32_SIZE;
    
    // Write header.
    *((uint8_t*)ptr)      = MAP32_TYPE;
    *((uint32_t*)(ptr+1)) = htonl(count);
}

bool minipack_is_map(void *ptr) {
    return minipack_is_fixmap(ptr) || minipack_is_map16(ptr) || minipack_is_map32(ptr);
}

size_t minipack_sizeof_map(uint32_t count)
{
    if(count <= FIXMAP_MAXSIZE) {
        return FIXMAP_SIZE;
    }
    else if(count <= MAP16_MAXSIZE) {
        return MAP16_SIZE;
    }
    return MAP32_SIZE;
}

size_t minipack_sizeof_map_elem(void *ptr)
{
    if(minipack_is_fixmap(ptr)) {
        return FIXMAP_SIZE;
    }
    else if(minipack_is_map16(ptr)) {
        return MAP16_SIZE;
    }
    else if(minipack_is_map32(ptr)) {
        return MAP32_SIZE;
    }
    else {
        return 0;
    }
}

uint32_t minipack_unpack_map(void *ptr, size_t *sz)
{
    if(minipack_is_fixmap(ptr)) {
        return (uint32_t)minipack_unpack_fixmap(ptr, sz);
    }
    else if(minipack_is_map16(ptr)) {
        return (uint32_t)minipack_unpack_map16(ptr, sz);
    }
    else if(minipack_is_map32(ptr)) {
        return minipack_unpack_map32(ptr, sz);
    }
    else {
        *sz = 0;
        return 0;
    }
}

void minipack_pack_map(void *ptr, uint32_t count, size_t *sz)
{
    if(count <= FIXMAP_MAXSIZE) {
        minipack_pack_fixmap(ptr, (uint8_t)count, sz);
        return;
    }
    else if(count <= MAP16_MAXSIZE) {
        minipack_pack_map16(ptr, (uint16_t)count, sz);
        return;
    }
    minipack_pack_map32(ptr, count, sz);
}

size_t minipack_sizeof_elem_and_data(void *ptr)
{
    size_t sz;
    
    // Integer.
    sz = minipack_sizeof_int_elem(ptr);
    if(sz > 0) return sz;
    
    // Unsigned Integer.
    sz = minipack_sizeof_uint_elem(ptr);
    if(sz > 0) return sz;
    
    // Float & Double
    if(minipack_is_float(ptr)) return minipack_sizeof_float();
    if(minipack_is_double(ptr)) return minipack_sizeof_double();

    // Nil & Boolean.
    if(minipack_is_nil(ptr)) return minipack_sizeof_nil();
    if(minipack_is_bool(ptr)) return minipack_sizeof_bool();

    // Raw
    uint32_t length = minipack_unpack_raw(ptr, &sz);
    if(sz > 0) return sz + length;
    
    // Map, Array and other data returns 0.
    return 0;
}


//==============================================================================
//
// Data Descriptor
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

// Creates a data descriptor with a given number of properties.
sky_data_descriptor *sky_data_descriptor_create(int64_t min_property_id, int64_t max_property_id)
{
    sky_data_descriptor *descriptor = NULL;

    // Add one property to account for the zero descriptor.
    int64_t property_count = (max_property_id - min_property_id) + 1;

    // Allocate memory for the descriptor and child descriptors in one block.
    size_t sz = sizeof(sky_data_descriptor) + (sizeof(sky_data_property_descriptor) * property_count);
    descriptor = calloc(sz, 1);
    descriptor->property_descriptors = (sky_data_property_descriptor*)(((void*)descriptor) + sizeof(sky_data_descriptor));
    descriptor->property_count = property_count;
    descriptor->property_zero_descriptor = NULL;
    
    // Initialize all property descriptors to noop.
    uint32_t i;
    for(i=0; i<property_count; i++) {
        int64_t property_id = min_property_id + (int64_t)i;
        descriptor->property_descriptors[i].property_id = property_id;
        descriptor->property_descriptors[i].set_func = sky_data_descriptor_set_noop;
        
        // Save a pointer to the descriptor that points to property zero.
        if(property_id == 0) {
            descriptor->property_zero_descriptor = &descriptor->property_descriptors[i];
        }
    }

    return descriptor;
}

// Removes a data descriptor reference from memory.
//
// descriptor - The data descriptor to free.
void sky_data_descriptor_free(sky_data_descriptor *descriptor)
{
    if(descriptor) {
        if(descriptor->action_property_descriptors != NULL) free(descriptor->action_property_descriptors);
        descriptor->action_property_descriptors = NULL;
        descriptor->action_property_descriptor_count = 0;
        
        free(descriptor);
    }
}


//--------------------------------------
// Value Management
//--------------------------------------

// Assigns a value to a struct through the data descriptor.
//
// descriptor  - The data descriptor.
// target      - The target struct.
// property_id - The property id to set.
// ptr         - A pointer to the memory location to read the value from.
// sz          - A pointer to where the number of bytes read are returned to.
//
// Returns 0 if successful, otherwise returns -1.
void sky_data_descriptor_set_value(sky_data_descriptor *descriptor,
                                   void *target, int64_t property_id,
                                   void *ptr, size_t *sz)
{
    // Find the property descriptor and call the set_func.
    sky_data_property_descriptor *property_descriptor = &descriptor->property_zero_descriptor[property_id];
    property_descriptor->set_func(target + property_descriptor->offset, ptr, sz);
}


//--------------------------------------
// Descriptor Management
//--------------------------------------

void sky_data_descriptor_set_data_sz(sky_data_descriptor *descriptor, uint32_t sz) {
    descriptor->data_sz = sz;
}

void sky_data_descriptor_set_timestamp_offset(sky_data_descriptor *descriptor, uint32_t offset) {
    descriptor->timestamp_descriptor.timestamp_offset = offset;
}

void sky_data_descriptor_set_ts_offset(sky_data_descriptor *descriptor, uint32_t offset) {
    descriptor->timestamp_descriptor.ts_offset = offset;
}

// Sets the data type and offset for a given property id.
//
// descriptor  - The data descriptor.
// property_id - The property id.
// offset      - The offset in the struct where the data should be set.
// data_type   - The data type to set.
//
// Returns 0 if successful, otherwise returns -1.
void sky_data_descriptor_set_property(sky_data_descriptor *descriptor,
                                      int64_t property_id,
                                      uint32_t offset,
                                      const char *data_type)
{
    // Retrieve property descriptor.
    sky_data_property_descriptor *property_descriptor = &descriptor->property_zero_descriptor[property_id];
    
    // Increment active property count if a new property is being set.
    if(property_descriptor->set_func == sky_data_descriptor_set_noop) {
        descriptor->active_property_count++;
    }
    
    // Set the offset and set_func function on the descriptor.
    property_descriptor->offset = offset;
    if(strlen(data_type) == 0) {
        property_descriptor->set_func = sky_data_descriptor_set_noop;
        property_descriptor->clear_func = NULL;
    }
    else if(strcmp(data_type, "string") == 0) {
        property_descriptor->set_func = sky_data_descriptor_set_string;
        property_descriptor->clear_func = sky_data_descriptor_clear_string;
    }
    else if(strcmp(data_type, "integer") == 0) {
        property_descriptor->set_func = sky_data_descriptor_set_int;
        property_descriptor->clear_func = sky_data_descriptor_clear_int;
    }
    else if(strcmp(data_type, "float") == 0) {
        property_descriptor->set_func = sky_data_descriptor_set_double;
        property_descriptor->clear_func = sky_data_descriptor_clear_double;
    }
    else if(strcmp(data_type, "boolean") == 0) {
        property_descriptor->set_func = sky_data_descriptor_set_boolean;
        property_descriptor->clear_func = sky_data_descriptor_clear_boolean;
    }

    // If descriptor is an action descriptor and it doesn't exist in our list then add it.
    if(property_id < 0) {
        uint32_t i;
        bool found = false;
        for(i=0; i<descriptor->action_property_descriptor_count; i++) {
            if(descriptor->action_property_descriptors[i]->property_id == property_id) {
                found = true;
                break;
            }
        }
        
        // If it's not in the list then add it.
        if(!found) {
            descriptor->action_property_descriptor_count++;
            descriptor->action_property_descriptors = realloc(descriptor->action_property_descriptors, sizeof(*descriptor->action_property_descriptors) * descriptor->action_property_descriptor_count);
            descriptor->action_property_descriptors[descriptor->action_property_descriptor_count-1] = property_descriptor;
        }
    }
}


//--------------------------------------
// Setters
//--------------------------------------

// Performs no operation on the data. Used for skipping over unused properties
// in the descriptor.
//
// target - Unused.
// value  - The memory location where a MessagePack value is encoded.
// sz     - A pointer to where the number of bytes read should be returned.
//
// Returns nothing.
void sky_data_descriptor_set_noop(void *target, void *value, size_t *sz)
{
    *sz = minipack_sizeof_elem_and_data(value);
}

// Reads a MessagePack string value from memory and sets the value to the
// given memory location.
//
// target - The place where the sky string should be written to.
// value  - The memory location where a MessagePack encoded raw is located.
// sz     - A pointer to where the number of bytes read should be returned.
//
// Returns nothing.
void sky_data_descriptor_set_string(void *target, void *value, size_t *sz)
{
    size_t _sz;
    sky_string *string = (sky_string*)target;
    string->length = minipack_unpack_raw(value, &_sz);
    string->data = (_sz > 0 ? value + _sz : NULL);
    *sz = _sz + string->length;
}

// Reads a MessagePack integer value from memory and sets the value to the
// given memory location.
//
// target - The place where the integer should be written to.
// value  - The memory location where a MessagePack encoded int is located.
// sz     - A pointer to where the number of bytes read should be returned.
//
// Returns nothing.
void sky_data_descriptor_set_int(void *target, void *value, size_t *sz)
{
    *((int32_t*)target) = (int32_t)minipack_unpack_int(value, sz);
}

// Reads a MessagePack double value from memory and sets the value to the
// given memory location.
//
// target - The place where the integer should be written to.
// value  - The memory location where a MessagePack encoded double is located.
// sz     - A pointer to where the number of bytes read should be returned.
//
// Returns nothing.
void sky_data_descriptor_set_double(void *target, void *value, size_t *sz)
{
    *((double*)target) = minipack_unpack_double(value, sz);
}

// Reads a MessagePack boolean value from memory and sets the value to the
// given memory location.
//
// target - The place where the integer should be written to.
// value  - The memory location where a MessagePack encoded boolean is located.
// sz     - A pointer to where the number of bytes read should be returned.
//
// Returns nothing.
void sky_data_descriptor_set_boolean(void *target, void *value, size_t *sz)
{
    *((bool*)target) = minipack_unpack_bool(value, sz);
}


//--------------------------------------
// Clear Functions
//--------------------------------------

// Clears a string.
//
// target - The location of the string data to reset.
//
// Returns nothing.
void sky_data_descriptor_clear_string(void *target)
{
    sky_string *string = (sky_string*)target;
    string->length = 0;
    string->data = NULL;
}

// Clears an integer.
//
// target - The location of the int data to reset.
//
// Returns nothing.
void sky_data_descriptor_clear_int(void *target)
{
    *((int32_t*)target) = 0;
}

// Clears a double.
//
// target - The location of the double data to reset.
//
// Returns nothing.
void sky_data_descriptor_clear_double(void *target)
{
    *((double*)target) = 0;
}

// Clears a boolean.
//
// target - The location of the boolean data to reset.
//
// Returns nothing.
void sky_data_descriptor_clear_boolean(void *target)
{
    *((bool*)target) = false;
}


//==============================================================================
//
// Cursor
//
//==============================================================================

void sky_cursor_set_data(sky_cursor *cursor);
void sky_cursor_clear_data(sky_cursor *cursor);

sky_cursor_set_ptr(sky_cursor *cursor, void *ptr, size_t sz)
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
    sky_cursor_clear_data(cursor);
}

void sky_cursor_next_event(sky_cursor *cursor)
{
    size_t event_length = 0;

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
        if(flag != EVENT_FLAG) badcursordata("eflag");
        ptr += sizeof(sky_event_flag_t);
        
        // Read timestamp.
        size_t sz;
        int64_t ts = minipack_unpack_int(ptr, &sz);
        if(sz == 0) badcursordata("timestamp");
        uint32_t timestamp = (ts >> SECONDS_BIT_OFFSET);

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
            
            // Clear old action data.
            uint32_t i;
            for(i=0; i<cursor->data_descriptor->action_property_descriptor_count; i++) {
                sky_data_property_descriptor *property_descriptor = cursor->data_descriptor->action_property_descriptors[i];
                property_descriptor->clear_func(cursor->data + property_descriptor->offset);
            }

            // Read msgpack map!
            uint32_t count = minipack_unpack_map(ptr, &sz);
            if(sz == 0) badcursordata("datamap");
            ptr += sz;

            // Loop over key/value pairs.
            for(i=0; i<count; i++) {
                // Read property id (key).
                int64_t property_id = minipack_unpack_int(ptr, &sz);
                if(sz == 0) badcursordata("key");
                ptr += sz;

                // Read property value and set it on the data object.
                sky_data_descriptor_set_value(cursor->data_descriptor, cursor->data, property_id, ptr, &sz);
                if(sz == 0) badcursordata("value");
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

void sky_cursor_clear_data(sky_cursor *cursor)
{
    if(cursor->data != NULL && cursor->data_descriptor != NULL && cursor->data_descriptor->data_sz > 0) {
        memset(cursor->data, 0, cursor->data_descriptor->data_sz);
    }
}

*/
import "C"

