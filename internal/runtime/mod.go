package runtime

import "fbp/internal/types"

// AbstractModule represents complex and atomic modules.
type AbstractModule interface {
	Run(in map[string]<-chan Msg, out map[string]chan<- Msg)
	Ports() (InPorts, OutPorts)
}

// ComplexModule is a composition of other modules.
type ComplexModule struct {
	in  InPorts
	out OutPorts
	wm  Workers // TODO: do we need it?
	net []Conn
}

func (cm ComplexModule) Ports() (InPorts, OutPorts) {
	return cm.in, cm.out
}

func (m ComplexModule) Run(in map[string]chan Msg, out map[string]chan Msg) {
	// TODO: pass in-out to connect-all?
	ConnectAll(m.net)
}

type InPorts Ports

type OutPorts Ports

type Workers map[string]AbstractModule

type Env map[string]AbstractModule

type Conn struct {
	Sender    <-chan Msg   // outPort
	Receivers []chan<- Msg // inPorts
}

type Ports map[string]types.Type

func NewModule(
	in InPorts,
	out OutPorts,
	wm Workers,
	net []Conn,
) ComplexModule {
	return ComplexModule{
		in:  in,
		out: out,
		wm:  wm,
		net: net,
	}
}

// AtomicModule represents module with native implementation.
type AtomicModule struct {
	in   InPorts
	out  OutPorts
	impl func(
		in map[string]<-chan Msg,
		out map[string]chan<- Msg,
	)
}

func (nm AtomicModule) Run(in map[string]<-chan Msg, out map[string]chan<- Msg) {
	nm.impl(in, out)
}

func (nm AtomicModule) Ports() (InPorts, OutPorts) {
	return nm.in, nm.out
}

func NewAtomicModule(
	in InPorts,
	out OutPorts,
	impl func(
		in map[string]<-chan Msg,
		out map[string]chan<- Msg,
	),
) AtomicModule {
	return AtomicModule{
		in:   in,
		out:  out,
		impl: impl,
	}
}
