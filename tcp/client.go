package tcp

import (
	"encoding/json"
	"fmt"
	packet2 "github.com/drizzleaio/socket/tcp/packet"
	"net"
)

type Client struct {
	conn         *Connection
	packetSystem *packet2.System
}

func Connect(ip, port string) *Client {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		fmt.Println("Failed to connect to server")
		return nil
	}

	client := &Client{
		packetSystem: packet2.NewPacketSystem(),
	}
	stream := NewStream(1024)
	stream.SetConnection(conn)
	client.conn = NewConnection(*stream, client.clientPacketHandler)

	return client
}

func (c *Client) Send(packet *packet2.DataPacket) {
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
