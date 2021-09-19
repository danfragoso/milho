package interpreter

import (
	"fmt"
	"strconv"
)

func __createSocket(params []Expression, session *Session) (Expression, error) {
	if len(params) != 3 {
		return nil, fmt.Errorf("createSocket: expected 3 parameters, got %d, parameters must be ('SERVER|'CLIENT, 'TCP|'UDP, PORT)", len(params))
	}

	sKind, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if sKind.Type() != SymbolExpr {
		return nil, fmt.Errorf("createSocket: expected socketType to be a symbol, got %s", sKind.Type())
	}

	if sKind.Value() != "SERVER" && sKind.Value() != "CLIENT" {
		return nil, fmt.Errorf("createSocket: expected socketType to be either SERVER or CLIENT, got %s", sKind.Value())
	}

	sProtocol, err := evaluate(params[1], session)
	if err != nil {
		return nil, err
	}

	if sProtocol.Type() != SymbolExpr {
		return nil, fmt.Errorf("createSocket: expected socketProtocol to be a symbol, got %s", sKind.Type())
	}

	if sProtocol.Value() != "TCP" && sProtocol.Value() != "UDP" {
		return nil, fmt.Errorf("createSocket: expected socketProtocol to be either TCP or UDP, got %s", sKind.Value())
	}

	sPort, err := evaluate(params[2], session)
	if err != nil {
		return nil, err
	}

	if sPort.Type() != NumberExpr {
		return nil, fmt.Errorf("createSocket: expected socketPort to be a numer, got %s", sKind.Type())
	}

	port, err := strconv.Atoi(sPort.Value())
	if err != nil {
		return nil, fmt.Errorf("createSocket: error parsing port, %s", err)
	}

	return createSocketExpression(socketType(sKind.Value()), socketProtocol(sProtocol.Value()), port)
}

func __writeSocket(params []Expression, session *Session) (Expression, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("writeSocket: expected 2 parameter, got %d", len(params))
	}

	sock, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if sock.Type() != SocketExpr {
		return nil, fmt.Errorf("writeSocket: expected socket to be a socket, got %s", sock.Type())
	}

	msg, err := evaluate(params[1], session)
	if err != nil {
		return nil, err
	}

	if msg.Type() != StringExpr {
		return nil, fmt.Errorf("writeSocket: expected msg to be a string, got %s", msg.Type())
	}

	return createNilExpression()
}
