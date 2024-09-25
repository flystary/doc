package tundrive

import (
	"context"
	"os"
	"sync"
	"sync/atomic"
	"tunnel/internal/config"
	"tunnel/internal/pool"
	"tunnel/internal/queue"
)

type TUNPacketMange interface {
	GetPacket() (pool.TUNPacket, error)
	DestructionPacket(pool.TUNPacket) error
}

type TUNDrive struct {
	ctx                context.Context
	lock               sync.Locker
	TunInfo            TUNInfo
	TunFIle            *os.File
	ReadPacketChannle  chan *pool.TUNPacket
	WritePacketChannle chan *pool.TUNPacket
	RXPacket           atomic.Uint64
	TXPacket           atomic.Uint64
	RXByte             atomic.Uint64
	TXByte             atomic.Uint64
	TXQueueManage      queue.QueueManage
}

func NewDefalutTunDrive(context context.Context, c config.Config) TUNDrive {
	return TUNDrive{
		ctx:                context,
		lock:               &sync.RWMutex{},
		ReadPacketChannle:  make(chan *pool.TUNPacket, 100000),
		WritePacketChannle: make(chan *pool.TUNPacket, 100000),
		TunInfo:            NewTunInfo(c),
	}
}

func (tun *TUNDrive) Start() {
	go func(t *TUNDrive) {
		tun.TunFIle = os.NewFile(uintptr(tun.TunInfo.FD), "/dev/net/tun")
		var readPacket pool.TUNPacket

		var err error
		var readNum int
		for {
			select {
			case <-t.ctx.Done():
				return
			default:
				readPacket, err = pool.PoolPacketManager.GetPacket()
				t.RXPacket.Add(1)
				if err != nil {
					return
				}
				readNum, err = t.TunFIle.Read(readPacket.RawByte)
				t.RXByte.Add(uint64(readNum))
				if err != nil {
					return
				}
				tun.ReadPacketChannle <- &readPacket
			}
		}
	}(tun)

	go func(t *TUNDrive) {
		var wirtePack *pool.TUNPacket
		var writeNum int
		for {
			select {
			case <-t.ctx.Done():
				return
			case wirtePack = <-t.WritePacketChannle:
				t.TXPacket.Add(1)
				writeNum, _ = t.TunFIle.Write(wirtePack.DecByte[8:])
				t.TXByte.Add(uint64(writeNum))
				pool.PoolPacketManager.DestructionPacket(wirtePack)
			}
		}
	}(tun)
}
