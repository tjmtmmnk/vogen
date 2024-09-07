package main

import (
	"context"
	"fmt"
	"go/ast"
	"go/types"
	"os"

	"golang.org/x/tools/go/packages"
)

func debug(x any) {
	fmt.Printf("%#v\n", x)
}

type structField struct {
	Name string
	Type string
}

type underlyingType struct {
	Name string
	Type string
}

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

	var underlyingTypes []*underlyingType
	var structFields []*structField

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok {
					for _, spec := range genDecl.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							switch t := typeSpec.Type.(type) {
							case *ast.StructType:
								for _, field := range t.Fields.List {
									if field.Names == nil || len(field.Names) == 0 {
										panic("field.Names is nil or empty")
									}
									name := field.Names[0]
									if ident, ok := field.Type.(*ast.Ident); ok {
										structFields = append(structFields, &structField{
											Name: name.Name,
											Type: ident.Name,
										})
									}
									if selector, ok := field.Type.(*ast.SelectorExpr); ok {
										ty := pkg.TypesInfo.TypeOf(selector)
										structFields = append(structFields, &structField{
											Name: name.Name,
											Type: ty.String(),
										})
									}
								}
							default:
								underlyingTypes = append(underlyingTypes, &underlyingType{
									Name: typeSpec.Name.Name,
									Type: getUnderlyingType(pkg.TypesInfo, t),
								})
							}
						}
					}
				}
			}
		}
		for _, v := range underlyingTypes {
			debug(v)
		}
		for _, v := range structFields {
			debug(v)
		}
	}
}

func getUnderlyingType(info *types.Info, spec ast.Expr) string {
	switch t := spec.(type) {
	case *ast.Ident:
		ty := info.TypeOf(t)
		return ty.Underlying().String()
	case *ast.SelectorExpr:
		ty := info.TypeOf(t)
		return ty.String() // time.Timeのような型はそのまま使う。Underlying()はtime.Timeの構造体の中身を指すため。
	}
	return ""
}
