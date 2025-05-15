package math

import (
	"github.com/sevenreup/duwa/src/object"
	"github.com/sevenreup/duwa/src/token"
	"github.com/shopspring/decimal"
)

// method=yochepa args=[number{number1}, number{number2}] return={number}
// This method returns the smaller of two numbers.
//
// `Example`
// ```
// Masamu.yochepa(5, 10) # returns 5
// ```
func methodMathMin(scope *object.Environment, tok token.Token, args ...object.Object) object.Object {
	if len(args) < 2 {
		panic("Masamu.yochepa requires at least two arguments")
	}

	if args[0].Type() != object.INTEGER_OBJ {
		return nil
	}

	if args[1].Type() != object.INTEGER_OBJ {
		return nil
	}

	number1 := args[0].(*object.Integer)
	number2 := args[1].(*object.Integer)

	if number1.Value.LessThan(number2.Value) {
		return number1
	}

	return number2
}

// method=sqrt args=[number{number1}] return={number}
// This method returns the square root of a number.
//
// `Example`
// ```
// Masamu.sqrt(25) # returns 5
// ```
func methodMathSqrt(scope *object.Environment, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 1 {
		panic("Masamu.sqrt requires one argument")
	}

	if args[0].Type() != object.INTEGER_OBJ {
		return nil
	}

	number := args[0].(*object.Integer)

	// NOTE: using pow(0.5) is the same as sqrt() but may not be accurate
	return &object.Integer{Value: number.Value.Pow(decimal.NewFromFloat(0.5))}
}

// method=round args=[number{number1}, number{number2}] return={number}
// This method rounds a number to a specified number of decimal places.
//
// `Example`
// ```
// Masamu.round(5.678, 2) # returns 5.68
// ```
func methodRound(scope *object.Environment, tok token.Token, args ...object.Object) object.Object {
	if len(args) < 2 {
		panic("Masamu.yochepa requires at least two arguments")
	}

	if args[0].Type() != object.INTEGER_OBJ {
		return nil
	}

	if args[1].Type() != object.INTEGER_OBJ {
		return nil
	}

	number := args[0].(*object.Integer)
	places := args[1].(*object.Integer)

	return &object.Integer{Value: number.Value.Round(int32(places.Value.IntPart()))}
}

// method=pansi args=[number{number1}] return={number}
// This method returns the largest integer less than or equal to a number.
//
// `Example`
// ```
// Masamu.pansi(5.678) # returns 5
// ```
func methodFloor(scope *object.Environment, tok token.Token, args ...object.Object) object.Object {
	if args[0].Type() != object.INTEGER_OBJ {
		return nil
	}
	number1 := args[0].(*object.Integer)
	return &object.Integer{Value: number1.Value.Floor()}
}

// library=Masamu
// This is the math module
// It contains functions for performing mathematical operations
// It is used to perform mathematical calculations
func Module() *object.LibraryModule {
	return object.NewBuiltInLibraryModule("Masamu", map[string]*object.LibraryFunction{
		"yochepa": object.NewBuiltin("yochepa", methodMathMin),
		"sqrt":    object.NewBuiltin("sqrt", methodMathSqrt),
		"round":   object.NewBuiltin("round", methodRound),
		"pansi":   object.NewBuiltin("pansi", methodFloor),
	})
}
