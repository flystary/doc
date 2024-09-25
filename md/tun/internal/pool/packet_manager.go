package pool

import (
	"fmt"
	"net"
	"sync"
)

type TUNPacket struct {
	DesIP      net.IP
	DesPort    int64
	Protocol   int
	EnableFEC  bool
	EnableKcp  bool
	PacketType int
	RawByte    []byte
	EncByte    []byte
	DecByte    []byte
}

var PacketPool = &sync.Pool{
	New: func() any {
		return make([]byte, 0, 1500)
	},
}

var PoolPacketManager ManagePakcetStruct

type ManagePakcetStruct struct {
}

func NewDefalutPacketManger() ManagePakcetStruct {
	return PoolPacketManager
}

func (p ManagePakcetStruct) GetPacket() (ret TUNPacket, err error) {
	var ok bool

	ret.DesIP = net.IP{}
	ret.DesPort = 0
	ret.Protocol = 0

	ret.DecByte, ok = PacketPool.Get().([]byte)
	if !ok {
		err = fmt.Errorf("get packe err")
		return
	}
	ret.RawByte, ok = PacketPool.Get().([]byte)
	if !ok {
		err = fmt.Errorf("get packe err")
		return
	}
	ret.EncByte, ok = PacketPool.Get().([]byte)
	if !ok {
		err = fmt.Errorf("get packe err")
		return
	}
	return
}

func (r ManagePakcetStruct) DestructionPacket(ret *TUNPacket) error {
	PacketPool.Put(ret.DecByte)
	PacketPool.Put(ret.EncByte)
	PacketPool.Put(ret.RawByte)
	return nil
}
