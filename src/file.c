#include <stdio.h>
#include <sys/stat.h>
#include <dirent.h>

#include "dbg.h"
#include "file.h"

//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// File Stats
//--------------------------------------

// Checks if a file exists.
//
// path - The path of the file.
//
// Returns true if the file at the given path exists. Otherwise returns false.
bool sky_file_exists(bstring path)
{
    struct stat buffer;
    int rc = stat(bdatae(path, ""), &buffer);
    return (rc == 0);
}

// Checks if a file is a directory.
//
// path - The path of the file.
//
// Returns true if the file at the given path is a directory. Otherwise
// returns false.
bool sky_file_is_dir(bstring path)
{
    struct stat buffer;
    int rc = stat(bdatae(path, ""), &buffer);
    if(rc == 0) {
        return S_ISDIR(buffer.st_mode);
    }
    else {
        return false;
    }
}

// Retrieves the size of a file, in bytes.
//
// path - The path of the file.
//
// Returns the size of the file at the given path in bytes.
off_t sky_file_get_size(bstring path)
{
    struct stat buffer;
    if(stat(bdatae(path, ""), &buffer) == 0) {
        return buffer.st_size;
    }
    else {
        return 0;
    }
}


//--------------------------------------
// File Copy
//--------------------------------------

// Copies a single file from the source path to the destination path.
//
// src  - The path of the file to copy.
// dest - The path where the copy should be placed.
//
// Returns 0 if successful, otherwise returns -1.
int sky_file_cp(bstring src, bstring dest)
{
    int rc;
    check(src != NULL, "Source path required");
    check(dest != NULL, "Destination path required");
    check(sky_file_exists(src), "Source file does not exist");
    check(!sky_file_is_dir(src), "Source file cannot be a directory");

    // Open source files.
    FILE *src_file = fopen(bdatae(src, ""), "r");
    check(src_file != NULL, "Unable to open source file for reading");

    // Open destination file.
    FILE *dest_file = fopen(bdatae(dest, ""), "w");
    check(dest_file != NULL, "Unable to open destination file for writing");
    
    // Read from source and write to destination until done.
    while(!feof(src_file)) {
        int buffer_size = 1024;
        char *buffer[buffer_size];
        
        size_t sz = fread(buffer, sizeof(char), buffer_size, src_file);
        fwrite(buffer, sizeof(char), sz, dest_file);
    }
    
    // Close files.
    fclose(src_file);
    fclose(dest_file);
    
    // Copy permissions to destination file.
    struct stat st;
    stat(bdatae(src, ""), &st);
    rc = chmod(bdatae(dest, ""), st.st_mode);
    check(rc == 0, "Unable to copy permissions");
    
    return 0;

error:
    return -1;
}

// Recursively copies a directory.
//
// src  - The path of the file or directory to copy.
// dest - The path where the copy should be placed.
//
// Returns 0 if successful, otherwise returns -1.
int sky_file_cp_r(bstring src, bstring dest)
{
    int rc;
    bstring ent_src = NULL;
    bstring ent_dest = NULL;
    check(src != NULL, "Source path required");
    check(dest != NULL, "Destination path required");
    check(sky_file_exists(src), "Source file does not exist");

    // If this is a directory then create a new dest directory and copy the
    // contents.
    if(sky_file_is_dir(src)) {
        // Create destination directory if it doesn't exist.
        if(!sky_file_exists(dest)) {
            struct stat st;
            rc = stat(bdatae(src, ""), &st);
            check(rc == 0, "Unable to stat source directory: %s", bdatae(src, ""));
            rc = mkdir(bdatae(dest, ""), st.st_mode);
            check(rc == 0, "Unable to create directory: %s", bdatae(dest, ""));
        }
        
        // Open directory.
        DIR *dir = opendir(bdatae(src, ""));
        check(dir != NULL, "Unable to open directory: %s", bdatae(src, ""));
        
        // Copy over contents of directory.
        struct dirent *ent;
        while((ent = readdir(dir))) {
            if(strcmp(ent->d_name, ".") != 0 && strcmp(ent->d_name, "..") != 0) {
                ent_src  = bformat("%s/%s", bdatae(src, ""), ent->d_name); check_mem(ent_src);
                ent_dest = bformat("%s/%s", bdatae(dest, ""), ent->d_name); check_mem(ent_dest);

                rc = sky_file_cp_r(ent_src, ent_dest);
                check(rc == 0, "Unable to copy: %s", bdatae(ent_src, ""));

                bdestroy(ent_src);
                bdestroy(ent_dest);
            }
        }
        
        // Close directory.
        closedir(dir);
    }
    // If this is a file then copy its contents.
    else {
        rc = sky_file_cp(src, dest);
        check(rc == 0, "Unable to copy file: %s", bdatae(src, ""));
    }
    
    return 0;

error:
    bdestroy(ent_src);
    bdestroy(ent_dest);
    return -1;
}


//--------------------------------------
// File Delete
//--------------------------------------

// Deletes a single file.
//
// path - The path of the file to delete.
//
// Returns 0 if successful, otherwise returns -1.
int sky_file_rm(bstring path)
{
    int rc;
    check(path != NULL, "Path required");
    check(!sky_file_is_dir(path), "File cannot be a directory");

    if(sky_file_exists(path)) {
        rc = remove(bdatae(path, ""));
        check(rc == 0, "Unable to delete file");
    }

    return 0;

error:
    return -1;
}

// Recursively deletes a file or directory.
//
// path  - The path of the file or directory to delete.
//
// Returns 0 if successful, otherwise returns -1.
int sky_file_rm_r(bstring path)
{
    int rc;
    bstring ent_path = NULL;
    check(path != NULL, "Path required");
    
    // If path doesn't exist then just ignore it.
    if(sky_file_exists(path)) {
        // If the file is a directory then delete its contents first and then
        // delete it.
        if(sky_file_is_dir(path)) {
            // Open directory.
            DIR *dir = opendir(bdatae(path, ""));
            check(dir != NULL, "Unable to open directory: %s", bdatae(path, ""));

            // Remove each file inside directory.
            struct dirent *ent;
            while((ent = readdir(dir))) {
                if(strcmp(ent->d_name, ".") != 0 && strcmp(ent->d_name, "..") != 0) {
                    ent_path  = bformat("%s/%s", bdatae(path, ""), ent->d_name); check_mem(ent_path);
                    rc = sky_file_rm_r(ent_path);
                    check(rc == 0, "Unable to delete: %s", bdatae(ent_path, ""));
                    bdestroy(ent_path);
                }
            }

            // Close directory.
            closedir(dir);

            // Remove directory.
            remove(bdatae(path, ""));
        }
        // If this is a file then delete it.
        else {
            rc = sky_file_rm(path);
            check(rc == 0, "Unable to delete file: %s", bdatae(path, ""));
        }
    }
    
    return 0;

error:
    bdestroy(ent_path);
    return -1;
}

