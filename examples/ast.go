package main

import (
	"go/parser"
	"go/token"
	"log"
)

func main() {
	src := []byte(`package main
import "fmt"
func main() {
	gogo func() { fmt.Println("vim-go") }()
	return
}
`)

	fset := token.NewFileSet()

	repo := "github.com/ggaaooppeenngg/gogo"

	_, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		log.Fatal(err)
	}
}
