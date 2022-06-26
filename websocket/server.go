package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/drizzleaio/socket/websocket/packet"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{}

type Server struct {
	connections []*Connection

	packetSystem *packet.System

	connectionHandler func(conn *websocket.Conn)
}

func NewServer(connectionHandler func(conn *websocket.Conn)) *Server {
	return &Server{
		connectionHandler: connectionHandler,
		packetSystem:      packet.NewPacketSystem(),
	}
}

func (s *Server) Serve(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	stream := NewStream(1024)
	stream.SetConnection(connection)
	defer stream.Close()

	c := NewConnection(*stream, s.packetHandler)

	s.connectionHandler(connection)

	s.connections = append(s.connections, c)

	select {
	case <-stream.closeWriter:
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

func (s *Server) packetHandler(msgType byte, data []byte) {
	handler, ok := s.packetSystem.Handlers[msgType]
	if !ok {
		fmt.Println("No handler for packet type", msgType)
		return
	}

	handler(data)
}
