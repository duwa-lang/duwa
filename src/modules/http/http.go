package http

import (
	"fmt"
	"net/http"
	"os"
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

// method=fayilo args=[HTTPServer{server}, mawu{route}, mawu{directory}] return={null}
// This method serves static files from a directory.
//
// `Example`
// ```
// http.fayilo(sevala, "/static/", "./public")
// ```
func methodServeFiles(scope *object.Environment, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 3 {
		panic("http.fayilo requires three arguments (server, route, directory)")
	}

	server, ok := args[0].(*HTTPServer)
	if !ok {
		panic("http.fayilo first argument must be an HTTP server")
	}

	if args[1].Type() != object.STRING_OBJ {
		panic("http.fayilo route must be a string")
	}

	if args[2].Type() != object.STRING_OBJ {
		panic("http.fayilo directory must be a string")
	}

	route := args[1].(*object.String).Value
	directory := args[2].(*object.String).Value

	// Check if directory exists
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		panic(fmt.Sprintf("Directory does not exist: %s", directory))
	}

	// Create file server handler
	fileServer := http.FileServer(http.Dir(directory))
	
	// Strip the route prefix when serving files
	if route == "/" {
		server.routes[route] = fileServer.ServeHTTP
	} else {
		server.routes[route] = http.StripPrefix(route, fileServer).ServeHTTP
	}

	return values.NULL
}

// method=fayilo_imodzi args=[HTTPServer{server}, mawu{route}, mawu{filepath}] return={null}
// This method serves a single file at a specific route.
//
// `Example`
// ```
// http.fayilo_imodzi(sevala, "/", "./index.html")
// ```
func methodServeFile(scope *object.Environment, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 3 {
		panic("http.fayilo_imodzi requires three arguments (server, route, filepath)")
	}

	server, ok := args[0].(*HTTPServer)
	if !ok {
		panic("http.fayilo_imodzi first argument must be an HTTP server")
	}

	if args[1].Type() != object.STRING_OBJ {
		panic("http.fayilo_imodzi route must be a string")
	}

	if args[2].Type() != object.STRING_OBJ {
		panic("http.fayilo_imodzi filepath must be a string")
	}

	route := args[1].(*object.String).Value
	filePath := args[2].(*object.String).Value

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		panic(fmt.Sprintf("File does not exist: %s", filePath))
	}

	// Create handler that serves the specific file
	server.routes[route] = func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "HEAD" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, filePath)
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
		"panga":        object.NewBuiltin("panga", methodCreateServer),  // create server
		"TENGA":        object.NewBuiltin("TENGA", methodGET),           // add GET route
		"yamba":        object.NewBuiltin("yamba", methodStartServer),   // start server
		"fayilo":       object.NewBuiltin("fayilo", methodServeFiles),   // serve directory
		"fayilo_imodzi": object.NewBuiltin("fayilo_imodzi", methodServeFile), // serve single file
	})
}
