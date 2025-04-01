package utils

import (
	"testing"

	"github.com/sevenreup/duwa/src/object"
	"github.com/sevenreup/duwa/src/values"
	"github.com/shopspring/decimal"
)

func TestIntegerObject(t *testing.T, obj object.Object, expected decimal.Decimal) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%s (%+s)", obj.String(), obj.String())
		return false
	}
	if !result.Value.Equal(expected) {
		t.Errorf("object has wrong value. got=%s, want=%s",
			result.Value.String(), expected.String())
		return false
	}
	return true
}

func TestStringObject(t *testing.T, obj object.Object, expected string) bool {
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

func TestBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
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

func TestLiteralExpression(
	t *testing.T,
	obj object.Object,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return TestIntegerObject(t, obj, decimal.NewFromInt(int64(v)))
	case int64:
		return TestIntegerObject(t, obj, decimal.NewFromInt(v))
	case bool:
		return TestBooleanObject(t, obj, v)
	case string:
		return TestStringObject(t, obj, v)
	}
	t.Errorf("type of exp not handled. got=%T", expected)
	return false
}

func TestNullObject(t *testing.T, obj object.Object) bool {
	if obj != values.NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func IsErrorObject(t *testing.T, obj object.Object, expected string) bool {
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

func IsNumberObject(t *testing.T, obj object.Object, expected int64) bool {
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
