package main

import (
	"context"
	"fmt"
	"go/ast"
	"go/types"
	"maps"
	"os"

	"golang.org/x/tools/go/packages"
)

func debug(x any) {
	fmt.Printf("%#v\n", x)
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
	for _, pkg := range pkgs {
		typeNameToUnderlyingType := map[string]types.Type{}
		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok {
					for _, spec := range genDecl.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							switch t := typeSpec.Type.(type) {
							case *ast.Ident, *ast.SelectorExpr:
								ty := pkg.TypesInfo.TypeOf(t)
								typeNameToUnderlyingType[typeSpec.Name.Name] = ty
							}
						}
					}
				}
			}
		}
		for k, v := range maps.All(typeNameToUnderlyingType) {
			fmt.Printf("%s\t%s\t%s\n", k, v.String(), v.Underlying())
		}
	}
}
