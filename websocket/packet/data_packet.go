package packet

import (
	"github.com/gorilla/websocket"
)

type DataPacket struct {
	Type   byte `json:"type"`
	Length int64
	Data   []byte `json:"data"`
}

func New(byteCode byte, data []byte) *DataPacket {
	return &DataPacket{
		Type:   byteCode,
		Length: int64(len(data)),
		Data:   data,
	}
}

func (packet *DataPacket) Write(conn *websocket.Conn) error {
	err := conn.WriteJSON(packet)
	if err != nil {
		return err
	}

	return nil
}