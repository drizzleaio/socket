package socket_library

import (
	"encoding/json"
	"fmt"
	"github.com/drizzleaio/socket/packet"
	"net"
)

type Client struct {
	conn         *Connection
	packetSystem *packet.System
}

func Connect(ip, port string) *Client {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		return nil
	}

	client := &Client{
		packetSystem: packet.NewPacketSystem(),
	}
	stream := NewStream(1024)
	stream.SetConnection(conn)
	client.conn = NewConnection(*stream, client.clientPacketHandler)

	return client
}

func (c *Client) Send(packet *packet.DataPacket) {
	c.conn.Send(packet)
}

func (c *Client) clientPacketHandler(msgType byte, data []byte) {
	handler, ok := c.packetSystem.Handlers[msgType]
	if !ok {
		fmt.Println("No handler for packet type", msgType)
		return
	}

	handler(data)
}

func AddClientHandler[T any](c *Client, id byte, handler func(packet *T)) {
	c.packetSystem.Handlers[id] = func(data []byte) { // Set the packet decoder for this ID
		out := new(T)                 // Create a new instance of the output type
		_ = json.Unmarshal(data, out) // Unmarshal the data into the output type
		handler(out)
	}
}
