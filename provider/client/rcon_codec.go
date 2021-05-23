// Implements a codec for https://developer.valvesoftware.com/wiki/Source_RCON_Protocol

package client

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"net/rpc"
)

var ServiceMethods = struct {
	Auth        string
	ExecCommand string
}{
	Auth:        "AUTH",
	ExecCommand: "EXECCOMMAND",
}

var packetTypeIds = struct {
	auth          int32
	authResponse  int32
	execCommand   int32
	responseValue int32
}{
	auth:          3,
	authResponse:  2,
	execCommand:   2,
	responseValue: 0,
}

// The difference betwee the reported size field, and the length of the body
const packetSizeOverhead = 10

type packetHeader struct {
	Size int32
	Id   int32
	Type int32
}

func readPacketHeader(reader io.Reader) (*packetHeader, error) {
	var p packetHeader
	if err := binary.Read(reader, binary.LittleEndian, &p.Size); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &p.Id); err != nil {
		return nil, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &p.Type); err != nil {
		return nil, err
	}
	return &p, nil
}

type packet struct {
	packetHeader
	Body string
}

func (p *packet) Write(writer io.Writer) error {
	if err := binary.Write(writer, binary.LittleEndian, p.Size); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, p.Id); err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, p.Type); err != nil {
		return err
	}
	if _, err := writer.Write([]byte(p.Body)); err != nil {
		return err
	}
	// Null terminate the body
	if err := binary.Write(writer, binary.LittleEndian, byte(0)); err != nil {
		return err
	}
	// Another null for the packet delimiter
	if err := binary.Write(writer, binary.LittleEndian, byte(0)); err != nil {
		return err
	}
	return nil
}

func newPacket(id int32, packetTypeId int32, body string) (*packet, error) {
	size := len(body) + packetSizeOverhead
	if size > math.MaxInt32 {
		return nil, errors.New("body too large for protocol")
	}
	return &packet{packetHeader{int32(size), id, packetTypeId}, body}, nil
}

type rconCodec struct {
	conn io.ReadWriteCloser

	// Some state we need to carry to adapt to the rpc interface
	prevAuthId  int32
	nextBodyLen int32
}

// newRconCodec returns a new rpc.ClientCodec using RCON on conn.
func newRconCodec(conn io.ReadWriteCloser) rpc.ClientCodec {
	return &rconCodec{
		conn: conn,
	}
}

func (c *rconCodec) WriteRequest(r *rpc.Request, param interface{}) error {
	if r.Seq > math.MaxInt32 {
		// Just error as it saves us implementing a mapping of
		// 64-bit seq to 32
		return errors.New("maximum sequence number exceeded: rcon only supports 32-bit sequence numbers: recreate the client to reset the counter")
	}
	packetId := int32(r.Seq)
	var packetTypeId int32
	switch r.ServiceMethod {
	case ServiceMethods.Auth:
		packetTypeId = packetTypeIds.auth
		c.prevAuthId = packetId
	case ServiceMethods.ExecCommand:
		packetTypeId = packetTypeIds.execCommand
	default:
		return fmt.Errorf("invalid method \"%s\"", r.ServiceMethod)
	}
	packet, err := newPacket(packetId, packetTypeId, param.(string))
	if err != nil {
		return err
	}
	return packet.Write(c.conn)
}

func (c *rconCodec) ReadResponseHeader(r *rpc.Response) error {
	p, err := readPacketHeader(c.conn)
	if err != nil {
		return err
	}
	c.nextBodyLen = p.Size - packetSizeOverhead
	switch p.Type {
	case packetTypeIds.authResponse:
		r.ServiceMethod = ServiceMethods.Auth
		r.Seq = uint64(c.prevAuthId)
		// ID == -1 is the error sentinel
		if p.Id == -1 {
			r.Error = "rcon auth failed"
		}
		return nil
	case packetTypeIds.responseValue:
		if p.Id == c.prevAuthId {
			// This appears to be the inital part of the 2 part response to an auth
			// Consume the empty body
			if err := c.ReadResponseBody(nil); err != nil {
				return err
			}
			// And process the next packet instead
			return c.ReadResponseHeader(r)
		} else {
			r.ServiceMethod = ServiceMethods.ExecCommand
			r.Seq = uint64(p.Id)
			return nil
		}
	default:
		return fmt.Errorf("unexpected packet type id %d", p.Id)
	}
}

func (c *rconCodec) ReadResponseBody(x interface{}) error {
	bodyLen := int(c.nextBodyLen)
	buffLen := bodyLen + 2 // 2 trailing null bytes
	c.nextBodyLen = 0
	buffer := make([]byte, buffLen)
	for bytesRead := 0; bytesRead < buffLen; {
		additionalBytesRead, err := c.conn.Read(buffer[bytesRead:])
		bytesRead += additionalBytesRead
		if err != nil {
			return err
		}
	}
	if x == nil {
		return nil
	}
	body := string(buffer[:bodyLen])
	*(x.(*string)) = body
	return nil
}

func (c *rconCodec) Close() error {
	return c.conn.Close()
}
