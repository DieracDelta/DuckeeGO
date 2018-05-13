package main

import (
	"github.com/otiai10/copy"
	"os"
	"strings"

	_ "fmt"

	"io/ioutil"

	"bytes"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
)

// argment is path to example program
var DEST = "./tmp/DuckieConcolic/"
var VERBOSE = false

func main() {

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
	typeMapping = make(map[string]string)

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

		if VERBOSE {
			_ = ast.Print(fset, uninstrumentedAST)
		}
		instrumentedAST := astutil.Apply(uninstrumentedAST, astutil.ApplyFunc(addInstrumentationPre), astutil.ApplyFunc(addInstrumentationPost))

		if VERBOSE {
			_ = ast.Print(fset, instrumentedAST)
		}
		var buf bytes.Buffer
		err = printer.Fprint(&buf, fset, instrumentedAST)
		if err != nil {
			panic(err)
		}
		// fmt.Println(buf.String())
		if strings.Contains(aGoFile.FilePath, "main") {

			aGoFile.FilePath = strings.Replace(aGoFile.FilePath, "main", "userMain", 1)

		}
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
											Value: "\"InstrumentedMainMethod\"",
										},
									},
								},
							},
						},
						&ast.ExprStmt{
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
						},
					},
				},
			},
		},
	}
	stuff.Imports = []*ast.ImportSpec{a, b}
	return stuff
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
