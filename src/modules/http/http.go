package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/sevenreup/duwa/src/object"
	"github.com/sevenreup/duwa/src/token"
	"github.com/sevenreup/duwa/src/values"
)

// HTTP Server object to hold the server state
type HTTPServer struct {
	server *http.Server
	routes map[string]func(http.ResponseWriter, *http.Request)
}

func (s *HTTPServer) String() string {
	return fmt.Sprintf("HTTP Server on %s", s.server.Addr)
}

func (s *HTTPServer) Type() object.ObjectType {
	return "HTTP_SERVER"
}

func (s *HTTPServer) Method(method string, args []object.Object) (object.Object, bool) {
	return nil, false
}

// method=yambisa args=[nambala{port}] return={HTTPServer}
// This method creates a new HTTP server on the specified port.
//
// `Example`
// ```
// sevala = http.yambisa(8080) # creates server on port 8080
// ```
func methodCreateServer(scope *object.Environment, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 1 {
		panic("http.yambisa requires one argument (port)")
	}

	if args[0].Type() != object.INTEGER_OBJ {
		panic("http.yambisa port must be a number")
	}

	port := args[0].(*object.Integer)
	portStr := ":" + strconv.FormatInt(port.Value.IntPart(), 10)

	server := &HTTPServer{
		server: &http.Server{Addr: portStr},
		routes: make(map[string]func(http.ResponseWriter, *http.Request)),
	}

	return server
}

// method=GET args=[HTTPServer{server}, mawu{path}] return={null}
// This method adds a GET route to the HTTP server.
//
// `Example`
// ```
// http.GET(sevala, "/hello")
// ```
func methodGET(scope *object.Environment, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 2 {
		panic("http.GET requires two arguments (server, path)")
	}

	server, ok := args[0].(*HTTPServer)
	if !ok {
		panic("http.GET first argument must be an HTTP server")
	}

	if args[1].Type() != object.STRING_OBJ {
		panic("http.GET path must be a string")
	}

	path := args[1].(*object.String).Value

	// Store a simple handler for now
	server.routes[path] = func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Hello from Duwa HTTP server! Path: %s", path)
	}

	return values.NULL
}

// method=yendesa args=[HTTPServer{server}] return={null}
// This method starts the HTTP server.
//
// `Example`
// ```
// http.yendesa(sevala) # starts the server
// ```
func methodStartServer(scope *object.Environment, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 1 {
		panic("http.yendesa requires one argument (server)")
	}

	server, ok := args[0].(*HTTPServer)
	if !ok {
		panic("http.yendesa argument must be an HTTP server")
	}

	// Set up the routes
	mux := http.NewServeMux()
	for path, handler := range server.routes {
		mux.HandleFunc(path, handler)
	}
	server.server.Handler = mux

	// Start server (blocking call to keep program alive)
	fmt.Printf("Starting HTTP server on %s\n", server.server.Addr)
	fmt.Printf("Server will keep running... Press Ctrl+C to stop\n")

	// Small delay to ensure the server is ready
	go func() {
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("Server is ready and listening for requests\n")
	}()

	if err := server.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Server error: %v\n", err)
	}

	return values.NULL
}

// library=http
// This is the HTTP module
// It contains functions for creating and managing HTTP servers
// It is used to create web servers and handle HTTP requests
func Module() *object.LibraryModule {
	return object.NewBuiltInLibraryModule("http", map[string]*object.LibraryFunction{
		"panga": object.NewBuiltin("panga", methodCreateServer), // create server
		"TENGA": object.NewBuiltin("TENGA", methodGET),          // add GET route
		"yamba": object.NewBuiltin("yamba", methodStartServer),  // start server
	})
}
