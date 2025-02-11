package evaluator

import (
	"path/filepath"
	"testing"

	"github.com/sevenreup/duwa/src/lexer"
	"github.com/sevenreup/duwa/src/object"
	"github.com/sevenreup/duwa/src/parser"
	"github.com/sevenreup/duwa/src/utils/environment"
	"github.com/sevenreup/duwa/src/values"
	"github.com/shopspring/decimal"
)

func testEval(input string) object.Object {
	l := lexer.New([]byte(input))
	p := parser.New(l)
	file := p.ParseFile()
	env := object.NewDefaultEnvironment()

	evaluatorInstance := Eval
	filename, _ := filepath.Abs("../../")
	environment.SetCompilationSettings(filename)
	object.RegisterEvaluator(evaluatorInstance)

	return Eval(file, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected decimal.Decimal) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if !result.Value.Equal(expected) {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%s, want=%s",
			result.Value, expected)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func testLiteralExpression(
	t *testing.T,
	obj object.Object,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerObject(t, obj, decimal.NewFromInt(int64(v)))
	case int64:
		return testIntegerObject(t, obj, decimal.NewFromInt(v))
	case bool:
		return testBooleanObject(t, obj, v)
	case string:
		return testStringObject(t, obj, v)
	}
	t.Errorf("type of exp not handled. got=%T", expected)
	return false
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != values.NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func isErrorObject(t *testing.T, obj object.Object, expected string) bool {
	err, ok := obj.(*object.Error)

	if !ok {
		t.Errorf("object is not Error. got=%T (%+v", obj, obj)
		return false
	}

	if err.Message != expected {
		t.Errorf("error has wrong message. got=%s, expected=%s", err.Message, expected)
		return false
	}

	return true
}

func isNumberObject(t *testing.T, obj object.Object, expected int64) bool {
	number, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("object is not Number. got=%T (%+v", obj, obj)
		return false
	}

	if number.Value.IntPart() != expected {
		t.Errorf("object has wrong value. got=%d, expected=%d", number.Value.IntPart(), expected)
		return false
	}

	return true
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"nambala n = 2; n++;", 3},
		{"nambala n = 2; n--;", 1},
		{"5 + 5 + 5 + 5- 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 +-50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 *-10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 +-10", 50},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"nambala x = 5; x += 1; x", 6},
		{"nambala x = 5; x -= 1; x", 4},
		{"nambala x = 5; x *= 2; x", 10},
		{"nambala x = 10; x /= 2; x", 5},
		{"nambala x = 0; x++; x", 1},
		{"nambala x = 6; x--; x", 5},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, decimal.NewFromInt(tt.expected))
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"zoona", true},
		{"bodza", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"1 >= 2", false},
		{"1 <= 2", true},
		{"zoona == zoona", true},
		{"bodza == bodza", true},
		{"zoona == bodza", false},
		{"zoona != bodza", true},
		{"bodza != zoona", true},
		{"zoona && zoona", true},
		{"bodza && bodza", false},
		{"zoona && bodza", false},
		{"zoona || zoona", true},
		{"bodza || bodza", false},
		{"zoona || bodza", true},
		{"(1 < 2) == zoona", true},
		{"(1 < 2) == bodza", false},
		{"(1 > 2) == zoona", false},
		{"(1 > 2) == bodza", true},
		{`("foo" == "bar") == bodza`, true},
		{`("foo" == "foo") == zoona`, true},
		{`0 <= 14 && 10 >= 10 && 18 <= 47`, true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!zoona", false},
		{"!bodza", true},
		{"!5", false},
		{"!!zoona", true},
		{"!!bodza", false},
		{"!!5", true}}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"ngati (zoona) { 10 }", 10},
		{"ngati (bodza) { 10 }", nil},
		{"ngati (1) { 10 }", 10},
		{"ngati (1 < 2) { 10 }", 10},
		{"ngati (1 > 2) { 10 }", nil},
		{"ngati (1 > 2) { 10 } kapena { 20 }", 20},
		{"ngati (1 < 2) { 10 } kapena { 20 }", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, decimal.NewFromInt(int64(integer)))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

// TODO: Test null return value
func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"bweza 10;", 10},
		{"bweza 10; 9;", 10},
		{"bweza 2 * 5; 9;", 10},
		{"9; bweza 2 * 5; 9;", 10},
		{`ndondomeko five() { bweza 5; };five();`, 5},
		{`
			ndondomeko five(n) { 
				ngati (n > 1) {
					lemba(n);
					bweza 5;
				}
				lemba(10);
				bweza 10; 
			};
			five(3);`,
			5,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, decimal.NewFromInt(tt.expected))
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + zoona;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + zoona; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-zoona",
			"unknown operator:-BOOLEAN",
		},
		{
			"zoona + bodza;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; zoona + bodza; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"ngati (10 > 1) { zoona + bodza; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
					ngati (10 > 1) {
 						ngati (10 > 1) {
 							bweza zoona + bodza;
 						}
 						bweza 1;
 					}
				`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello"- "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			`{"dzina": "Maliko"}[ndondomeko d(x) { x }];`,
			"1:20:: runtime error: unusable as map key: FUNCTION",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestAssignmentStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"nambala a = 5; a;", 5},
		{"nambala a = 5 * 5; a;", 25},
		{"nambala a = 5; nambala b = a; b;", 5},
		{"nambala a = 5; nambala b = a; nambala c = a + b + 5; c;", 15},
		{`mawu a = "b"; a;`, "b"},
		{`mawu a = "5"; mawu b = a; mawu c = a + b + "5"; c;`, "555"},
	}
	for _, tt := range tests {
		testLiteralExpression(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "ndondomeko phatikizaZiwiri(x) { x + 2; };"
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}
	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"ndondomeko identity(x) { x; }; identity(5);", 5},
		{"ndondomeko  identity(x) { bweza x; }; identity(5);", 5},
		{"ndondomeko double(x){ x * 2; }; double(5);", 10},
		{"ndondomeko  add(x, y) { x + y; }; add(5, 5);", 10},
		{"ndondomeko add(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"ndondomeko zina(x) { x; }(5)", 5},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), decimal.NewFromInt(tt.expected))
	}
}

func TestWhileExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`pamene (bodza) { }`, nil},
		{`nambala n = 0; pamene (n < 10) { n = n + 1 }; n`, 10},
		{`nambala n = 0; pamene (n < 10) { n++ }; n`, 10},
		{"nambala n = 10; pamene (n > 0) { n = n - 1 }; n", 0},
		{"nambala n = 10; pamene (n > 0) { n-- }; n", 0},
		{"nambala n = 0; pamene (n < 10) { n = n + 1 }", nil},
		{"nambala n = 10; pamene (n > 0) { n = n - 1 }", nil},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		number, ok := tt.expected.(int)

		if ok {
			isNumberObject(t, result, int64(number))
		} else {
			if result != nil {
				t.Errorf("object is not NULL. got=%T (nil)", number)
			}
		}
	}
}

func TestForExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`x = 10; za (x = y; x > 0; x = x - 1) { x }`, "identifier not found: y"},
		{`za (x = 0; x < 10; x = x + 1) { y }`, "identifier not found: y"},
		{`bar = zoona; za (x = 0; x < 10; x = x + 1) { y; print(bar) }`, "identifier not found: y"},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		number, ok := tt.expected.(int64)

		if ok {
			isNumberObject(t, result, number)
		} else {
			isErrorObject(t, result, tt.expected.(string))
		}
	}
}

func TestClosures(t *testing.T) {
	input := `
 	ndondomeko newAdder(x) {
 		ndondomeko temp(y) { x + y };
 	};
 	nambala addTwo = newAdder(2);
 	addTwo(2);`
	testIntegerObject(t, testEval(input), decimal.NewFromInt(4))
}

func TestStrings(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			input:    `"Hello" + " " + "World!"`,
			expected: "Hello World!",
		},
		{
			input:    `"Hello" + " " + "World" + "!"`,
			expected: `Hello World!`,
		},
		{
			input:    `"Hello" + 3`,
			expected: "Hello3",
		},
		{
			input:    `"Hello" + 3 + " " + "World" + "!"`,
			expected: "Hello3 World!",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		str, ok := evaluated.(*object.String)
		if !ok {
			t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
		}
		if str.Value != tt.expected {
			t.Errorf("String has wrong value. got=%q (%+v)", str.Value, tt.expected)
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}
	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}
	testIntegerObject(t, result.Elements[0], decimal.NewFromInt(1))
	testIntegerObject(t, result.Elements[1], decimal.NewFromInt(4))
	testIntegerObject(t, result.Elements[2], decimal.NewFromInt(6))
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"nambala i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"nambala[] myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"nambala[] myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"nambala[] myArray = [1, 2, 3]; nambala i = myArray[0]; myArray[i]",
			2,
		},
		{
			"nambala[] myArray = [1, 2, 3]; myArray[0] = 5; myArray[0]",
			5,
		},
		{
			"nambala[] myArray = [1, 2, 3]; nambala temp = myArray[0]; myArray[0] = 5; myArray[0] = temp; myArray[0]",
			1,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, decimal.NewFromInt(int64(integer)))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

// TODO: Test Array errors

func TestMethodCalls(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"nambala[] myArray = [1, 2, 3];myArray.kutalika();",
			3,
		},
		{
			"nambala[] myArray = [1, 2, 3];myArray.chotsaKumbuyo();myArray[1];",
			2,
		},
		{
			"nambala[] myArray = [1, 2, 3];myArray.chotsaKutsogolo();myArray[0];",
			2,
		},
		{
			"nambala[] myArray = [1, 2, 3];myArray.Kankha(8);myArray[3];",
			8,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, decimal.NewFromInt(int64(integer)))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `mawu two = "two";
	{
		"one": 10- 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		zoona: 5,
		bodza: 6
	}`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Map)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}
	expected := map[object.MapKey]int64{
		(&object.String{Value: "one"}).MapKey():                  1,
		(&object.String{Value: "two"}).MapKey():                  2,
		(&object.String{Value: "three"}).MapKey():                3,
		(&object.Integer{Value: decimal.NewFromInt(4)}).MapKey(): 4,
		values.TRUE.MapKey():                                     5,
		values.FALSE.MapKey():                                    6,
	}
	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}
	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}
		testLiteralExpression(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`mawu key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{zoona: 5}[zoona]`,
			5,
		},
		{
			`{bodza: 5}[bodza]`,
			5,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testLiteralExpression(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestClasses(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`
			kalasi Munthu {
				ndondomeko zaka() { 
					bweza 10;
				}
			}
			Munthu maliko = Munthu();
			maliko.zaka();
			`,
			10,
		},
		{
			`
			kalasi Munthu {
				numbala zaka = 2;
				ndondomeko yikaZaka() { 
					zaka = 10;
				}
			}
			Munthu maliko = Munthu();
			maliko.yikaZaka();
			maliko.zaka;
			`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, decimal.NewFromInt(int64(integer)))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestImport(t *testing.T) {
	tests := []struct {
		input string
	}{
		{`tenga fmt kuchokera "fmt";`},
		{`tenga { export1, export2 } kuchokera "fmt";`},
	}
	for _, tt := range tests {
		testEval(tt.input)
	}
}
