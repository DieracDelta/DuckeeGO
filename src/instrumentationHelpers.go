package main

import (
	"fmt"
	// for rewriting
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	// "strings"
	// "reflect"
)

var typeMapping map[string]string

var concolicIntTypeString = "CONCOLIC_INT"
var concolicBoolTypeString = "CONCOLIC_BOOL"
var concolicMapTypeString = "CONCOLIC_MAP"
var dummyTypeString = "DUMMY_TYPE"

func binExprType(tok token.Token) string {
	switch tok {
	case token.ADD:
		return concolicIntTypeString
	case token.SUB:
		return concolicIntTypeString
	case token.MUL:
		return concolicIntTypeString
	case token.QUO:
		return concolicIntTypeString
	case token.REM:
		return concolicIntTypeString
	case token.AND:
		return concolicIntTypeString
	case token.OR:
		return concolicIntTypeString
	case token.XOR:
		return concolicIntTypeString
	case token.SHL:
		return concolicIntTypeString
	case token.SHR:
		return concolicIntTypeString
	case token.AND_NOT:
		return concolicBoolTypeString
	case token.LAND:
		return concolicBoolTypeString
	case token.LOR:
		return concolicBoolTypeString
	case token.EQL:
		return concolicIntTypeString
	case token.LSS:
		return concolicIntTypeString
	case token.GTR:
		return concolicIntTypeString
	case token.NEQ:
		return concolicIntTypeString
	case token.GEQ:
		return concolicIntTypeString
	case token.LEQ:
		return concolicIntTypeString
	default:
		return dummyTypeString
	}
}

// ===================== PRE =====================

func instrumentBasicLitPre(curNode *astutil.Cursor) {
	castedNode := curNode.Node().(*ast.BasicLit)
	if castedNode.Kind == token.INT {
		augNode :=
			ast.CallExpr{
				Fun: &ast.Ident{
					Name: "concolicTypes.MakeConcolicIntConst",
				},
				Args: []ast.Expr{castedNode},
			}
		curNode.Replace(&augNode)
	} else if castedNode.Kind == token.STRING {
		// TODO
	}
}

func instrumentAssignStmtPre(curNode *astutil.Cursor) bool {
	castedNode := curNode.Node().(*ast.AssignStmt)

	switch castedNode.Lhs[0].(type) {
	case *ast.Ident:
		lhsName := castedNode.Lhs[0].(*ast.Ident).Name
		if _, ok := typeMapping[lhsName]; !ok {
			rhs := castedNode.Rhs[0]
			switch rhs.(type) {
			case *ast.BinaryExpr:
				binTypeString := binExprType(rhs.(*ast.BinaryExpr).Op)
				if binTypeString != dummyTypeString {
					typeMapping[lhsName] = binTypeString
				}

			case *ast.CompositeLit:
				castedRhs := rhs.(*ast.CompositeLit)
				switch castedRhs.Type.(type) {
				case *ast.MapType:
					t := castedRhs.Type.(*ast.MapType)
					_, okKey := t.Key.(*ast.Ident)
					_, okVal := t.Value.(*ast.Ident)
					if okKey && okVal {
						k := t.Key.(*ast.Ident).Name
						v := t.Value.(*ast.Ident).Name
						if k == "int" && v == "int" {
							typeMapping[lhsName] = concolicMapTypeString
						}
					}
				default:
				}
			}
		}

	default:
	}
	return true
}

func instrumentCompositeLitPre(curNode *astutil.Cursor) {
	castedNode := curNode.Node().(*ast.CompositeLit)

	switch castedNode.Type.(type) {
	case *ast.MapType:
		t := castedNode.Type.(*ast.MapType)
		_, okKey := t.Key.(*ast.Ident)
		_, okVal := t.Value.(*ast.Ident)
		if okKey && okVal {
			k := t.Key.(*ast.Ident).Name
			v := t.Value.(*ast.Ident).Name
			if k == "int" && v == "int" {
				newCompLit := &ast.CompositeLit{Type: &ast.MapType{Key: &ast.Ident{Name: "int"}, Value: &ast.Ident{Name: "int"}}}
				newNode := &ast.CallExpr{Fun: &ast.Ident{Name: "concolicTypes.MakeConcolicMapConst"}, Args: []ast.Expr{newCompLit}}
				curNode.Replace(newNode)
			} else {
				curNode.Replace(castedNode)
			}
		}
	default:
	}
}

func instrumentFuncDeclPre(curNode *astutil.Cursor) {
	castedNode := curNode.Node().(*ast.FuncDecl)
	ast.Print(token.NewFileSet(), castedNode)

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
				// TODO fix this
				Name: "concolicTypes.SymStack.SetArgsPopped",
			},
		},
	}
	castedNode.Body.List = append([]ast.Stmt{poppedStatement}, castedNode.Body.List...)

	castedType := castedNode.Type
	newCastedType := &ast.FuncType{
		Func: token.NoPos,
		Results: &ast.FieldList{
			List: []*ast.Field{},
		},
		Params: &ast.FieldList{
			List: []*ast.Field{},
		},
	}

	// TOOD OUTPUT PARAMS
	if castedType.Results != nil && castedType.Results.List != nil {
		for _, aParam := range castedType.Results.List {
			// supposed to look like
			// i = makeConcolicIntVar(cv, "i")
			switch aParam.Type.(*ast.Ident).Name {
			case "concolicTypes.ConcolicString":
				// newCastedType.Results.List = append([]&ast.Field{&ast.Ident{Name: "string"}}, newCastedType.Results...)
				newCastedType.Results.List =
					append(
						[]*ast.Field{
							&ast.Field{
								Type: &ast.Ident{Name: "string"},
								// 		Type: "string",
							},
						},
						newCastedType.Results.List...)
			case "concolicTypes.ConcolicInt":
				// newCastedType.Results.List = append([]ast.Field{&ast.Ident{Name: "int"}}, newCastedType.Results...)
				// aParam.Type = &ast.Ident{Name: "int"}
				newCastedType.Results.List =
					append(
						[]*ast.Field{
							&ast.Field{
								Type: &ast.Ident{Name: "int"},
							},
						},
						newCastedType.Results.List...)
			case "concolicTypes.ConcolicBool":
				newCastedType.Results.List =
					append(
						[]*ast.Field{
							&ast.Field{
								Type: &ast.Ident{Name: "bool"},
							},
						},
						newCastedType.Results.List...)
			default:
				fmt.Printf(aParam.Type.(*ast.Ident).Name + "\r\n")
				// fmt.("WE DON'T SUPPORT THIS TYPE!")
				// if the type is wrong, it's all wrong, so move onto next parameter
				break

			}
		}
	}

	for index1, aParam := range castedType.Params.List {
		aParam = castedType.Params.List[len(castedType.Params.List)-1-index1]
		newParam := &ast.Field{
			// Type:  aParam.Type,
			Names: []*ast.Ident{},
		}
		for index2, aName := range aParam.Names {
			aName = aParam.Names[len(aParam.Names)-1-index2]
			var z3DataType string
			var goDataType string
			var concolName string
			canInstrument := true
			switch aParam.Type.(type) {
			case *ast.Ident:
				switch aParam.Type.(*ast.Ident).Name {
				// case "string":
				// 	fallthrough
				// TODO does this support string correctly
				case "int":
					z3DataType = "Int"
					goDataType = "int"
					concolName = "Int"
				case "bool":
					z3DataType = "Bool"
					goDataType = "bool"
					concolName = "Bool"
				default:
					canInstrument = false
					fmt.Printf(aParam.Type.(*ast.Ident).Name + "\r\n")
					// fmt.("WE DON'T SUPPORT THIS TYPE!")
					// if the type is wrong, it's all wrong, so move onto next parameter
					break

				}
			case *ast.MapType:
				kCast, kOk := aParam.Type.(*ast.MapType).Key.(*ast.Ident)
				vCast, vOk := aParam.Type.(*ast.MapType).Value.(*ast.Ident)
				if kOk && vOk && kCast.Name == "int" && vCast.Name == "int" {
					z3DataType = "Array"
					goDataType = "map[int]int"
					concolName = "Map"
				} else {
					canInstrument = false
				}
			}

			if canInstrument {
				newParam.Names = append(newParam.Names, &ast.Ident{Name: aName.Name + "Val"})
			} else {
				newParam.Names = append(newParam.Names, &ast.Ident{Name: aName.Name})
			}

			// I CHANGED THIS
			// aParam.Type = &ast.Ident{Name: strings.ToLower(methodPiece)}

			newParam.Type = &ast.Ident{Name: goDataType}
			newCastedType.Params.List = append([]*ast.Field{newParam}, newCastedType.Params.List...)
			// fmt.Print("hidbasdf\r\n")
			// ast.Print(token.NewFileSet(), aParam.Type)
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
							Name: "concolicTypes.MakeConcolic" + concolName,
						},
						Args: []ast.Expr{
							&ast.Ident{
								Name: aName.Name + "Val",
							},
							&ast.Ident{
								Name: "concolicTypes.SymStack.PopArg().(z3." + z3DataType + ")",
							},
						},
					},
				},
			}
			castedNode.Body.List = append([]ast.Stmt{&newNode}, castedNode.Body.List...)

			// set each parameter to a different Name
			// aName.Name = aName.Name + "Val"
			castedNode.Type = newCastedType
			curNode.Replace(castedNode)
		}
	}
}

// ===================== POST =====================

func instrumentBinaryExprPost(curNode *astutil.Cursor) {
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
		addedNode.Name = "ConcIntEq"
	case token.LSS:
		addedNode.Name = "ConcIntLT"
	case token.GTR:
		addedNode.Name = "ConcIntGT"
	case token.NEQ:
		addedNode.Name = "ConcIntNE"
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

func instrumentUnaryExprPost(curNode *astutil.Cursor) {
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

func instrumentBasicLitPost(curNode *astutil.Cursor) {

}

func instrumentCompositeLitPost(curNode *astutil.Cursor) {

}

func instrumentAssignStmtPost(curNode *astutil.Cursor) {
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
		// TODO iterate through all
		switch castedNode.Lhs[0].(type) {
		case *ast.IndexExpr:
			lhs := castedNode.Lhs[0].(*ast.IndexExpr)
			switch lhs.X.(type) {
			case *ast.Ident:
				lhsName := lhs.X.(*ast.Ident).Name
				if res, ok := typeMapping[lhsName]; ok && res == concolicMapTypeString {
					// newNode := &ast.Ident{Name: "hi"}
					newNode := &ast.AssignStmt{Lhs: []ast.Expr{&ast.Ident{Name: "_"}},
						Rhs: []ast.Expr{&ast.CallExpr{Fun: &ast.Ident{Name: lhsName + ".ConcMapPut"}, Args: []ast.Expr{lhs.Index, castedNode.Rhs[0]}}},
						Tok: token.ASSIGN,
					}
					curNode.Replace(newNode)
				}
			}
		case *ast.Ident:
			switch castedNode.Rhs[0].(type) {
			case *ast.IndexExpr:
				rhs := castedNode.Rhs[0].(*ast.IndexExpr)
				switch rhs.X.(type) {
				case *ast.Ident:
					rhsName := rhs.X.(*ast.Ident).Name
					if res, ok := typeMapping[rhsName]; ok && res == concolicMapTypeString {
						// newNode := &ast.Ident{Name: "hi"}
						newNode := &ast.AssignStmt{Lhs: []ast.Expr{&ast.Ident{Name: "_"}},
							Rhs: []ast.Expr{&ast.CallExpr{Fun: &ast.Ident{Name: rhsName + ".ConcMapGet"}, Args: []ast.Expr{rhs.Index}}},
							Tok: token.ASSIGN,
						}
						curNode.Replace(newNode)
					}
				}
			}
		default:
			/*
				switch castedNode.Rhs[0].(type) {
				case *ast.CallExpr:
					switch castedNode.Rhs[0].(*ast.CallExpr).Fun.(type) {
					case *ast.FuncLit:
						blah := castedNode.Rhs[0].(*ast.CallExpr).Fun.(*ast.FuncLit).Type.Results.List[0].Type.(*ast.Ident)
						switch blah.Name {
						case "concolicTypes.ConcolicString":
							blah.Name = "string"
						case "concolicTypes.ConcolicInt":
							blah.Name = "int"
						case "concolicTypes.ConcolicBool":
							blah.Name = "bool"
						default:
							// fmt.Printf(aParam.Type.(*ast.Ident).Name + "\r\n")
							// fmt.("WE DON'T SUPPORT THIS TYPE!")
							// if the type is wrong, it's all wrong, so move onto next parameter
							break

						}
						// ast.Print(token.NewFileSet(), castedNode)
						curNode.Replace(castedNode)
					case *ast.Ident:
					default:
					}
				case *ast.IndexExpr:

				case *ast.FuncLit:
				case *ast.CompositeLit:
				default:
				}*/
		}

		// ast.Print(token.NewFileSet(), castedNode.Rhs[0])
	}

	/*
		replacementNode :=
			ast.AssignStmt{
				Tok: curNode.Node().(*ast.AssignStmt).Tok,
				Lhs: castedNode.Lhs,
				Rhs: castedNode.Rhs,
			}
		curNode.Replace(&replacementNode)
	*/
}

func instrumentIncDecStmtPost(curNode *astutil.Cursor) {
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

func instrumentIdentPost(curNode *astutil.Cursor) {
	// switch curNode.Parent().(type) {
	// case *ast.FuncDecl:
	// default:

	// TODO bad idea?
	if _, ok := curNode.Parent().(*ast.Field); !ok {
		castedNode := curNode.Node().(*ast.Ident)
		var varOrConst, concType string
		switch castedNode.Name {
		case "int":
			castedNode.Name = "concolicTypes.ConcolicInt"
			curNode.Replace(castedNode)
		case "bool":
			castedNode.Name = "concolicTypes.ConcolicBool"
			curNode.Replace(castedNode)
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

}

func instrumentIfStmtPost(curNode *astutil.Cursor) {
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

func instrumentFuncDeclPost(curNode *astutil.Cursor) {

}

func instrumentCallExprPost(curNode *astutil.Cursor) bool {
	castedNode := curNode.Node().(*ast.CallExpr)
	switch castedNode.Fun.(type) {
	// case *ast.SelectorExpr:
	case *ast.Ident:
		if castedNode.Fun == nil || castedNode.Fun.(*ast.Ident).Obj == nil || castedNode.Fun.(*ast.Ident).Obj.Decl == nil {
			break
		}
		objectified := castedNode.Fun.(*ast.Ident).Obj.Decl.(*ast.FuncDecl)

		var objectifiedNode *ast.FieldList
		if objectified != nil && objectified.Type != nil {
			objectifiedNode = objectified.Type.Results
			// ast.Print(token.NewFileSet(), objectifiedNode)
		} else {
			objectified = nil
		}

		parNode := curNode.Parent()
		daName, aval := getName(&parNode)
		var daVal string
		switch aval {
		case "int":
			aval = "Int"
			fallthrough
		case "bool":
			aval = "Bool"
			fallthrough
		case "string":
			aval = "String"
			daVal = "concolicTypes.Concolic" + aval
		default:
			daVal = aval

		}
		newNode := &ast.CallExpr{

			Fun: &ast.FuncLit{
				Type: &ast.FuncType{
					Results: objectifiedNode,
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{&ast.Ident{Name: daName + "Val"}},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{castedNode},
							// TODO := or = and actualy make it right with right type
						},
						&ast.DeclStmt{
							Decl: &ast.GenDecl{
								Tok:   token.VAR,
								Specs: []ast.Spec{&ast.TypeSpec{Name: &ast.Ident{Name: daName}, Type: &ast.Ident{Name: daVal}}},
							},
							// TODO fix the typing
						},
					},
				},
			},
		}
		newNode.Fun.(*ast.FuncLit).Body.List = append(
			[]ast.Stmt{&ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.Ident{Name: "concolicTypes.SymStack.SetArgsPushed"},
				},
			},
			},
			newNode.Fun.(*ast.FuncLit).Body.List...)
		// switching order of args
		// newNode.Fun.(*ast.FuncLit).Body.List,
		// &ast.ExprStmt{
		// 	X: &ast.CallExpr{
		// 		Fun: &ast.Ident{Name: "symStack.SetArgsPushed"},
		// 	},
		// },
		// )

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
						// aNameNode.Name += "Val"
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

			// TODO was in here at one point
			// for i, aParam := range newNode.Fun.(*ast.FuncLit).Type.Results.List {
			// 	// supposed to look like
			// 	// i = makeConcolicIntVar(cv, "i")
			// 	switch aParam.Type.(*ast.Ident).Name {
			// 	case "string":
			// 		newNode.Fun.(*ast.FuncLit).Type.Results.List[i].Type = &ast.Ident{Name: "concolicTypes.ConcolicString"}
			// 		newNode.Fun.(*ast.FuncLit).Type.Results.List[i].Type = &ast.Ident{Name: "concolicTypes.ConcolicString"}
			// 	case "int":
			// 		// fmt.Printf("mother of gawd")
			// 		newNode.Fun.(*ast.FuncLit).Type.Results.List[i].Type = &ast.Ident{Name: "concolicTypes.ConcolicInt", NamePos: token.NoPos}

			// 		// ast.Print(token.NewFileSet(), newNode)
			// 	case "bool":
			// 		newNode.Fun.(*ast.FuncLit).Type.Results.List[i].Type = &ast.Ident{Name: "concolicTypes.ConcolicBool"}
			// 	default:
			// 		fmt.Printf(aParam.Type.(*ast.Ident).Name + "\r\n")
			// 		// fmt.("WE DON'T SUPPORT THIS TYPE!")
			// 		// if the type is wrong, it's all wrong, so move onto next parameter
			// 		break

			daName, _ := getName(&parNode)
			afterIf := instrumentParentOfCallExpr(curNode)
			if afterIf != nil {
				newNode.Fun.(*ast.FuncLit).Body.List = append(newNode.Fun.(*ast.FuncLit).Body.List, afterIf)

			}
			newNode.Fun.(*ast.FuncLit).Body.List =
				append(
					newNode.Fun.(*ast.FuncLit).Body.List,
					&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: daName}}})
			// 	}
			// }
			for _, aParam := range objectifiedNode.List {
				switch aParam.Type.(type) {
				case *ast.Ident:
					switch aParam.Type.(*ast.Ident).Name {
					case "concolicTypes.ConcolicString":
						aParam.Type.(*ast.Ident).Name = "string"
						// aParam.Type.(*ast.Ident).Obj.Name = "string"
						// aParam.Type.(*ast.Ident).NamePos = token.NoPos
					case "concolicTypes.ConcolicInt":
						aParam.Type.(*ast.Ident).Name = "int"
						// aParam.Type.(*ast.Ident).Obj.Name = "int"
						// aParam.Type.(*ast.Ident).NamePos = token.NoPos
					case "concolicTypes.ConcolicBool":
						aParam.Type.(*ast.Ident).Name = "bool"
						// aParam.Type.(*ast.Ident).Obj.Name = "bool"
						// aParam.Type.(*ast.Ident).NamePos = token.NoPos
						// aParam.Type.(*ast.Ident).Obj.
					default:
						fmt.Printf(aParam.Type.(*ast.Ident).Name + "\r\n")
						// fmt.("WE DON'T SUPPORT THIS TYPE!")
						// if the type is wrong, it's all wrong, so move onto next parameter

					}

				}
			}
			curNode.Replace(newNode)
		}

	case *ast.SelectorExpr:
		castedChild := castedNode.Fun.(*ast.SelectorExpr)
		switch castedChild.X.(type) {
		case *ast.Ident:
			castedGrandChild := castedChild.X.(*ast.Ident)
			if castedGrandChild.Name == "concolicTypes" {
				switch castedChild.Sel.Name {
				case "MakeFuzzyInt":
					castedChild.Sel.Name = "MakeConcolicIntVar"
					// castedNode.Fun.(*ast.SelectorExpr).Sel.(*ast.Ident).Name = "MakeConcolicIntVar"
				case "MakeFuzzyString":
					castedChild.Sel.Name = "MakeConcolicStringVar"
				case "MakeFuzzyBool":
					castedChild.Sel.Name = "MakeConcolicBoolVar"
				case "MakeFuzzyMapIntInt":
					castedChild.Sel.Name = "MakeConcolicMapVar"
				default:
					return true
				}
				castedNode.Args = []ast.Expr{castedNode.Args[0]}
				curNode.Replace(castedNode)
			}
		}
	case *ast.FuncLit:
	case *ast.CallExpr:
	default:
		// ast.Print(token.NewFileSet(), curNode.Node())
		panic("not supported!")
	}
	// TODO
	return true

}

func instrumentReturnStmtPost(curNode *astutil.Cursor) {
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
	castedNode.Name.Name = "InstrumentedMainMethod"
}

// right now only correctly implemented for a single node
func getName(parNode *ast.Node) (string, string) {
	switch (*parNode).(type) {
	case *ast.AssignStmt:
		castedParentNode := (*parNode).(*ast.AssignStmt)
		actualName := "_"
		actualType := ""
		for _, val := range castedParentNode.Lhs {
			switch val.(type) {
			case *ast.Ident:
				actualName = val.(*ast.Ident).Name
				// decl := val.(*ast.Ident).Obj.Decl
				// ast.Print(token.NewFileSet(), decl)
				// 				switch decl.(type){
				// case
				// 				}

			}
		}
		if theRhs, ok1 := castedParentNode.Rhs[0].(*ast.CallExpr); ok1 {
			if theIdent, ok2 := theRhs.Fun.(*ast.Ident); ok2 {
				if theDecl, ok3 := theIdent.Obj.Decl.(*ast.FuncDecl); ok3 {
					for _, aResult := range theDecl.Type.Results.List {
						if castedResult, ok4 := aResult.Type.(*ast.Ident); ok4 {
							actualType = castedResult.Name
						}
					}
				}

			}

		}
		return actualName, actualType
	default:
		// probably shoudlnt augment/shouldn't hit this
		return "", ""
	}
}

func instrumentParentOfCallExpr(curNode *astutil.Cursor) *ast.IfStmt {
	// castedNode := curNode.Node().(*ast.CallExpr)
	parentNode := curNode.Parent()
	switch parentNode.(type) {
	case *ast.AssignStmt:
		castedParentNode := parentNode.(*ast.AssignStmt)
		actualName := "_"
		for _, val := range castedParentNode.Lhs {
			switch val.(type) {
			case *ast.Ident:
				actualName = val.(*ast.Ident).Name
				// castedParentNode.Lhs[index] = &ast.Ident{Name: val.(*ast.Ident).Name + "Val"}
			}
		}
		// nextNode := &ast.BlockStmt{
		// 	List: []ast.Stmt{
		// 		castedParentNode,
		// 	},
		// }
		// if castedParentNode.Tok == token.DEFINE {
		// 	declStatement := &ast.DeclStmt{
		// 		Decl: &ast.GenDecl{
		// 			Tok: token.VAR,
		// 			Specs: []ast.Spec{
		// 				&ast.TypeSpec{
		// 					Name: &ast.Ident{
		// 						Name: actualName,
		// 					},
		// 				},
		// 			},
		// 		},
		// 	}
		// 	nextNode.List = append(nextNode.List, declStatement)
		// }
		// TODO if  statement
		ifStmnt := &ast.IfStmt{
			// TODO this is hax
			Cond: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "concolicTypes.SymStack",
					},
					Sel: &ast.Ident{
						Name: "AreArgsPushed",
					},
				},
			},
			// TODO types "make concolic int string OR BOOL "
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.Ident{Name: actualName},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.Ident{
									Name: "concolicTypes.MakeConcolicIntConst",
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: actualName + "Val",
									},
								},
							},
						},
					},
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "concolicTypes.SymStack",
								},
								Sel: &ast.Ident{
									Name: "ClearArgs",
								},
							},
						},
					},
				},
			},
			Else: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.Ident{Name: actualName},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.Ident{
									// TODO
									Name: "concolicTypes.MakeConcolicInt",
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: actualName + "Val",
									},
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "concolicTypes.SymStack",
											},
											Sel: &ast.Ident{
												Name: "PopReturn().",
											},
										},
										Args: []ast.Expr{
											&ast.SelectorExpr{
												X: &ast.Ident{
													Name: "z3",
												},
												Sel: &ast.Ident{
													// TODO types
													Name: "Int",
												},
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
		return ifStmnt
		// nextNode.List = append(nextNode.List, ifStmnt)
		// curNode.InsertAfter(ifStmnt)
		// TOOD gotta put this somewhere boi

		// queueOfThings.stage1.Push(curNode.Parent().GetId().Id, nextNode)

		// pre := func(curNode *astutil.Cursor) bool {
		// 	return true
		// }

		// bruh := true
		// castedParentNode. = nextNode

		// post := func(cn *astutil.Cursor) bool {
		// 	if cn.Node() == parentNode {
		// 		// doesn't traverse childrens so we gucci
		// 		// TODO . replace
		// 		fmt.Print("FUCK ME")
		// 		cn.Replace(nextNode)
		// 		bruh = false
		// 	}
		// 	return true
		// }
		// astutil.Apply(parentNode, nil, astutil.ApplyFunc(post))
		// TODO
	default:
		return nil
	}
}

func instrumentDeclParentCheckPre(curNode *astutil.Cursor) {
	parNode := curNode.Parent()
	switch parNode.(type) {
	case *ast.File:
		typeMapping = make(map[string]string)
	}
}
