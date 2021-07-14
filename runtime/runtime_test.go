package runtime

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"github.com/ggaaooppeenngg/gogo/go/ast"
	"github.com/ggaaooppeenngg/gogo/go/parser"
	"github.com/ggaaooppeenngg/gogo/go/printer"
	"github.com/ggaaooppeenngg/gogo/go/token"
)

var config = printer.Config{
	Mode:     printer.UseSpaces | printer.TabIndent | printer.SourcePos,
	Tabwidth: 8,
}

func TestMemAlloc(t *testing.T) {
	pointer := mal(32)
	a1 := (*int64)(unsafe.Pointer(pointer + 8))
	*a1 = 1
	var data []int64
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	sh.Data = pointer
	sh.Len = 4
	sh.Cap = 4
	fmt.Println(data[1])
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
func Hi() {
}


func TestSwitch(t *testing.T) {
	buf := GoGoBuf{}
	gosave(&buf)
	stack := mal(1024)
	sp := stack + 1024 - 4*8
	fmt.Println(sp)
	newBuf := GoGoBuf{
		PC: funcPC(Hi),
		SP: sp,
	}
	gogo(&newBuf)
}

func TestSchedule(t *testing.T) {
	NewProc(func() {
//line testdata/main.go2:10
		for i := 0; i < 100; i++ {
			fmt.Println("vim-go", "123")
			time.Sleep(time.Second)
		}
	})
//line testdata/main.go2:14
	NewProc(func() {
//line testdata/main.go2:16
		for i := 0; i < 100; i++ {
		fmt.Println("vim-go2", "321")
			time.Sleep(time.Second)
		}
	})
	time.Sleep(100 * time.Second)

}
