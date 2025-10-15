package duwa

import (
	"testing"

	"github.com/duwa-lang/duwa/src/object"
	"github.com/shopspring/decimal"
)

func TestCallE_Success(t *testing.T) {
	env := object.NewDefaultEnvironment()
	duwa := New(env)

	// Load test file
	duwa.Run(`
		ndondomeko add(a, b) {
			bweza a + b;
		}
	`)

	// Test successful function call
	args := []object.Object{
		&object.Integer{Value: decimal.NewFromInt(5)},
		&object.Integer{Value: decimal.NewFromInt(3)},
	}

	result, err := duwa.CallE("add", args)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	intResult, ok := result.(*object.Integer)
	if !ok {
		t.Fatalf("Expected Integer result, got: %T", result)
	}

	expected := int64(8)
	if intResult.Value.IntPart() != expected {
		t.Errorf("Expected %d, got %d", expected, intResult.Value.IntPart())
	}
}

func TestCallE_FunctionNotFound(t *testing.T) {
	env := object.NewDefaultEnvironment()
	duwa := New(env)

	// Test calling non-existent function
	args := []object.Object{}
	result, err := duwa.CallE("nonExistentFunction", args)

	if err == nil {
		t.Fatal("Expected error for non-existent function, got nil")
	}

	if result != nil {
		t.Fatalf("Expected nil result, got: %v", result)
	}

	expectedError := "function not found: nonExistentFunction"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestCallE_ArgumentCountMismatch(t *testing.T) {
	env := object.NewDefaultEnvironment()
	duwa := New(env)

	// Load test file
	duwa.Run(`
		ndondomeko multiply(a, b) {
			bweza a * b;
		}
	`)

	// Test calling with wrong number of arguments
	args := []object.Object{
		&object.Integer{Value: decimal.NewFromInt(5)},
	}

	result, err := duwa.CallE("multiply", args)

	if err == nil {
		t.Fatal("Expected error for argument count mismatch, got nil")
	}

	if result != nil {
		t.Fatalf("Expected nil result, got: %v", result)
	}

	// The error message should contain "argument count mismatch"
	if len(err.Error()) == 0 {
		t.Error("Expected non-empty error message")
	}
}

func TestCallE_WithStringResult(t *testing.T) {
	env := object.NewDefaultEnvironment()
	duwa := New(env)

	// Load test file
	duwa.Run(`
		ndondomeko greet(name) {
			bweza "Moni " + name;
		}
	`)

	// Test successful function call with string
	args := []object.Object{
		&object.String{Value: "Duwa"},
	}

	result, err := duwa.CallE("greet", args)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	strResult, ok := result.(*object.String)
	if !ok {
		t.Fatalf("Expected String result, got: %T", result)
	}

	expected := "Moni Duwa"
	if strResult.Value != expected {
		t.Errorf("Expected '%s', got '%s'", expected, strResult.Value)
	}
}

func TestEnvironmentCallE_Success(t *testing.T) {
	env := object.NewDefaultEnvironment()
	duwa := New(env)

	// Load test file
	duwa.Run(`
		ndondomeko subtract(a, b) {
			bweza a - b;
		}
	`)

	// Test successful function call via environment
	args := []object.Object{
		&object.Integer{Value: decimal.NewFromInt(10)},
		&object.Integer{Value: decimal.NewFromInt(4)},
	}

	result, err := env.CallE("subtract", args)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	intResult, ok := result.(*object.Integer)
	if !ok {
		t.Fatalf("Expected Integer result, got: %T", result)
	}

	expected := int64(6)
	if intResult.Value.IntPart() != expected {
		t.Errorf("Expected %d, got %d", expected, intResult.Value.IntPart())
	}
}

func TestEnvironmentCallE_Error(t *testing.T) {
	env := object.NewDefaultEnvironment()

	// Test calling non-existent function via environment
	args := []object.Object{}
	result, err := env.CallE("missingFunc", args)

	if err == nil {
		t.Fatal("Expected error for non-existent function, got nil")
	}

	if result != nil {
		t.Fatalf("Expected nil result, got: %v", result)
	}
}
