package packet

import "encoding/json"

type Packet struct {
	Type byte
	Data interface{}
}

func (p *Packet) Marshal() *DataPacket {
	jsonBytes, err := json.Marshal(p.Data)
	if err != nil {
		return nil
	}

	return New(p.Type, jsonBytes)
}
