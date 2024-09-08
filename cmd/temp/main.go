package main

import (
	"context"
	"fmt"
	"go/ast"
	"go/types"
	"log"
	"os"
	"os/exec"
	"regexp"
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
}

var regex = regexp.MustCompile(`^(.+)/(\w+\.\w+)$`)

func main() {
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
	pkgs, err := packages.Load(cfg, "./sample/address.go")
	if err != nil {
		panic(err)
	}

	dataList := make([]*TemplateData, 0, len(pkgs))

	for _, pkg := range pkgs {
		data := &TemplateData{
			Structs:                  make([]*Struct, 0),
			TypeNameToUnderlyingType: make(map[string]string),
		}

		for _, syntax := range pkg.Syntax {
			for _, decl := range syntax.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok {
					for _, spec := range genDecl.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							switch t := typeSpec.Type.(type) {
							case *ast.StructType:
								structFields := make([]*StructField, len(t.Fields.List))

								for i, field := range t.Fields.List {
									if field.Names == nil || len(field.Names) == 0 {
										panic("field.Names is nil or empty")
									}
									name := field.Names[0]
									if ident, ok := field.Type.(*ast.Ident); ok {
										structFields[i] = &StructField{
											Name: name.Name,
											Type: ident.Name,
										}
									} else if selector, ok := field.Type.(*ast.SelectorExpr); ok {
										ty := pkg.TypesInfo.TypeOf(selector)
										typeName, importPath := extractTypeAndImportPath(ty.String())
										structFields[i] = &StructField{
											Name: name.Name,
											Type: typeName,
										}
										if importPath != "" {
											data.ImportPaths = append(data.ImportPaths, importPath)
										}
									}
								}
								data.Structs = append(data.Structs, &Struct{
									Name:   typeSpec.Name.Name,
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
			}
		}

		dataList = append(dataList, data)
		data.PackageName = pkg.Name
		data.ConstructorPrefix = "New"
		generateConstructor("sample/address.go", *data, map[string]bool{})
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

func generateConstructor(filename string, data TemplateData, constructorReturnsError map[string]bool) {
	tmpl, err := template.New("constructor").Funcs(template.FuncMap{
		"camelCase": toCamelCase,
		"pascalCase": toPascalCase,
		"shouldReturnError": func() bool {
			return len(constructorReturnsError) > 0
		},
		"constructorReturnsError": func(typ string) bool {
			return constructorReturnsError[data.ConstructorPrefix+typ]
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
		"isDefinedType": func(typeName string) bool {
			underlying, ok := data.TypeNameToUnderlyingType[typeName]
			if !ok {
				return false
			}
			return underlying != typeName
		},
		"uniqueImportPaths": func() []string {
			seen := map[string]bool{}
			uniqueImportPaths := make([]string, 0, len(data.ImportPaths))
			for _, path := range data.ImportPaths {
				if _, ok := seen[path]; !ok {
					uniqueImportPaths = append(uniqueImportPaths, path)
					seen[path] = true
				}
			}
			return uniqueImportPaths
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
     {{if constructorReturnsError .Type}}
       tempVarByVogen{{.Name}}, err := {{$prefix | pascalCase}}{{$structName | pascalCase}}{{$field.Name | pascalCase}}({{.Name | camelCase}})
       if err != nil {
        return nil, err
       }
     {{else}}
       tempVarByVogen{{.Name}} := {{$prefix | pascalCase}}{{$structName | pascalCase}}{{$field.Name | pascalCase}}({{.Name | camelCase}})
     {{end}}
  {{end}}
    return &{{$structName}}{
     {{range $index, $field := .Fields}}
       {{.Name}}: {{if isDefinedType .Type}}tempVarByVogen{{.Name}}{{else}}{{.Name | camelCase}}{{end}},
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

	//runGoImports(outputFilename)
	//runGoFmt(outputFilename)

	log.Printf("successfully in %s\n", outputFilename)
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
