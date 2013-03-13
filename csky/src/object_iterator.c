#include <stdlib.h>
#include "sky/object_iterator.h"

//==============================================================================
//
// Functions
//
//==============================================================================

void sky_object_iterator_next(sky_object_iterator *iterator);

sky_object_iterator *sky_object_iterator_create(leveldb_t *db)
{
    sky_object_iterator *iterator = calloc(1, sizeof(sky_object_iterator));

    // Initialize LevelDB iterator.
    leveldb_readoptions_t* readoptions = leveldb_readoptions_create();
    iterator->leveldb_iterator = leveldb_create_iterator(db, readoptions);
    leveldb_readoptions_destroy(readoptions);
    leveldb_iter_seek_to_first(iterator->leveldb_iterator);

    // Move cursor to initial object.
    sky_object_iterator_next(iterator);

    return iterator;
}

void sky_object_iterator_free(sky_object_iterator *iterator)
{
    if(iterator->leveldb_iterator) leveldb_iter_destroy(iterator->leveldb_iterator);
    iterator->leveldb_iterator = NULL;
    free(iterator);
}

void sky_object_iterator_next(sky_object_iterator *iterator)
{
    // Move to next object.
    if(leveldb_iter_valid(iterator->leveldb_iterator)) {
        // Retrieve the object data for this object.
        size_t data_length;
        void *data = (void*)leveldb_iter_value(iterator->leveldb_iterator, &data_length);
        
        // Set the pointer on the cursor.
        sky_cursor_set_ptr(&iterator->cursor, data, data_length);

        // Move to the next object.
        leveldb_iter_next(iterator->leveldb_iterator);
        if(leveldb_iter_valid(iterator->leveldb_iterator)) {
            iterator->eof = false;
        }
    }
    else {
        iterator->eof = true;
    }
}

bool sky_object_iterator_eof(sky_object_iterator *iterator)
{
    return iterator->eof;
}
