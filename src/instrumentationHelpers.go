package main

import (
	_ "fmt"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
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
		// because of z3 we don't really support strings so shouldn't need to incldue anything here
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
		Func:    token.NoPos,
		Results: castedType.Results,
		Params: &ast.FieldList{
			List: []*ast.Field{},
		},
	}

	for index1, aParam := range castedType.Params.List {
		aParam = castedType.Params.List[len(castedType.Params.List)-1-index1]
		newParam := &ast.Field{
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
				// TODO this doesn't quite this support string correctly; but z3 doesn't either
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
					// fmt.("WE DON'T SUPPORT THIS TYPE!")
					// if the type is wrong, move onto next parameter
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
				switch concolName {
				case "Int":
					typeMapping[aName.Name] = concolicIntTypeString
				case "Bool":
					typeMapping[aName.Name] = concolicBoolTypeString
				case "Map":
					typeMapping[aName.Name] = concolicMapTypeString
				}
			} else {
				newParam.Names = append(newParam.Names, &ast.Ident{Name: aName.Name})
			}

			newParam.Type = &ast.Ident{Name: goDataType}
			newCastedType.Params.List = append([]*ast.Field{newParam}, newCastedType.Params.List...)
			// add in concolic constructors before the _ thingies
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

			castedNode.Type = newCastedType
			curNode.Replace(castedNode)
		}
	}
}

// ===================== POST =====================

func instrumentBinaryExprPost(curNode *astutil.Cursor) {
	castedNode := curNode.Node().(*ast.BinaryExpr)

	addedNode := &ast.Ident{
		Name: "",
	}

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
	// !, ^, and - are the only unary expressions we support
	castedNode := curNode.Node().(*ast.UnaryExpr)
	switch castedNode.Op {
	case token.SUB:
		switch castedNode.X.(type) {
		case *ast.Ident:
			curNode.Parent().(*ast.AssignStmt).Lhs = append([]ast.Expr{&ast.Ident{Name: "_"}}, curNode.Parent().(*ast.AssignStmt).Lhs...)
			newNode := &ast.SelectorExpr{
				X: castedNode.X,
				Sel: &ast.Ident{
					Name: "ConcIntMul(concolicTypes.MakeConcolicIntConst(-1))",
				},
			}

			curNode.InsertAfter(newNode)
		case *ast.CallExpr:
			castedChild := castedNode.X.(*ast.CallExpr)
			castedChild.Args[0].(*ast.BasicLit).Value = "-" + castedNode.X.(*ast.CallExpr).Args[0].(*ast.BasicLit).Value

		}
		curNode.Replace(castedNode.X)
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
	// for the sake of completeness
}

func instrumentCompositeLitPost(curNode *astutil.Cursor) {
	// for the sake of completeness
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
		switch castedNode.Lhs[0].(type) {
		case *ast.IndexExpr:
			lhs := castedNode.Lhs[0].(*ast.IndexExpr)
			switch lhs.X.(type) {
			case *ast.Ident:
				lhsName := lhs.X.(*ast.Ident).Name
				if res, ok := typeMapping[lhsName]; ok && res == concolicMapTypeString {
					newNode := &ast.AssignStmt{
						Lhs: []ast.Expr{&ast.Ident{Name: "_"}},
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
						newNode := &ast.AssignStmt{Lhs: []ast.Expr{&ast.Ident{Name: "_"}},
							Rhs: []ast.Expr{&ast.CallExpr{Fun: &ast.Ident{Name: rhsName + ".ConcMapGet"}, Args: []ast.Expr{rhs.Index}}},
							Tok: token.ASSIGN,
						}
						curNode.Replace(newNode)
					}
				}
			}
		default:
		}
	}
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
			// TODO it might be better to SHA this
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
	} else {
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
			},
		}

	}
}

func instrumentFuncDeclPost(curNode *astutil.Cursor) {
	// here for completeness sake
}

func instrumentCallExprPost(curNode *astutil.Cursor) bool {
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
			objectifiedNode = nil
		}

		parNode := curNode.Parent()
		daName, aval := getName(&parNode)
		switch aval {
		case "int":
			aval = "concolicTypes.ConcolicInt"
		case "bool":
			aval = "concolicTypes.ConcolicBool"
		default:

		}

		var newNode *ast.CallExpr
		newFiledList := &ast.FieldList{
			List: []*ast.Field{},
		}

		if objectifiedNode != nil {
			for _, aField := range objectifiedNode.List {
				nameToConvert := aField.Type.(*ast.Ident).Name
				var convertedName string
				switch nameToConvert {
				case "string":
					convertedName = "concolicTypes.ConcolicString"
				case "bool":
					convertedName = "concolicTypes.ConcolicBool"
				case "int":
					convertedName = "concolicTypes.ConcolicInt"
				case "map":
					// TODO HALP CHRIS
					convertedName = "concolicTypes.ConcolicMap"
				}

				newFiledList.List = append(newFiledList.List, &ast.Field{
					Type: &ast.Ident{
						Name: convertedName,
					},
				})

			}

		}

		if daName != "" {

			newNode = &ast.CallExpr{

				Fun: &ast.FuncLit{
					Type: &ast.FuncType{
						Results: newFiledList,
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{&ast.Ident{Name: daName + "Val"}},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{castedNode},
								// TODO := or = and make the type for clearly defined make it right with right type
							},
							&ast.DeclStmt{
								Decl: &ast.GenDecl{
									Tok:   token.VAR,
									Specs: []ast.Spec{&ast.TypeSpec{Name: &ast.Ident{Name: daName}, Type: &ast.Ident{Name: aval}}},
								},
								// TODO fix the typing
							},
						},
					},
				},
			}
		} else {
			newNode = &ast.CallExpr{

				Fun: &ast.FuncLit{
					Type: &ast.FuncType{
						Results: newFiledList,
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: castedNode,
							},
							// TODO := or = and actualy make it right with right type
						},
					},
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

		if objectified != nil {
			declaration := castedNode.Fun.(*ast.Ident).Obj.Decl.(*ast.FuncDecl)
			paramList := declaration.Type.Params.List

			for index, aParam := range paramList {
				containedInMap := false
				if _, ok := typeMapping[castedNode.Args[index].(*ast.Ident).Name]; ok {
					containedInMap = true
				}
				if aParam.Type.(*ast.Ident).Name == "string" || aParam.Type.(*ast.Ident).Name == "int" || aParam.Type.(*ast.Ident).Name == "bool" || (len(aParam.Type.(*ast.Ident).Name) >= 3 && aParam.Type.(*ast.Ident).Name[0:3] == "map" && containedInMap) {
					for range aParam.Names {
						newNode.Fun.(*ast.FuncLit).Body.List = append(
							[]ast.Stmt{
								&ast.ExprStmt{
									X: &ast.CallExpr{
										Fun: &ast.Ident{
											Name: "concolicTypes.SymStack.PushArg",
										},
										Args: []ast.Expr{
											&ast.Ident{
												// TODO this is a hack-- fix it by replacing ident with selectorstatement
												Name: castedNode.Args[index].(*ast.Ident).Name + ".Z3Expr",
											},
										},
									},
								},
							}, newNode.Fun.(*ast.FuncLit).Body.List...)
					}
					castedNode.Args[index] = &ast.SelectorExpr{
						X: castedNode.Args[index],
						Sel: &ast.Ident{
							Name: "Value",
						},
					}
				}
			}

			daName, releType := getName(&parNode)
			afterIf := instrumentParentOfCallExpr(curNode, releType)
			if afterIf != nil && len(newNode.Fun.(*ast.FuncLit).Body.List) > 0 {
				newNode.Fun.(*ast.FuncLit).Body.List = append(newNode.Fun.(*ast.FuncLit).Body.List, afterIf)

			}
			if daName != "" {
				newNode.Fun.(*ast.FuncLit).Body.List =
					append(
						newNode.Fun.(*ast.FuncLit).Body.List,
						&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: daName}}})
			} else {
			}
			if objectifiedNode != nil && objectifiedNode.List != nil {
				for _, aParam := range objectifiedNode.List {
					switch aParam.Type.(type) {
					case *ast.Ident:
						switch aParam.Type.(*ast.Ident).Name {
						case "concolicTypes.ConcolicString":
							aParam.Type.(*ast.Ident).Name = "string"
						case "concolicTypes.ConcolicInt":
							aParam.Type.(*ast.Ident).Name = "int"
						case "concolicTypes.ConcolicBool":
							aParam.Type.(*ast.Ident).Name = "bool"
						default:
							// if the type is wrong, so move onto next parameter

						}

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

// TODO thinking about other types
func instrumentParentOfCallExpr(curNode *astutil.Cursor, rType string) *ast.IfStmt {
	var releType string
	switch rType {
	case "string":
		releType = "String"
	case "int":
		releType = "Int"
	case "map":
		releType = "Map"
	case "bool":
		releType = "Bool"

	}
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
			}
		}
		ifStmnt := &ast.IfStmt{
			// TODO this is hax; make it cleaner with more selector statemetns
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
									Name: "concolicTypes.MakeConcolic" + releType + "Const",
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
									Name: "concolicTypes.MakeConcolic" + releType,
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
													Name: releType,
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

// 0 is LHS
// 1 is RHS
// 2 is not found
func whichSideAmIOn(curNode *astutil.Cursor) int {
	castedCurNode := curNode.Node().(*ast.IndexExpr)
	if _, ok := curNode.Parent().(*ast.AssignStmt); !ok {
		return 1

	}
	castedParNode, _ := curNode.Parent().(*ast.AssignStmt)
	oldVal := castedCurNode.Lbrack
	castedCurNode.Lbrack = 6969696969696969
	curNode.Replace(castedCurNode)
	rVal := 2
	// explore lhs
	for _, theExpr := range castedParNode.Lhs {
		if ca, ok := theExpr.(*ast.IndexExpr); ok {
			if ca.Lbrack == 6969696969696969 {
				rVal = 0
				goto finished
			}
		}
	}
	// explore rhs
	for _, theExpr := range castedParNode.Rhs {
		if ca, ok := theExpr.(*ast.IndexExpr); ok {
			if ca.Lbrack == 6969696969696969 {
				rVal = 1
				goto finished
			}
		}
	}
finished:
	castedCurNode.Lbrack = oldVal
	curNode.Replace(castedCurNode)
	return rVal
}

func instrumentIndexExprPre(curNode *astutil.Cursor) {
	castedNode := curNode.Node().(*ast.IndexExpr)
	switch whichSideAmIOn(curNode) {
	case 0:
	case 1:
		// rhs turns into GET
		replacementNode :=
			ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   castedNode.X,
					Sel: &ast.Ident{Name: "ConcMapGet"},
				},
				Args: []ast.Expr{
					castedNode.Index,
				},
			}
		curNode.Replace(&replacementNode)
	default:
		panic("o no")
	}
}
