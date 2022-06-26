package socket_library

import (
	"fmt"
	"testing"
)

type TestPacket struct {
	Name string
	User uint8
}

type Product struct {
	Name     string    `json:"name"`
	Sku      string    `json:"sku"`
	Color    string    `json:"color"`
	Variants []Variant `json:"variants"`
}

type Variant struct {
	ID   string `json:"id"`
	Size string `json:"size"`
}

func onTestPacket(data *TestPacket) {
	fmt.Println(data.Name, data.User)
}

func onPacket(data *Product) {
	fmt.Println(data.Name)
}

func TestSocket(t *testing.T) {
	client := Connect("localhost", "5000")
	if client == nil {
		t.Error("Client is nil")
	}
	AddClientHandler(client, 1, onPacket)

	for {

	}
}
