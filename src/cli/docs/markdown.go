package docs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sevenreup/duwa/src/cli/utils"
)

// Markdown generation templates
const (
	typeTemplate = `# {{.Name}}
{{if .Alternatives}}
**Alternatives:** {{range .Alternatives}}{{.}} {{end}}
{{end}}

{{.Doc}}

## Methods

{{range .Methods}}
### {{.Name}}

{{.Doc}}

**Arguments:**
{{range .Arguments}}
- {{.Name}}: {{.Type}}
{{end}}

**Returns:** {{.RetunType}}

{{end}}
`

	libraryTemplate = `# {{.Name}} Library

{{.Doc}}

## Methods

{{range .Methods}}
### {{.Name}}

{{.Doc}}

**Arguments:**
{{range .Arguments}}
- {{.Name}}: {{.Type}}
{{end}}

**Returns:** {{.RetunType}}

{{end}}
`

	builtinsTemplate = `# Built-in Functions
{{range .}}
{{range .Functions}}
## {{.Name}}

{{.Doc}}

**Arguments:**
{{range .Arguments}}
- {{.Name}}: {{.Type}}
{{end}}

**Returns:** {{.RetunType}}

{{end}}
{{end}}
`
)

func (p *Parser) generateMarkdown() error {
	if err := os.MkdirAll(p.config.OutputDir, 0755); err != nil {
		return err
	}

	// Generate documentation for types
	typeTempl := template.Must(template.New("type").Parse(typeTemplate))
	for _, t := range p.types {
		filename := filepath.Join(p.config.OutputDir, fmt.Sprintf("type_%s.md", strings.ToLower(t.Name)))
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		if err := typeTempl.Execute(file, t); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}

	// Generate documentation for libraries
	libTempl := template.Must(template.New("library").Parse(libraryTemplate))
	for _, lib := range p.libraries {
		filename := filepath.Join(p.config.OutputDir, fmt.Sprintf("library_%s.md", strings.ToLower(lib.Name)))
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		if err := libTempl.Execute(file, lib); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}

	// Generate documentation for built-in functions
	filename := filepath.Join(p.config.OutputDir, "builtins.md")
	return utils.WriteFileFromTemplate(filename, builtinsTemplate, p.builtins)
}
