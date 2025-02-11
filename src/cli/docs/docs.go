package docs

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ArgumentData struct {
	Name string
	Type string
}

type FunctionInfo struct {
	Name      string
	Arguments []ArgumentData
	RetunType string `json:"returns"`
	Doc       string
}

type StructInfo struct {
	Name         string
	Alternatives []string
	Doc          string
	Methods      []FunctionInfo
	SourceFile   string // Track which file this type came from
}

type LibraryInfo struct {
	Name       string
	Doc        string
	Methods    []FunctionInfo
	SourceFile string
}

type BuiltinInfo struct {
	Functions  []FunctionInfo
	SourceFile string
}

type ParserConfig struct {
	RootDir             string
	ExcludePatterns     []string
	IncludePatterns     []string
	MinMethods          int
	IncludeUndocumented bool
	OutputDir           string
	GenerateFormat      string
}

type Parser struct {
	config    ParserConfig
	types     []StructInfo
	libraries []LibraryInfo
	builtins  []BuiltinInfo
}

func NewParser(config ParserConfig) *Parser {
	return &Parser{
		config:    config,
		types:     make([]StructInfo, 0),
		libraries: make([]LibraryInfo, 0),
		builtins:  make([]BuiltinInfo, 0),
	}
}

func parseLibraryLine(line string) (string, error) {
	libraryRegex := regexp.MustCompile(`library=(\w+)`)
	libraryMatch := libraryRegex.FindStringSubmatch(line)
	if len(libraryMatch) < 2 {
		return "", fmt.Errorf("library name not found")
	}
	return libraryMatch[1], nil
}

// isLibraryMap checks if a map declaration matches the library function signature
func isLibraryMap(expr ast.Expr) bool {
	compositeLit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return false
	}
	mapType, ok := compositeLit.Type.(*ast.MapType)
	if !ok {
		return false
	}

	// Check key type is string
	if _, ok := mapType.Key.(*ast.Ident); !ok || mapType.Key.(*ast.Ident).Name != "string" {
		return false
	}

	// Check value type is *object.LibraryFunction
	starExpr, ok := mapType.Value.(*ast.StarExpr)
	if !ok {
		return false
	}

	selector, ok := starExpr.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	return selector.Sel.Name == "LibraryFunction"
}

// parseBuiltinFunction parses a builtin function declaration
func parseBuiltinFunction(funcDoc *doc.Func) (FunctionInfo, error) {
	if funcDoc.Doc == "" {
		return FunctionInfo{}, fmt.Errorf("no documentation")
	}

	docLine, documentation := parseDocumentation(funcDoc.Doc)
	if !strings.Contains(docLine, "type=builtin-func") {
		return FunctionInfo{}, fmt.Errorf("not a builtin function")
	}

	name, args, returnType, err := parseMethodLine(docLine)
	if err != nil {
		return FunctionInfo{}, err
	}

	return FunctionInfo{
		Name:      name,
		Arguments: args,
		RetunType: returnType,
		Doc:       documentation,
	}, nil
}

// analyzeLibraries processes a single file and returns found libraries
func (p *Parser) analyzeLibraries(filename string, file *ast.File) ([]LibraryInfo, error) {
	var libraries []LibraryInfo

	ast.Inspect(file, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok {
			if genDecl.Doc == nil {
				return true
			}

			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for i, val := range valueSpec.Values {
						if isLibraryMap(val) {
							docText := genDecl.Doc.Text()
							firstLine, documentation := parseDocumentation(docText)
							libraryName, err := parseLibraryLine(firstLine)
							if err != nil {
								continue
							}

							library := LibraryInfo{
								Name:       libraryName,
								Doc:        documentation,
								Methods:    []FunctionInfo{},
								SourceFile: filename,
							}

							// Get the variable name for this library
							if i < len(valueSpec.Names) {
								libraryVarName := valueSpec.Names[i].Name
								// Now find all methods associated with this library
								methods := p.findLibraryMethods(file, libraryVarName)
								library.Methods = methods
							}

							libraries = append(libraries, library)
						}
					}
				}
			}
		}
		return true
	})

	return libraries, nil
}

// findLibraryMethods finds all functions that belong to a library
func (p *Parser) findLibraryMethods(file *ast.File, libraryVarName string) []FunctionInfo {
	var methods []FunctionInfo

	ast.Inspect(file, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			if strings.HasPrefix(funcDecl.Name.Name, "method") && funcDecl.Doc != nil {
				methodLine, documentation := parseDocumentation(funcDecl.Doc.Text())
				name, args, returnType, err := parseMethodLine(methodLine)
				if err == nil {
					functionInfo := FunctionInfo{
						Name:      name,
						Arguments: args,
						RetunType: returnType,
						Doc:       documentation,
					}
					methods = append(methods, functionInfo)
				}
			}
		}
		return true
	})

	return methods
}

// analyzeFile processes a single file and returns found types and libraries
func (p *Parser) analyzeFile(filename string) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("error parsing file %s: %v", filename, err)
	}

	// Parse libraries first
	libraries, err := p.analyzeLibraries(filename, file)
	if err != nil {
		return err
	}
	p.libraries = append(p.libraries, libraries...)

	// Parse types using existing code
	pkg := &ast.Package{
		Name:  file.Name.Name,
		Files: map[string]*ast.File{filename: file},
	}

	docPkg := doc.New(pkg, ".", doc.AllDecls)
	for _, t := range docPkg.Types {
		if t.Doc == "" && !p.config.IncludeUndocumented {
			continue
		}

		if t.Doc != "" {
			typeLine, documentation := parseDocumentation(t.Doc)
			typeName, alternatives, err := parseStructLine(typeLine)
			if err != nil {
				continue
			}

			structInfo := StructInfo{
				Name:         typeName,
				Alternatives: alternatives,
				Doc:          documentation,
				Methods:      []FunctionInfo{},
				SourceFile:   filename,
			}

			for _, method := range t.Methods {
				if method.Doc != "" {
					methodLine, documentation := parseDocumentation(method.Doc)
					name, args, returnType, err := parseMethodLine(methodLine)
					if err == nil {
						functionInfo := FunctionInfo{
							Name:      name,
							Arguments: args,
							RetunType: returnType,
							Doc:       documentation,
						}
						structInfo.Methods = append(structInfo.Methods, functionInfo)
					}
				}
			}

			if len(structInfo.Methods) >= p.config.MinMethods {
				p.types = append(p.types, structInfo)
			}
		}
	}

	var currentBuiltins []FunctionInfo
	for _, fn := range docPkg.Funcs {
		if strings.HasPrefix(fn.Name, "BuiltIn") {
			if builtinFunc, err := parseBuiltinFunction(fn); err == nil {
				currentBuiltins = append(currentBuiltins, builtinFunc)
			}
		}
	}

	if len(currentBuiltins) > 0 {
		p.builtins = append(p.builtins, BuiltinInfo{
			Functions:  currentBuiltins,
			SourceFile: filename,
		})
	}

	return nil
}

// Parse scans the directory and processes all valid files
func (p *Parser) Parse() error {
	return filepath.Walk(p.config.RootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			if !p.shouldProcessFile(path) {
				return nil
			}

			err := p.analyzeFile(path)
			if err != nil {
				fmt.Printf("Warning: Error processing file %s: %v\n", path, err)
			}
		}

		return nil
	})
}

// GetTypes returns all found types
func (p *Parser) GetTypes() []StructInfo {
	return p.types
}

// GetLibraries returns all found libraries
func (p *Parser) GetLibraries() []LibraryInfo {
	return p.libraries
}

func (p *Parser) GetBuiltins() []BuiltinInfo {
	return p.builtins
}

// shouldProcessFile checks if a file should be processed based on include/exclude patterns
func (p *Parser) shouldProcessFile(filename string) bool {
	// Check exclude patterns
	for _, pattern := range p.config.ExcludePatterns {
		matched, err := filepath.Match(pattern, filepath.Base(filename))
		if err == nil && matched {
			return false
		}
	}

	// If include patterns are specified, file must match at least one
	if len(p.config.IncludePatterns) > 0 {
		for _, pattern := range p.config.IncludePatterns {
			matched, err := filepath.Match(pattern, filepath.Base(filename))
			if err == nil && matched {
				return true
			}
		}
		return false
	}

	return true
}

// parseStructLine parses the first line of struct documentation
func parseStructLine(line string) (string, []string, error) {
	typeRegex := regexp.MustCompile(`type=(\w+)`)
	typeMatch := typeRegex.FindStringSubmatch(line)
	if len(typeMatch) < 2 {
		return "", nil, fmt.Errorf("type name not found")
	}
	typeName := typeMatch[1]

	altRegex := regexp.MustCompile(`alternative=(\w+)`)
	altMatch := altRegex.FindStringSubmatch(line)
	var alternatives []string
	if len(altMatch) >= 2 {
		alternatives = strings.Split(altMatch[1], ",")
	}

	return typeName, alternatives, nil
}

// parseMethodLine parses the method signature line
func parseMethodLine(line string) (string, []ArgumentData, string, error) {
	methodRegex := regexp.MustCompile(`method=(\w+)`)
	methodMatch := methodRegex.FindStringSubmatch(line)
	if len(methodMatch) < 2 {
		return "", nil, "", fmt.Errorf("method name not found")
	}
	methodName := methodMatch[1]

	var args []ArgumentData
	argsRegex := regexp.MustCompile(`args=\[(.*?)\]`)
	argsMatch := argsRegex.FindStringSubmatch(line)
	if len(argsMatch) >= 2 && argsMatch[1] != "" {
		argParts := strings.Split(argsMatch[1], ",")
		for _, arg := range argParts {
			argRegex := regexp.MustCompile(`(\w+)\{(\w+)\}`)
			argMatch := argRegex.FindStringSubmatch(strings.TrimSpace(arg))
			if len(argMatch) >= 3 {
				args = append(args, ArgumentData{
					Type: argMatch[1],
					Name: argMatch[2],
				})
			}
		}
	}

	returnRegex := regexp.MustCompile(`return=\{(.*?)\}`)
	returnMatch := returnRegex.FindStringSubmatch(line)
	if len(returnMatch) < 2 {
		return "", nil, "", fmt.Errorf("return type not found")
	}
	returnType := returnMatch[1]

	return methodName, args, returnType, nil
}

func parseDocumentation(docText string) (string, string) {
	lines := strings.Split(strings.TrimSpace(docText), "\n")
	if len(lines) == 0 {
		return "", ""
	}
	return lines[0], strings.TrimSpace(strings.Join(lines[1:], "\n"))
}

func GenerateDocs() {
	// Example usage
	fmt.Println("Generating documentation...")
	config := ParserConfig{
		RootDir:             "./src",
		ExcludePatterns:     []string{"*_test.go", "*.generated.go"},
		IncludePatterns:     []string{"*.go"},
		MinMethods:          0,
		IncludeUndocumented: true,
		OutputDir:           "../duwa-site/content/generated",
		GenerateFormat:      "json",
	}

	parser := NewParser(config)
	err := parser.Parse()
	if err != nil {
		fmt.Printf("Error parsing directory: %v\n", err)
		return
	}

	if config.GenerateFormat == "json" {
		parser.generateJson()
	} else {
		parser.generateMarkdown()
	}
	fmt.Println("Documentation generated successfully!")
}
