package websocket

import (
	"log"
	"testing"
)

type Product struct {
	Name     string    `json:"name"`
	Site     string    `json:"site"`
	Sku      string    `json:"sku"`
	Color    string    `json:"color"`
	Variants []Variant `json:"variants"`
}

type Variant struct {
	ID   string `json:"id"`
	Size string `json:"size"`
}

func onTestPacket(data *Product) {
	log.Println("Received product:", data.Sku)
}

func TestWs(t *testing.T) {
	c := Connect("wss", "localhost", "/ws")
	AddClientHandler(c, 1, onTestPacket)

	for {

	}

}
