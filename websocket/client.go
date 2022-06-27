package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/drizzleaio/socket/websocket/packet"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

type Client struct {
	conn         *Connection
	packetSystem *packet.System
}

func Connect(scheme, host, path string) *Client {
	u := url.URL{Scheme: scheme, Host: host, Path: path}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	client := &Client{
		packetSystem: packet.NewPacketSystem(),
	}
	stream := NewStream(1024)
	stream.SetConnection(c)
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
