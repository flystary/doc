package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
	"tunnel/internal/pool"
	"tunnel/internal/queue"
)

func main() {
	QueueManage := queue.NewQueueManage()
	ctx, _ := context.WithCancel(context.Background())
	netPool := pool.NewDefalutPacketManger()

	// 创建一个协议的队列
	QueueManage.AddQueue(ctx, queue.QUEUE_TYPE_HASH, nil)
	QueueManage.AddQueue(ctx, queue.QUEUE_TYPE_PROTOCOL, nil)

	// 创建一个五元组hash队列

	QueueManage.Start()

	wg := sync.WaitGroup{}
	endNum := 0

	go func(qm queue.QueueManage) {
		var packet pool.TUNPacket
		for {
			select {
			case packet = <-qm.ReadChan:
				wg.Done()
				endNum++
				netPool.DestructionPacket(packet)
				// fmt.Printf("read index:%d %+v\n", binary.LittleEndian.Uint64(packet.RawByte[:]), packet)
			}
		}
	}(QueueManage)

	// sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// <-sigs

	// var err error
	ip := net.ParseIP("192.168.1.2")

	ticker := time.NewTicker(time.Second * 1)
	go func() {
		for _ = range ticker.C {
			status, _ := json.Marshal(QueueManage.GetStaus())
			fmt.Printf("queue status : %s \n", status)
		}
	}()

	for index := 0; index < 10000000; index++ {
		wg.Add(1)
		packet, _ := netPool.GetPacket()
		packet.DesIP = ip
		packet.DesPort = int64(index % 100)
		packet.Protocol = 0x6
		ubyte := []byte{0, 0, 0, 0, 0, 0, 0, 0}
		binary.LittleEndian.PutUint64(ubyte[:], uint64(index))
		packet.RawByte = ubyte
		QueueManage.WirteChan <- packet

	}

	wg.Wait()
	status, _ := json.Marshal(QueueManage.GetStaus())
	fmt.Printf("queue status : %s \n", status)
	// println("ssss")
	// sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// <-sigs

}
