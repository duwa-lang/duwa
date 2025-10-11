package examples

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/duwa-lang/duwa/src/evaluator"
	"github.com/duwa-lang/duwa/src/object"
	"github.com/duwa-lang/duwa/src/parser"
	"github.com/duwa-lang/duwa/src/utils/environment"
)

func testEval(input string) object.Object {
	p := parser.NewParser()
	file := p.ParseFile([]byte(input))
	env := object.NewDefaultEnvironment()

	evaluatorInstance := evaluator.Eval
	filename, _ := filepath.Abs("../../")
	environment.SetCompilationSettings(filename)
	object.RegisterEvaluator(evaluatorInstance)

	return evaluator.Eval(file, env)
}

func TestExamples(t *testing.T) {
	// read all files, ending with .duwa
	files, err := filepath.Glob("*.duwa")
	if err != nil {
		t.Fatalf("Failed to read example files: %v", err)
	}

	var failedFiles []string

	for _, file := range files {
		// read the file content
		content, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("Failed to read file %s: %v", file, err)
			failedFiles = append(failedFiles, file)
			continue
		}

		// skip files containing "khonso.landira"
		if strings.Contains(string(content), "khonso.landira") {
			t.Logf("Skipping file %s (contains khonso.landira)", file)
			continue
		}

		// evaluate the content
		result := testEval(string(content))
		if _, ok := result.(*object.Error); ok {
			failedFiles = append(failedFiles, file)
		}
	}

	if len(failedFiles) > 0 {
		t.Errorf("The following example files failed: %v", failedFiles)
	}
}
