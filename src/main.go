package main

// z3.stuff
// import "github.com/aclements/go-z3/z3"
import "os"

import "fmt"

// import "io/ioutil"
// for rewriting
import "reflect"
import "bytes"
import "go/parser"
import "go/ast"
import "go/token"
import "go/printer"
import "golang.org/x/tools/go/ast/astutil"

// import "concolicTypes"

// import "reflect"

// argment is path to example program

func main() {
	if false {
		fmt.Print("mr duck\r\n")
	}
	fset := token.NewFileSet()
	// TODO add more files  by including more args
	filePath := os.Args[1]

	uninstrumentedAST, err := parser.ParseFile(fset, filePath, nil, 0)

	if err != nil {
		panic(err)
	}

	ast.Print(fset, uninstrumentedAST)
	instrumentedAST := astutil.Apply(uninstrumentedAST, astutil.ApplyFunc(addInstrumentationPre), astutil.ApplyFunc(addInstrumentationPost))

	// concolicExecute(instrumentedAST)
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, instrumentedAST)
	fmt.Println(buf.String())
}

// case *ast.BinaryExpr:
// case *ast.BasicLit:
// if curNode.Kind == token.INT {
// 	// implement replacement
// 	fmt.Printf("quack quack %s\r\n", curNode.Value)
// }

func addInstrumentationPre(curNode *astutil.Cursor) bool {
	// switch curNode.Node().(type) {
	return true

}

// TODO moo. add function lookup functionality
// TODO not really used ever tho so who cares
func containsIntType(curNode *ast.Node) bool {
	switch (*curNode).(type) {
	case *ast.BinaryExpr:
		ducker1 := ((*curNode).(*ast.BinaryExpr).X).(ast.Node)
		ducker2 := ((*curNode).(*ast.BinaryExpr).Y).(ast.Node)
		return containsIntType(&ducker1) || containsIntType(&ducker2)
	case *ast.BasicLit:
		return (*curNode).(*ast.BasicLit).Kind == token.INT
	default:
		return false
	}
}

func addInstrumentationPost(curNode *astutil.Cursor) bool {
	fmt.Println(reflect.TypeOf(curNode.Node()))
	switch curNode.Node().(type) {
	// the idea is to find a binary expression
	// then check if it contains an int type (or function that returns int type)
	// replace with the node with callexpr if it does
	// case *ast.
	case *ast.BinaryExpr:
		castedNode := curNode.Node().(*ast.BinaryExpr)

		// TODO add switch to determine the function you use
		addedNode := &ast.Ident{
			Name: "",
		}
		// if stuff gets assigned

		switch castedNode.Op {
		case token.ADD:
			addedNode.Name = "Add"
		case token.SUB:
			addedNode.Name = "Sub"
		case token.MUL:
			addedNode.Name = "Mul"
		case token.QUO:
			addedNode.Name = "Div"
		case token.REM:
			addedNode.Name = "Rem"
		default:
			panic("unsupported operation!!")
		}

		// depending on what it is, may not need to use int or w/e
		// actNode := curNode.Node()
		// if containsIntType(&actNode) {
		replacementNode :=
			ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   castedNode.X,
					Sel: addedNode,
				},
				Args: []ast.Expr{castedNode.Y},
			}

		curNode.Replace(&replacementNode)
		// }
	case *ast.BasicLit:
		castedNode := curNode.Node().(*ast.BasicLit)
		if castedNode.Kind == token.INT {
			identifier := getIdentifier(curNode)
			// if it's not a declaration of some sort lul
			var theElt []ast.Expr
			// TODO dry this out
			if identifier == "" {
				theElt = []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: "\"\"",
					},
					&ast.BasicLit{
						Kind:  token.INT,
						Value: "true",
					},
				}
			} else {
				theElt = []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: "\"" + identifier + "\"",
					},
					&ast.BasicLit{
						Kind:  token.INT,
						Value: "false",
					},
				}

			}
			bruh :=
				ast.CompositeLit{
					Type: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "concolicTypes",
						},
						Sel: &ast.Ident{
							Name: "ConcolicInt",
						},
					},
					Elts: []ast.Expr{
						castedNode,
						&ast.CompositeLit{
							Type: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "symTypes",
								},
								Sel: &ast.Ident{
									Name: "SymInt",
								},
							},
							Elts: theElt,
						},
					},
				}
			curNode.Replace(&bruh)
			// TODO implement replacement
		} else if castedNode.Kind == token.STRING {

		}

	case *ast.AssignStmt:
		castedNode := curNode.Node().(*ast.AssignStmt)
		addedNode := &ast.Ident{
			Name: "",
		}
		switch castedNode.Tok {
		case token.ADD_ASSIGN:
			addedNode.Name = "Add"
		case token.SUB_ASSIGN:
			addedNode.Name = "Sub"
		case token.MUL_ASSIGN:
			addedNode.Name = "Mul"
		case token.QUO_ASSIGN:
			addedNode.Name = "Div"
		case token.REM_ASSIGN:
			addedNode.Name = "Rem"
		default:
			return true
		}

		replacementNode :=
			ast.AssignStmt{
				Tok: token.ASSIGN,
				Lhs: castedNode.Lhs,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						// TODO assert about x len
						Fun: &ast.SelectorExpr{
							X:   castedNode.Lhs[0],
							Sel: addedNode,
						},
						Args: castedNode.Rhs,
					},
				},
			}
		curNode.Replace(&replacementNode)

	case *ast.FuncType:
	case *ast.BlockStmt:
	case *ast.Ident:
		castedNode := curNode.Node().(*ast.Ident)
		switch castedNode.Name {
		case "int":
			castedNode.Name = "concolicTypes.ConcolicInt"
		case "bool":
			castedNode.Name = "concolicTypes.ConcolicBool"
		case "true":
			fallthrough
		case "false":
			// TODO dry this out (combine it with other else ifs/put into method)
			identifier := getIdentifier(curNode)
			var theElt []ast.Expr
			if identifier == "" {
				theElt = []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: "\"\"",
					},
					&ast.BasicLit{
						Kind:  token.INT,
						Value: "true",
					},
				}
			} else {
				theElt = []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: "\"" + identifier + "\"",
					},
					&ast.BasicLit{
						Kind:  token.INT,
						Value: "false",
					},
				}

			}
			bruh :=
				ast.CompositeLit{
					Type: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "concolicTypes",
						},
						Sel: &ast.Ident{
							Name: "ConcolicBool",
						},
					},
					Elts: []ast.Expr{
						castedNode,
						&ast.CompositeLit{
							Type: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "symTypes",
								},
								Sel: &ast.Ident{
									Name: "SymBool",
								},
							},
							Elts: theElt,
						},
					},
				}
			curNode.Replace(&bruh)
		}
	default:
	}
	return true
}

func getIdentifier(curNode *astutil.Cursor) string {
	index := curNode.Index()
	parentNode := curNode.Parent()
	switch parentNode.(type) {
	case *ast.File:
		break
	case *ast.FuncDecl:
		break
	case *ast.AssignStmt:
		castedParentNode := parentNode.(*ast.AssignStmt)
		return castedParentNode.Lhs[index].(*ast.Ident).Name
	case *ast.ValueSpec:
		castedParentNode := parentNode.(*ast.ValueSpec)
		return castedParentNode.Names[index].Name
	default:
		break
	}

	return ""
}

func concolicExecute(instrumentedFile ast.Node) {

}
