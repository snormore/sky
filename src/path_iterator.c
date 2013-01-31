#include <stdlib.h>
#include <assert.h>

#include "path_iterator.h"
#include "mem.h"
#include "dbg.h"


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

// Creates a reference to a path iterator.
// 
// Returns a reference to the new path iterator if successful. Otherwise returns
// null.
sky_path_iterator *sky_path_iterator_create()
{
    sky_path_iterator *iterator = sky_path_iterator_alloc(); check_mem(iterator);
    sky_path_iterator_init(iterator);
    return iterator;
    
error:
    sky_path_iterator_free(iterator);
    return NULL;
}

// Allocates memory for a path iterator.
// 
// Returns a reference to the new path iterator if successful. Otherwise returns
// null.
sky_path_iterator *sky_path_iterator_alloc()
{
    return malloc(sizeof(sky_path_iterator));
}

// Initializes a path iterator.
// 
// iterator - The iterator to initialize.
void sky_path_iterator_init(sky_path_iterator *iterator)
{
    memset(iterator, 0, sizeof(sky_path_iterator));
}

// Uninitializes a path iterator.
// 
// iterator - The iterator.
void sky_path_iterator_uninit(sky_path_iterator *iterator)
{
    iterator->tablet = NULL;
    if(iterator->leveldb_iterator) leveldb_iter_destroy(iterator->leveldb_iterator);
    iterator->leveldb_iterator = NULL;
}

// Removes a path iterator reference from memory.
//
// iterator - The path iterator to free.
void sky_path_iterator_free(sky_path_iterator *iterator)
{
    if(iterator) {
        sky_path_iterator_uninit(iterator);
        free(iterator);
    }
}


//--------------------------------------
// Source
//--------------------------------------

// Assigns a tablet as the source.
// 
// iterator - The iterator.
// tablet   - The data file to iterate over.
//
// Returns 0 if successful, otherwise returns -1.
int sky_path_iterator_set_tablet(sky_path_iterator *iterator, sky_tablet *tablet)
{
    int rc;
    assert(iterator != NULL);
    iterator->tablet = tablet;
    iterator->eof = false;

    // Initialize LevelDB iterator.
    iterator->leveldb_iterator = leveldb_create_iterator(tablet->leveldb_db, tablet->readoptions);
    check(iterator->leveldb_iterator != NULL, "Unable to create LevelDB iterator");
    leveldb_iter_seek_to_first(iterator->leveldb_iterator);

    // Move cursor to initial path.
    rc = sky_path_iterator_next(iterator);
    check(rc == 0, "Unable to move to initial path");

    return 0;
    
error:
    return -1;
}


//--------------------------------------
// Iteration
//--------------------------------------

// Moves the iterator to point to the next path.
// 
// iterator - The iterator.
//
// Returns 0 if successful, otherwise returns -1.
int sky_path_iterator_next(sky_path_iterator *iterator)
{
    int rc;
    assert(iterator != NULL);
    assert(iterator->tablet != NULL);
    assert(iterator->leveldb_iterator != NULL);

    // Move to next path.
    if(leveldb_iter_valid(iterator->leveldb_iterator)) {
        // Retrieve the path data for this object.
        size_t data_length;
        void *data = (void*)leveldb_iter_value(iterator->leveldb_iterator, &data_length);
        
        // Set the pointer on the cursor.
        rc = sky_cursor_set_ptr(&iterator->cursor, data, data_length);
        check(rc == 0, "Unable to set cursor pointer");

        // Move to the next object.
        leveldb_iter_next(iterator->leveldb_iterator);
        
        if(leveldb_iter_valid(iterator->leveldb_iterator)) {
            iterator->eof = false;
        }
    }
    else {
        iterator->eof = true;
    }

    return 0;
    
error:
    return -1;
}

// Checks if the path iterator is at the end.
// 
// iterator - The iterator.
//
// Returns true if at the end, otherwise returns false.
bool sky_path_iterator_eof(sky_path_iterator *iterator)
{
    assert(iterator != NULL);
    return iterator->eof;
}


