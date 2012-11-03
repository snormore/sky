#include <stdio.h>
#include <stdlib.h>
#include <getopt.h>
#include <signal.h>

#include "bstring.h"
#include "dbg.h"
#include "server.h"
#include "version.h"


//==============================================================================
//
// Definitions
//
//==============================================================================

#define SKY_DEFAULT_DATA_PATH "/usr/local/sky/data"

typedef struct {
    bstring path;
    int port;
} skyd_options;



//==============================================================================
//
// Forward Declarations
//
//==============================================================================

int skyd_server_create(bstring path, int port, sky_server **ret);

skyd_options *skyd_options_parse(int argc, char **argv);

void skyd_options_free(skyd_options *options);


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Main
//--------------------------------------

int main(int argc, char **argv)
{
    int rc;

    // Parse command line options.
    skyd_options *options = skyd_options_parse(argc, argv);

    // Create server.
    sky_server *server = NULL;
    rc = skyd_server_create(options->path, options->port, &server);
    check(rc == 0, "Unable to create server");
    
    // Display status.
    printf("Sky Server v%s\n", SKY_VERSION);
    printf("Listening on 0.0.0.0:%d, CTRL+C to stop\n", server->port);
    
    // Signal handlers.

    // Run server.
    sky_server_start(server);
    while(true) {
        sky_server_accept(server);
    }

    // Clean up.
    sky_server_stop(server);
    sky_server_free(server);
    skyd_options_free(options);

    return 0;

error:
    return 1;
}


//--------------------------------------
// Server
//--------------------------------------

// Creates a server and registers message handlers.
//
// path - The path to the data directory.
// port - The port to run the server on.
// ret  - A pointer where the server should be returned to.
//
// Return 0 if successful, otherwise returns -1.
int skyd_server_create(bstring path, int port, sky_server **ret)
{
    int rc;
    sky_server *server = NULL;
    check(blength(path) > 0, "Path required");

    // Create server.
    server = sky_server_create(path);
    check_mem(server);
    
    // Assign a different port if one is provided.
    if(port > 0) {
        server->port = port;
    }
    
    // Register the default message handlers on the server.
    rc = sky_server_add_default_message_handlers(server);
    check(rc == 0, "Unable to register default message handlers");

    // Ignore sigpipe.
    signal(SIGPIPE, SIG_IGN);

    // Return the server reference.
    *ret = server;
    return 0;

error:
    *ret = NULL;
    sky_server_free(server);
    return -1;
}


//--------------------------------------
// Command Line Options
//--------------------------------------

// Parses the command line options.
//
// argc - The number of arguments.
// argv - An array of argument strings.
//
// Returns a pointer to an Options struct.
skyd_options *skyd_options_parse(int argc, char **argv)
{
    skyd_options *options = calloc(1, sizeof(*options));
    check_mem(options);
    
    // Command line options.
    struct option long_options[] = {
        {"port", optional_argument, 0, 'p'},
        {0, 0, 0, 0}
    };

    // Parse command line options.
    while(1) {
        int option_index = 0;
        int c = getopt_long(argc, argv, "p:", long_options, &option_index);
        
        // Check for end of options.
        if(c == -1) {
            break;
        }
        
        // Parse each option.
        switch(c) {
            case 'i': {
                options->port = atoi(optarg);
                break;
            }
        }
    }
    
    argc -= optind;
    argv += optind;

    // Default the path if one is not specified.
    if(argc < 1) {
        options->path = bfromcstr(SKY_DEFAULT_DATA_PATH); check_mem(options->path);
    }
    else {
        options->path = bfromcstr(argv[0]); check_mem(options->path);
    }

    // Default port.
    if(options->port < 0) {
        fprintf(stderr, "Error: Invalid port number.\n\n");
        exit(1);
    }

    return options;
    
error:
    exit(1);
}

// Frees an Options struct from memory.
//
// Returns nothing.
void skyd_options_free(skyd_options *options)
{
    if(options) {
        bdestroy(options->path);
        free(options);
    }
}

