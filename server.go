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

func New(port string) *Server {
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

func AddHandler[T any](s *packet.System, id byte, handler func(packet *T)) {
	s.Handlers[id] = func(data []byte) { // Set the packet decoder for this ID
		out := new(T)                 // Create a new instance of the output type
		_ = json.Unmarshal(data, out) // Unmarshal the data into the output type
		handler(out)
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
