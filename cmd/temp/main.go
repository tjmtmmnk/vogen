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

func main() {
	ctx := context.Background()
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cfg := &packages.Config{
		Context: ctx,
		Mode:    packages.LoadAllSyntax,
		Dir:     wd,
	}
	pkgs, err := packages.Load(cfg, "./sample")
	if err != nil {
		panic(err)
	}
	for _, pkg := range pkgs {
		typeMap := map[string]types.TypeAndValue{}
		for k, v := range maps.All(pkg.TypesInfo.Types) {
			if t, ok := k.(*ast.Ident); ok {
				typeMap[t.Name] = v
			}
		}
		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok {
					for _, spec := range genDecl.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							if t, ok := typeMap[typeSpec.Name.Name]; ok {
								if named, ok := t.Type.(*types.Named); ok {
									fmt.Println(named, named.Underlying())
								}
							}
						}
					}
				}
			}
		}
	}
}
