#include <inttypes.h>
#include <stdbool.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/time.h>
#include <zmq.h>

#include "bstring.h"
#include "server.h"
#include "message_header.h"
#include "add_event_message.h"
#include "next_actions_message.h"
#include "add_action_message.h"
#include "get_action_message.h"
#include "get_actions_message.h"
#include "add_property_message.h"
#include "get_property_message.h"
#include "get_properties_message.h"
#include "multi_message.h"
#include "sky_zmq.h"
#include "dbg.h"


//==============================================================================
//
// Forward Declarations
//
//==============================================================================

int sky_server_get_table(sky_server *server, bstring name, sky_table **ret);

int sky_server_open_table(sky_server *server, bstring name, bstring path, sky_table **ret);

int sky_server_create_servlets(sky_server *server, sky_table *table);

int sky_server_stop_servlets(sky_server *server);


//==============================================================================
//
// Global Variables
//
//==============================================================================

// A counter to track the next available server id. Creating multiple servers
// is not thread safe, however, there's probably not a good reason to span
// multiple servers across multiple threads in the same process.
int32_t next_server_id = 0;


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

// Creates a reference to a server instance.
//
// path - The data directory path.
//
// Returns a reference to the server.
sky_server *sky_server_create(bstring path)
{
    sky_server *server = NULL;
    server = calloc(1, sizeof(sky_server)); check_mem(server);
    server->id = next_server_id++;
    server->path = bstrcpy(path);
    if(path) check_mem(server->path);
    server->port = SKY_DEFAULT_PORT;
    server->context = zmq_ctx_new();
    
    return server;

error:
    sky_server_free(server);
    return NULL;
}

// Frees the servlets on a server instance.
//
// server - The server.
//
// Returns nothing.
void sky_server_free_servlets(sky_server *server)
{
    if(server) {
        uint32_t i;
        for(i=0; i<server->servlet_count; i++) {
            sky_servlet_free(server->servlets[i]);
            server->servlets[i] = NULL;
        }
        free(server->servlets);
        server->servlets = NULL;
        server->servlet_count = 0;
    }
}

// Frees the tables on a server instance.
//
// server - The server.
//
// Returns nothing.
void sky_server_free_tables(sky_server *server)
{
    if(server) {
        uint32_t i;
        for(i=0; i<server->table_count; i++) {
            sky_table_free(server->tables[i]);
            server->tables[i] = NULL;
        }
        free(server->tables);
        server->tables = NULL;
        server->table_count = 0;
    }
}

// Frees the message handlers on a server instance.
//
// server - The server.
//
// Returns nothing.
void sky_server_free_message_handlers(sky_server *server)
{
    if(server) {
        uint32_t i;
        for(i=0; i<server->message_handler_count; i++) {
            sky_message_handler_free(server->message_handlers[i]);
            server->message_handlers[i] = NULL;
        }
        free(server->message_handlers);
        server->message_handlers = NULL;
        server->message_handler_count = 0;
    }
}

// Frees a server instance from memory.
//
// server - The server object to free.
//
// Returns nothing.
void sky_server_free(sky_server *server)
{
    if(server) {
        if(server->path) bdestroy(server->path);
        sky_server_free_tables(server);
        sky_server_free_servlets(server);
        sky_server_free_message_handlers(server);
        zmq_ctx_destroy(server->context);

        free(server);
    }
}



//--------------------------------------
// State
//--------------------------------------

// Starts a server. Once a server is started, it can accept messages over TCP
// on the bind address and port number specified by the server object.
//
// server - The server to start.
//
// Returns 0 if successful, otherwise returns -1.
int sky_server_start(sky_server *server)
{
    int rc;

    check(server != NULL, "Server required");
    check(server->state == SKY_SERVER_STATE_STOPPED, "Server already running");
    check(server->port > 0, "Port required");

    // Initialize socket info.
    server->sockaddr = calloc(1, sizeof(struct sockaddr_in));
    check_mem(server->sockaddr);
    server->sockaddr->sin_addr.s_addr = INADDR_ANY;
    server->sockaddr->sin_port = htons(server->port);
    server->sockaddr->sin_family = AF_INET;

    // Create socket.
    server->socket = socket(AF_INET, SOCK_STREAM, 0);
    check(server->socket != -1, "Unable to create a socket");
    
    // Set socket for reuse.
    int optval = 1;
    rc = setsockopt(server->socket, SOL_SOCKET, SO_REUSEADDR, &optval, sizeof(optval));
    check(rc == 0, "Unable to set socket for reuse");
    
    // Bind socket.
    rc = bind(server->socket, (struct sockaddr*)server->sockaddr, sizeof(struct sockaddr_in));
    check(rc == 0, "Unable to bind socket");
    
    // Listen on socket.
    rc = listen(server->socket, SKY_LISTEN_BACKLOG);
    check(rc != -1, "Unable to listen on socket");
    
    // Update server state.
    server->state = SKY_SERVER_STATE_RUNNING;
    
    return 0;

error:
    sky_server_stop(server);
    return -1;
}

// Stops a server. This actions closes the TCP socket and in-process messages
// will be aborted.
//
// server - The server to stop.
//
// Returns 0 if successful, otherwise returns -1.
int sky_server_stop(sky_server *server)
{
    int rc;
    check(server != NULL, "Server required");
    
    // Close socket if open.
    if(server->socket > 0) {
        close(server->socket);
    }
    server->socket = 0;

    // Clear socket info.
    if(server->sockaddr) {
        free(server->sockaddr);
    }
    server->sockaddr = NULL;
    
    // Send a shutdown signal to all servlets and wait for response.
    rc = sky_server_stop_servlets(server);
    check(rc == 0, "Unable to stop servlets");

    // Update server state.
    server->state = SKY_SERVER_STATE_STOPPED;
    
    return 0;

error:
    return -1;
}

// Sends a shutdown signal to all servlets and waits for their confirmation
// before returning.
//
// server - The server.
//
// Returns 0 if successful, otherwise returns -1.
int sky_server_stop_servlets(sky_server *server)
{
    int rc;
    void *pull_socket = NULL;
    void *push_socket = NULL;
    check(server != NULL, "Server required");

    // Ignore if there are no servlets.
    if(server->servlet_count == 0) {
        return 0;
    }

    // Create pull socket.
    pull_socket = zmq_socket(server->context, ZMQ_PULL);
    check(pull_socket != NULL, "Unable to create server shutdown socket");

    // Bind pull socket.
    rc = zmq_bind(pull_socket, SKY_SERVER_SHUTDOWN_URI);
    check(rc == 0, "Unable to bind server shutdown socket");

    // Send a shutdown message to each servlet.
    uint32_t i;
    for(i=0; i<server->servlet_count; i++) {
        // Create push socket.
        push_socket = zmq_socket(server->context, ZMQ_PUSH);
        check(push_socket != NULL, "Unable to create shutdown push socket");
    
        // Connect to servlet.
        rc = zmq_connect(push_socket, bdata(server->servlets[i]->uri));
        check(rc == 0, "Unable to connect to servlet for shutdown");

        // Send NULL worklet for shutdown.
        void *ptr = NULL;
        rc = sky_zmq_send_ptr(push_socket, &ptr);
        check(rc == 0, "Unable to send worklet message");
        
        // Close socket.
        rc = zmq_close(push_socket);
        check(rc == 0, "Unable to close server shutdown push socket");
        push_socket = NULL;
    }

    // Read in one pull message for every push message sent.
    for(i=0; i<server->servlet_count; i++) {
        // Receive servlet ref on shutdown.
        sky_servlet *servlet = NULL;
        rc = sky_zmq_recv_ptr(pull_socket, (void**)&servlet);
        check(rc == 0, "Server unable to receive shutdown response");
    }

    // Clean up socket.
    rc = zmq_close(pull_socket);
    check(rc == 0, "Unable to close server shutdown pull socket");

    // Clean up servlets.
    free(server->servlets);
    server->servlets = NULL;
    server->servlet_count = 0;

    return 0;

error:
    if(pull_socket) zmq_close(pull_socket);
    if(push_socket) zmq_close(push_socket);
    return -1;
}

// Accepts a connection on a running server. Once a connection is accepted then
// the message is parsed and processed.
//
// server - The server.
//
// Returns 0 if successful, otherwise returns -1.
int sky_server_accept(sky_server *server)
{
    int rc;
    check(server != NULL, "Server required");

    // Accept the next connection.
    int sockaddr_size = sizeof(struct sockaddr_in);
    int socket = accept(server->socket, (struct sockaddr*)server->sockaddr, (socklen_t *)&sockaddr_size);
    check(socket != -1, "Unable to accept connection");

    // Wrap socket in a buffered file reference.
    FILE *input = fdopen(socket, "r");
    FILE *output = fdopen(dup(socket), "w");
    check(input != NULL, "Unable to open buffered socket input");
    check(output != NULL, "Unable to open buffered socket output");
    
    // Process message.
    rc = sky_server_process_message(server, input, output);
    check(rc == 0, "Unable to process message");
    
    return 0;

error:
    fclose(input);
    fclose(output);
    return -1;
}


//--------------------------------------
// Message Processing
//--------------------------------------

// Processes a single message.
//
// server - The server.
// input  - The input stream.
// output - The output stream.
//
// Returns 0 if successful, otherwise returns -1.
int sky_server_process_message(sky_server *server, FILE *input, FILE *output)
{
    int rc;
    sky_message_header *header = NULL;
    
    // Parse message header.
    header = sky_message_header_create(); check_mem(header);
    rc = sky_message_header_unpack(header, input);
    check(rc == 0, "Unable to unpack message header");

    // Ignore the table if this is a multi message.
    if(biseqcstr(header->name, "multi") == 1) {
        rc = sky_server_process_multi_message(server, input, output);
        check(rc == 0, "Unable to process multi message");
    }
    else {
        // Retrieve appropriate message handler by name.
        sky_message_handler *handler = NULL;
        rc = sky_server_get_message_handler(server, header->name, &handler);
        check(rc == 0, "Unable to get message handler");

        // Open table.
        sky_table *table = NULL;
        rc = sky_server_get_table(server, header->table_name, &table);
        check(rc == 0, "Unable to open table");

        // If the handler exists then use it to process the message.
        if(handler != NULL) {
            rc = handler->process(server, header, table, input, output);
        }
        // Parse appropriate message type.
        else {
            fclose(input);
            fclose(output);
            sentinel("Invalid message type");
        }
        check(rc == 0, "Unable to process message: %s", bdata(header->name));
    }
    
    sky_message_header_free(header);
    return 0;

error:
    sky_message_header_free(header);
    return -1;
}


//--------------------------------------
// Message Handlers
//--------------------------------------

// Retrieves a message handler from the server by name.
//
// server - The server.
// name   - The name of the message handler.
// ret    - A pointer to where the handler should be returned.
//
// Returns 0 if successful, otherwise returns -1.
int sky_server_get_message_handler(sky_server *server, bstring name,
                                   sky_message_handler **ret)
{
    check(server != NULL, "Server required");
    check(blength(name), "Message name required");
    check(ret != NULL, "Return pointer required");
    
    // Initialize return value.
    *ret = NULL;
    
    // Make sure a handler with the same name doesn't exist.
    uint32_t i;
    for(i=0; i<server->message_handler_count; i++) {
        if(biseq(server->message_handlers[i]->name, name) == 1) {
            *ret = server->message_handlers[i];
            break;
        }
    }

    return 0;

error:
    *ret = NULL;
    return -1;
}

// Adds a message handler to the server.
//
// server  - The server.
// handler - The message handler to add.
//
// Returns 0 if successful, otherwise returns -1.
int sky_server_add_message_handler(sky_server *server,
                                   sky_message_handler *handler)
{
    int rc;
    check(server != NULL, "Server required");
    check(handler != NULL, "Message handler required");
    
    // Make sure a handler with the same name doesn't exist.
    sky_message_handler *existing_handler = NULL;
    rc = sky_server_get_message_handler(server, handler->name, &existing_handler);
    check(rc == 0, "Unable to get existing handler");
    check(existing_handler == NULL, "Message handler '%s' already exists on server", bdata(handler->name));
    
    // Append handler to server's list of message handlers.
    server->message_handler_count++;
    server->message_handlers = realloc(server->message_handlers, server->message_handler_count * sizeof(*server->message_handlers));
    check_mem(server->message_handlers);
    server->message_handlers[server->message_handler_count-1] = handler;
    
    return 0;

error:
    return -1;
}

// Removes a message handler to the server.
//
// server  - The server.
// handler - The message handler to remove.
//
// Returns 0 if successful, otherwise returns -1.
int sky_server_remove_message_handler(sky_server *server,
                                      sky_message_handler *handler)
{
    check(server != NULL, "Server required");
    check(handler != NULL, "Message handler required");
    
    // Remove handler from server's list of message handlers.
    uint32_t i,j;
    for(i=0; i<server->message_handler_count; i++) {
        sky_message_handler *handler = server->message_handlers[i];
        if(server->message_handlers[i] == handler) {
            for(j=i; j<server->message_handler_count-1; j++) {
                server->message_handlers[j] = server->message_handlers[j+1];
            }
            server->message_handlers[server->message_handler_count-1] = NULL;
            server->message_handler_count--;
            break;
        }
    }
    
    return 0;

error:
    return -1;
}

// Add the standard message handlers to the server.
//
// server  - The server.
//
// Returns 0 if successful, otherwise returns -1.
int sky_server_add_default_message_handlers(sky_server *server)
{
    int rc;
    sky_message_handler *handler = NULL;
    check(server != NULL, "Server required");
    
    // 'Add Event' message.
    handler = sky_add_event_message_handler_create(); check_mem(handler);
    rc = sky_server_add_message_handler(server, handler);
    check(rc == 0, "Unable to add message handler");

    // 'Next Actions' message.
    handler = sky_next_actions_message_handler_create(); check_mem(handler);
    rc = sky_server_add_message_handler(server, handler);
    check(rc == 0, "Unable to add message handler");

    // 'Add Action' message.
    handler = sky_add_action_message_handler_create(); check_mem(handler);
    rc = sky_server_add_message_handler(server, handler);
    check(rc == 0, "Unable to add message handler");

    // 'Get Action' message.
    handler = sky_get_action_message_handler_create(); check_mem(handler);
    rc = sky_server_add_message_handler(server, handler);
    check(rc == 0, "Unable to add message handler");

    // 'Get Actions' message.
    handler = sky_get_actions_message_handler_create(); check_mem(handler);
    rc = sky_server_add_message_handler(server, handler);
    check(rc == 0, "Unable to add message handler");

    // 'Add Property' message.
    handler = sky_add_property_message_handler_create(); check_mem(handler);
    rc = sky_server_add_message_handler(server, handler);
    check(rc == 0, "Unable to add message handler");

    // 'Get Property' message.
    handler = sky_get_property_message_handler_create(); check_mem(handler);
    rc = sky_server_add_message_handler(server, handler);
    check(rc == 0, "Unable to add message handler");

    return 0;

error:
    return -1;
}


//--------------------------------------
// Table management
//--------------------------------------

// Retrieves a reference to a table by name. If the table is not found then
// it will be opened automatically.
//
// server - The server.
// name   - The table name.
// ret    - A pointer to where the table reference should be returned.
//
// Returns 0 if successful, otherwise returns -1.
int sky_server_get_table(sky_server *server, bstring name, sky_table **ret)
{
    int rc;
    sky_table *table = NULL;
    bstring path = NULL;
    check(server != NULL, "Server required");
    check(blength(name) > 0, "Table name required");
    
    // Initialize return values.
    *ret = NULL;
    
    // Determine the path to the table.
    path = bformat("%s/%s", bdata(server->path), bdata(name));
    check_mem(path);

    // Loop over tables to see if it's open yet.
    uint32_t i;
    for(i=0; i<server->table_count; i++) {
        if(biseq(server->tables[i]->path, path) == 1) {
            *ret = server->tables[i];
            break;
        }
    }

    // If the table is not yet opened then open it.
    if(*ret == NULL) {
        rc = sky_server_open_table(server, name, path, ret);
        check(rc == 0, "Unable to open table");
    }

    bdestroy(path);
    return 0;

error:
    bdestroy(path);
    sky_table_free(table);
    *ret = NULL;
    return -1;
}

// Opens a table and attaches it to the server. Once a table is opened it
// cannot be closed until the server shuts down.
//
// server - The server.
// name   - The table name.
// path   - The table path.
// ret    - A pointer to where the table reference should be returned.
//
// Returns 0 if successful, otherwise returns -1.
int sky_server_open_table(sky_server *server, bstring name, bstring path,
                          sky_table **ret)
{
    int rc;
    sky_table *table = NULL;
    check(server != NULL, "Server required");
    check(blength(path) > 0, "Table path required");
    check(ret != NULL, "Return pointer required");
    
    // Initialize return values.
    *ret = NULL;
    
    // Create the table.
    table = sky_table_create(); check_mem(table);
    table->name = bstrcpy(name); check_mem(table->name);
    rc = sky_table_set_path(table, path);
    check(rc == 0, "Unable to set table path");

    // Open the table.
    rc = sky_table_open(table);
    check(rc == 0, "Unable to open table");
    
    // Append the table to the list of open tables.
    server->table_count++;
    server->tables = realloc(server->tables, server->table_count * sizeof(*server->tables));
    check_mem(server->tables);
    server->tables[server->table_count-1] = table;

    // Open servlets for the table's tablets.
    rc = sky_server_create_servlets(server, table);
    check(rc == 0, "Unable to create servlets for table");

    // Return table.
    *ret = table;

    return 0;

error:
    sky_table_free(table);
    *ret = NULL;
    return -1;
}

//--------------------------------------
// Servlet Management
//--------------------------------------

// Retrieves a servlet that processes a given tablet.
//
// server    - The server.
// tablet    - The tablet.
// servlet   - A pointer to where the servlet should be returned.
//
// Returns 0 if successful, otherwise returns -1.
int sky_server_get_tablet_servlet(sky_server *server, sky_tablet *tablet,
                                  sky_servlet **servlet)
{
    check(server != NULL, "Server required");
    check(tablet != NULL, "Tablet required");
    check(servlet != NULL, "Servlet return pointer required");
    
    // Loop over all servlets and find the one associated with the tablet.
    uint32_t i;
    for(i=0; i<server->servlet_count; i++) {
        if(server->servlets[i]->tablet == tablet) {
            *servlet = server->servlets[i];
            break;
        }
    }
    
    return 0;

error:
    *servlet = NULL;
    return -1;
}

// Retrieves a list of servlets associated with a given table.
//
// server   - The server.
// table    - The table.
// servlets - A pointer to where the servlets should be returned.
// count    - A pointer to where the number of servlets should be returned.
//
// Returns 0 if successful, otherwise returns -1.
int sky_server_get_table_servlets(sky_server *server, sky_table *table,
                                  sky_servlet ***servlets, uint32_t *count)
{
    check(server != NULL, "Server required");
    check(table != NULL, "Table required");
    check(servlets != NULL, "Servlets return pointer required");
    check(count != NULL, "Servlet count return pointer required");
    
    // Allocate array.
    *count = table->tablet_count;
    *servlets = calloc(*count, sizeof(*servlets));
    check_mem(*servlets);
    
    // Loop over all servlets and find ones associated with the table.
    uint32_t i, index=0;
    for(i=0; i<server->servlet_count; i++) {
        sky_servlet *servlet = server->servlets[i];
        if(servlet->tablet->table == table) {
            (*servlets)[index] = servlet;
            index++;
        }
    }
    
    return 0;

error:
    *count = 0;
    free(*servlets);
    return -1;
}

// Creates a set of servlets for a given table. Each servlet is responsible
// for a single tablet on the table.
//
// server - The server.
// table  - The table to create servlets against.
//
// Returns 0 if successful, otherwise returns -1.
int sky_server_create_servlets(sky_server *server, sky_table *table)
{
    int rc;
    sky_servlet *servlet = NULL;
    check(server != NULL, "Server required");
    check(table, "Table required");
    
    // Allocate additional space for the new servlets.
    uint32_t new_servlet_count = table->tablet_count;
    server->servlets = realloc(server->servlets, (server->servlet_count + new_servlet_count) * sizeof(*server->servlets));
    check_mem(server->servlets);
    memset(&server->servlets[server->servlet_count], 0, sizeof(*server->servlets) * new_servlet_count);
    
    // Loop over tablets and create one servlet for each one.
    uint32_t i;
    for(i=0; i<new_servlet_count; i++) {
        // Create the servlet.
        sky_tablet *tablet = table->tablets[i];
        servlet = sky_servlet_create(server, tablet); check_mem(servlet);

        // Start the servlet.
        rc = sky_servlet_start(servlet);
        check(rc == 0, "Unable to start servlet");
        
        // Append the servlet to the server.
        server->servlets[server->servlet_count] = servlet;
        server->servlet_count++;

        servlet = NULL;
    }

    return 0;

error:
    sky_servlet_free(servlet);
    return -1;
}


//--------------------------------------
// Multi Message
//--------------------------------------

// Parses and process a multi message.
//
// server - The server.
// input  - The input file stream.
// output - The output file stream.
//
// Returns 0 if successful, otherwise returns -1.
int sky_server_process_multi_message(sky_server *server, FILE *input,
                                     FILE *output)
{
    int rc;
    check(server != NULL, "Server required");
    check(input != NULL, "Input required");
    check(output != NULL, "Output stream required");
    
    debug("Message received: [MULTI]");
    
    // Parse message.
    sky_multi_message *message = sky_multi_message_create(); check_mem(message);
    rc = sky_multi_message_unpack(message, input);
    check(rc == 0, "Unable to parse MULTI message");

    // Start time.
    struct timeval tv;
    gettimeofday(&tv, NULL);
    int64_t t0 = (tv.tv_sec*1000) + (tv.tv_usec/1000);

    // Process message.
    uint32_t i;
    for(i=0; i<message->message_count; i++) {
        rc = sky_server_process_message(server, input, output);
        check(rc == 0, "Unable to process child message");
    }

    // End time.
    gettimeofday(&tv, NULL);
    int64_t t1 = (tv.tv_sec*1000) + (tv.tv_usec/1000);
    printf("MULTI: %d messages processed in %.3f seconds\n", message->message_count, ((float)(t1-t0))/1000);
    
    return 0;

error:
    return -1;
}
