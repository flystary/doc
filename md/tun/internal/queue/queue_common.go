package queue

import (
	"encoding/json"
	"sync/atomic"
)

type RateInfo struct {
	Name        string        `json:"name"`
	Packet      atomic.Uint64 `json:"packet"`
	Byte        atomic.Uint64 `json:"byte"`
	FinshPacket uint64        `json:"finsh_packet"`
	DropPacket  uint64        `json:"drop_packet"`
	FinshByte   atomic.Uint64 `json:"finsh_byte"`
	DropPByte   atomic.Uint64 `json:"drop_byte"`
}

type QueueManageStatus struct {
	Name        string     `json:"name"`
	QueueStatus []RateInfo `json:"queue_status"`
}

func NewRateInfo(name string) (ret RateInfo) {
	return RateInfo{
		Name:        name,
		Packet:      atomic.Uint64{},
		Byte:        atomic.Uint64{},
		FinshPacket: 0,
		DropPacket:  0,
		FinshByte:   atomic.Uint64{},
		DropPByte:   atomic.Uint64{},
	}
}

func (r RateInfo) MarshalJSON() (res []byte, err error) {

	return json.Marshal(struct {
		Name        string `json:"name"`
		Packet      uint64 `json:"packet"`
		Byte        uint64 `json:"byte"`
		FinshPacket uint64 `json:"finsh_packet"`
		FinshByte   uint64 `json:"finsh_byte"`
	}{
		Name:        r.Name,
		Packet:      r.Packet.Load(),
		Byte:        r.Byte.Load(),
		FinshPacket: r.FinshPacket,
		FinshByte:   r.FinshByte.Load(),
	})
}
