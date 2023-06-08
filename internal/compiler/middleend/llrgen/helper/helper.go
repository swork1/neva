package helper

import (
	"github.com/nevalang/neva/internal/compiler"
	ts "github.com/nevalang/neva/pkg/types"
)

type Helper struct {
	ts.Helper
}

/* --- COMPONENTS  --- */

// MainComponent returns entity of kind "component" with main component type params and io
func (h Helper) MainComponent(nodes map[string]compiler.Node, net []compiler.Connection) compiler.Entity {
	return compiler.Entity{
		Kind: compiler.ComponentEntity,
		Component: compiler.Component{
			TypeParams: []ts.Param{
				h.ParamWithNoConstr("t"),
			},
			IO: compiler.IO{
				In: map[string]compiler.Port{
					"start": {
						Type: h.Rec(nil), // TODO any?
					},
				},
				Out: map[string]compiler.Port{
					"exit": {
						Type: h.Inst("int"),
					},
				},
			},
			Nodes: nodes,
			Net:   net,
		},
	}
}

func (h Helper) Node(instance compiler.Instance) compiler.Node {
	return compiler.Node{
		Instance: instance,
	}
}

func (h Helper) NodeWithStaticPorts(
	instance compiler.Node,
	ports map[compiler.RelPortAddr]compiler.EntityRef,
) compiler.Node {
	return compiler.Node{
		Instance:      instance,
		StaticInports: ports,
	}
}

func (h Helper) NodeInstance(pkg, entity string, args ...ts.Expr) compiler.Instance {
	return compiler.Instance{
		Ref: compiler.EntityRef{
			Pkg:  pkg,
			Name: entity,
		},
		TypeArgs: args,
	}
}

func (h Helper) InstanceWithDI(pkg, entity string, di map[string]compiler.Instance, args ...ts.Expr) compiler.Instance {
	return compiler.Instance{
		Ref: compiler.EntityRef{
			Pkg:  pkg,
			Name: entity,
		},
		TypeArgs:    args,
		ComponentDI: di,
	}
}

/* --- MESSAGES  --- */

func (h Helper) MsgEntity(exported bool, v compiler.MsgValue) compiler.Entity {
	return compiler.Entity{
		Exported: exported,
		Kind:     compiler.MsgEntity,
		Msg: compiler.Msg{
			Value: v,
		},
	}
}

func (h Helper) MsgWithRefEntity(exported bool, ref *compiler.EntityRef) compiler.Entity {
	return compiler.Entity{
		Exported: exported,
		Kind:     compiler.MsgEntity,
		Msg: compiler.Msg{
			Ref: ref,
		},
	}
}

func (h Helper) IntMsgValue(v int) compiler.MsgValue {
	return compiler.MsgValue{
		TypeExpr: h.Inst("int"),
		Int:      v,
	}
}

func (h Helper) IntMsg(exported bool, v int) compiler.Entity {
	return h.MsgEntity(
		exported,
		h.IntMsgValue(v),
	)
}

func (h Helper) IntVecMsgEntity(exported bool, vv []compiler.Msg) compiler.Entity {
	return h.MsgEntity(exported, compiler.MsgValue{
		TypeExpr: h.Inst("vec", h.Inst("int")),
		Vec:      vv,
	})
}

/* --- OTHER  --- */

func (h Helper) Imports(ss ...string) map[string]string {
	m := make(map[string]string, len(ss))
	for _, s := range ss {
		m[s] = s
	}
	return m
}