package mapper

/*
#cgo LDFLAGS: -llmdb
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <inttypes.h>
#include <string.h>
#include <lmdb.h>

// The number of microseconds per second.
#define USEC_PER_SEC        1000000

// A bit-mask to extract the microseconds from a Sky timestamp.
#define USEC_MASK           0xFFFFF

// The number of bits that seconds are shifted over in a timestamp.
#define SECONDS_BIT_OFFSET  20

#define badcursordata(MSG, PTR) do {\
    fprintf(stderr, "Cursor pointing at invalid raw event data [" MSG "]: %p\n", PTR); \
    cursor->next_event->timestamp = 0; \
    return false; \
} while(0)

#define debug(M, ...) fprintf(stderr, "DEBUG %s:%d: " M "\n", __FILE__, __LINE__, ##__VA_ARGS__)

typedef struct {
  char *data;
  int64_t sz;
} sky_string;

typedef struct sky_cursor sky_cursor;

typedef struct {
  bool eos;
  bool eof;
  int64_t timestamp;
} sky_event;

struct sky_cursor {
    sky_event *event;
    sky_event *next_event;

    bool blocking_eos;
    uint32_t max_timestamp;
    uint32_t session_idle_in_sec;

    void *key_prefix;
    int64_t key_prefix_sz;
    MDB_cursor* lmdb_cursor;
};


bool sky_cursor_next_event(sky_cursor *cursor);

void sky_cursor_update_eos(sky_cursor *cursor);


// Creates a new cursor instance.
sky_cursor *sky_cursor_new(MDB_cursor *lmdb_cursor, void *key_prefix, int64_t key_prefix_sz)
{
    sky_cursor *cursor = calloc(1, sizeof(sky_cursor));
    if(cursor == NULL) debug("[malloc] Unable to allocate cursor.");

    cursor->lmdb_cursor = lmdb_cursor;
    cursor->key_prefix = key_prefix;
    cursor->key_prefix_sz = key_prefix_sz;

    return cursor;
}

// Removes a cursor from memory.
void sky_cursor_free(sky_cursor *cursor)
{
    if(cursor) {
        if(cursor->key_prefix != NULL) free(cursor->key_prefix);
        cursor->key_prefix = NULL;
        cursor->key_prefix_sz = 0;
        free(cursor);
    }
}

// Sets up object after cursor has already been positioned.
bool sky_cursor_iter_object(sky_cursor *cursor, MDB_val *key, MDB_val *data)
{
    if(cursor->key_prefix_sz > 0 && (key->mv_size < cursor->key_prefix_sz || memcmp(cursor->key_prefix, key->mv_data, cursor->key_prefix_sz) != 0)) {
        return false;
    }
    // fprintf(stderr, "\nOBJ (%.*s) [%d]\n", (int)key->mv_size, (char*)key->mv_data, (int)key->mv_size);

    // Clear the data object if set.
    cursor->session_idle_in_sec = 0;
    cursor->blocking_eos = false;
    cursor->next_event->eof = false;
    // TODO: memset(cursor->event, 0, cursor->event_sz);

    // Read first event into "next" event.
    // TODO: if(!sky_cursor_read(cursor, cursor->next_event, data->mv_data)) {
    // TODO:     return false;
    // TODO: }

    // Move "next" event to current event and put the next event in buffer.
    return sky_cursor_next_event(cursor);
}

// Moves the cursor to point to the first object. If a prefix is set then
// move to the first object that with the given prefix.
bool sky_cursor_init(sky_cursor *cursor)
{
    int rc;
    MDB_val key, data;

    if(cursor->key_prefix_sz == 0) {
        if((rc = mdb_cursor_get(cursor->lmdb_cursor, &key, &data, MDB_FIRST)) != 0) {
            if(rc != MDB_NOTFOUND) fprintf(stderr, "MDB_FIRST error: %s\n", mdb_strerror(rc));
            return false;
        }

    } else {
        key.mv_data = cursor->key_prefix;
        key.mv_size = cursor->key_prefix_sz;
        if((rc = mdb_cursor_get(cursor->lmdb_cursor, &key, &data, MDB_SET_RANGE)) != 0) {
            if(rc != MDB_NOTFOUND) fprintf(stderr, "MDB_SET_RANGE error: %s\n", mdb_strerror(rc));
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
    if(cursor->event->eof || (cursor->event->eos && cursor->blocking_eos)) {
        return false;
    }
    cursor->blocking_eos = true;

    // Read the next event.
    if(!cursor->next_event->eof) {
        MDB_val key, data;
        int rc = mdb_cursor_get(cursor->lmdb_cursor, &key, &data, MDB_NEXT_DUP);
        if(rc != 0) {
            // Clear next event if there isn't one.
            // TODO: memset(cursor->next_event, 0, cursor->event_sz);
            cursor->next_event->eof = true;

            if(rc != MDB_NOTFOUND) {
                printf("lmdb cursor error: %s\n", mdb_strerror(rc));
            }
        } else {
            cursor->next_event->eof = false;
            // TODO: if(!sky_cursor_read(cursor, cursor->next_event, data.mv_data)) {
            // TODO:     return true;
            // TODO: }
        }
    }

    // Update eos.
    sky_cursor_update_eos(cursor);

    return true;
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
    cursor->blocking_eos = false;
}

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
    "unsafe"
    "github.com/szferi/gomdb"
)

func sky_cursor_new(lmdb_cursor *mdb.Cursor, prefix string) *C.sky_cursor {
    return C.sky_cursor_new(lmdb_cursor.MdbCursor(), unsafe.Pointer(C.CString(prefix)), C.int64_t(len(prefix)))
}

func sky_cursor_free(cursor *C.sky_cursor) {
    C.sky_cursor_free(cursor)
}
