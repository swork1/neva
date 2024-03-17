// Package parser implements source code parsing.
// It uses parser (and lexer) generated by ANTLR4 from neva.g4 grammar file.
package parser

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/antlr4-go/antlr/v4"

	"github.com/nevalang/neva/internal/compiler"
	generated "github.com/nevalang/neva/internal/compiler/parser/generated"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

type treeShapeListener struct {
	*generated.BasenevaListener
	file src.File
}

type Parser struct {
	isDebug bool
}

func (p Parser) ParseModules(
	rawMods map[src.ModuleRef]compiler.RawModule,
) (map[src.ModuleRef]src.Module, *compiler.Error) {
	parsedMods := make(map[src.ModuleRef]src.Module, len(rawMods))

	for modRef, rawMod := range rawMods {
		parsedPkgs, err := p.ParsePackages(modRef, rawMod.Packages)
		if err != nil {
			return nil, compiler.Error{
				Err:      errors.New("Parsing error"),
				Location: &src.Location{ModRef: modRef},
			}.Wrap(err)
		}

		parsedMods[modRef] = src.Module{
			Manifest: rawMod.Manifest,
			Packages: parsedPkgs,
		}
	}

	return parsedMods, nil
}

func (p Parser) ParsePackages(
	modRef src.ModuleRef,
	rawPkgs map[string]compiler.RawPackage,
) (
	map[string]src.Package,
	*compiler.Error,
) {
	packages := make(map[string]src.Package, len(rawPkgs))

	for pkgName, pkgFiles := range rawPkgs {
		parsedFiles, err := p.ParseFiles(modRef, pkgName, pkgFiles)
		if err != nil {
			return nil, compiler.Error{
				Location: &src.Location{PkgName: pkgName},
			}.Wrap(err)
		}

		packages[pkgName] = parsedFiles
	}

	return packages, nil
}

func (p Parser) ParseFiles(
	modRef src.ModuleRef,
	pkgName string,
	files map[string][]byte,
) (map[string]src.File, *compiler.Error) {
	result := make(map[string]src.File, len(files))

	for fileName, fileBytes := range files {
		loc := src.Location{
			ModRef:   modRef,
			PkgName:  pkgName,
			FileName: fileName,
		}
		parsedFile, err := p.parseFile(loc, fileBytes)
		if err != nil {
			return nil, compiler.Error{Location: &loc}.Wrap(err)
		}
		result[fileName] = parsedFile
	}

	return result, nil
}

func (p Parser) parseFile(
	loc src.Location,
	bb []byte,
) (f src.File, err *compiler.Error) {
	defer func() {
		if e := recover(); e != nil {
			compilerErr, ok := e.(*compiler.Error)
			if ok {
				err = compiler.Error{Location: &loc}.Wrap(compilerErr)
				return
			}
			err = &compiler.Error{
				Err: fmt.Errorf(
					"%v: %v",
					e,
					string(debug.Stack()),
				),
				Location: &loc,
			}
		}
	}()

	input := antlr.NewInputStream(string(bb))
	lexer := generated.NewnevaLexer(input)
	lexerErrors := &CustomErrorListener{}
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(lexerErrors)
	tokenStream := antlr.NewCommonTokenStream(lexer, 0)

	parserErrors := &CustomErrorListener{}
	prsr := generated.NewnevaParser(tokenStream)
	prsr.RemoveErrorListeners()
	prsr.AddErrorListener(parserErrors)
	if p.isDebug {
		prsr.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	}
	prsr.BuildParseTrees = true

	tree := prsr.Prog()
	listener := &treeShapeListener{}

	antlr.ParseTreeWalkerDefault.Walk(listener, tree)

	if len(lexerErrors.Errors) > 0 {
		return src.File{}, &compiler.Error{
			Err:      errors.Join(lexerErrors.Errors...),
			Location: &src.Location{},
			Meta:     &src.Meta{},
		}
	}
	if len(parserErrors.Errors) > 0 {
		return src.File{}, &compiler.Error{
			Err:      errors.Join(parserErrors.Errors...),
			Location: &src.Location{},
			Meta:     &src.Meta{},
		}
	}

	return listener.file, nil
}

func New(isDebug bool) Parser {
	return Parser{isDebug: isDebug}
}
