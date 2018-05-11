package main

import (
	"fmt"
	// for rewriting
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	// "reflect"
)

func instrumentBinaryExpr(curNode *astutil.Cursor) {
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
	replacementNode :=
		ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   castedNode.X,
				Sel: addedNode,
			},
			Args: []ast.Expr{castedNode.Y},
		}

	curNode.Replace(&replacementNode)
}

func instrumentUnaryExpr(curNode *astutil.Cursor) {
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
}

func instrumentBasicLit(curNode *astutil.Cursor) {
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
				// &ast.Ident{
				// 	Name: "cv",
				// },
				&ast.Ident{
					Name: "\"" + identifier + "\"",
				},
			}
		}
		augNode :=
			ast.CallExpr{
				Fun: &ast.Ident{
					Name: "concolicTypes.MakeConcolic" + "Int" + varOrConst,
				},
				Args: theArgs,
			}
		curNode.Replace(&augNode)
	} else if castedNode.Kind == token.STRING {
		// TODO
	}
}

func instrumentAssignStmt(curNode *astutil.Cursor) {
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
		return
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
}

func instrumentIncDecStmt(curNode *astutil.Cursor) {
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

	augNode := ast.CallExpr{
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
				Args: []ast.Expr{&augNode},
			},
		},
	}
	curNode.Replace(&replacementNode)
}

func instrumentIdent(curNode *astutil.Cursor) {
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
				// &ast.Ident{
				// 	Name: "cv",
				// },
				&ast.Ident{
					Name: "\"" + identifier + "\"",
				},
			}
		}
		augNode :=
			ast.CallExpr{
				Fun: &ast.Ident{
					Name: "concolicTypes.MakeConcolic" + concType + varOrConst,
				},
				Args: theArgs,
			}
		curNode.Replace(&augNode)

	}

}

func instrumentIfStmt(curNode *astutil.Cursor) {
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
						// &ast.BasicLit{
						// 	Kind:  token.STRING,
						// 	Value: "currPathConstrs",
						// },
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
							// &ast.BasicLit{
							// 	Kind:  token.STRING,
							// 	Value: "currPathConstrs",
							// },
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
}

func instrumentFuncDecl(curNode *astutil.Cursor) {
	castedNode := curNode.Node().(*ast.FuncDecl)
	// don't instrument main (I'm assuming main *could* have args)
	// just being safe
	if castedNode.Name.Name == "main" {
		instrumentMainMethod(curNode)
		return
	}

	// add setargspopped statement right after all vars do their thing
	poppedStatement := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.Ident{
				Name: "concolicTypes.SymStack.SetArgsPopped",
			},
		},
	}
	castedNode.Body.List = append([]ast.Stmt{poppedStatement}, castedNode.Body.List...)

	castedType := castedNode.Type

	for _, aParam := range castedType.Results.List {
		// supposed to look like
		// i = makeConcolicIntVar(cv, "i")
		switch aParam.Type.(*ast.Ident).Name {
		case "concolicTypes.ConcolicString":
			aParam.Type = &ast.Ident{Name: "string"}
		case "concolicTypes.ConcolicInt":
			aParam.Type = &ast.Ident{Name: "int"}
		case "concolicTypes.ConcolicBool":
			aParam.Type = &ast.Ident{Name: "bool"}
		default:
			fmt.Printf(aParam.Type.(*ast.Ident).Name + "\r\n")
			// fmt.("WE DON'T SUPPORT THIS TYPE!")
			// if the type is wrong, it's all wrong, so move onto next parameter
			break

		}
	}

	for index1, aParam := range castedType.Params.List {
		aParam = castedType.Params.List[len(castedType.Params.List)-1-index1]
		for index2, aName := range aParam.Names {
			aName = aParam.Names[len(aParam.Names)-1-index2]
			// supposed to look like
			// i = makeConcolicIntVar(cv, "i")
			var methodPiece string
			switch aParam.Type.(*ast.Ident).Name {
			case "string":
				fallthrough
				// TODO does this support string correctly
			case "concolicTypes.ConcolicString":
				methodPiece = "String"
			case "int":
				fallthrough
			case "concolicTypes.ConcolicInt":
				methodPiece = "Int"
			case "bool":
				fallthrough
			case "concolicTypes.ConcolicBool":
				methodPiece = "Bool"
			default:
				fmt.Printf(aParam.Type.(*ast.Ident).Name + "\r\n")
				// fmt.("WE DON'T SUPPORT THIS TYPE!")
				// if the type is wrong, it's all wrong, so move onto next parameter
				break

			}
			// TODO add in concolic constructors before the _ thingies
			// add "_ = y" for example
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
							Name: "concolicTypes.MakeConcolic" + methodPiece,
						},
						Args: []ast.Expr{
							&ast.Ident{
								Name: aName.Name,
							},
							&ast.Ident{
								Name: "concolicTypes.SymStack.PopArg().(z3." + methodPiece + ")",
							},
						},
					},
				},
			}
			castedNode.Body.List = append([]ast.Stmt{&newNode}, castedNode.Body.List...)

			// set each parameter to a different Name
			aName.Name = aName.Name + "Val"
		}
	}
	// example:
	// ruckkerduck( vw * ConcreteValues, curPathConstrs []z3.Bool)

	// replacing argument values

	// newFuncArgs := []*ast.Field{
	// 	&ast.Field{
	// 		Names: []*ast.Ident{
	// 			&ast.Ident{
	// 				Name: "cv",
	// 			},
	// 		},
	// 		Type: &ast.Ident{
	// 			Name: "* concolicTypes.ConcreteValues",
	// 		},
	// 	},
	// 	&ast.Field{
	// 		Names: []*ast.Ident{
	// 			&ast.Ident{
	// 				Name: "currPathConstrs",
	// 			},
	// 		},
	// 		Type: &ast.Ident{
	// 			Name: "*[]z3.Bool",
	// 		},
	// 	},
	// }
	// castedType.Params.List = newFuncArgs
}

func instrumentCallExpr(curNode *astutil.Cursor) {
	castedNode := curNode.Node().(*ast.CallExpr)
	switch castedNode.Fun.(type) {
	case *ast.Ident:
		if castedNode.Fun == nil || castedNode.Fun.(*ast.Ident).Obj == nil || castedNode.Fun.(*ast.Ident).Obj.Decl == nil {
			break
		}
		objectified := castedNode.Fun.(*ast.Ident).Obj.Decl.(*ast.FuncDecl)

		var objectifiedNode *ast.FieldList
		if objectified != nil && objectified.Type != nil {
			objectifiedNode = objectified.Type.Results
		} else {
			objectified = nil
		}
		fmt.Printf("SHISTSHITSHTI")

		newNode := &ast.CallExpr{
			Fun: &ast.FuncLit{
				Type: &ast.FuncType{
					Results: objectifiedNode,
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ReturnStmt{
							Results: []ast.Expr{castedNode},
						},
					},
				},
			},
		}
		if objectified != nil {
			declaration := castedNode.Fun.(*ast.Ident).Obj.Decl.(*ast.FuncDecl)
			// name := declaration.Name.Name
			paramList := declaration.Type.Params.List
			for _, aParam := range paramList {
				// TODO value
				if aParam.Type.(*ast.Ident).Name == "string" || aParam.Type.(*ast.Ident).Name == "int" || aParam.Type.(*ast.Ident).Name == "bool" {
					for _, aNameNode := range aParam.Names {
						newNode.Fun.(*ast.FuncLit).Body.List = append(
							[]ast.Stmt{
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Fun: &ast.Ident{
											Name: "concolicTypes.SymStack.PushArg",
										},
										Args: []ast.Expr{
											&ast.Ident{
												Name: aNameNode.Name + ".Z3Expr",
											},
										},
									},
								},
							}, newNode.Fun.(*ast.FuncLit).Body.List...)
						aNameNode.Name += "Val"
						// TODO might fuck some things up
						aNameNode.Obj = nil
					}
				}
			}

			for i, _ := range castedNode.Args {
				castedNode.Args[i] = &ast.SelectorExpr{
					X: castedNode.Args[i],
					Sel: &ast.Ident{
						Name: "Value",
					},
				}

			}

			newNode.Fun.(*ast.FuncLit).Body.List = append(
				[]ast.Stmt{&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.Ident{Name: "concolicTypes.SymStack.SetArgsPushed"},
					},
				},
				},
				newNode.Fun.(*ast.FuncLit).Body.List...)

			for _, aParam := range objectifiedNode.List {
				// supposed to look like
				// i = makeConcolicIntVar(cv, "i")
				switch aParam.Type.(*ast.Ident).Name {
				case "concolicTypes.ConcolicString":
					aParam.Type = &ast.Ident{Name: "string"}
				case "concolicTypes.ConcolicInt":
					aParam.Type = &ast.Ident{Name: "int"}
				case "concolicTypes.ConcolicBool":
					aParam.Type = &ast.Ident{Name: "bool"}
				default:
					fmt.Printf(aParam.Type.(*ast.Ident).Name + "\r\n")
					// fmt.("WE DON'T SUPPORT THIS TYPE!")
					// if the type is wrong, it's all wrong, so move onto next parameter
					break

				}
			}
			// for _, aParam := range objectifiedNode.List {
			// 	switch aParam.Type.(type) {
			// 	case *ast.Ident:
			// 		switch aParam.Type.(*ast.Ident).Name {
			// 		case "concolicTypes.ConcolicString":
			// 			aParam.Type.(*ast.Ident).Name = "string"
			// 			aParam.Type.(*ast.Ident).Obj.Name = "string"
			// 			aParam.Type.(*ast.Ident).NamePos = token.NoPos
			// 		case "concolicTypes.ConcolicInt":
			// 			aParam.Type.(*ast.Ident).Name = "int"
			// 			aParam.Type.(*ast.Ident).Obj.Name = "int"
			// 			aParam.Type.(*ast.Ident).NamePos = token.NoPos
			// 		case "concolicTypes.ConcolicBool":
			// 			aParam.Type.(*ast.Ident).Name = "bool"
			// 			aParam.Type.(*ast.Ident).Obj.Name = "bool"
			// 			aParam.Type.(*ast.Ident).NamePos = token.NoPos
			// 			// aParam.Type.(*ast.Ident).Obj.
			// 		default:
			// 			fmt.Printf(aParam.Type.(*ast.Ident).Name + "\r\n")
			// 			// fmt.("WE DON'T SUPPORT THIS TYPE!")
			// 			// if the type is wrong, it's all wrong, so move onto next parameter

			// 		}

			// 	}
			// }

		}
		curNode.Replace(newNode)
	case *ast.SelectorExpr:
	case *ast.FuncLit:
	default:
		ast.Print(token.NewFileSet(), curNode.Node())
		panic("not supported!")
	}
	// TODO

}

func instrumentReturnStmt(curNode *astutil.Cursor) {
	castedNode := curNode.Node().(*ast.ReturnStmt)
	newNode := &ast.BlockStmt{
		List: []ast.Stmt{
			castedNode,
		},
	}
	for index, val := range castedNode.Results {
		symStackStmt :=
			[]ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.Ident{
							Name: "concolicTypes.SymStack.PushReturn",
						},
						// TODO add all arguments
						// TODO enable type checking for return
						// currently asssuming returns concolic whatever
						Args: []ast.Expr{
							&ast.SelectorExpr{
								X: val,
								Sel: &ast.Ident{
									Name: "Z3Expr",
								},
							},
						},
					},
				},
			}
		newNode.List = append(symStackStmt, newNode.List...)
		castedNode.Results[index] = &ast.SelectorExpr{
			X: val,
			Sel: &ast.Ident{
				Name: "Value",
			},
		}
	}
	curNode.Replace(newNode)
}

// add in a handler here and rename the method
func instrumentMainMethod(curNode *astutil.Cursor) {
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
	castedNode.Name.Name = "instrumentedMainMethod"
}
