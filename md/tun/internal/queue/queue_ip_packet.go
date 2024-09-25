package queue

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"tunnel/internal/pool"
)

type QueueIPPacketManage struct {
	ctx       context.Context
	num       int
	Queue     []*QueueIPPacket
	WriteChan chan pool.TUNPacket
}

type QueueIPPacket struct {
	Name      string
	ctx       context.Context
	Rate      RateInfo
	Hook      QueueHook
	OutPacket chan pool.TUNPacket
	InPacket  chan pool.TUNPacket
}

func (ippacket *QueueIPPacket) Start() {
	go func(q *QueueIPPacket) {
		var packet pool.TUNPacket
		var err error
		for {
			select {
			case <-q.ctx.Done():
				return
			case packet = <-q.InPacket:
				if q.Hook == nil {
					q.Rate.FinshPacket++
					q.Rate.FinshByte.Add(uint64(len(packet.RawByte)))
					q.OutPacket <- packet
					return
				}
				err = q.Hook(&packet)
				if err == nil {
					q.Rate.FinshPacket++
					q.Rate.FinshByte.Add(uint64(len(packet.RawByte)))
					q.OutPacket <- packet
					return
				}
				if !errors.Is(err, ErrTakOver) {
					q.Rate.DropPacket++
					q.Rate.DropPByte.Add(uint64(len(packet.RawByte)))

				} else {
					q.Rate.FinshPacket++
					q.Rate.FinshByte.Add(uint64(len(packet.RawByte)))
				}
				return
			}
		}
	}(ippacket)
}

func NewQueueIPPacket(ctx context.Context, num int, Hook QueueHook, wrCh chan pool.TUNPacket) (q QueueIPPacketManage, rdCh chan pool.TUNPacket) {
	q = QueueIPPacketManage{
		ctx:       ctx,
		Queue:     make([]*QueueIPPacket, num),
		num:       num,
		WriteChan: wrCh,
	}
	rdCh = make(chan pool.TUNPacket, 10000)

	childCtx, _ := context.WithCancel(ctx)

	for index := 0; index < q.num; index++ {
		q.Queue[index] = &QueueIPPacket{
			ctx:       childCtx,
			InPacket:  make(chan pool.TUNPacket, 10000),
			OutPacket: rdCh,
			Hook:      Hook,
			Rate:      NewRateInfo(fmt.Sprintf("IPPakcet_%d", index)),
			Name:      fmt.Sprintf("queue_ippakcet_%d", index),
		}
	}
	return
}

func (q QueueIPPacketManage) setTunPakcet(Packet pool.TUNPacket) (err error) {
	var value uint64
	value = 0
	value += uint64(Packet.DesPort)
	value += uint64(Packet.Protocol)
	value += binary.LittleEndian.Uint64(Packet.DesIP)
	index := value % uint64(q.num)

	q.Queue[index].Rate.Packet.Add(1)
	q.Queue[index].Rate.Byte.Add(uint64(len(Packet.RawByte)))
	q.Queue[index].InPacket <- Packet
	return
}

func (q QueueIPPacketManage) Start() {
	for _, queue := range q.Queue {
		queue.Start()
	}

	go func(q QueueIPPacketManage) {
		var packet pool.TUNPacket
		for {
			select {
			case <-q.ctx.Done():
				return
			case packet = <-q.WriteChan:
				q.setTunPakcet(packet)
			}
		}
	}(q)
}

func (q QueueIPPacketManage) GetStaus() (res QueueManageStatus) {
	res.Name = "QueueIPPacketManage"
	res.QueueStatus = make([]RateInfo, 0, len(q.Queue))
	for _, queue := range q.Queue {
		res.QueueStatus = append(res.QueueStatus, queue.Rate)
	}
	return
}
