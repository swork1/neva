package portgen

import (
	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/runtime"
)

type PortGen struct{}

// TODO use map as input?
func (p PortGen) Ports(addrs []runtime.PortAddr) map[runtime.PortAddr]chan core.Msg {
	ports := make(map[runtime.PortAddr]chan core.Msg, len(addrs))

	for i := range addrs {
		ports[addrs[i]] = make(chan core.Msg)
	}

	return ports
}

func New() PortGen {
	return PortGen{}
}