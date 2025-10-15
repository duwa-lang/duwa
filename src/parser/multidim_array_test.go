package parser

import (
	"testing"

	"github.com/duwa-lang/duwa/src/ast"
)

func TestMultiDimensionalArrayLiterals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int // expected number of elements in outer array
	}{
		{
			name:     "2D array - 2x2",
			input:    "[[1, 2], [3, 4]]",
			expected: 2,
		},
		{
			name:     "2D array - 3x3",
			input:    "[[1, 2, 3], [4, 5, 6], [7, 8, 9]]",
			expected: 3,
		},
		{
			name:     "3D array",
			input:    "[[[1, 2], [3, 4]], [[5, 6], [7, 8]]]",
			expected: 2,
		},
		{
			name:     "Jagged array",
			input:    "[[1], [2, 3], [4, 5, 6]]",
			expected: 3,
		},
		{
			name:     "Empty nested arrays",
			input:    "[[], [1], []]",
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			file := p.ParseFile([]byte(tt.input))
			checkParserErrors(t, p)

			if len(file.Statements) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(file.Statements))
			}

			stmt, ok := file.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("expected ExpressionStatement, got %T", file.Statements[0])
			}

			array, ok := stmt.Expression.(*ast.ArrayLiteral)
			if !ok {
				t.Fatalf("expected ArrayLiteral, got %T", stmt.Expression)
			}

			if len(array.Elements) != tt.expected {
				t.Errorf("expected %d outer elements, got %d", tt.expected, len(array.Elements))
			}

			// Verify that elements are also arrays (for the first element)
			if len(array.Elements) > 0 {
				innerArray, ok := array.Elements[0].(*ast.ArrayLiteral)
				if !ok {
					t.Errorf("expected inner element to be ArrayLiteral, got %T", array.Elements[0])
				}
				if innerArray != nil && len(innerArray.Elements) == 0 && tt.name != "Empty nested arrays" {
					t.Errorf("expected inner array to have elements")
				}
			}
		})
	}
}

func TestMultiDimensionalArrayIndexAccess(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "2D array access",
			input: "arr[0][1]",
		},
		{
			name:  "3D array access",
			input: "matrix[0][1][2]",
		},
		{
			name:  "4D array access",
			input: "hyper[0][1][2][3]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			file := p.ParseFile([]byte(tt.input))
			checkParserErrors(t, p)

			if len(file.Statements) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(file.Statements))
			}

			stmt, ok := file.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("expected ExpressionStatement, got %T", file.Statements[0])
			}

			// The result should be an IndexExpression
			_, ok = stmt.Expression.(*ast.IndexExpression)
			if !ok {
				t.Fatalf("expected IndexExpression, got %T", stmt.Expression)
			}
		})
	}
}

func TestMultiDimensionalArrayTypeDeclarations(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "2D array type",
			input: "nambala[][] matrix = [[1, 2], [3, 4]];",
		},
		{
			name:  "3D array type",
			input: "nambala[][][] cube = [[[1, 2]], [[3, 4]]];",
		},
		{
			name:  "4D array type",
			input: "nambala[][][][] hyper = [[[[1]]]];",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			file := p.ParseFile([]byte(tt.input))
			checkParserErrors(t, p)

			if len(file.Statements) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(file.Statements))
			}

			stmt, ok := file.Statements[0].(*ast.VariableDeclarationStatement)
			if !ok {
				t.Fatalf("expected VariableDeclarationStatement, got %T", file.Statements[0])
			}

			if stmt.Identifier == nil {
				t.Fatalf("expected identifier to be set")
			}

			if stmt.Value == nil {
				t.Fatalf("expected value to be set")
			}

			// Verify the value is an array literal
			_, ok = stmt.Value.(*ast.ArrayLiteral)
			if !ok {
				t.Fatalf("expected value to be ArrayLiteral, got %T", stmt.Value)
			}
		})
	}
}

func TestNestedArrayModification(t *testing.T) {
	input := "matrix[0][1] = 42;"

	p := NewParser()
	file := p.ParseFile([]byte(input))
	checkParserErrors(t, p)

	if len(file.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(file.Statements))
	}

	// Just verify it parses without errors
	// The actual assignment is tested in runtime tests
	if len(p.Errors()) > 0 {
		t.Fatalf("parser had errors: %v", p.Errors())
	}
}
