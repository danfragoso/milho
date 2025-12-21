package interpreter

import (
	"bufio"
	"fmt"
	"net"
	"sync"

	"github.com/danfragoso/milho/mir"
)

// TCP socket state management
var (
	tcpSockets     = make(map[int64]*tcpSocketInstance)
	tcpSocketMutex sync.RWMutex
	nextSocketID   int64 = 1
)

type tcpSocketInstance struct {
	ID     int64
	Conn   net.Conn
	Reader *bufio.Reader
	Writer *bufio.Writer
}

// __tcpConnect creates a TCP connection to host:port
func __tcpConnect(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("tcp-connect: expected 2 parameters (host, port), got %d", len(params))
	}

	hostExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if hostExpr.Type() != mir.StringExpr {
		return nil, fmt.Errorf("tcp-connect: host must be a string, got %s", hostExpr.Type())
	}

	host := hostExpr.(*mir.StringExpression).Val

	portExpr, err := evaluate(params[1], session)
	if err != nil {
		return nil, err
	}

	if portExpr.Type() != mir.NumberExpr {
		return nil, fmt.Errorf("tcp-connect: port must be a number, got %s", portExpr.Type())
	}

	port := portExpr.(*mir.NumberExpression).Numerator

	// Connect to TCP server
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("tcp-connect: failed to connect to %s: %s", address, err)
	}

	tcpSocketMutex.Lock()
	socketID := nextSocketID
	nextSocketID++

	instance := &tcpSocketInstance{
		ID:     socketID,
		Conn:   conn,
		Reader: bufio.NewReader(conn),
		Writer: bufio.NewWriter(conn),
	}

	tcpSockets[socketID] = instance
	tcpSocketMutex.Unlock()

	return mir.CreateNumberExpression(socketID, 1)
}

// __tcpWrite writes data to a TCP socket
func __tcpWrite(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("tcp-write: expected 2 parameters (socketID, data), got %d", len(params))
	}

	socketIDExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if socketIDExpr.Type() != mir.NumberExpr {
		return nil, fmt.Errorf("tcp-write: socketID must be a number, got %s", socketIDExpr.Type())
	}

	socketID := socketIDExpr.(*mir.NumberExpression).Numerator

	dataExpr, err := evaluate(params[1], session)
	if err != nil {
		return nil, err
	}

	if dataExpr.Type() != mir.StringExpr {
		return nil, fmt.Errorf("tcp-write: data must be a string, got %s", dataExpr.Type())
	}

	data := dataExpr.(*mir.StringExpression).Val

	tcpSocketMutex.RLock()
	instance, exists := tcpSockets[socketID]
	tcpSocketMutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("tcp-write: socket %d not found", socketID)
	}

	// Write data to socket
	_, err = instance.Writer.WriteString(data)
	if err != nil {
		return nil, fmt.Errorf("tcp-write: failed to write: %s", err)
	}

	// Flush the buffer
	err = instance.Writer.Flush()
	if err != nil {
		return nil, fmt.Errorf("tcp-write: failed to flush: %s", err)
	}

	return mir.CreateNilExpression()
}

// __tcpRead reads data from a TCP socket until delimiter
func __tcpRead(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) < 1 || len(params) > 2 {
		return nil, fmt.Errorf("tcp-read: expected 1-2 parameters (socketID, [delimiter]), got %d", len(params))
	}

	socketIDExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if socketIDExpr.Type() != mir.NumberExpr {
		return nil, fmt.Errorf("tcp-read: socketID must be a number, got %s", socketIDExpr.Type())
	}

	socketID := socketIDExpr.(*mir.NumberExpression).Numerator

	// Default delimiter is newline
	delimiter := byte('\n')
	if len(params) == 2 {
		delimExpr, err := evaluate(params[1], session)
		if err != nil {
			return nil, err
		}

		if delimExpr.Type() != mir.StringExpr {
			return nil, fmt.Errorf("tcp-read: delimiter must be a string, got %s", delimExpr.Type())
		}

		delimStr := delimExpr.(*mir.StringExpression).Val
		if len(delimStr) > 0 {
			delimiter = delimStr[0]
		}
	}

	tcpSocketMutex.RLock()
	instance, exists := tcpSockets[socketID]
	tcpSocketMutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("tcp-read: socket %d not found", socketID)
	}

	// Read until delimiter
	data, err := instance.Reader.ReadString(delimiter)
	if err != nil {
		return nil, fmt.Errorf("tcp-read: failed to read: %s", err)
	}

	return mir.CreateStringExpression(data)
}

// __tcpReadLine reads a line from a TCP socket (until \n)
func __tcpReadLine(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("tcp-read-line: expected 1 parameter (socketID), got %d", len(params))
	}

	socketIDExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if socketIDExpr.Type() != mir.NumberExpr {
		return nil, fmt.Errorf("tcp-read-line: socketID must be a number, got %s", socketIDExpr.Type())
	}

	socketID := socketIDExpr.(*mir.NumberExpression).Numerator

	tcpSocketMutex.RLock()
	instance, exists := tcpSockets[socketID]
	tcpSocketMutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("tcp-read-line: socket %d not found", socketID)
	}

	// Read until newline
	data, err := instance.Reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("tcp-read-line: failed to read: %s", err)
	}

	return mir.CreateStringExpression(data)
}

// __tcpClose closes a TCP socket
func __tcpClose(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("tcp-close: expected 1 parameter (socketID), got %d", len(params))
	}

	socketIDExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if socketIDExpr.Type() != mir.NumberExpr {
		return nil, fmt.Errorf("tcp-close: socketID must be a number, got %s", socketIDExpr.Type())
	}

	socketID := socketIDExpr.(*mir.NumberExpression).Numerator

	tcpSocketMutex.Lock()
	instance, exists := tcpSockets[socketID]
	if !exists {
		tcpSocketMutex.Unlock()
		return nil, fmt.Errorf("tcp-close: socket %d not found", socketID)
	}

	// Close connection
	closeErr := instance.Conn.Close()
	if closeErr != nil {
		tcpSocketMutex.Unlock()
		return nil, fmt.Errorf("tcp-close: failed to close: %s", closeErr)
	}

	delete(tcpSockets, socketID)
	tcpSocketMutex.Unlock()

	return mir.CreateNilExpression()
}
