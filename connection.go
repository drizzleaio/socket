package socket_library

import "github.com/drizzleaio/socket/packet"

type Connection struct {
	stream Stream

	messageHandler func(msgType byte, message []byte)
}

func NewConnection(stream Stream, messageHandler func(msgType byte, message []byte)) *Connection {
	connection := &Connection{
		stream,
		messageHandler,
	}
	go connection.handleMessages()

	return connection
}

func (c *Connection) handleMessages() {
	for {
		select {
		case message := <-c.stream.Incoming:
			c.messageHandler(message.Type, message.Data)
		}
	}
}

func (c *Connection) Send(dataPacket *packet.DataPacket) {
	c.stream.Outgoing <- dataPacket
}
