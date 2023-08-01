package shared

import (
	"fmt"

	ts "github.com/nevalang/neva/pkg/types"
)

type HLFile struct {
	Imports  map[string]string
	Entities map[string]Entity
}

type Entity struct {
	Exported  bool
	Kind      EntityKind
	Msg       HLMsg
	Type      ts.Def // FIXME https://github.com/nevalang/neva/issues/186
	Interface Interface
	Component Component
}

type EntityKind uint8

const (
	ComponentEntity EntityKind = iota + 1
	MsgEntity
	TypeEntity
	InterfaceEntity
)

func (e EntityKind) String() string {
	switch e {
	case ComponentEntity:
		return "component"
	case MsgEntity:
		return "msg"
	case TypeEntity:
		return "type"
	case InterfaceEntity:
		return "interface"
	default:
		return "unknown"
	}
}

type Component struct {
	Interface Interface
	Nodes     map[string]Node
	Net       []Connection
}

type Interface struct {
	Params []ts.Param
	IO     IO
}

type Node struct {
	Ref         EntityRef
	TypeArgs    []ts.Expr
	ComponentDI map[string]Node
}

type EntityRef struct {
	Pkg  string // "" for local entities (alias, namespace)
	Name string
}

func (e EntityRef) String() string {
	if e.Pkg == "" {
		return e.Name
	}
	return fmt.Sprintf("%s.%s", e.Pkg, e.Name)
}

type HLMsg struct {
	Ref   *EntityRef // if nil then use value
	Value MsgValue
}

type MsgValue struct {
	TypeExpr ts.Expr          // type of the message
	Bool     bool             // only for messages with `bool`  type
	Int      int              // only for messages with `int` type
	Float    float64          // only for messages with `float` type
	Str      string           // only for messages with `str` type
	Vec      []HLMsg          // only for types with `vec` type
	Map      map[string]HLMsg // only for types with `map` type
}

type IO struct {
	In, Out map[string]Port
}

type Port struct {
	Type  ts.Expr
	IsArr bool
}

type Connection struct {
	SenderSide    SenderConnectionSide
	ReceiverSides []PortConnectionSide
}

// SenderConnectionSide can have outport or message as a source of data
type SenderConnectionSide struct {
	MsgRef *EntityRef // if not nil then port addr must not be used
	PortConnectionSide
}

type PortConnectionSide struct {
	PortAddr  ConnPortAddr
	Selectors []Selector
}

type Selector struct {
	RecField string // "" means use ArrIdx
	ArrIdx   int
}

type ConnPortAddr struct {
	Node string
	RelPortAddr
}

type RelPortAddr struct {
	Name string
	Idx  uint8
}
