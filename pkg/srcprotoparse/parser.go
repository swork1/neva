// Package srcproto parse implements parsing of bytes of source code to bytes of src protobuf.
package srcprotoparse

import (
	"context"
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/nevalang/neva/internal/compiler/src"
	"github.com/nevalang/neva/pkg/srcproto"
)

var p = parser.New(false)

func ParseFile(ctx context.Context, bb []byte) (*srcproto.File, error) {
	parsedSrcFile, err := p.ParseFile(ctx, bb)
	if err != nil {
		return nil, fmt.Errorf("parse file: %w", err)
	}

	protoSrcFile, err := fileToProto(parsedSrcFile)
	if err != nil {
		return nil, fmt.Errorf("cast file: %w", err)
	}

	return protoSrcFile, nil
}

func fileToProto(file src.File) (*srcproto.File, error) {
	entityMap := make(map[string]*srcproto.Entity)

	for k, v := range file.Entities {
		entity, err := entityToProto(v)
		if err != nil {
			return nil, fmt.Errorf("convert entity: %v", err)
		}

		entityMap[k] = entity
	}

	return &srcproto.File{
		Imports:  map[string]string{},
		Entities: entityMap,
	}, nil
}

// TODO cast entity values themselves
func entityToProto(entity src.Entity) (*srcproto.Entity, error) {
	//nolint:nosnakecase
	switch entity.Kind {
	case src.ComponentEntity:
		protoComponent, err := componentToProto(entity.Component)
		if err != nil {
			return nil, fmt.Errorf("convert to proto: %w", err)
		}

		return &srcproto.Entity{
			Kind:      srcproto.EntityKind_ENTITY_KIND_COMPONENT,
			Component: protoComponent,
		}, nil
	case src.ConstEntity:
		return &srcproto.Entity{Kind: srcproto.EntityKind_ENTITY_KIND_CONST}, nil
	case src.TypeEntity:
		return &srcproto.Entity{Kind: srcproto.EntityKind_ENTITY_KIND_TYPE_DEF}, nil
	case src.InterfaceEntity:
		return &srcproto.Entity{Kind: srcproto.EntityKind_ENTITY_KIND_INTERFACE}, nil
	default:
		return nil, errors.New("unknown entity kind")
	}
}

func componentToProto(component src.Component) (*srcproto.Component, error) {
	nodesProto := make(map[string]*srcproto.Node)
	for key, node := range component.Nodes {
		nodeProto, err := ConvertNodeToProto(node)
		if err != nil {
			return nil, fmt.Errorf("convert node: %v", err)
		}
		nodesProto[key] = nodeProto
	}

	connectionsProto := make([]*srcproto.Connection, 0, len(component.Net))
	for _, connection := range component.Net {
		connectionProto, err := ConvertConnectionToProto(connection)
		if err != nil {
			return nil, fmt.Errorf("convert connection: %v", err)
		}
		connectionsProto = append(connectionsProto, connectionProto)
	}

	interfaceProto, err := ConvertInterfaceToProto(component.Interface)
	if err != nil {
		return nil, fmt.Errorf("convert interface: %v", err)
	}

	return &srcproto.Component{
		Interface:   interfaceProto,
		Nodes:       nodesProto,
		Connections: connectionsProto,
	}, nil
}

func ConvertInterfaceToProto(interfaceValue src.Interface) (*srcproto.Interface, error) {
	typeParamsProto := make([]*srcproto.TypeParam, 0, len(interfaceValue.TypeParams))
	for _, param := range interfaceValue.TypeParams {
		typeParamsProto = append(typeParamsProto, &srcproto.TypeParam{
			Name:   param.Name,
			Constr: param.Constr.String(),
		})
	}

	ioProto, err := ConvertIOToProto(interfaceValue.IO)
	if err != nil {
		return nil, fmt.Errorf("convert IO: %v", err)
	}

	return &srcproto.Interface{
		TypeParams: typeParamsProto,
		Io:         &ioProto,
	}, nil
}

func ConvertNodeToProto(node src.Node) (*srcproto.Node, error) {
	componentDIsProto := make(map[string]*srcproto.Node, len(node.ComponentDI))
	for key, nodeDI := range node.ComponentDI {
		nodeDIProto, err := ConvertNodeToProto(nodeDI)
		if err != nil {
			return nil, fmt.Errorf("convert node DI: %v", err)
		}
		componentDIsProto[key] = nodeDIProto
	}

	entityRefProto, err := ConvertEntityRefToProto(node.EntityRef)
	if err != nil {
		return nil, fmt.Errorf("convert entity ref: %v", err)
	}

	typeArgs := make([]string, 0, len(node.TypeArgs))
	for _, arg := range node.TypeArgs {
		typeArgs = append(typeArgs, arg.String())
	}

	return &srcproto.Node{
		EntityRef:    entityRefProto,
		TypeArgs:     typeArgs,
		ComponentDis: componentDIsProto,
	}, nil
}

func ConvertEntityRefToProto(entityRef src.EntityRef) (*srcproto.EntityRef, error) {
	return &srcproto.EntityRef{
		Pkg:  entityRef.Pkg,
		Name: entityRef.Name,
	}, nil
}

func ConvertConstToProto(constValue src.Const) (srcproto.Const, error) {
	refProto, err := ConvertEntityRefToProto(*constValue.Ref)
	if err != nil {
		return srcproto.Const{}, fmt.Errorf("convert entity ref: %v", err)
	}

	valueProto, err := ConvertMsgToProto(*constValue.Value)
	if err != nil {
		return srcproto.Const{}, fmt.Errorf("convert msg: %v", err)
	}

	return srcproto.Const{
		Ref:   refProto,
		Value: valueProto,
	}, nil
}

func ConvertMsgToProto(msg src.Msg) (*srcproto.Msg, error) {
	vecsProto := make([]*srcproto.Msg, 0, len(msg.Vec))
	for _, v := range msg.Vec {
		vecProto, err := ConvertMsgToProto(*v.Value)
		if err != nil {
			return nil, fmt.Errorf("convert vec msg: %v", err)
		}
		vecsProto = append(vecsProto, vecProto)
	}

	mapProto := make(map[string]*srcproto.Msg, len(msg.Map))
	for key, v := range msg.Map {
		mapValueProto, err := ConvertMsgToProto(*v.Value) // TODO handle references
		if err != nil {
			return nil, fmt.Errorf("convert map msg: %v", err)
		}
		mapProto[key] = mapValueProto
	}

	return &srcproto.Msg{
		TypeExpr: msg.TypeExpr.String(),
		Bool:     msg.Bool,
		Int:      int64(msg.Int),
		Float:    msg.Float,
		Str:      msg.Str,
		Vecs:     vecsProto,
		Map:      mapProto,
	}, nil
}

func ConvertIOToProto(io src.IO) (srcproto.IO, error) {
	insProto := make([]*srcproto.Port, 0, len(io.In))
	for _, port := range io.In {
		portProto, err := ConvertPortToProto(port)
		if err != nil {
			return srcproto.IO{}, fmt.Errorf("convert in port: %v", err)
		}
		insProto = append(insProto, portProto)
	}

	outsProto := make([]*srcproto.Port, 0, len(io.Out))
	for _, port := range io.Out {
		portProto, err := ConvertPortToProto(port)
		if err != nil {
			return srcproto.IO{}, fmt.Errorf("convert out port: %v", err)
		}
		outsProto = append(outsProto, portProto)
	}

	return srcproto.IO{
		Ins:  insProto,
		Outs: outsProto,
	}, nil
}

func ConvertPortToProto(port src.Port) (*srcproto.Port, error) {
	return &srcproto.Port{
		TypeExpr: port.TypeExpr.String(),
		IsArray:  port.IsArray,
	}, nil
}

func ConvertConnectionToProto(connection src.Connection) (*srcproto.Connection, error) {
	senderSideProto, err := ConvertSenderConnectionSideToProto(connection.SenderSide)
	if err != nil {
		return nil, fmt.Errorf("convert sender side: %v", err)
	}

	receiverSidesProto := make([]*srcproto.ReceiverConnectionSide, 0, len(connection.ReceiverSides))
	for _, receiverSide := range connection.ReceiverSides {
		receiverSideProto, err := ConvertReceiverConnectionSideToProto(receiverSide)
		if err != nil {
			return nil, fmt.Errorf("convert receiver side: %v", err)
		}
		receiverSidesProto = append(receiverSidesProto, receiverSideProto)
	}

	return &srcproto.Connection{
		SenderSide:    senderSideProto,
		ReceiverSides: receiverSidesProto,
	}, nil
}

func ConvertSenderConnectionSideToProto(
	senderConnectionSide src.SenderConnectionSide,
) (*srcproto.SenderConnectionSide, error) {
	if senderConnectionSide.ConstRef != nil {
		return &srcproto.SenderConnectionSide{
			Selectors: senderConnectionSide.Selectors,
			ConstRef: &srcproto.EntityRef{
				Pkg:  senderConnectionSide.ConstRef.Pkg,
				Name: senderConnectionSide.ConstRef.Name,
			},
		}, nil
	}

	portAddrProto, err := ConvertPortAddrToProto(*senderConnectionSide.PortAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert port address: %v", err)
	}

	return &srcproto.SenderConnectionSide{
		PortAddr:  portAddrProto,
		Selectors: senderConnectionSide.Selectors,
	}, nil
}

func ConvertReceiverConnectionSideToProto(
	receiverConnectionSide src.ReceiverConnectionSide,
) (*srcproto.ReceiverConnectionSide, error) {
	portAddrProto, err := ConvertPortAddrToProto(receiverConnectionSide.PortAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert port address: %v", err)
	}

	return &srcproto.ReceiverConnectionSide{
		PortAddr:  portAddrProto,
		Selectors: receiverConnectionSide.Selectors,
	}, nil
}

func ConvertPortAddrToProto(portAddr src.PortAddr) (*srcproto.PortAddr, error) {
	result := &srcproto.PortAddr{
		Node: portAddr.Node,
		Port: portAddr.Port,
	}
	if portAddr.Idx == nil {
		return result, nil
	}
	result.Idx = int32(*portAddr.Idx)
	return result, nil
}
