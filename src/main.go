package main

// z3.stuff
// import "github.com/aclements/go-z3/z3"
import (
	"github.com/otiai10/copy"
	"os"

	"fmt"

	"io/ioutil"

	// for rewriting
	"bytes"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"reflect"
)

// import "concolicTypes"

// import "reflect"

// argment is path to example program
var DEST = "/tmp/DuckieConcolic/"

func main() {
	if false {
		fmt.Print("mr duck\r\n")
	}

	fileConfigPath := os.Args[1]

	jsonFile, err := os.Open(fileConfigPath)
	if err != nil {
		panic(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = jsonFile.Close()
	if err != nil {
		panic(err)
	}

	var configData ConfigData
	err = json.Unmarshal(byteValue, &configData)
	if err != nil {
		panic(err)
	}

	err = copy.Copy(configData.ProjectPath, DEST)
	if err != nil {
		panic(err)
	}

	for _, aGoFile := range configData.ConfigData {
		fset := token.NewFileSet()
		// TODO add more files  by including more args
		filePath := configData.ProjectPath + aGoFile.FilePath
		uninstrumentedAST, err := parser.ParseFile(fset, filePath, nil, 0)

		if err != nil {
			panic(err)
		}

		_ = ast.Print(fset, uninstrumentedAST)
		instrumentedAST := astutil.Apply(uninstrumentedAST, astutil.ApplyFunc(addInstrumentationPre), astutil.ApplyFunc(addInstrumentationPost))

		var buf bytes.Buffer
		err = printer.Fprint(&buf, fset, instrumentedAST)
		if err != nil {
			panic(err)
		}
		// fmt.Println(buf.String())
		_ = os.Remove(DEST + aGoFile.FilePath)
		tmpFile, _ := os.Create(DEST + aGoFile.FilePath)
		_, _ = tmpFile.WriteString(buf.String())
		_ = tmpFile.Close()
	}

	fset := token.NewFileSet()
	var buf bytes.Buffer
	err = printer.Fprint(&buf, fset, constructMain(configData))
	if err != nil {
		panic(err)
	}
	_ = os.Remove(DEST + "main.go")
	tmpFile, _ := os.Create(DEST + "main.go")
	_, _ = tmpFile.WriteString(buf.String())
	_ = tmpFile.Close()

}

func constructMain(configData ConfigData) *ast.File {
	stuff := &ast.File{
		Name: &ast.Ident{
			Name: configData.Package,
		},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"reflect\"",
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"~/school/6.858/DuckDuckGo/src/concolicTypes\"",
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "main",
				},
				Type: &ast.FuncType{},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "h",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "new",
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "Handler",
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "method",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   &ast.Ident{Name: "reflect"},
												Sel: &ast.Ident{Name: "ValueOf"},
											},
											Args: []ast.Expr{
												&ast.Ident{Name: "h"},
											},
										},
										Sel: &ast.Ident{Name: "MethodByName"},
									},
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: "\"DONTWORRYABOUTITBRO\"",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, aThing := range configData.ConfigData {
		for _, aFunc := range aThing.Functions {
			node1 :=
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "method",
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "reflect"},
										Sel: &ast.Ident{Name: "ValueOf"},
									},
									Args: []ast.Expr{
										&ast.Ident{Name: "h"},
									},
								},
								Sel: &ast.Ident{Name: "MethodByName"},
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: "\"" + aFunc.Name + "\"",
								},
							},
						},
					},
				}
			node2 := &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.Ident{Name: "concolicTypes.concolicExec"},
					Args: []ast.Expr{
						&ast.Ident{
							Name: "method",
						},
						&ast.BasicLit{
							Kind:  token.INT,
							Value: "100",
						},
					},
				},
			}

			stuff.Decls[2].(*ast.FuncDecl).Body.List = append(stuff.Decls[2].(*ast.FuncDecl).Body.List, node1)
			stuff.Decls[2].(*ast.FuncDecl).Body.List = append(stuff.Decls[2].(*ast.FuncDecl).Body.List, node2)

		}

	}
	return stuff
}

func addInstrumentationPre(curNode *astutil.Cursor) bool {
	// TODO don't really need anything in here yet
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
		case token.AND:
			addedNode.Name = "And"
		case token.OR:
			addedNode.Name = "Or"
		case token.XOR:
			addedNode.Name = "Xor"
		case token.SHL:
			addedNode.Name = "Shl"
		// TODO add support for andnot
		// case token.AND_NOT:
		// 		addedNode.Name = "AndNot"
		case token.LAND:
			addedNode.Name = "LAnd"
		case token.LOR:
			addedNode.Name = "LOr"
		case token.EQL:
			addedNode.Name = "Equals"
		case token.LSS:
			addedNode.Name = "Lss"
		case token.GTR:
			addedNode.Name = "Gtr"
		case token.NEQ:
			addedNode.Name = "NEq"
		case token.GEQ:
			addedNode.Name = "GEq"
		case token.LEQ:
			addedNode.Name = "LEq"
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
	case *ast.UnaryExpr:
		// ! on bools is the only case I can think of
		castedNode := curNode.Node().(*ast.UnaryExpr)
		switch castedNode.Op {
		case token.NOT:
			replacemenetNode := ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: castedNode.X,
					Sel: &ast.Ident{
						Name: "Not",
					},
				},
			}
			curNode.Replace(&replacemenetNode)
		}
		// ruckkerduck( vw * ConcreteValues, curPathConstrs []z3.Bool)
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
						&ast.CallExpr{
							Fun: &ast.Ident{
								Name: "cv.getIntValue",
							},
							Args: []ast.Expr{
								castedNode,
							},
						},
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
		case token.AND_ASSIGN:
			addedNode.Name = "And"
		case token.OR_ASSIGN:
			addedNode.Name = "Or"
		case token.XOR_ASSIGN:
			addedNode.Name = "XOr"
		case token.SHL_ASSIGN:
			addedNode.Name = "Shl"
		case token.SHR_ASSIGN:
			addedNode.Name = "Shr"
		// case token.AND_NOT_ASSIGN:
		// 	addedNode.Name = ""
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

	case *ast.IncDecStmt:
		castedNode := curNode.Node().(*ast.IncDecStmt)
		addedNode := &ast.Ident{
			Name: "",
		}
		switch castedNode.Tok {
		case token.INC:
			addedNode.Name = "Add"
		case token.DEC:
			addedNode.Name = "Sub"
		}

		regNode := &ast.BasicLit{
			Kind:  token.INT,
			Value: "1",
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
					regNode,
					&ast.CompositeLit{
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "symTypes",
							},
							Sel: &ast.Ident{
								Name: "SymInt",
							},
						},
						Elts: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: "\"\"",
							},
							&ast.BasicLit{
								Kind:  token.INT,
								Value: "true",
							},
						},
					},
				},
			}

		replacementNode := ast.AssignStmt{
			Tok: token.ASSIGN,
			Lhs: []ast.Expr{castedNode.X},
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   castedNode.X,
						Sel: addedNode,
					},
					Args: []ast.Expr{&bruh},
				},
			},
		}
		curNode.Replace(&replacementNode)

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
	case *ast.IfStmt:
		// TODO worry about this
	case *ast.FuncDecl:
		castedNode := curNode.Node().(*ast.FuncDecl)
		castedNode.Recv = &ast.FieldList{
			List: []*ast.Field{
				&ast.Field{
					Names: []*ast.Ident{
						&ast.Ident{
							Name: "h",
						},
					},
					Type: &ast.Ident{
						Name: "Handler",
					},
				},
			},
		}
		castedType := castedNode.Type

		for _, aParam := range castedType.Params.List {
			// theType := aParam.Type.(*ast.Ident).Name
			for _, aName := range aParam.Names {
				// i = makeConcolicIntVar(cv, "i")
				var randoType string
				switch aParam.Type.(*ast.Ident).Name {
				case "string":
					randoType = "String"
				case "int":
					randoType = "Int"
				case "bool":
					randoType = "Bool"
				case "concolicTypes.ConcolicString":
					randoType = "String"
				case "concolicTypes.ConcolicInt":
					randoType = "Int"
				case "concolicTypes.ConcolicBool":
					randoType = "Bool"
				default:
					fmt.Printf(aParam.Type.(*ast.Ident).Name)
					panic("WTF WE DON'T SUPPORT THIS TYPE!")

				}
				newNode := ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: aName.Name,
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{
								Name: "concolicTypes.MakeConcolic" + randoType + "Var",
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "cv",
								},
								&ast.Ident{
									Name: "\"" + aName.Name + "\"",
								},
							},
						},
					},
				}
				castedNode.Body.List = append([]ast.Stmt{&newNode}, castedNode.Body.List...)

			}

		}

		// ruckkerduck( vw * ConcreteValues, curPathConstrs []z3.Bool)
		newFuncArgs := []*ast.Field{
			&ast.Field{
				Names: []*ast.Ident{
					&ast.Ident{
						Name: "cw",
					},
				},
				Type: &ast.Ident{
					Name: "ConcolicTypes.ConcreteValues",
				},
			},
			&ast.Field{
				Names: []*ast.Ident{
					&ast.Ident{
						Name: "curPathConstrs",
					},
				},
				Type: &ast.Ident{
					Name: "*[]z3.Bool",
				},
			},
		}
		castedType.Params.List = newFuncArgs

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
