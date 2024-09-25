package negotication

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"sync/atomic"
	"time"
	"tunnel/internal/config"
)

const (
	WAN_FENG_TYPE_HEARTBEAT = iota
	WAN_FENG_TYPE_CLIENT_ID
	WAN_FENG_TYPE_RSA
	WAN_FENG_TYPE_KEY
	WAN_FENG_TYPE_REDUNDANCY
	WAN_FENG_TYPE_IP
	WAN_FENG_TYPE_PAYLOAD
	WAN_FENG_TYPE_HEARTBEAT_ECHO
	WAN_FENG_TYPE_HEARTBEAT_ECHO_REPLY
)

type Negotication struct {
	IsRunning            atomic.Bool
	TransportLocalIp     string
	NegotiationServePort int
	NegotiationServeIP   string
	laddr                *net.TCPAddr
	raddr                *net.TCPAddr
	Conn                 *net.TCPConn
	ctx                  context.Context
	CancelFunc           context.CancelFunc
}

type NegoticationData struct {
	Type   uint8
	Length uint32
	Body   []byte
}

func NewNegotication(config config.Config) (r Negotication, err error) {
	r.NegotiationServePort = config.GetInt("negotiation_server_port", 0)
	if r.NegotiationServePort == 0 {
		err = fmt.Errorf("negotiation_server_port is not set")
		return
	}
	r.NegotiationServeIP = config.GetString("negotiation_server_ip", "")
	if len(r.NegotiationServeIP) == 0 {
		err = fmt.Errorf("negotiation_server_ip is not set")
		return
	}

	r.TransportLocalIp = config.GetString("transport_local_ip", "0.0.0.0")
	if len(r.NegotiationServeIP) == 0 {
		err = fmt.Errorf("negotiation_server_ip is not set")
		return
	}
	r.IsRunning.Store(false)
	r.laddr = &net.TCPAddr{
		IP:   net.ParseIP(r.TransportLocalIp),
		Port: 0,
	}
	r.raddr = &net.TCPAddr{
		IP:   net.ParseIP(r.NegotiationServeIP),
		Port: r.NegotiationServePort,
	}

	r.Conn, err = net.DialTCP("tcp", r.laddr, r.raddr)
	return
}

func (n *Negotication) SetCtx(ctx context.Context) error {
	n.ctx, n.CancelFunc = context.WithCancel(ctx)
	return nil
}

func (n *Negotication) SendHearteat() {
	ticker := time.NewTicker(time.Second * 1)
	sendData := NegoticationData{}
	sendData.Type = WAN_FENG_TYPE_HEARTBEAT
	sendData.Length = 0

	ret := sendData.GetSendByte()
	for {
		select {
		case <-n.ctx.Done():
			return
		case <-ticker.C:
			n.Conn.Write(ret)
		}
	}
}

func (n *Negotication) SendData(d NegoticationData) (int, error) {

	ret := d.GetSendByte()
	return n.Conn.Write(ret)
}

func (d NegoticationData) GetSendByte() (ret []byte) {
	bodyLen := make([]byte, 0, 4)

	bodyLen = binary.BigEndian.AppendUint32(bodyLen, uint32(len(d.Body)))
	ret = make([]byte, 0, 1024)
	ret = append(ret, d.Type)
	ret = append(ret, bodyLen...)
	ret = append(ret, d.Body...)
	return
}

func GetNegoticationData(b []byte) (ret NegoticationData, err error) {
	if len(b) < 3 {
		err = fmt.Errorf("byte len less 3")
		return
	}
	ret.Type = b[0]
	ret.Length = binary.BigEndian.Uint32(b[1:5])
	ret.Body = b[5 : ret.Length+5]
	return
}
