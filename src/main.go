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
	// "reflect"
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

		astutil.AddImport(fset, uninstrumentedAST, "concolicTypes")
		astutil.AddImport(fset, uninstrumentedAST, "github.com/aclements/go-z3/z3")

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
	mainFile := constructMain(configData)
	astutil.AddImport(fset, mainFile, "reflect")
	astutil.AddImport(fset, mainFile, "concolicTypes")
	err = printer.Fprint(&buf, fset, mainFile)
	if err != nil {
		panic(err)
	}
	_ = os.Remove(DEST + "main.go")
	tmpFile, _ := os.Create(DEST + "main.go")
	_, _ = tmpFile.WriteString(buf.String())
	_ = tmpFile.Close()

}

func constructMain(configData ConfigData) *ast.File {
	a := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: "\"concolicTypes\"",
		},
	}
	b := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: "\"reflect\"",
		},
	}
	d := &ast.TypeSpec{
		Name: &ast.Ident{
			Name: "Handler",
		},
		Type: &ast.StructType{
			Fields: &ast.FieldList{},
		},
	}

	stuff := &ast.File{
		Name: &ast.Ident{
			Name: configData.Package,
		},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok:   token.IMPORT,
				Specs: []ast.Spec{a},
			},
			&ast.GenDecl{
				Tok:   token.IMPORT,
				Specs: []ast.Spec{b},
			},
			&ast.GenDecl{
				Tok:   token.TYPE,
				Specs: []ast.Spec{d},
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
					Fun: &ast.Ident{Name: "concolicTypes.ConcolicExec"},
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

			stuff.Decls[3].(*ast.FuncDecl).Body.List = append(stuff.Decls[3].(*ast.FuncDecl).Body.List, node1)
			stuff.Decls[3].(*ast.FuncDecl).Body.List = append(stuff.Decls[3].(*ast.FuncDecl).Body.List, node2)

		}
	}
	stuff.Imports = []*ast.ImportSpec{a, b}
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
	// fmt.Println(reflect.TypeOf(curNode.Node()))
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
			addedNode.Name = "ConcIntAdd"
		case token.SUB:
			addedNode.Name = "ConcIntSub"
		case token.MUL:
			addedNode.Name = "ConcIntMul"
		case token.QUO:
			addedNode.Name = "ConcIntDiv"
		case token.REM:
			addedNode.Name = "ConcIntMod"
		case token.AND:
			addedNode.Name = "ConcIntAnd"
		case token.OR:
			addedNode.Name = "ConcIntOr"
		case token.XOR:
			addedNode.Name = "ConcIntXOr"
		case token.SHL:
			addedNode.Name = "ConcIntSHL"
		case token.SHR:
			addedNode.Name = "ConcIntSHR"
		case token.AND_NOT:
			addedNode.Name = "ConcBoolAndNot"
		case token.LAND:
			addedNode.Name = "ConcBoolAnd"
		case token.LOR:
			addedNode.Name = "ConcBoolOr"
		case token.EQL:
			addedNode.Name = "ConcEq"
		case token.LSS:
			addedNode.Name = "ConcIntLT"
		case token.GTR:
			addedNode.Name = "ConcIntGT"
		case token.NEQ:
			addedNode.Name = "ConcNE"
		case token.GEQ:
			addedNode.Name = "ConcIntGE"
		case token.LEQ:
			addedNode.Name = "ConcIntLE"
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
						Name: "ConcBoolNot",
					},
				},
			}
			curNode.Replace(&replacemenetNode)
		case token.XOR:
			replacemenetNode := ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: castedNode.X,
					Sel: &ast.Ident{
						Name: "ConcIntNot",
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
			var theArgs []ast.Expr
			var varOrConst string
			if identifier == "" {
				varOrConst = "Const"
				theArgs = []ast.Expr{
					castedNode,
				}
			} else {
				varOrConst = "Var"
				theArgs = []ast.Expr{
					&ast.Ident{
						Name: "cv",
					},
					&ast.Ident{
						Name: "\"" + identifier + "\"",
					},
				}
			}
			bruh :=
				ast.CallExpr{
					Fun: &ast.Ident{
						Name: "concolicTypes.MakeConcolic" + "Int" + varOrConst,
					},
					Args: theArgs,
				}
			curNode.Replace(&bruh)
		} else if castedNode.Kind == token.STRING {

		}

	case *ast.AssignStmt:
		castedNode := curNode.Node().(*ast.AssignStmt)
		addedNode := &ast.Ident{
			Name: "",
		}
		switch castedNode.Tok {
		case token.ADD_ASSIGN:
			addedNode.Name = "ConcIntAdd"
		case token.SUB_ASSIGN:
			addedNode.Name = "ConcIntSub"
		case token.MUL_ASSIGN:
			addedNode.Name = "ConcIntMul"
		case token.QUO_ASSIGN:
			addedNode.Name = "ConcIntDiv"
		case token.REM_ASSIGN:
			addedNode.Name = "ConcIntMod"
		case token.AND_ASSIGN:
			addedNode.Name = "ConcIntAnd"
		case token.OR_ASSIGN:
			addedNode.Name = "ConcIntOr"
		case token.XOR_ASSIGN:
			addedNode.Name = "ConcIntXOr"
		case token.SHL_ASSIGN:
			addedNode.Name = "ConcIntSHL"
		case token.SHR_ASSIGN:
			addedNode.Name = "ConcIntSHR"
		case token.AND_NOT_ASSIGN:
			addedNode.Name = "ConcBoolAndNot"
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
			addedNode.Name = "ConcIntAdd"
		case token.DEC:
			addedNode.Name = "ConcIntSub"
		}

		regNode := &ast.BasicLit{
			Kind:  token.INT,
			Value: "1",
		}

		bruh := ast.CallExpr{
			Fun: &ast.Ident{
				Name: "concolicTypes.MakeConcolicIntConst",
			},
			Args: []ast.Expr{regNode},
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
		var varOrConst, concType string
		switch castedNode.Name {
		case "int":
			castedNode.Name = "concolicTypes.ConcolicInt"
		case "bool":
			castedNode.Name = "concolicTypes.ConcolicBool"
		case "true":
			fallthrough
		case "false":
			concType = "Bool"
			identifier := getIdentifier(curNode)
			var theArgs []ast.Expr
			if identifier == "" {
				varOrConst = "Const"
				theArgs = []ast.Expr{
					castedNode,
				}
			} else {
				varOrConst = "Var"
				theArgs = []ast.Expr{
					&ast.Ident{
						Name: "cv",
					},
					&ast.Ident{
						Name: "\"" + identifier + "\"",
					},
				}
			}
			bruh :=
				ast.CallExpr{
					Fun: &ast.Ident{
						Name: "concolicTypes.MakeConcolic" + concType + varOrConst,
					},
					Args: theArgs,
				}
			curNode.Replace(&bruh)

		}
	case *ast.IfStmt:
		castedNode := curNode.Node().(*ast.IfStmt)
		cond := castedNode.Cond
		castedNode.Cond = &ast.SelectorExpr{
			X: cond,
			Sel: &ast.Ident{
				Name: "Value",
			},
		}
		castedNode.Body.List = append(
			[]ast.Stmt{
				// TODO it might be better to SHA shit here tbh
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "concolicTypes",
							},
							Sel: &ast.Ident{
								Name: "AddPositivePathConstr",
							},
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: "currPathConstrs",
							},
							&ast.SelectorExpr{
								X: cond,
								Sel: &ast.Ident{
									Name: "Z3Expr",
								},
							},
						},
					},
				},
			},
			castedNode.Body.List...)
		if castedNode.Else != nil {
			castedNode.Else = &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "concolicTypes",
								},
								Sel: &ast.Ident{
									Name: "AddNegativePathConstr",
								},
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: "currPathConstrs",
								},
								&ast.SelectorExpr{
									X: cond,
									Sel: &ast.Ident{
										Name: "Z3Expr",
									},
								},
							},
						},
					},
					castedNode.Else,
				},
			}
		}

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
					fallthrough
				case "concolicTypes.ConcolicString":
					randoType = "String"
				case "int":
					fallthrough
				case "concolicTypes.ConcolicInt":
					randoType = "Int"
				case "bool":
					fallthrough
				case "concolicTypes.ConcolicBool":
					randoType = "Bool"
				default:
					fmt.Printf(aParam.Type.(*ast.Ident).Name + "\r\n")
					panic("WTF WE DON'T SUPPORT THIS TYPE!")

				}
				newNode2 := ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "_",
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.Ident{
							Name: aName.Name,
						},
					},
				}
				castedNode.Body.List = append([]ast.Stmt{&newNode2}, castedNode.Body.List...)
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
						Name: "cv",
					},
				},
				Type: &ast.Ident{
					Name: "* concolicTypes.ConcreteValues",
				},
			},
			&ast.Field{
				Names: []*ast.Ident{
					&ast.Ident{
						Name: "currPathConstrs",
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

func isValid(node *ast.Node) bool {
	return true
}
