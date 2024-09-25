package transport

import (
	"context"
	"fmt"
	"net"
	"tunnel/internal/config"
)

const (
	TransportProtocolTCP = iota
	TransportProtocolUDP
	TransportProtocolETP
)

type TransportStruct struct {
	Protocol   int
	ServerIP   net.IP
	ServerPort int
	LocalIp    net.IP
	ctx        context.Context
	Conn       net.Conn
}

func NewTransport(ctx context.Context, config config.Config) (ret TransportStruct, err error) {
	ret.ctx = ctx
	ret.ServerIP = net.ParseIP(config.GetString("negotiation_server_ip", ""))
	if ret.ServerIP.IsUnspecified() {
		err = fmt.Errorf("negotiation_server_ip:%s IsUnspecified ", config.GetString("negotiation_server_ip", ""))
		return
	}

	ret.LocalIp = net.ParseIP(config.GetString("transport_local_ip", ""))

	ret.ServerPort = config.GetInt("negotiation_server_port", 0)
	if ret.ServerPort <= 0 {
		err = fmt.Errorf("negotiation_server_port:%d is error ", ret.ServerPort)
		return
	}
	Protocol := config.GetString("transport_protocol", "")
	switch Protocol {
	case "tcp":
		ret.Protocol = TransportProtocolTCP
	case "udp":
		ret.Protocol = TransportProtocolUDP
	case "etp":
		ret.Protocol = TransportProtocolETP
	default:
		err = fmt.Errorf("transport_protocol:'%s' error", Protocol)
		return
	}

	return
}

func (c *TransportStruct) CreateConn() (err error) {

	if c.Protocol == TransportProtocolUDP {
		laddr := &net.UDPAddr{
			IP:   c.LocalIp,
			Port: 0,
		}
		raddr := &net.UDPAddr{
			IP:   c.ServerIP,
			Port: c.ServerPort,
		}
		c.Conn, err = net.DialUDP("udp", laddr, raddr)
		return

	}
	if c.Protocol == TransportProtocolTCP {
		laddr := &net.TCPAddr{
			IP:   c.LocalIp,
			Port: 0,
		}
		raddr := &net.TCPAddr{
			IP:   c.ServerIP,
			Port: c.ServerPort,
		}
		c.Conn, err = net.DialTCP("tcp", laddr, raddr)
		return
	}

	return
}
