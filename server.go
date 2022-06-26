package socket_library

import (
	"encoding/json"
	"fmt"
	"github.com/drizzleaio/socket/packet"
	"net"
)

type Server struct {
	connections []*Connection

	packetSystem *packet.System
}

func NewServer(port string) *Server {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil
	}

	server := &Server{
		packetSystem: packet.NewPacketSystem(),
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				panic(err)
			}
			stream := NewStream(1024)
			stream.SetConnection(conn)

			c := NewConnection(*stream, server.packetHandler)

			server.connections = append(server.connections, c)
		}
	}()

	return server
}

func AddServerHandler[T any](s *Server, id byte, handler func(packet *T)) {
	s.packetSystem.Handlers[id] = func(data []byte) {
		out := new(T)
		_ = json.Unmarshal(data, out)
		handler(out)
	}
}

func (s *Server) Emit(p *packet.DataPacket) {
	for _, c := range s.connections {
		c.Send(p)
	}
}

func (s *Server) packetHandler(msgType byte, data []byte) {
	handler, ok := s.packetSystem.Handlers[msgType]
	if !ok {
		fmt.Println("No handler for packet type", msgType)
		return
	}

	handler(data)
}
