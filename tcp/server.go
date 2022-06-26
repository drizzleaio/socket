package tcp

import (
	"encoding/json"
	"fmt"
	packet2 "github.com/drizzleaio/socket/tcp/packet"
	"net"
)

type Server struct {
	connections []*Connection

	listener     net.Listener
	packetSystem *packet2.System
}

func NewServer(port string, connectionHandler func(conn net.Conn)) *Server {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil
	}

	server := &Server{
		listener:     listener,
		packetSystem: packet2.NewPacketSystem(),
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

			connectionHandler(conn)
			server.connections = append(server.connections, c)
		}
	}()

	return server
}

func (s *Server) Destroy() {
	err := s.listener.Close()
	if err != nil {
		return
	}
}

func AddServerHandler[T any](s *Server, id byte, handler func(packet *T)) {
	s.packetSystem.Handlers[id] = func(data []byte) {
		out := new(T)
		_ = json.Unmarshal(data, out)
		handler(out)
	}
}

func (s *Server) Emit(p *packet2.DataPacket) {
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
