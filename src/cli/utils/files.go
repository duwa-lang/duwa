package utils

import (
	"os"
	"text/template"
)

func WriteFileFromTemplate(filePath string, textTemplate string, data any) error {
	builtinsTempl := template.Must(template.New("builtins").Parse(textTemplate))
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return builtinsTempl.Execute(file, data)
}
