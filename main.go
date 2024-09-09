package main

import (
	"context"
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"log"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"text/template"
	"unicode"

	"golang.org/x/tools/go/packages"
)

func debug(x any) {
	fmt.Printf("%#v\n", x)
}

type Struct struct {
	Name   string
	Fields []*StructField
}

type StructField struct {
	Name string
	Type string
}

type TemplateData struct {
	ConstructorPrefix        string
	PackageName              string
	Structs                  []*Struct
	ImportPaths              []string
	TypeNameToUnderlyingType map[string]string
	ConstructorReturnsError  map[string]bool
}

var regex = regexp.MustCompile(`^(.+)/(\w+\.\w+)$`)

func main() {
	filePath := flag.String("path", "", "Source file path")
	structNames := flag.String("structs", "", "Comma-separated list of struct names to generate")
	prefix := flag.String("prefix", "", "Prefix for constructor function")
	factory := flag.Bool("factory", false, "Generate factory functions")
	flag.Parse()

	if *filePath == "" || *structNames == "" || *prefix == "" {
		log.Fatalf("Usage: vogen -path <FilePath> -structs <StructName1,StructName2,...> -prefix <Prefix>")
	}

	targetStructs := strings.Split(*structNames, ",")

	ctx := context.Background()
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cfg := &packages.Config{
		Context: ctx,
		Mode:    packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Dir:     wd,
	}
	pkgs, err := packages.Load(cfg, *filePath)
	if err != nil {
		panic(err)
	}

	if len(pkgs) == 0 {
		log.Fatalf("no packages found")
	}

	pkg := pkgs[0]
	data := &TemplateData{
		PackageName:              pkg.Name,
		ConstructorPrefix:        *prefix,
		Structs:                  make([]*Struct, 0),
		TypeNameToUnderlyingType: make(map[string]string),
		ConstructorReturnsError:  make(map[string]bool),
	}

	for _, syntax := range pkg.Syntax {
		for _, decl := range syntax.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok {
				for _, spec := range genDecl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						switch t := typeSpec.Type.(type) {
						case *ast.StructType:
							structFields := make([]*StructField, len(t.Fields.List))
							structName := typeSpec.Name.Name

							if !slices.Contains(targetStructs, structName) {
								continue
							}

							for i, field := range t.Fields.List {
								if field.Names == nil || len(field.Names) == 0 {
									panic("field.Names is nil or empty")
								}
								name := field.Names[0].Name
								if ident, ok := field.Type.(*ast.Ident); ok {
									structFields[i] = &StructField{
										Name: name,
										Type: ident.Name,
									}
								} else if selector, ok := field.Type.(*ast.SelectorExpr); ok {
									ty := pkg.TypesInfo.TypeOf(selector)
									typeName, importPath := extractTypeAndImportPath(ty.String())
									structFields[i] = &StructField{
										Name: name,
										Type: typeName,
									}
									if importPath != "" {
										data.ImportPaths = append(data.ImportPaths, importPath)
									}
								}
							}
							data.Structs = append(data.Structs, &Struct{
								Name:   structName,
								Fields: structFields,
							})
						default:
							typeName, importPath := extractTypeAndImportPath(getUnderlyingType(pkg.TypesInfo, t))
							if importPath != "" {
								data.ImportPaths = append(data.ImportPaths, importPath)
							}
							data.TypeNameToUnderlyingType[typeSpec.Name.Name] = typeName
						}
					}
				}
			}
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				data.ConstructorReturnsError[funcDecl.Name.Name] = false

				list := funcDecl.Type.Results.List
				// (type, error)のみ対応
				if funcDecl.Type.Results != nil && len(list) == 2 {
					if ident, ok := list[1].Type.(*ast.Ident); ok && ident.Name == "error" {
						data.ConstructorReturnsError[funcDecl.Name.Name] = true
					}
				}
			}
		}
	}

	generateConstructor(*filePath, data)

	if *factory {
		generateFactory(*filePath, data)
	}
}

func extractTypeAndImportPath(s string) (string, string) {
	if !regex.MatchString(s) {
		return s, ""
	}
	matches := regex.FindStringSubmatch(s)
	qualifiedTypeName := matches[2]
	sp := strings.Split(qualifiedTypeName, ".")
	if len(sp) != 2 {
		return s, ""
	}
	importPath := matches[1] + "/" + sp[0]
	return qualifiedTypeName, importPath
}

func getUnderlyingType(info *types.Info, expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.SelectorExpr:
		ty := info.TypeOf(t)
		return ty.String() // time.Timeのような型はそのまま使う。Underlying()はtime.Timeの構造体の中身を指すため。
	case *ast.StarExpr:
		ty := info.TypeOf(t)
		if pointer, ok := ty.(*types.Pointer); ok {
			if named, ok := pointer.Elem().(*types.Named); ok {
				return "*" + named.Obj().Name()
			}
		}
		return ty.Underlying().String()
	default:
		ty := info.TypeOf(t)
		return ty.Underlying().String()
	}
}

func toCamelCase(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func toPascalCase(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func uniqueImportPaths(data *TemplateData) []string {
	seen := map[string]bool{}
	uniqueImportPaths := make([]string, 0, len(data.ImportPaths))
	for _, path := range data.ImportPaths {
		if _, ok := seen[path]; !ok {
			uniqueImportPaths = append(uniqueImportPaths, path)
			seen[path] = true
		}
	}
	return uniqueImportPaths
}

func generateConstructor(filename string, data *TemplateData) {
	tmpl, err := template.New("constructor").Funcs(template.FuncMap{
		"camelCase":  toCamelCase,
		"pascalCase": toPascalCase,
		"shouldReturnError": func() bool {
			return len(data.ConstructorReturnsError) > 0
		},
		"constructorReturnsError": func(constructorName string) bool {
			return data.ConstructorReturnsError[constructorName]
		},
		"isPointer": func(typ string) bool {
			return strings.HasPrefix(typ, "*")
		},
		"getUnderlyingType": func(typeName string) string {
			underlying, ok := data.TypeNameToUnderlyingType[typeName]
			if !ok {
				return typeName
			}
			return underlying
		},
		"hasConstructor": func(constructorName string) bool {
			_, ok := data.ConstructorReturnsError[constructorName]
			return ok
		},
		"uniqueImportPaths": func() []string {
			return uniqueImportPaths(data)
		},
	}).Parse(`// Code generated by vogen DO NOT EDIT.
package {{.PackageName}}
import (
  {{range uniqueImportPaths}}
	"{{.}}"
  {{end}}
)
{{$prefix := .ConstructorPrefix}}
{{range .Structs}}
  {{$structName := .Name}}
  func {{$prefix}}{{$structName}}({{range $index, $field := .Fields}}{{if $index}}, {{end}}{{$field.Name | camelCase}} {{ getUnderlyingType $field.Type }}{{end}}) {{if shouldReturnError}}(*{{$structName}}, error){{else}}*{{$structName}}{{end}} {
   {{range $index, $field := .Fields}}
     {{$constructorName := printf "%s%s%s" ($prefix | pascalCase) ($structName | pascalCase) ($field.Name | pascalCase)}}
     {{if hasConstructor $constructorName}}
       {{if constructorReturnsError $constructorName}}
         tempVarByVogen{{.Name}}, err := {{$constructorName}}({{.Name | camelCase}})
         if err != nil {
          return nil, err
         }
       {{else}}
         tempVarByVogen{{.Name}} := {{$constructorName}}({{.Name | camelCase}})
       {{end}}
     {{end}}
  {{end}}
    return &{{$structName}}{
     {{range $index, $field := .Fields}}
       {{$constructorName := printf "%s%s%s" ($prefix | pascalCase) ($structName | pascalCase) ($field.Name | pascalCase)}}
       {{.Name}}: {{if hasConstructor $constructorName}}tempVarByVogen{{.Name}}{{else}}{{.Name | camelCase}}{{end}},
     {{end}}
    }{{if shouldReturnError}}, nil{{end}}
  }
{{end}}
`)
	if err != nil {
		log.Fatalf("failed to parse template: %v", err)
	}

	baseFileName := strings.TrimSuffix(filename, ".go")
	outputFilename := baseFileName + "_vo_gen.go"
	f, err := os.Create(outputFilename)
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer f.Close()

	err = tmpl.Execute(f, data)
	if err != nil {
		log.Fatalf("failed to execute template: %v", err)
	}

	runGoImports(outputFilename)
	runGoFmt(outputFilename)

	log.Printf("successfully in %s\n", outputFilename)
}

func generateFactory(filename string, data *TemplateData) {
	tmpl, err := template.New("factory").Funcs(template.FuncMap{
		"pascalCase": toPascalCase,
		"camelCase":  toCamelCase,
		"uniqueImportPaths": func() []string {
			return uniqueImportPaths(data)
		},
	}).Parse(`// Code generated by vogen DO NOT EDIT.
package {{.PackageName}}
import (
  {{range uniqueImportPaths}}
	"{{.}}"
  {{end}}
  "testing"
)
{{range .Structs}}
  {{$structName := .Name}}
  type {{$structName}}Setter struct {
    {{range $index, $field := .Fields}}
  	{{$field.Name}} *{{$field.Type}}
    {{end}}
  }
  
  func Build{{$structName}}(t *testing.T, s *{{$structName}}Setter) *{{$structName}} {
    obj := &{{$structName}}{}
    {{range $index, $field := .Fields}}
  	{{$constructorName := printf "Build%s%s" ($structName | pascalCase) ($field.Name | pascalCase)}}
  	if s.{{$field.Name}} == nil {
        obj.{{$field.Name}} = {{$constructorName}}(t)
      } else {
        obj.{{$field.Name}} = *s.{{$field.Name}}
      }
    {{end}}
    return obj
  }
{{end}}
`)
	if err != nil {
		log.Fatalf("failed to parse template: %v", err)
	}

	baseFileName := strings.TrimSuffix(filename, ".go")
	outputFilename := baseFileName + "_factory_gen.go"
	f, err := os.Create(outputFilename)
	defer f.Close()

	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		log.Fatalf("failed to execute template: %v", err)
	}

	runGoImports(outputFilename)
	runGoFmt(outputFilename)

	log.Printf("generate factory successfully in %s\n", outputFilename)
}

func runGoImports(filename string) {
	cmd := exec.Command("go", "run", "golang.org/x/tools/cmd/goimports", "-w", filename)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to run goimports: %v", err)
	}
}

func runGoFmt(filename string) {
	cmd := exec.Command("go", "run", "cmd/gofmt", "-w", filename)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to run gofmt: %v", err)
	}
}
