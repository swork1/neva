package analyzer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/src"
)

var (
	ErrMainComponentWithTypeParams     = errors.New("main component can't have type parameters")
	ErrMainComponentNodes              = errors.New("something wrong with main component's nodes")
	ErrEntityNotFoundByNodeRef         = errors.New("entity not found by node ref")
	ErrMainComponentInportsCount       = errors.New("main component must have one inport")
	ErrMainComponentOutportsCount      = errors.New("main component must have exactly one outport")
	ErrMainComponentWithoutEnterInport = errors.New("main component must have 'enter' inport")
	ErrMainComponentWithoutExitOutport = errors.New("main component must have 'exit' outport")
	ErrMainPortIsArray                 = errors.New("main component's ports cannot not be arrays")
	ErrMainComponentPortTypeNotAny     = errors.New("main component's ports must be of type any")
	ErrMainComponentNodeNotComponent   = errors.New("main component's nodes must be components only")
)

func (a Analyzer) analyzeMainComponent(cmp src.Component, pkg src.Package, scope src.Scope) error {
	if len(cmp.Interface.TypeParams) != 0 {
		return fmt.Errorf("%w: %v", ErrMainComponentWithTypeParams, cmp.Interface.TypeParams)
	}

	if err := a.analyzeMainComponentIO(cmp.Interface.IO); err != nil {
		return fmt.Errorf("main component io: %w", err)
	}

	if err := a.analyzeMainComponentNodes(cmp.Nodes, pkg, scope); err != nil {
		return fmt.Errorf("%w: %v", ErrMainComponentNodes, err)
	}

	return nil
}

func (a Analyzer) analyzeMainComponentIO(io src.IO) error {
	if len(io.Out) != 1 {
		return fmt.Errorf("%w: %v", ErrMainComponentOutportsCount, io.Out)
	}
	if len(io.In) != 1 {
		return fmt.Errorf("%w: %v", ErrMainComponentInportsCount, io.In)
	}

	enterInport, ok := io.In["enter"]
	if !ok {
		return ErrMainComponentWithoutEnterInport
	}
	if err := a.analyzeMainComponentPort(enterInport); err != nil {
		return fmt.Errorf("enter inport: %w", err)
	}

	exitInport, ok := io.Out["exit"]
	if !ok {
		return ErrMainComponentWithoutExitOutport
	}
	if err := a.analyzeMainComponentPort(exitInport); err != nil {
		return fmt.Errorf("exit outport: %w", err)
	}

	return nil
}

func (a Analyzer) analyzeMainComponentPort(port src.Port) error {
	if port.IsArray {
		return ErrMainPortIsArray
	}
	if !(src.Scope{}).IsTopType(port.TypeExpr) {
		return ErrMainComponentPortTypeNotAny
	}
	return nil
}

func (Analyzer) analyzeMainComponentNodes(nodes map[string]src.Node, pkg src.Package, scope src.Scope) error {
	for nodeName, node := range nodes {
		nodeEntity, _, err := scope.Entity(node.EntityRef)
		if err != nil {
			return fmt.Errorf(
				"%w: node name %v: entity ref %v: %v",
				ErrEntityNotFoundByNodeRef,
				nodeName,
				node.EntityRef,
				err,
			)
		}
		if nodeEntity.Kind != src.ComponentEntity {
			return fmt.Errorf("%w: %v: %v", ErrMainComponentNodeNotComponent, nodeName, node.EntityRef)
		}
	}
	return nil
}