package socket_library

import (
	"fmt"
	"github.com/drizzleaio/socket/packet"
	"net"
	"testing"
)

type TestPacket struct {
	Name string
	User uint8
}

func onTestPacket(data *TestPacket) {
	fmt.Println(data.Name, data.User)
}

func TestSocket(t *testing.T) {
	server := New("5000")
	if server == nil {
		t.Error("Server is nil")
	}
	AddHandler(server.packetSystem, 0, onTestPacket)

	// Connect to a server
	conn, _ := net.Dial("tcp", "localhost:5000")

	// Create a stream
	stream := NewStream(1024)
	stream.SetConnection(conn)

	test := packet.Packet{
		Type: 0,
		Data: TestPacket{
			Name: "test",
			User: 0,
		},
	}

	stream.Outgoing <- test.Marshal()

	for {

	}
}
