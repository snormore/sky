#ifndef _tests_functional_server_helpers_h
#define _tests_functional_server_helpers_h

#include <sys/socket.h>
#include <sys/types.h>
#include <netinet/in.h>
#include <netdb.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <errno.h>
#include <arpa/inet.h> 
#include <pthread.h>

#include <sky.h>
#include <dbg.h>
#include "../minunit.h"


//==============================================================================
//
// Definitions
//
//==============================================================================

typedef struct {
    sky_server *server;
    uint32_t message_count;
} server_options;

void *run_server(void *_options);


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Server
//--------------------------------------

sky_server *start_server(uint32_t message_count, pthread_t *thread)
{
    struct tagbstring ROOT = bsStatic(".");
    sky_server *server = sky_server_create(&ROOT);
    server->port = TEST_PORT;
    sky_server_add_default_message_handlers(server);
    sky_server_start(server);

    signal(SIGPIPE, SIG_IGN);

    server_options *options = calloc(1, sizeof(*options));
    options->server = server;
    options->message_count = message_count;
    pthread_create(thread, NULL, run_server, (void*)options);
    return NULL;
}

void *run_server(void *_options)
{
    server_options *options = (server_options*)_options;
    uint32_t i;
    for(i=0; i<options->message_count; i++) {
        sky_server_accept(options->server);
    }
    // HACK:
    usleep(1000);
    sky_server_stop(options->server);
    sky_server_free(options->server);
    free(_options);
    return NULL;
}


//--------------------------------------
// Messages
//--------------------------------------

int _send_msg(char *path)
{
    int sz=0, rc;
    int input_len=0, output_len=0;
    char buffer[1024];
    char input_msg[1024];
    char output_msg[1024];
    struct sockaddr_in addr;

    // Read input message.
    FILE *input_file = fopen(path, "r");
    check(input_file != NULL, "Unable to open input message path");
    input_len = fread(input_msg, sizeof(char), sizeof(input_msg), input_file);
    check(input_len > 0, "Unable to read input message");
    fclose(input_file);

    // Setup connection info.
    memset(buffer, 0, sizeof(buffer));
    memset(&addr, '0', sizeof(addr));
    int sock = socket(AF_INET, SOCK_STREAM, 0);
    addr.sin_family = AF_INET;
    addr.sin_port = htons(TEST_PORT); 
    inet_pton(AF_INET, "127.0.0.1", &addr.sin_addr);

    // Connect to server.
    rc = connect(sock, (struct sockaddr*)&addr, sizeof(addr));
    check(rc == 0, "Unable to connect to server");

    // Send message.
    sz = write(sock, input_msg, input_len);
    check(sz > 0, "Unable to send input message");

    // Read output message.
    output_len = read(sock, output_msg, sizeof(output_msg));
    check(sz > 0, "Unable to recv output message");
    
    // Write output to file.
    FILE *output_file = fopen("tmp/output", "w");
    check(output_file != NULL, "Unable to open output message path: tmp/output");
    sz = fwrite(output_msg, sizeof(char), output_len, output_file);
    check(sz == output_len, "Unable to write output message bytes: written:%d, exp:%d", sz, output_len);
    fclose(output_file);
    
    return 0;

error:
    return -1;
}

// Sends the contents of a given path to the server.
#define send_msg(PATH) do {\
    int rc = _send_msg(PATH); \
    mu_assert_int_equals(rc, 0); \
} while(0)

// Asserts the contents of an output message.
#define mu_assert_msg(EXP_FILENAME) mu_assert_file("tmp/output", EXP_FILENAME)

#endif