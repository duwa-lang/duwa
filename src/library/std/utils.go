package std

import (
	"path/filepath"
	"strings"

	"github.com/sevenreup/duwa/src/utils/environment"
)

func GetFilePath(path string) (string, error) {
	if !strings.HasSuffix(path, ".duwa") {
		path += ".duwa"
	}
	return filepath.Abs(filepath.Join(environment.CompilerEnvironment.SDKPath, "./src/builtins/std", path))
}
