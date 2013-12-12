package hashmap

/*
#include <stdio.h>
#include <stdlib.h>
#include <inttypes.h>

#define HASHMAP_BUCKET_COUNT 256

typedef enum {
    int_value_type,
    hashmap_value_type,
} sky_hashmap_elem_type_e;

typedef struct sky_hashmap sky_hashmap;

typedef struct {
  int64_t key;
  sky_hashmap_elem_type_e value_type;
  union {
    int64_t int_value;
    sky_hashmap *hashmap_value;
  } value;
} sky_hashmap_elem;

typedef struct {
    sky_hashmap_elem *elements;
    int64_t count;
} sky_hashmap_bucket;

struct sky_hashmap {
    sky_hashmap_bucket buckets[HASHMAP_BUCKET_COUNT];
};


// Creates a new hashmap instance.
sky_hashmap *sky_hashmap_new()
{
    sky_hashmap *hashmap = calloc(1, sizeof(sky_hashmap));
    return hashmap;
}

// Frees a hashmap instance.
void sky_hashmap_free(sky_hashmap *hashmap)
{
    // TODO
}

// Retrieves the bucket and element in the hashmap for a key.
inline void sky_hashmap_find(sky_hashmap *hashmap, int64_t key, sky_hashmap_bucket **bucket, sky_hashmap_elem **element)
{
    int64_t i;
    int64_t hash = key % HASHMAP_BUCKET_COUNT;
    *bucket = &hashmap->buckets[hash];
    for (i=0; i<(*bucket)->count; i++) {
        if((*bucket)->elements[i].key == key) {
            *element = &((*bucket)->elements[i]);
            return;
        }
    }
    *element = NULL;
}

// Retrieves the bucket and element in the hashmap for a key.
// The element is created if it's not found.
inline void sky_hashmap_find_or_create(sky_hashmap *hashmap, int64_t key, sky_hashmap_elem_type_e value_type, sky_hashmap_elem **elem)
{
    sky_hashmap_bucket *bucket;
    sky_hashmap_find(hashmap, key, &bucket, elem);

    // Create new element if it didn't exist.
    if((*elem) == NULL) {
        bucket->count++;
        bucket->elements = realloc(bucket->elements, (sizeof(*bucket->elements) * bucket->count));
        *elem = &bucket->elements[bucket->count-1];
        (*elem)->key = key;
        (*elem)->value_type = value_type;
        (*elem)->value.int_value = 0;
    } else if((*elem)->value_type != value_type) {
        // If the existing type is different that the new one then override it.
        if((*elem)->value_type == hashmap_value_type) {
            sky_hashmap_free((*elem)->value.hashmap_value);
        }
        (*elem)->value.hashmap_value = NULL;
        (*elem)->value_type = value_type;
    }
}

// Retrieves an int value for a given key.
// Creates the key if it doesn't already exist.
int64_t sky_hashmap_get(sky_hashmap *hashmap, int64_t key)
{
    sky_hashmap_elem *elem;
    sky_hashmap_find_or_create(hashmap, key, int_value_type, &elem);
    return elem->value.int_value;
}

// Sets an int value for a given key.
void sky_hashmap_set(sky_hashmap *hashmap, int64_t key, int64_t value)
{
    sky_hashmap_elem *elem;
    sky_hashmap_find_or_create(hashmap, key, int_value_type, &elem);
    elem->value.int_value = value;
}

// Retrieves a hashmap value for a given key.
// Creates the key if it doesn't already exist.
sky_hashmap *sky_hashmap_submap(sky_hashmap *hashmap, int64_t key)
{
    sky_hashmap_elem *elem;
    sky_hashmap_find_or_create(hashmap, key, hashmap_value_type, &elem);
    if(elem->value.hashmap_value == NULL) {
        elem->value.hashmap_value = sky_hashmap_new();
    }
    return elem->value.hashmap_value;
}


// FOR BENCHMARKING ONLY.
sky_hashmap *sky_hashmap_benchmark(sky_hashmap *hashmap, int64_t n)
{
    int64_t i;
    for (i=0; i<n; i++) {
        sky_hashmap_set(hashmap, i % 256, 100);
    }
}
*/
import "C"
import "unsafe"

// Hashmap wraps the underlying C struct.
type Hashmap struct {
    C *C.sky_hashmap
}

// New creates a new Hashmap instance.
func New() *Hashmap {
    return &Hashmap{C.sky_hashmap_new()}
}

// Free releases the underlying C memory.
func (h *Hashmap) Free() {
    if h.C != nil {
        C.free(unsafe.Pointer(h.C))
        h.C = nil
    }
}

// Get retrieves the int value for a given key.
func (h *Hashmap) Get(key int64) int64 {
    return int64(C.sky_hashmap_get(h.C, C.int64_t(key)))
}

// Set sets an int value for a given key.
func (h *Hashmap) Set(key int64, value int64) {
    C.sky_hashmap_set(h.C, C.int64_t(key), C.int64_t(value))
}

// Submap retrieves the hashmap value for a given key.
func (h *Hashmap) Submap(key int64) *Hashmap {
    return &Hashmap{C.sky_hashmap_submap(h.C, C.int64_t(key))}
}

// benchmark runs set N number of times within the C context.
func benchmark(h *Hashmap, n int64) {
    C.sky_hashmap_benchmark(h.C, C.int64_t(n))
}

