package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	// Define the flag for the struct name
	structName := flag.String("name", "", "Name of the struct to search for")
	flag.Parse()

	if *structName == "" {
		log.Fatal("Please provide a struct name using the -name flag")
	}

	// Get the current directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Walk through all .go files in the directory
	found := false
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".go" {
			// Parse the .go file
			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return err
			}

			// Search for the struct with the specified name
			for _, decl := range node.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok {
					for _, spec := range genDecl.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							if structType, ok := typeSpec.Type.(*ast.StructType); ok {
								if typeSpec.Name.Name == *structName {
									found = true
									generateParser(typeSpec, structType)
								}
							}
						}
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	if !found {
		log.Fatalf("Struct with name %s not found", *structName)
	}
}

func generateParser(typeSpec *ast.TypeSpec, structType *ast.StructType) {
	filename := fmt.Sprintf("%s_jsonparse.go", strings.ToLower(typeSpec.Name.Name))
	fmt.Println("Generating parser for struct", typeSpec.Name.Name)
	fmt.Println("Generated file:", filename)

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Write the package name
	fmt.Fprintf(f, "package main\n\n")

	// Write the imports
	// TODO:
	// Change to encoding/json/jsontext when proposal is accepted:
	// - https://github.com/golang/go/issues/71497
	fmt.Fprintf(f, "import (\n")
	fmt.Fprintf(f, "\t\"github.com/go-json-experiment/json/jsontext\"")
	fmt.Fprintf(f, ")\n")
}

var fnTemplate = template.Must(template.New("fn").Parse(`
func (s *{{.StructName}}) ParseJSON(data []byte) error {
}
`))
