// Implements a client for https://developer.valvesoftware.com/wiki/Source_RCON_Protocol
// Differs to https://github.com/gtaylor/factorio-rcon
// mainly in that it is cocurrency-safe, handles parallel/interleaved calls
// and is based on stdlib net/rpc client

package client

import (
	"io"
	"net"
	"net/rpc"
)

type RCON struct {
	c *rpc.Client
}

func NewClient(conn io.ReadWriteCloser) *rpc.Client {
	return rpc.NewClientWithCodec(newRconCodec(conn))
}

func Dial(address string) (*RCON, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	rcon := &RCON{c: NewClient(conn)}
	return rcon, nil
}

func (r *RCON) Close() error {
	return r.c.Close()
}

func (r *RCON) Execute(command string) (string, error) {
	var response string
	err := r.c.Call(ServiceMethods.ExecCommand, command, &response)
	return response, err
}

func (r *RCON) Authenticate(password string) (err error) {
	return r.c.Call(ServiceMethods.Auth, password, nil)
}
