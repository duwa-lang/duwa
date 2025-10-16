package parser

import (
	"testing"

	"github.com/duwa-lang/duwa/src/ast"
)

// TestIfExpressionWithIdentifier tests if expressions with simple identifiers as conditions
func TestIfExpressionWithIdentifier(t *testing.T) {
	tests := []struct {
		input      string
		identifier string
		shouldFail bool
	}{
		{input: `ngati(started){}`, identifier: "started", shouldFail: false},
		{input: `ngati (started) {}`, identifier: "started", shouldFail: false},
		{input: `ngati(isReady){}`, identifier: "isReady", shouldFail: false},
		{input: `ngati (x) { lemba("yes") }`, identifier: "x", shouldFail: false},
	}

	for _, tt := range tests {
		p := NewParser()
		file := p.ParseFile([]byte(tt.input))

		if tt.shouldFail {
			if len(p.Errors()) == 0 {
				t.Errorf("Expected parsing to fail for: %s", tt.input)
			}
			continue
		}

		checkParserErrors(t, p)

		if len(file.Statements) != 1 {
			t.Fatalf("file.Statements does not contain 1 statement. got=%d for input: %s",
				len(file.Statements), tt.input)
		}

		stmt, ok := file.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("file.Statements[0] is not ast.ExpressionStatement. got=%T for input: %s",
				file.Statements[0], tt.input)
		}

		exp, ok := stmt.Expression.(*ast.IfExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T for input: %s",
				stmt.Expression, tt.input)
		}

		if !testIdentifier(t, exp.Condition, tt.identifier) {
			t.Errorf("Failed to parse condition identifier for: %s", tt.input)
		}
	}
}

// TestIfExpressionWithMethodCall tests if expressions with method calls as conditions
func TestIfExpressionWithMethodCall(t *testing.T) {
	tests := []struct {
		input      string
		object     string
		method     string
		shouldFail bool
	}{
		{input: `ngati(zenera.isKeyPressed(zenera.KEY_SPACE)){}`, object: "zenera", method: "isKeyPressed", shouldFail: false},
		{input: `ngati (zenera.isKeyPressed(zenera.KEY_SPACE)) {}`, object: "zenera", method: "isKeyPressed", shouldFail: false},
		{input: `ngati(player.isAlive()){}`, object: "player", method: "isAlive", shouldFail: false},
		{input: `ngati (game.checkCollision(x, y)) { bweza zoona; }`, object: "game", method: "checkCollision", shouldFail: false},
	}

	for _, tt := range tests {
		p := NewParser()
		file := p.ParseFile([]byte(tt.input))

		if tt.shouldFail {
			if len(p.Errors()) == 0 {
				t.Errorf("Expected parsing to fail for: %s", tt.input)
			}
			continue
		}

		checkParserErrors(t, p)

		if len(file.Statements) != 1 {
			t.Fatalf("file.Statements does not contain 1 statement. got=%d for input: %s",
				len(file.Statements), tt.input)
		}

		stmt, ok := file.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("file.Statements[0] is not ast.ExpressionStatement. got=%T for input: %s",
				file.Statements[0], tt.input)
		}

		exp, ok := stmt.Expression.(*ast.IfExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T for input: %s",
				stmt.Expression, tt.input)
		}

		methodCall, ok := exp.Condition.(*ast.MethodExpression)
		if !ok {
			t.Fatalf("exp.Condition is not ast.MethodExpression. got=%T for input: %s",
				exp.Condition, tt.input)
		}

		if !testIdentifier(t, methodCall.Left, tt.object) {
			t.Errorf("Failed to parse method object for: %s", tt.input)
		}

		if !testIdentifier(t, methodCall.Method, tt.method) {
			t.Errorf("Failed to parse method name for: %s", tt.input)
		}
	}
}

// TestIfExpressionWithLogicalOperators tests if expressions with logical operators (&&, ||)
func TestIfExpressionWithLogicalOperators(t *testing.T) {
	tests := []struct {
		input      string
		operator   string
		shouldFail bool
	}{
		{
			input:      `ngati(zenera.isKeyPressed(zenera.KEY_SPACE) || zenera.isKeyPressed(zenera.KEY_ENTER)){}`,
			operator:   "||",
			shouldFail: false,
		},
		{
			input:      `ngati (zenera.isKeyPressed(zenera.KEY_SPACE) || zenera.isKeyPressed(zenera.KEY_ENTER)) {}`,
			operator:   "||",
			shouldFail: false,
		},
		{
			input:      `ngati(x > 5 && y < 10){}`,
			operator:   "&&",
			shouldFail: false,
		},
		{
			input:      `ngati (started && !paused) {}`,
			operator:   "&&",
			shouldFail: false,
		},
		{
			input:      `ngati(player.isAlive() && player.hasAmmo()){}`,
			operator:   "&&",
			shouldFail: false,
		},
		{
			input:      `ngati (zoona || bodza) { lemba("test") }`,
			operator:   "||",
			shouldFail: false,
		},
	}

	for _, tt := range tests {
		p := NewParser()
		file := p.ParseFile([]byte(tt.input))

		if tt.shouldFail {
			if len(p.Errors()) == 0 {
				t.Errorf("Expected parsing to fail for: %s", tt.input)
			}
			continue
		}

		checkParserErrors(t, p)

		if len(file.Statements) != 1 {
			t.Fatalf("file.Statements does not contain 1 statement. got=%d for input: %s",
				len(file.Statements), tt.input)
		}

		stmt, ok := file.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("file.Statements[0] is not ast.ExpressionStatement. got=%T for input: %s",
				file.Statements[0], tt.input)
		}

		exp, ok := stmt.Expression.(*ast.IfExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T for input: %s",
				stmt.Expression, tt.input)
		}

		infixExp, ok := exp.Condition.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp.Condition is not ast.InfixExpression. got=%T for input: %s",
				exp.Condition, tt.input)
		}

		if infixExp.Operator != tt.operator {
			t.Errorf("Expected operator %s, got %s for input: %s",
				tt.operator, infixExp.Operator, tt.input)
		}
	}
}

// TestIfExpressionWithComplexConditions tests complex nested conditions
func TestIfExpressionWithComplexConditions(t *testing.T) {
	tests := []struct {
		input       string
		description string
		shouldFail  bool
	}{
		{
			input:       `ngati((x > 5 && y < 10) || z == 0){}`,
			description: "Grouped conditions with mixed operators",
			shouldFail:  false,
		},
		{
			input:       `ngati (player.x >= 0 && player.x <= width && player.y >= 0 && player.y <= height) {}`,
			description: "Multiple chained && operations",
			shouldFail:  false,
		},
		{
			input:       `ngati(arr[0] > 0 || arr[1] > 0){}`,
			description: "Array access in condition",
			shouldFail:  false,
		},
		{
			input:       `ngati (!started && ready) {}`,
			description: "Negation with logical AND",
			shouldFail:  false,
		},
		{
			input:       `ngati(game.isOver() || time == 0){}`,
			description: "Method call with OR operator",
			shouldFail:  false,
		},
	}

	for _, tt := range tests {
		p := NewParser()
		file := p.ParseFile([]byte(tt.input))

		if tt.shouldFail {
			if len(p.Errors()) == 0 {
				t.Errorf("Expected parsing to fail for %s: %s", tt.description, tt.input)
			}
			continue
		}

		checkParserErrors(t, p)

		if len(file.Statements) != 1 {
			t.Fatalf("file.Statements does not contain 1 statement for %s. got=%d",
				tt.description, len(file.Statements))
		}

		stmt, ok := file.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("file.Statements[0] is not ast.ExpressionStatement for %s. got=%T",
				tt.description, file.Statements[0])
		}

		_, ok = stmt.Expression.(*ast.IfExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.IfExpression for %s. got=%T",
				tt.description, stmt.Expression)
		}

		// Just verify it parses, the condition structure is tested in other tests
		t.Logf("✓ Successfully parsed: %s", tt.description)
	}
}

// TestIfExpressionWithoutSpaces tests expressions without spaces
func TestIfExpressionWithoutSpaces(t *testing.T) {
	tests := []struct {
		input       string
		description string
		shouldFail  bool
	}{
		{input: `ngati(x>5){}`, description: "No spaces in comparison", shouldFail: false},
		{input: `ngati(x==y){}`, description: "No spaces in equality", shouldFail: false},
		{input: `ngati(x!=y){}`, description: "No spaces in inequality", shouldFail: false},
		{input: `ngati(x>=5&&y<=10){}`, description: "No spaces with && operator", shouldFail: false},
		{input: `ngati(x<5||y>10){}`, description: "No spaces with || operator", shouldFail: false},
	}

	for _, tt := range tests {
		p := NewParser()
		file := p.ParseFile([]byte(tt.input))

		if tt.shouldFail {
			if len(p.Errors()) == 0 {
				t.Errorf("Expected parsing to fail for %s: %s", tt.description, tt.input)
			}
			continue
		}

		checkParserErrors(t, p)

		if len(file.Statements) != 1 {
			t.Fatalf("file.Statements does not contain 1 statement for %s. got=%d",
				tt.description, len(file.Statements))
		}

		t.Logf("✓ Successfully parsed: %s", tt.description)
	}
}

// TestWhileExpressionWithComplexConditions tests while loops with complex conditions
func TestWhileExpressionWithComplexConditions(t *testing.T) {
	tests := []struct {
		input       string
		description string
		shouldFail  bool
	}{
		{input: `pamene(running){}`, description: "While with simple identifier", shouldFail: false},
		{input: `pamene(game.isRunning()){}`, description: "While with method call", shouldFail: false},
		{input: `pamene(x > 0 && y < 100){}`, description: "While with logical AND", shouldFail: false},
		{input: `pamene(x > 0 || y > 0){}`, description: "While with logical OR", shouldFail: false},
	}

	for _, tt := range tests {
		p := NewParser()
		file := p.ParseFile([]byte(tt.input))

		if tt.shouldFail {
			if len(p.Errors()) == 0 {
				t.Errorf("Expected parsing to fail for %s: %s", tt.description, tt.input)
			}
			continue
		}

		checkParserErrors(t, p)

		if len(file.Statements) != 1 {
			t.Fatalf("file.Statements does not contain 1 statement for %s. got=%d",
				tt.description, len(file.Statements))
		}

		stmt, ok := file.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("file.Statements[0] is not ast.ExpressionStatement for %s. got=%T",
				tt.description, file.Statements[0])
		}

		_, ok = stmt.Expression.(*ast.WhileExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.WhileExpression for %s. got=%T",
				tt.description, stmt.Expression)
		}

		t.Logf("✓ Successfully parsed: %s", tt.description)
	}
}
