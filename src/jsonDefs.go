package main

import "go/ast"

type ConfigData struct {
	ConfigData  []GoFile `json:"Data"`
	ProjectPath string   `json:"ProjectPath"`
	Package     string   `json:"Package"`
}

type GoFile struct {
	FilePath  string       `json:"RelFilePath"`
	Functions []GoFunction `json:"Functions"`
}

type GoFunction struct {
	Name   string   `json:"Name"`
	Object string   `json:"Object"`
	Input  []string `json:"Input"`
	Output []string `json:"Output"`
	Node   *ast.Node
}
