package websocket

import (
	"fmt"
	"github.com/drizzleaio/socket/websocket/packet"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"testing"
)

type TestPacket struct {
	Name string
	User uint8
}

func onTestPacket(data *TestPacket) {
	fmt.Println(data.Name, data.User)
}

func TestWs(t *testing.T) {
	go func() {
		server := NewServer(func(conn *websocket.Conn) {
			fmt.Println("New connection! " + conn.RemoteAddr().String())
		})
		AddServerHandler(server, 0, onTestPacket)

		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			server.Serve(w, r)
		})
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	con := NewStream(1024)
	go con.SetConnection(c)
	//

	fmt.Println("Connected!")
	p := packet.Packet{
		Type: 0,
		Data: TestPacket{
			Name: "Test",
			User: 1,
		},
	}
	con.Outgoing <- p.Marshal()

	for {
		select {
		case msg := <-con.Incoming:
			log.Printf("received: %v", msg.Data)
		}
	}
}
