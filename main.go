package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	var structName string
	var dir string

	// Get the current directory
	curdir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	// Define the flag for the struct name
	flag.StringVar(&structName, "name", "", "Name of the struct to search for")
	flag.StringVar(&dir, "dir", curdir, "Directory to search for the struct")
	flag.Parse()

	if structName == "" {
		return errors.New("please provide a struct name using the -name flag")
	}

	// Walk through all .go files in the directory
	file, typeSpec, err := findTypeInDir(dir, structName)
	if err != nil {
		return fmt.Errorf("find type in directory: %w", err)
	}

	// Ensure is struct type
	if _, ok := typeSpec.Type.(*ast.StructType); !ok {
		return fmt.Errorf("type '%s' is not a struct", structName)
	}

	// Create an *ast.File to hold the generated code
	outfile := &ast.File{
		Name: file.Name,
	}

	outname := fmt.Sprintf("%s_json.go", strings.ToLower(structName))
	fset := token.NewFileSet()
	fset.AddFile(outname, fset.Base(), 0)

	// Generate imports
	generateImports(outfile, []string{
		"bytes",
		"github.com/go-json-experiment/json/jsontext",
		"github.com/paskozdilar/go-gen-json/jsonutil",
	})

	// Generate field aliases
	aliasSpecs := &ast.GenDecl{Tok: token.TYPE}
	outfile.Decls = append(outfile.Decls, aliasSpecs)
	prefix := fmt.Sprintf("field_%s", structName)

	aliasMap := make(map[string]ast.Expr)
	if err := generateAliases(aliasSpecs, typeSpec, prefix, aliasMap); err != nil {
		return fmt.Errorf("generate field aliases: %w", err)
	}
	// for k, v := range aliasMap {
	// 	fmt.Printf("aliasMap[%s] = %s\n", k, v)
	// }

	// Generate top-level method
	generateTopLevelMethod(outfile, typeSpec, "ParseJSON")

	// Generate per-alias functions
	generateSubMethods(outfile, typeSpec, "parseJSON", aliasMap)

	// Print outfile AST into a file
	f, err := os.Create(filepath.Join(dir, outname))
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()
	if err := format.Node(f, fset, outfile); err != nil {
		f.Close()
		os.Remove(f.Name())
		return fmt.Errorf("print file: %w", err)
	}

	fmt.Printf("Generated file: %s\n", f.Name())
	return nil
}

func findTypeInDir(dir, structName string) (*ast.File, *ast.TypeSpec, error) {
	// Read files in the directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("read directory: %s", err)
	}

	for _, entry := range entries {
		// Skip non-go files
		if !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}
		path := filepath.Join(dir, entry.Name())

		// Find the struct type in the file
		if file, typeSpec, ok := findTypeInFile(path, structName); ok {
			return file, typeSpec, nil
		}
	}

	return nil, nil, fmt.Errorf("type %s not found in %s", structName, dir)
}

func findTypeInFile(path, structName string) (*ast.File, *ast.TypeSpec, bool) {
	// Parse file
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
	if err != nil {
		return nil, nil, false
	}

	// Find the struct type in top-level declarations
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, s := range gd.Specs {
			ts, ok := s.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if ts.Name.Name == structName {
				return file, ts, true
			}
		}
	}

	return nil, nil, false
}

func generateImports(outfile *ast.File, paths []string) {
	decl := &ast.GenDecl{
		Tok: token.IMPORT,
	}
	for _, path := range paths {
		decl.Specs = append(decl.Specs, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("%q", path),
			},
		})
	}
	outfile.Decls = append(outfile.Decls, decl)
}

func generateAliases(aliasSpecs *ast.GenDecl, typeSpec *ast.TypeSpec, prefix string, aliasMap map[string]ast.Expr) error {
	// Generate top level alias
	aliasSpecs.Specs = append(aliasSpecs.Specs, &ast.TypeSpec{
		Name:   ast.NewIdent(prefix),
		Type:   typeSpec.Name,
		Assign: 1, // any non-zero value will cause the type to be aliased
	})

	// Generate field aliases
	fields := typeSpec.Type.(*ast.StructType).Fields.List
	for _, field := range fields {
		// Skip unexported fields
		if !field.Names[0].IsExported() {
			continue
		}
		if err := generateFieldAlias(aliasSpecs, field, prefix, aliasMap); err != nil {
			return err
		}
	}
	return nil
}

func generateFieldAlias(aliasSpecs *ast.GenDecl, field *ast.Field, prefix string, aliasMap map[string]ast.Expr) error {
	// Use the json tag if available, otherwise use the field name
	var tag string
	if field.Tag != nil {
		tag = reflect.StructTag(strings.Trim(field.Tag.Value, "`")).Get("json")
	}
	if tag == "" {
		tag = field.Names[0].Name
	}

	// Add field name to prefix
	name := fmt.Sprintf("%s_%s", prefix, field.Names[0].Name)
	aliasMap[name] = field.Type

	// We can remove some of the below cases, since struct type can only
	// be one of:
	switch t := field.Type.(type) {
	case *ast.Ident:
		// First we check the basic type, like int or string
		aliasSpecs.Specs = append(aliasSpecs.Specs, &ast.TypeSpec{
			Name:   ast.NewIdent(name),
			Type:   t,
			Assign: 1, // any non-zero value will cause the type to be aliased
		})
	case *ast.ArrayType:
		// Then we check the array type
		aliasSpecs.Specs = append(aliasSpecs.Specs, &ast.TypeSpec{
			Name:   ast.NewIdent(name),
			Type:   t,
			Assign: 1, // any non-zero value will cause the type to be aliased
		})
	default:
		fmt.Printf("skipping field %s of type %s\n", name, field.Type)
		// case *ast.StarExpr:
		// Then we check the pointer type
		// generatePointerAlias(aliasSpecs, t, name)

		// aliasSpecs.Specs = append(aliasSpecs.Specs, &ast.TypeSpec{
		// 	Name: ast.NewIdent(name),
		// 	Type: t,
		// })

		// case *ast.Ident:
		// case *ast.StructType:
		// case *ast.ArrayType:
		// case *ast.BadExpr:
		// case *ast.Ellipsis:
		// case *ast.BasicLit:
		// case *ast.FuncLit:
		// case *ast.CompositeLit:
		// case *ast.ParenExpr:
		// case *ast.SelectorExpr:
		// case *ast.IndexExpr:
		// case *ast.IndexListExpr:
		// case *ast.SliceExpr:
		// case *ast.TypeAssertExpr:
		// case *ast.CallExpr:
		// case *ast.StarExpr:
		// case *ast.UnaryExpr:
		// case *ast.BinaryExpr:
		// case *ast.KeyValueExpr:
		// case *ast.ArrayType:
		// case *ast.StructType:
		// case *ast.FuncType:
		// case *ast.InterfaceType:
		// case *ast.MapType:
		// case *ast.ChanType:
	}
	return nil
}

func generateTopLevelMethod(outfile *ast.File, typeSpec *ast.TypeSpec, name string) {
	// FuncDecl declares both functions and methods
	decl := &ast.FuncDecl{}

	// Add the method body to the method declaration at the end
	defer func() {
		outfile.Decls = append(outfile.Decls, decl)
	}()

	// Method name is a simple Identifier
	decl.Name = ast.NewIdent(name)

	// Method receives a pointer to the struct type, named "target"
	decl.Recv = &ast.FieldList{List: []*ast.Field{{
		Names: []*ast.Ident{ast.NewIdent("target")},
		Type: &ast.StarExpr{
			X:    &ast.Ident{Name: typeSpec.Name.Name},
			Star: 1,
		},
	}}}

	// Method type is defined by params and results
	decl.Type = &ast.FuncType{}

	// Method has a single parameter, a byte slice named "data"
	decl.Type.Params = &ast.FieldList{List: []*ast.Field{{
		Names: []*ast.Ident{ast.NewIdent("b")},
		Type: &ast.ArrayType{
			Elt: &ast.Ident{Name: "byte"},
		},
	}}}

	// Method returns a single unnamed value of type error
	decl.Type.Results = &ast.FieldList{List: []*ast.Field{{
		Type: &ast.Ident{Name: "error"},
	}}}

	// Method body contains a block statement with a list of sub-statements
	decl.Body = &ast.BlockStmt{}

	// The substatements are the following:

	// d := jsontext.NewDecoder(bytes.NewReader(b))
	decl.Body.List = append(decl.Body.List, &ast.AssignStmt{
		Lhs: []ast.Expr{ast.NewIdent("d")},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{&ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "jsontext"},
				Sel: ast.NewIdent("NewDecoder"),
			},
			Args: []ast.Expr{&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "bytes"},
					Sel: ast.NewIdent("NewReader"),
				},
				Args: []ast.Expr{ast.NewIdent("b")},
			}},
		}},
	})

	// return parseJSON_<structName>(target, d)
	decl.Body.List = append(decl.Body.List, &ast.ReturnStmt{
		Results: []ast.Expr{&ast.CallExpr{
			Fun:  &ast.Ident{Name: fmt.Sprintf("parseJSON_%s", typeSpec.Name.Name)},
			Args: []ast.Expr{ast.NewIdent("target"), ast.NewIdent("d")},
		}},
	})
}

func generateSubMethods(outfile *ast.File, typeSpec *ast.TypeSpec, prefix string, aliasMap map[string]ast.Expr) error {
	// decl := &ast.FuncDecl{}
	// name := fmt.Sprintf("%s_%s", prefix, typeSpec.Name.Name)
	// decl.Name = ast.NewIdent(name)
	return nil
}
