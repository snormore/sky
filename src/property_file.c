#include <stdlib.h>
#include <inttypes.h>
#include <assert.h>

#include "dbg.h"
#include "mem.h"
#include "bstring.h"
#include "file.h"
#include "property_file.h"
#include "minipack.h"

//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

// Creates a reference to an property file.
// 
// Returns a reference to the new property file if successful. Otherwise returns
// null.
sky_property_file *sky_property_file_create()
{
    sky_property_file *property_file = calloc(sizeof(sky_property_file), 1);
    check_mem(property_file);
    return property_file;
    
error:
    sky_property_file_free(property_file);
    return NULL;
}

// Removes an property file reference from memory.
//
// property_file - The property file to free.
void sky_property_file_free(sky_property_file *property_file)
{
    if(property_file) {
        if(property_file->path) bdestroy(property_file->path);
        property_file->path = NULL;
        sky_property_file_unload(property_file);
        free(property_file);
    }
}


//--------------------------------------
// Path
//--------------------------------------

// Retrieves the file path of an property file.
//
// property_file - The property file.
// path          - A pointer to where the path should be returned.
//
// Returns 0 if successful, otherwise returns -1.
int sky_property_file_get_path(sky_property_file *property_file, bstring *path)
{
    assert(property_file != NULL);

    *path = bstrcpy(property_file->path);
    if(property_file->path) check_mem(*path);

    return 0;

error:
    *path = NULL;
    return -1;
}

// Sets the file path of an property file.
//
// property_file - The property file.
// path          - The file path to set.
//
// Returns 0 if successful, otherwise returns -1.
int sky_property_file_set_path(sky_property_file *property_file, bstring path)
{
    assert(property_file != NULL);

    if(property_file->path) bdestroy(property_file->path);
    
    property_file->path = bstrcpy(path);
    if(path) check_mem(property_file->path);

    return 0;

error:
    property_file->path = NULL;
    return -1;
}


//--------------------------------------
// Persistence
//--------------------------------------

// Loads properties from file.
//
// property_file - The property file to load.
//
// Returns 0 if successful, otherwise returns -1.
int sky_property_file_load(sky_property_file *property_file)
{
    int rc;
    size_t sz;
    uint32_t count = 0;
    FILE *file = NULL;
    sky_property **properties = NULL;
    assert(property_file != NULL);
    assert(property_file->path != NULL);

    // Unload any properties currently in memory.
    rc = sky_property_file_unload(property_file);
    check(rc == 0, "Unable to unload property file");

    // Read in properties file if it exists.
    if(sky_file_exists(property_file->path)) {
        file = fopen(bdata(property_file->path), "r");
        check(file, "Failed to open property file: %s",  bdata(property_file->path));

        // Read property count.
        count = minipack_fread_array(file, &sz);
        check(sz != 0, "Unable to read properties array at byte: %ld", ftell(file));

        // Allocate properties.
        if(count > 0) {
            properties = malloc(sizeof(sky_property*) * count);
            check_mem(properties);
        }

        // Read properties.
        uint32_t i;
        for(i=0; i<count; i++) {
            sky_property *property = sky_property_create(); check_mem(property);
            
            rc = sky_property_unpack(property, file);
            check(rc == 0, "Unable to unpack property");

            // Append to array.
            property->property_file = property_file;
            properties[i] = property;
        }

        // Close the file.
        fclose(file);
    }

    // Store property list on property file.
    property_file->properties = properties;
    property_file->property_count = count;

    return 0;

error:
    if(file) fclose(file);
    if(properties) free(properties);
    return -1;
}

// Saves properties to file.
//
// property_file - The property file to save.
//
// Returns 0 if successful, otherwise returns -1.
int sky_property_file_save(sky_property_file *property_file)
{
    int rc;
    size_t sz;
    FILE *file = NULL;
    assert(property_file != NULL);
    assert(property_file->path != NULL);

    // Open file.
    file = fopen(bdata(property_file->path), "w");
    check(file, "Failed to open property file: %s", bdata(property_file->path));

    // Write property array.
    rc = minipack_fwrite_array(file, property_file->property_count, &sz);
    check(rc == 0, "Unable to write properties array");

    // Write properties.
    uint32_t i;
    for(i=0; i<property_file->property_count; i++) {
        sky_property *property = property_file->properties[i];

        rc = sky_property_pack(property, file);
        check(rc == 0, "Unable to pack property");
    }

    // Close the file.
    fclose(file);

    return 0;

error:
    if(file) fclose(file);
    return -1;
}

// Unloads the properties in the property file from memory.
//
// property_file - The property file to save.
//
// Returns 0 if successful, otherwise returns -1.
int sky_property_file_unload(sky_property_file *property_file)
{
    if(property_file) {
        // Release properties.
        if(property_file->properties) {
            uint32_t i=0;
            for(i=0; i<property_file->property_count; i++) {
                sky_property *property = property_file->properties[i];
                sky_property_free(property);
                property_file->properties[i] = NULL;
            }
            free(property_file->properties);
            property_file->properties = NULL;
        }
        
        property_file->property_count = 0;
    }
    
    return 0;
}


//--------------------------------------
// Property Management
//--------------------------------------

// Retrieves an property by id.
//
// property_file - The property file to search.
// property_id   - The id of the property.
// ret           - A pointer to where the property should be returned to.
//
// Returns 0 if successful, otherwise returns -1.
int sky_property_file_find_by_id(sky_property_file *property_file,
                                 sky_property_id_t property_id,
                                 sky_property **ret)
{
    assert(property_file != NULL);
    
    // Initialize return values.
    *ret = NULL;
    
    // Loop over properties to find matching name.
    uint32_t i;
    for(i=0; i<property_file->property_count; i++) {
        if(property_file->properties[i]->id == property_id) {
            *ret = property_file->properties[i];
            break;
        }
    }
    
    return 0;
}

// Retrieves a property with a given name.
//
// property_file - The property file to search.
// name          - The name of the property.
// ret           - A pointer to where the property should be returned to.
//
// Returns 0 if successful, otherwise returns -1.
int sky_property_file_find_by_name(sky_property_file *property_file,
                                   bstring name, sky_property **ret)
{
    assert(property_file != NULL);
    assert(name != NULL);
    
    // Initialize return value.
    *ret = NULL;
    
    // Loop over properties to find matching name.
    uint32_t i;
    for(i=0; i<property_file->property_count; i++) {
        if(biseq(property_file->properties[i]->name, name) == 1) {
            *ret = property_file->properties[i];
            break;
        }
    }
    
    return 0;
}


// Adds an property to an property file.
//
// property_file - The property file.
// property      - The property to add to the property file.
//
// Returns 0 if successful, otherwise returns -1.
int sky_property_file_add_property(sky_property_file *property_file,
                                   sky_property *property)
{
    uint32_t i;
    assert(property_file != NULL);
    assert(property != NULL);
    check(property->id == 0, "Property ID must be zero");
    check(property->property_file == NULL, "Property must not be attached to a property file");
    
    // Make sure an property with that name doesn't already exist.
    sky_property *_property;
    int rc = sky_property_file_find_by_name(property_file, property->name, &_property);
    check(rc == 0 && _property == NULL, "Property already exists with the same name");
    
    // If this is an object property then find the next positive id.
    if(property->type == SKY_PROPERTY_TYPE_OBJECT) {
        int64_t max_id = 0;
        for(i=0; i<property_file->property_count; i++) {
            if(property_file->properties[i]->id > max_id) {
                max_id = property_file->properties[i]->id;
            }
        }
        check(max_id < INT8_MAX, "No additional object property ids available");
        property->id = max_id + 1;
    }
    // If this is an action property then find the next negative id.
    else if(property->type == SKY_PROPERTY_TYPE_ACTION) {
        int64_t min_id = 0;
        for(i=0; i<property_file->property_count; i++) {
            if(property_file->properties[i]->id < min_id) {
                min_id = property_file->properties[i]->id;
            }
        }
        check(min_id > INT8_MIN, "No additional action property ids available");
        property->id = min_id - 1;
    }
    else {
        sentinel("Invalid property type: %d", property->type);
    }

    // Link to property file.
    property->property_file = property_file;

    // Append property to list.
    property_file->property_count++;
    property_file->properties = realloc(property_file->properties, sizeof(sky_property*) * property_file->property_count);
    check_mem(property_file->properties);
    property_file->properties[property_file->property_count-1] = property;
    
    return 0;

error:
    return -1;
}


//--------------------------------------
// Data Descriptors
//--------------------------------------

// Initializes a data descriptor based on a property file. This function
// automatically sets up the minimum and maximum property identifiers.
//
// property_file - The property file.
// ret           - A pointer to where the descriptor should be returned.
//
// Returns 0 if successful, otherwise returns -1.
int sky_property_file_create_data_descriptor(sky_property_file *property_file,
                                             sky_data_descriptor **ret)
{
    uint32_t i;
    assert(property_file != NULL);
    assert(ret != NULL);

    // Determine minimum and maximum property identifiers.
    sky_property_id_t min_property_id = 0;
    sky_property_id_t max_property_id = 0;
    for(i=0; i<property_file->property_count; i++) {
        sky_property *property = property_file->properties[i];
        
        if(property->id < min_property_id) {
            min_property_id = property->id;
        }
        if(property->id > max_property_id) {
            max_property_id = property->id;
        }
    }
    
    // Initialize the descriptor.
    sky_data_descriptor *descriptor = sky_data_descriptor_create(min_property_id, max_property_id);
    check_mem(descriptor);
    
    // Return descriptor.
    *ret = descriptor;
    
    return 0;

error:
    sky_data_descriptor_free(descriptor);
    *ret = NULL;
    return -1;
}

