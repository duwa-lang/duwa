package evaluator

import (
	"testing"

	"github.com/duwa-lang/duwa/src/object"
)

func TestClassInstanceIsolation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Basic instance property isolation",
			input: `
				kalasi Droplet {
					nambala x = 0;
					nambala y = 0;
				}

				droplet1 = Droplet();
				droplet2 = Droplet();

				droplet1.x = 10;
				droplet1.y = 20;

				droplet2.x + droplet2.y
			`,
			expected: "0", // droplet2 should still have x=0, y=0
		},
		{
			name: "Multiple instances with different values",
			input: `
				kalasi Counter {
					nambala count = 0;
				}

				c1 = Counter();
				c2 = Counter();
				c3 = Counter();

				c1.count = 5;
				c2.count = 10;
				c3.count = 15;

				c1.count + c2.count + c3.count
			`,
			expected: "30",
		},
		{
			name: "Instance property modification doesn't affect new instances",
			input: `
				kalasi Box {
					nambala value = 100;
				}

				b1 = Box();
				b1.value = 50;

				b2 = Box();

				b2.value
			`,
			expected: "100", // b2 should have the original default value
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)

			if evaluated == nil {
				t.Fatalf("Eval returned nil")
			}

			integer, ok := evaluated.(*object.Integer)
			if !ok {
				t.Fatalf("Expected Integer, got %T (%+v)", evaluated, evaluated)
			}

			if integer.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, integer.String())
			}
		})
	}
}
