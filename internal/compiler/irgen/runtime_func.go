package irgen

import (
	"errors"
	"strings"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
	"github.com/nevalang/neva/internal/runtime/ir"
)

func getFuncRef(component src.Component, nodeTypeArgs []ts.Expr) (string, error) {
	args, ok := component.Directives[compiler.ExternDirective]
	if !ok {
		return "", nil
	}

	if len(args) == 1 {
		return args[0], nil
	}

	if len(nodeTypeArgs) == 0 || nodeTypeArgs[0].Inst == nil {
		return "", nil
	}

	firstTypeArg := nodeTypeArgs[0].Inst.Ref.String()
	for _, arg := range args {
		parts := strings.Split(arg, " ")
		if firstTypeArg == parts[0] {
			return parts[1], nil
		}
	}

	return "", errors.New("type argument mismatches runtime func directive")
}

func getCfgMsg(node src.Node, scope src.Scope) (*ir.Msg, *compiler.Error) {
	args, ok := node.Directives[compiler.BindDirective]
	if !ok {
		return nil, nil
	}

	entity, location, err := scope.Entity(compiler.ParseEntityRef(args[0]))
	if err != nil {
		return nil, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
		}
	}

	return getIRMsgBySrcRef(entity.Const, scope.WithLocation(location))
}
