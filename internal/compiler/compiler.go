package compiler

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/src"
	"github.com/nevalang/neva/pkg/ir"
)

type Compiler struct {
	repo     Repository
	parser   Parser
	analyzer Analyzer
	irgen    IRGenerator
}

type (
	Repository interface {
		ByPath(context.Context, string) (map[string]RawPackage, error)
		Save(context.Context, string, *ir.Program) error
	}
	Parser interface {
		ParseProg(context.Context, map[string]RawPackage) (src.Program, error)
	}
	Analyzer interface {
		Analyze(src.Program) (src.Program, error)
	}
	IRGenerator interface {
		Generate(context.Context, src.Program) (*ir.Program, error)
	}

	RawPackage map[string][]byte // File -> content.
)

func (c Compiler) Compile(ctx context.Context, srcPath, dstPath string) (*ir.Program, error) {
	rawProg, err := c.repo.ByPath(ctx, srcPath)
	if err != nil {
		return nil, fmt.Errorf("repo by path: %w", err)
	}

	parsedProg, err := c.parser.ParseProg(ctx, rawProg)
	if err != nil {
		return nil, fmt.Errorf("parse prog: %w", err)
	}

	analyzedProg, err := c.analyzer.Analyze(parsedProg)
	if err != nil {
		return nil, fmt.Errorf("analyze: %w", err)
	}

	irProg, err := c.irgen.Generate(ctx, analyzedProg)
	if err != nil {
		return nil, fmt.Errorf("generate IR: %w", err)
	}

	if dstPath == "" {
		return irProg, nil
	}

	if err := c.repo.Save(ctx, dstPath, irProg); err != nil {
		return nil, fmt.Errorf("repo save: %w", err)
	}

	return irProg, nil
}

func New(
	repo Repository,
	parser Parser,
	analyzer Analyzer,
	irgen IRGenerator,
) Compiler {
	return Compiler{
		repo:     repo,
		parser:   parser,
		analyzer: analyzer,
		irgen:    irgen,
	}
}
