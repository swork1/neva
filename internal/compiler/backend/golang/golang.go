package golang

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"strings"
	"text/template"
	"unicode"

	"github.com/nevalang/neva/internal/compiler"
)

//go:embed tmpl/main.go.tmpl runtime
var Efs embed.FS // TODO make private

type Backend struct{}

var ErrExecTmpl = errors.New("execute template")

func (b Backend) GenerateTarget(ctx context.Context, prog compiler.LLProgram) ([]byte, error) {
	tmpl, err := template.New("main.go.tmpl").Funcs(template.FuncMap{
		"getMsg":           b.getMsg,
		"getPorts":         b.getPortsFunc(prog.Ports),
		"getPortChVarName": b.getPortChVarName,
		"getConnComment":   b.getConnComment,
	}).ParseFS(Efs, "tmpl/main.go.tmpl")
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, prog); err != nil {
		return nil, errors.Join(ErrExecTmpl, err)
	}

	return buf.Bytes(), nil
}

var ErrUnknownMsgType = errors.New("unknown msg type")

func (b Backend) getMsg(msg compiler.LLMsg) (string, error) {
	switch msg.Type {
	case compiler.LLIntMsg:
		return fmt.Sprintf("runtime.NewIntMsg(%d)", msg.Int), nil
	}
	return "", fmt.Errorf("%w: %v", ErrUnknownMsgType, msg.Type)
}

func (b Backend) getConnComment(conn compiler.LLConnection) string {
	s := b.fmtPortAddr(conn.SenderSide.PortAddr) + " -> "

	for _, rcvr := range conn.ReceiverSides {
		s += b.fmtPortAddr(rcvr.PortAddr)
	}

	return "// " + s
}

func (b Backend) fmtPortAddr(addr compiler.LLPortAddr) string {
	return fmt.Sprintf("%s.%s[%d]", addr.Path, addr.Name, addr.Idx)
}

func (b Backend) getPortChVarName(addr compiler.LLPortAddr) string {
	path := b.handleSpecialChars(addr.Path)
	port := addr.Name
	if path != "" {
		port = b.uppercaseFirstLetter(addr.Name)
	}
	return fmt.Sprintf("%s%s%dPort", path, port, addr.Idx)
}

func (b Backend) getPortsFunc(ports map[compiler.LLPortAddr]uint8) func(path, port string) string {
	return func(path, port string) string {
		var s string
		for addr := range ports {
			if addr.Path == path && addr.Name == port {
				s = s + b.getPortChVarName(addr) + ","
			}
		}
		return s
	}
}

func (b Backend) handleSpecialChars(portPath string) string {
	var (
		buffer          bytes.Buffer
		shouldUppercase bool
	)

	for i := 0; i < len(portPath); i++ {
		if portPath[i] == '.' || portPath[i] == '/' {
			shouldUppercase = true
			continue
		}
		s := string(portPath[i])
		if shouldUppercase {
			s = strings.ToUpper(s)
			shouldUppercase = false
		}
		buffer.WriteString(s)
	}

	return buffer.String()
}

func (b Backend) uppercaseFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	bb := []byte(s)
	bb[0] = byte(unicode.ToUpper(rune(bb[0])))
	return string(bb)
}
