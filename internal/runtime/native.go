package runtime

type NativeModule struct {
	in   InportsInterface
	out  OutportsInterface
	impl func(NodeIO)
}

func (a NativeModule) SpawnWorker(env map[string]Module) (NodeIO, error) {
	io := createIO(a.in, a.out)
	go a.impl(io)
	return io, nil
}

func (a NativeModule) Interface() ModuleInterface {
	return ModuleInterface{
		In:  a.in,
		Out: a.out,
	}
}

func NewNativeModule(
	in InportsInterface,
	out OutportsInterface,
	impl func(NodeIO),
) NativeModule {
	return NativeModule{
		in:   in,
		out:  out,
		impl: impl,
	}
}

func createIO(in InportsInterface, out OutportsInterface) NodeIO {
	inports := make(map[string]chan Msg, len(in))
	outports := make(map[string]chan Msg, len(in))

	for k := range in {
		inports[k] = make(chan Msg)
	}
	for k := range out {
		inports[k] = make(chan Msg)
	}

	return NodeIO{inports, outports}
}