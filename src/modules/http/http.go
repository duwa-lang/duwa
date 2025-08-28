package http

import (
	"net/http"
	"strconv"

	"github.com/sevenreup/duwa/src/object"
	"github.com/sevenreup/duwa/src/token"
)

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

// library=http
// This is the HTTP module
// It contains functions for creating and managing HTTP servers
// It is used to create web servers and handle HTTP requests
func Module() *object.LibraryModule {
	return object.NewBuiltInLibraryModule("http", map[string]*object.LibraryFunction{
		"panga": object.NewBuiltin("panga", methodCreateServer), // create server
	})
}
