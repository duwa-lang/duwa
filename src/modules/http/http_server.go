package http

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sevenreup/duwa/src/object"
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
	return "http.Server"
}

func (s *HTTPServer) Method(method string, args []object.Object) (object.Object, bool) {
	switch method {
	case "yamba":
		return s.methodStartServer(args)
	case "TENGA":
		return s.methodGET(args)
	case "fayilo":
		return s.methodServeFiles(args)
	case "fayilo_imodzi":
		return s.methodServeFile(args)
	}
	return nil, false
}

// method=GET args=[HTTPServer{server}, mawu{path}] return={null}
// This method adds a GET route to the HTTP server.
//
// `Example`
// ```
// http.GET(sevala, "/hello")
// ```
func (s *HTTPServer) methodGET(args []object.Object) (object.Object, bool) {
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

	return values.NULL, true
}

// method=fayilo args=[HTTPServer{server}, mawu{route}, mawu{directory}] return={null}
// This method serves static files from a directory.
//
// `Example`
// ```
// http.fayilo(sevala, "/static/", "./public")
// ```
func (s *HTTPServer) methodServeFiles(args []object.Object) (object.Object, bool) {
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

	return values.NULL, true
}

// method=fayilo_imodzi args=[HTTPServer{server}, mawu{route}, mawu{filepath}] return={null}
// This method serves a single file at a specific route.
//
// `Example`
// ```
// http.fayilo_imodzi(sevala, "/", "./index.html")
// ```
func (s *HTTPServer) methodServeFile(args []object.Object) (object.Object, bool) {
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

	return values.NULL, true
}

// method=yamba args=[HTTPServer{server}] return={null}
// This method starts the HTTP server.
//
// `Example`
// ```
// http.yamba(sevala) # starts the server
// ```
func (s *HTTPServer) methodStartServer(args []object.Object) (object.Object, bool) {
	if len(args) != 1 {
		panic("http.yamba requires one argument (server)")
	}

	server, ok := args[0].(*HTTPServer)
	if !ok {
		panic("http.yamba argument must be an HTTP server")
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

	return values.NULL, true
}
