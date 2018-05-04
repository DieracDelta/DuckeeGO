package main

// z3.stuff
// import "github.com/aclements/go-z3/z3"
import "os"

import "fmt"

// import "io/ioutil"
import "go/parser"
import "go/ast"
import "go/token"

// argment is path to example program
func main() {
	fset := token.NewFileSet()
	// TODO add more files  by including more args
	fileName := os.Args[1]

	parsedFile, err := parser.ParseFile(fset, fileName, nil, 0)

	if err != nil {
		panic(err)
	}

	// ast.Print(fset, parsedFile)
	add_instrumentation(fset, parsedFile)

	concolic_execute(parsedFile)
}

func add_instrumentation(fset *token.FileSet, parsedFile *ast.File) {
	queue := parsedFile.Decls
	for len(queue) > 0 {
		curNode := queue[0]
		queue = queue[1:]

		switch curNode := curNode.(type) {
		case *ast.GenDecl:
			switch curNode.Tok {
			case token.CONST:
			case token.TYPE:
			case token.VAR:
			case token.IMPORT:
			}
			fmt.Print(curNode)
		case *ast.FuncDecl:
			fmt.Print(curNode)
		}
	}
}

func concolic_execute(instrumenetedFile *ast.File) {
}
