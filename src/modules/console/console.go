package console

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/duwa-lang/duwa/src/object"
	"github.com/duwa-lang/duwa/src/token"
	"github.com/duwa-lang/duwa/src/values"
)

// method=lemba args=[...] return={null}
// This method prints the arguments to the console
//
// `Example`
// ```
// khonso.lemba("Hello, World!") # prints "Hello, World!"
// ```
func methodConsolePrint(env *object.Environment, tok token.Token, args ...object.Object) object.Object {
	values := make([]string, 0)

	for _, value := range args {
		values = append(values, value.String())
	}

	libPrint(env, values)

	return nil
}

// method=landira args=[mawu{zolemba}] return={mawu}
// This method reads a string from the console
//
// `Example`
// ```
// khonso.landira() # returns the string entered by the user
// ```
func methodConsoleRead(scope *object.Environment, tok token.Token, args ...object.Object) object.Object {
	if len(args) == 1 {
		prompt := args[0].(*object.String).Value

		fmt.Print(prompt)
	}

	val, err := scope.Console.Read()

	if err != nil {
		return values.NULL
	}

	return &object.String{Value: val}
}

// method=fufuta args=[] return={null}
// This method clears the console
//
// `Example`
// ```
// khonso.fufuta() # clears the console
// ```
func methodConsoleClear(scope *object.Environment, tok token.Token, args ...object.Object) object.Object {
	err := scope.Console.Clear()
	if err != nil {
		return values.NULL
	}

	return values.NULL
}

func libPrint(env *object.Environment, values []string) {
	if len(values) > 0 {
		str := make([]string, 0)

		str = append(str, values...)
		strRaw, _ := strconv.Unquote(`"` + strings.Join(str, " ") + `"`)
		env.Logger.Info(strRaw)
	}
}

// type=builtin-func method=lemba args=[any{valueToPrint}] return={null}
// The lemba function prints a value to the console.
func BuiltInPrint(env *object.Environment, tok token.Token, args ...object.Object) object.Object {
	if len(args) > 0 {
		str := make([]string, 0)

		for _, value := range args {
			str = append(str, value.String())
		}
		env.Logger.Info(strings.Join(str, " "))
	}

	return nil
}

// type=builtin-func method=lembanzr args=[any{valueToPrint}] return={null}
// The lembanzr function prints a value to the console and adds a newline.
func BuiltInPrintLine(env *object.Environment, tok token.Token, args ...object.Object) object.Object {
	if len(args) > 0 {
		str := make([]string, 0)

		for _, value := range args {
			str = append(str, value.String())
		}

		env.Logger.Info(strings.Join(str, " ") + "\n")
	} else {
		env.Logger.Info("\n")
	}

	return nil
}

// library=khonso
// This is the console module
// It contains functions that interact with the console
// It is used to read and write to the console
func Module() *object.LibraryModule {
	return object.NewBuiltInLibraryModule("khonso", map[string]*object.LibraryFunction{
		"lemba":   object.NewBuiltin("lemba", methodConsolePrint),
		"fufuta":  object.NewBuiltin("fufuta", methodConsoleClear),
		"landira": object.NewBuiltin("landira", methodConsoleRead),
	})
}

func Builtins() map[string]*object.LibraryFunction {
	return map[string]*object.LibraryFunction{
		"lemba":    object.NewBuiltin("lemba", BuiltInPrintLine),
		"lembanzr": object.NewBuiltin("lembanzr", BuiltInPrintLine),
	}
}
