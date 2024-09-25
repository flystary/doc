package queue

import (
	"context"
	"errors"
	"tunnel/internal/pool"
)

const IPPROTO_DEFALUT = 0x0
const IPPROTO_TCP = 0x6
const IPPROTO_UDP = 0x11
const IPPROTO_ICMP = 0x1

type QueueProtocolManage struct {
	ctx     context.Context
	Queue   map[int]*QueueProtocol
	WriteCh chan pool.TUNPacket
}

type QueueProtocol struct {
	Name      string
	Rate      RateInfo
	ctx       context.Context
	Hook      QueueHook
	OutPacket chan pool.TUNPacket
	InPacket  chan pool.TUNPacket
}

func (protocol *QueueProtocol) Start() {
	go func(q *QueueProtocol) {
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
	}(protocol)
}

func NewQueueProtocol(ctx context.Context, Hook QueueHook, wrCh chan pool.TUNPacket) (q QueueProtocolManage, rdCh chan pool.TUNPacket) {

	childCtx, _ := context.WithCancel(ctx)
	q = QueueProtocolManage{
		ctx:     ctx,
		Queue:   make(map[int]*QueueProtocol, 5),
		WriteCh: wrCh,
	}
	rdCh = make(chan pool.TUNPacket, 10000)

	q.Queue[IPPROTO_DEFALUT] = &QueueProtocol{
		ctx:       childCtx,
		InPacket:  make(chan pool.TUNPacket, 10000),
		Hook:      Hook,
		OutPacket: rdCh,
		Name:      "QueueProtocol_DEFALUT",
		Rate:      NewRateInfo("QueueProtocol_DEFALUT"),
	}
	// q.Queue[IPPROTO_DEFALUT].Start()

	q.Queue[IPPROTO_TCP] = &QueueProtocol{
		ctx:       childCtx,
		InPacket:  make(chan pool.TUNPacket, 10000),
		Hook:      Hook,
		OutPacket: rdCh,
		Name:      "QueueProtocol_TCP",
		Rate:      NewRateInfo("QueueProtocol_TCP"),
	}
	// q.Queue[IPPROTO_TCP].Start()

	q.Queue[IPPROTO_UDP] = &QueueProtocol{
		ctx:       childCtx,
		InPacket:  make(chan pool.TUNPacket, 10000),
		Hook:      Hook,
		OutPacket: rdCh,
		Name:      "QueueProtocol_UDP",
		Rate:      NewRateInfo("QueueProtocol_UDP"),
	}
	// q.Queue[IPPROTO_UDP].Start()

	q.Queue[IPPROTO_ICMP] = &QueueProtocol{
		ctx:       childCtx,
		InPacket:  make(chan pool.TUNPacket, 10000),
		Hook:      Hook,
		OutPacket: rdCh,
		Name:      "QueueProtocol_ICMP",
		Rate:      NewRateInfo("QueueProtocol_ICMP"),
	}

	return
}

func (q QueueProtocolManage) setTunPakcet(Packet pool.TUNPacket) (err error) {
	var queue *QueueProtocol
	var ok bool
	queue, ok = q.Queue[Packet.Protocol]
	if ok {
		queue = q.Queue[IPPROTO_DEFALUT]
	}

	queue.InPacket <- Packet
	queue.Rate.Packet.Add(1)
	queue.Rate.Byte.Add(uint64(len(Packet.RawByte)))
	return
}

func (q QueueProtocolManage) Start() {
	for _, queue := range q.Queue {
		queue.Start()
	}
	go func(q QueueProtocolManage) {
		var packet pool.TUNPacket
		for {
			select {
			case <-q.ctx.Done():
				return
			case packet = <-q.WriteCh:
				q.setTunPakcet(packet)
			}
		}
	}(q)
}

func (q QueueProtocolManage) GetStaus() (res QueueManageStatus) {
	res.Name = "QueueProtocolManage"
	res.QueueStatus = make([]RateInfo, 0, len(q.Queue))
	for _, queue := range q.Queue {
		res.QueueStatus = append(res.QueueStatus, queue.Rate)
	}
	return
}
