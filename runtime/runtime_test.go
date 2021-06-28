package runtime

import (
	"log"
	"os"
	"testing"

	"github.com/ggaaooppeenngg/gogo/go/ast"
	"github.com/ggaaooppeenngg/gogo/go/parser"
	"github.com/ggaaooppeenngg/gogo/go/printer"
	"github.com/ggaaooppeenngg/gogo/go/token"
)

var config = printer.Config{
	Mode:     printer.UseSpaces | printer.TabIndent | printer.SourcePos,
	Tabwidth: 8,
}

func TestAST(t *testing.T) {
	src := []byte(`package main
import "fmt"
func main() {
	gogo func() { fmt.Println("vim-go") }()
	return
}
`)

	fset := token.NewFileSet()

	funcName := "gogo.NewProc"

	pf, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		log.Fatal(err)
	}
	for i, decl := range pf.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if ok {
			stmtList := funcDecl.Body.List
			for i := range stmtList {
				goStmt, ok := stmtList[i].(*ast.GoGoStmt)
				if ok {
					expr := &ast.CallExpr{Fun: &ast.Ident{Name: funcName}, Args: []ast.Expr{goStmt.Call}}
					stmtList[i] = &ast.ExprStmt{X: expr}
				}
			}
			pf.Decls[i] = funcDecl
		}
		genDecl, ok := decl.(*ast.GenDecl)
		if ok {
			if genDecl.Tok == token.IMPORT {
				importDecl := genDecl
				importDecl.Specs = append(importDecl.Specs, &ast.ImportSpec{
					Name: &ast.Ident{Name: "gogo"},
					Path: &ast.BasicLit{Value: "\"github.com/ggaaooppeenngg/gogo\""},
				})
				pf.Decls[i] = importDecl
			}
		}
	}
	config.Fprint(os.Stdout, fset, pf)
}
