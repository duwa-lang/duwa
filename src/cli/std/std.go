package std

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/duwa-lang/duwa/src/cli/utils"
)

const (
	builtinsTemplate = `package std

var BuiltinFiles = []string{
	{{- range .Files }}
	"{{ . }}",
	{{- end }}
}
`
)

type StdGeneratorInfo struct {
	Files []string
}

func GenerateStdTypes() {
	fmt.Println("Generating std types")
	var files []string
	root := "./src/builtins/std"

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".duwa" {
			relativePath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			// Remove the file extension
			relativePath = relativePath[:len(relativePath)-len(filepath.Ext(relativePath))]
			files = append(files, relativePath)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	filename := filepath.Join("./src/library/std", "includes.go")

	utils.WriteFileFromTemplate(filename, builtinsTemplate, StdGeneratorInfo{Files: files})
}
