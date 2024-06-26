// Package runtime implements environment for dataflow programs execution.
package runtime

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type Runtime struct {
	connector  Connector
	funcRunner FuncRunner
}

var ErrNilDeps = errors.New("runtime deps nil")

func New(connector Connector, funcRunner FuncRunner) Runtime {
	return Runtime{
		connector:  connector,
		funcRunner: funcRunner,
	}
}

var (
	ErrStartPortNotFound = errors.New("start port not found")
	ErrExitPortNotFound  = errors.New("stop port not found")
	ErrConnector         = errors.New("connector")
	ErrFuncRunner        = errors.New("func runner")
)

func (r Runtime) Run(ctx context.Context, prog Program) error {
	enter := prog.Ports[PortAddr{Path: "in", Port: "start"}]
	if enter == nil {
		return ErrStartPortNotFound
	}

	exit := prog.Ports[PortAddr{Path: "out", Port: "stop"}]
	if exit == nil {
		return ErrExitPortNotFound
	}

	funcRun, err := r.funcRunner.Run(prog.Funcs)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFuncRunner, err)
	}

	cancelableCtx, cancel := context.WithCancel(ctx)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		funcRun(
			context.WithValue(
				cancelableCtx,
				"cancel", //nolint:staticcheck // SA1029
				cancel,
			),
		)
		wg.Done()
	}()

	go func() {
		r.connector.Connect(cancelableCtx, prog.Connections)
		wg.Done()
	}()

	go func() {
		enter <- &baseMsg{}
	}()

	go func() {
		<-exit
		cancel()
	}()

	wg.Wait()

	return nil
}
