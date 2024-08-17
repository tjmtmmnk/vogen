package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"unicode"
)

type Field struct {
	Name string
	Type string
}

type TemplateData struct {
	PackageName string
	StructName  string
	Fields      []Field
}

func main() {
	// Define flags for source file and struct names
	sourceFile := flag.String("source", "", "Source file name")
	structNames := flag.String("structs", "", "Comma-separated list of struct names to generate constructors for")
	flag.Parse()

	// Check if both flags are provided
	if *sourceFile == "" || *structNames == "" {
		log.Fatalf("Usage: go run gen.go -source <FileName> -structs <StructName1,StructName2,...>")
	}

	filename := *sourceFile
	targetStructs := strings.Split(*structNames, ",")

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("failed to parse file: %v", err)
	}

	typeMap := make(map[string]string) // 型情報を保持するマップ
	constructors := []TemplateData{}   // 生成するコンストラクタの情報を保持する変数

	// ASTを解析して構造体やdefined typeを収集
	for _, f := range node.Decls {
		if genDecl, ok := f.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					// defined typeのチェック
					if ident, ok := typeSpec.Type.(*ast.Ident); ok {
						typeMap[typeSpec.Name.Name] = ident.Name
					}
					// 構造体の解析
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						for _, targetStruct := range targetStructs {
							if typeSpec.Name.Name == targetStruct {
								fields := []Field{}
								for _, field := range structType.Fields.List {
									fieldType := exprToString(field.Type)
									for _, name := range field.Names {
										fields = append(fields, Field{Name: name.Name, Type: fieldType})
									}
								}

								// 構造体情報を変数に設定
								constructors = append(constructors, TemplateData{
									PackageName: node.Name.Name,
									StructName:  typeSpec.Name.Name,
									Fields:      fields,
								})
							}
						}
					}
				}
			}
		}
	}

	// ターゲット構造体のコンストラクタを生成
	for _, constructor := range constructors {
		generateConstructor(constructor, typeMap)
	}
}

// 型情報の文字列化
func exprToString(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return exprToString(v.X) + "." + v.Sel.Name
	case *ast.StarExpr:
		return "*" + exprToString(v.X)
	case *ast.ArrayType:
		return "[]" + exprToString(v.Elt)
	default:
		return ""
	}
}

// camelCaseに変換
func toCamelCase(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// コンストラクタのコードを生成
func generateConstructor(data TemplateData, typeMap map[string]string) {
	tmpl, err := template.New("constructor").Funcs(template.FuncMap{
		"camelCase": toCamelCase,
		"isDefinedType": func(typ string) bool {
			_, defined := typeMap[typ]
			return defined
		},
		"getBaseType": func(typ string) string {
			if baseType, defined := typeMap[typ]; defined {
				return baseType
			}
			return typ
		},
	}).Parse(`// Auto-generated constructor for {{.StructName}}
package {{.PackageName}}
func New{{.StructName}}({{range $index, $field := .Fields}}{{if $index}}, {{end}}{{$field.Name | camelCase}} {{if isDefinedType $field.Type}}{{getBaseType $field.Type}}{{else}}{{$field.Type}}{{end}}{{end}}) *{{.StructName}} {
    return &{{.StructName}}{
        {{range .Fields}}
        {{.Name}}: {{if isDefinedType .Type}}New{{.Type}}({{.Name | camelCase}}){{else}}{{.Name | camelCase}}{{end}},
        {{end}}
    }
}
`)
	if err != nil {
		log.Fatalf("failed to parse template: %v", err)
	}

	outputFilename := strings.ToLower(data.StructName) + "_constructor_gen.go"
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

	log.Printf("Constructor for %s generated successfully in %s\n", data.StructName, outputFilename)
}

// goimportsを実行
func runGoImports(filename string) {
	cmd := exec.Command("go", "run", "golang.org/x/tools/cmd/goimports", "-w", filename)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to run goimports: %v", err)
	}
}

// go fmtを実行
func runGoFmt(filename string) {
	cmd := exec.Command("go", "run", "cmd/gofmt", "-w", filename)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to run gofmt: %v", err)
	}
}
