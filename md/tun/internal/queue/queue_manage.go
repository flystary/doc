package queue

import (
	"context"
	"fmt"
	"tunnel/internal/pool"
)

const (
	QUEUE_TYPE_PROTOCOL = iota
	QUEUE_TYPE_HASH
	QUEUE_TYPE_IP_PACKET
)

type QueueHook func(*pool.TUNPacket) error

var ErrTakOver = fmt.Errorf("packet  take over")

type QueueManageInface interface {
	Start()
	GetStaus() (res QueueManageStatus)
	// SetChan(WriteChan, ReatChan chan TUNPacket)
}

type QueueManage struct {
	Table     []QueueManageInface
	ReadChan  chan pool.TUNPacket
	WirteChan chan pool.TUNPacket
	lastRdCh  *chan pool.TUNPacket
	lastWrCh  *chan pool.TUNPacket
}

func NewQueueManage() (r QueueManage) {
	// QueueHook = func(t pool.TUNPacket) error {}
	r.Table = make([]QueueManageInface, 0, 255)
	r.WirteChan = make(chan pool.TUNPacket, 10000)
	r.lastRdCh = &r.ReadChan
	r.lastWrCh = &r.WirteChan
	return
}

func (q *QueueManage) AddQueue(ctx context.Context, QTyep int, HOOK QueueHook) (err error) {
	var queue QueueManageInface
	switch QTyep {
	case QUEUE_TYPE_PROTOCOL:
		queue, q.ReadChan = NewQueueProtocol(ctx, HOOK, *q.lastWrCh)
		q.lastWrCh = &q.ReadChan
	case QUEUE_TYPE_HASH:
		queue, q.ReadChan = NewQueueHash(ctx, 1, HOOK, *q.lastWrCh)
		q.lastWrCh = &q.ReadChan
	case QUEUE_TYPE_IP_PACKET:
		queue, q.ReadChan = NewQueueIPPacket(ctx, 1, HOOK, *q.lastWrCh)
		q.lastWrCh = &q.ReadChan
	default:
		return fmt.Errorf("type is error")
	}

	q.Table = append(q.Table, queue)
	return
}

func (q *QueueManage) Start() {
	for _, queue := range q.Table {
		queue.Start()
	}
}

func (q *QueueManage) GetStaus() (res []QueueManageStatus) {
	res = make([]QueueManageStatus, 0, len(q.Table))
	for _, queue := range q.Table {
		add := queue.GetStaus()
		res = append(res, add)
	}
	return
}
