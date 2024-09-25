package abyssclinet

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"net"
	"time"
	"tunnel/internal/config"
	"tunnel/internal/fec"
	"tunnel/internal/kcp"
	"tunnel/internal/pool"
	"tunnel/internal/queue"
	"tunnel/internal/transport"
	tundrive "tunnel/internal/tun_drive"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"tunnel/internal/negotication"
	wfrsa "tunnel/internal/wf_rsa"
)

var errNegoticationFinsh = errors.New("Negotication finsh")

const MagicStr = "7x"

type Client struct {
	Ras          *wfrsa.RSAInfo
	Negotication negotication.Negotication
	ClientID     string
	ConfigPath   string
	Config       config.Config
	CancelFunc   context.CancelFunc
	Fec          fec.FEC
	KCP          kcp.KCP
	TUNDrive     tundrive.TUNDrive
	ServerIP     net.IP
	ctx          context.Context
	Transport    transport.TransportStruct
	IngressQueue queue.QueueManage
	EngressQueue queue.QueueManage
}

var DecPakcetFun = func(c Client) func(*pool.TUNPacket) error {
	return func(t *pool.TUNPacket) (err error) {
		t.EncByte, err = c.Ras.Aes256GcmDec(t.DecByte)
		return
	}
}

var EncPakcetFun = func(c Client) func(*pool.TUNPacket) error {
	return func(t *pool.TUNPacket) (err error) {
		t.DecByte = c.Ras.Aes256GcmEnc(t.DecByte)
		return
	}
}

func NewClinet(ctx context.Context, ConfigPath string) (r Client, err error) {
	// var ctx context.Context
	// ctx, r.CancelFunc = context.WithCancel(ctx)

	r.ctx = ctx
	r.ConfigPath = ConfigPath
	r.Config, err = config.NewConfigPath(r.ConfigPath)
	if err != nil {
		return Client{}, err
	}

	r.Negotication, err = negotication.NewNegotication(r.Config)

	r.ClientID = r.Config.GetString("client_id", "")
	if len(r.ClientID) == 0 {
		err = fmt.Errorf("client_id is empty")
		return Client{}, err
	}

	r.Ras, err = wfrsa.NewRSA(r.Config)
	if err != nil {
		return Client{}, err
	}

	r.Fec, err = fec.NewFec(ctx, r.Config)
	if err != nil {
		return Client{}, err
	}

	r.KCP, err = kcp.NewKCP(ctx, r.Config)
	if err != nil {
		return Client{}, err
	}

	r.Transport, err = transport.NewTransport(ctx, r.Config)
	if err != nil {
		return Client{}, err
	}

	r.TUNDrive = tundrive.NewDefalutTunDrive(ctx, r.Config)

	r.IngressQueue = queue.NewQueueManage()
	// 添加队列检查magic等信息
	r.IngressQueue.AddQueue(ctx, queue.QUEUE_TYPE_IP_PACKET, func(packet *pool.TUNPacket) (err error) {
		var length int16
		if len(packet.RawByte) <= 8 {
			return fmt.Errorf("len less 8")
		}
		length = int16(packet.RawByte[2]<<8) + int16(packet.RawByte[3])
		if length != int16(len(packet.RawByte)) {
			return fmt.Errorf("len err")
		}

		if subtle.ConstantTimeCompare(packet.RawByte[:2], []byte(MagicStr)) != 1 {
			return fmt.Errorf("magic err")
		}
		if subtle.ConstantTimeCompare(packet.RawByte[4:8], r.ServerIP) != 1 {
			return fmt.Errorf("ServerIP err")
		}
		return
	})

	// 添加解密队列信息
	r.IngressQueue.AddQueue(ctx, queue.QUEUE_TYPE_IP_PACKET, func(packet *pool.TUNPacket) (err error) {
		packet.DecByte, err = r.Ras.Aes256GcmDec(packet.RawByte[8:])
		if err != nil {
			return err
		}
		packet.PacketType = int(packet.RawByte[0] >> 4)
		packet.EnableFEC = ((packet.RawByte[0] & 0x80) == 0x80)

		if packet.PacketType == negotication.WAN_FENG_TYPE_HEARTBEAT_ECHO_REPLY {
			return queue.ErrTakOver
		}
		if packet.PacketType == negotication.WAN_FENG_TYPE_PAYLOAD {
			r.TUNDrive.WritePacketChannle <- packet
		}
		return
	})

	r.EngressQueue = queue.NewQueueManage()

	// 分析数据tun卡的数据包信息
	r.EngressQueue.AddQueue(ctx, queue.QUEUE_TYPE_IP_PACKET, func(packet *pool.TUNPacket) (err error) {
		packeteth := gopacket.NewPacket(packet.RawByte, layers.LayerTypeIPv4, gopacket.Default)

		return
	})

	return
}

func (c Client) StartNegotication() (err error) {
	if !c.Negotication.IsRunning.CompareAndSwap(false, true) {
		return fmt.Errorf("Negotication is running")
	}

	defer c.Negotication.IsRunning.Store(false)

	readByte := make([]byte, 2048)

	var NegoticationData negotication.NegoticationData

	// 开始发送协商第一个包
	sendData := negotication.NegoticationData{
		Type:   negotication.WAN_FENG_TYPE_CLIENT_ID,
		Length: uint32(len(c.ClientID)),
		Body:   []byte(c.ClientID),
	}
	_, err = c.Negotication.SendData(sendData)

	if err != nil {
		c.Negotication.CancelFunc()
		return
	}

	for {
		_, err = c.Negotication.Conn.Read(readByte)
		if err != nil {
			c.Negotication.CancelFunc()
			return
		}
		NegoticationData, err = negotication.GetNegoticationData(readByte)
		if err != nil {
			c.Negotication.CancelFunc()
			return
		}
		err = c.getData(NegoticationData)

		// errNegoticationFinsh 表示完成协商
		if errors.Is(err, errNegoticationFinsh) {
			return nil
		}
		c.Negotication.CancelFunc()
		return
	}
}

func (c Client) StartTransport() (err error) {

	// 协商完成，启动tun卡
	if err = c.TUNDrive.TunInfo.CreateTunDrive(); err != nil {
		c.Negotication.CancelFunc()
		return
	}

	// tun卡启动
	if err = c.TUNDrive.TunInfo.RunTunDrive(); err != nil {
		c.Negotication.CancelFunc()
		return
	}

	// 创建连接
	if err = c.Transport.CreateConn(); err != nil {
		c.Negotication.CancelFunc()
		return
	}

	// 开启tun卡的数据传输
	c.TUNDrive.Start()

	// 开启心跳
	go c.StartHeartbeat()

	return
}

func (c Client) StartHeartbeat() {
	ticker := time.NewTicker(time.Second * 1)
	sendData := negotication.NegoticationData{}
	sendData.Type = negotication.WAN_FENG_TYPE_HEARTBEAT
	sendData.Length = 0

	ret := sendData.GetSendByte()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.Transport.Conn.Write(ret)
		}
	}
}

func (c Client) getData(NegoticationData negotication.NegoticationData) (err error) {
	var sendND negotication.NegoticationData
	switch NegoticationData.Type {
	case negotication.WAN_FENG_TYPE_HEARTBEAT:

	case negotication.WAN_FENG_TYPE_CLIENT_ID:
		sendND.Type = negotication.WAN_FENG_TYPE_RSA
		sendND.Body, err = c.Ras.GetPublicKeyByte()
		if err != nil {
			err = fmt.Errorf("WAN_FENG_TYPE_CLIENT_ID GetPublicKeyByte err:%s ", err.Error())
			return
		}
		sendND.Length = uint32(len(sendND.Body))
		_, err = c.Negotication.SendData(sendND)
		if err != nil {
			err = fmt.Errorf("WAN_FENG_TYPE_CLIENT_ID SendData err:%s ", err.Error())
			return
		}
	case negotication.WAN_FENG_TYPE_RSA:
		c.Ras.SetSeverPublicKey(NegoticationData.Body)
		sendND.Type = negotication.WAN_FENG_TYPE_KEY
		sendND.Length = 1
		sendND.Body = []byte{byte(c.Ras.EncType)}
		_, err = c.Negotication.SendData(sendND)
		if err != nil {
			err = fmt.Errorf("WAN_FENG_TYPE_RSA SendData err:%s ", err.Error())
			return
		}
	case negotication.WAN_FENG_TYPE_KEY:
		if err = c.negotication_chekc_key(NegoticationData); err != nil {
			err = fmt.Errorf("WAN_FENG_TYPE_KEY negotication_chekc_key err:%s ", err.Error())
			return
		}
		sendND.Type = negotication.WAN_FENG_TYPE_REDUNDANCY
		sendND.Length = 2
		sendND.Body = append(sendND.Body, byte(c.Fec.FecMode))
		if c.KCP.Enable {
			sendND.Body = append(sendND.Body, 1)
		} else {
			sendND.Body = append(sendND.Body, 0)
		}
		_, err = c.Negotication.SendData(sendND)
		if err != nil {
			err = fmt.Errorf("WAN_FENG_TYPE_KEY SendData err:%s ", err.Error())
			return
		}
	case negotication.WAN_FENG_TYPE_REDUNDANCY:
		if NegoticationData.Body[0] != byte(c.Fec.FecMode) {
			err = fmt.Errorf("WAN_FENG_TYPE_KEY WAN_FENG_TYPE_REDUNDANCY fec is not same FecMode:%d config:%d ", int(NegoticationData.Body[0]), c.Fec.FecMode)
			return
		}
		enable_kcp := 0
		if c.KCP.Enable {
			enable_kcp = 1
		}
		if NegoticationData.Body[1] != byte(enable_kcp) {
			err = fmt.Errorf("WAN_FENG_TYPE_KEY WAN_FENG_TYPE_REDUNDANCY kcp is not same enable_kcp:%d config:%d ", int(NegoticationData.Body[1]), enable_kcp)
			return
		}
		sendND.Type = negotication.WAN_FENG_TYPE_IP
		sendND.Length = 0
		_, err = c.Negotication.SendData(sendND)
		if err != nil {
			err = fmt.Errorf("WAN_FENG_TYPE_REDUNDANCY  SendData err:%s ", err.Error())
			return
		}
	case negotication.WAN_FENG_TYPE_IP:
		if err = c.negotication_check_ip(NegoticationData); err != nil {
			err = fmt.Errorf("WAN_FENG_TYPE_IP  negotication_check_ip  err:%s ", err.Error())
			return
		}

		return errNegoticationFinsh
	}

	return
}

func (c Client) negotication_chekc_key(data negotication.NegoticationData) (err error) {
	if c.Ras.EncType == wfrsa.ENC_NONE {
		return
	}
	ras_size := c.Ras.Prvkey.Size()
	if c.Ras.EncType == wfrsa.ENC_AES_GCM {
		if len(data.Body) != ras_size*3 {
			return
		}

		c.Ras.ServerEncInfo.Key, err = c.Ras.RsaDecrypt(data.Body[0:ras_size])
		if err != nil {
			return err
		}

		c.Ras.ServerEncInfo.IV, err = c.Ras.RsaDecrypt(data.Body[ras_size : ras_size*2])
		if err != nil {
			return err
		}

		c.Ras.ServerEncInfo.ADD, err = c.Ras.RsaDecrypt(data.Body[ras_size*2 : ras_size*3])
		if err != nil {
			return err
		}
	} else if c.Ras.EncType == wfrsa.ENC_AES_GCM {
		if len(data.Body) != ras_size*2 {
			return
		}
		c.Ras.ServerEncInfo.Key, err = c.Ras.RsaEncrypt(data.Body[0:ras_size])
		if err != nil {
			return err
		}

		c.Ras.ServerEncInfo.IV, err = c.Ras.RsaEncrypt(data.Body[ras_size : ras_size*2])
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("EncType is err EncType:%d", c.Ras.EncType)
	}
	return
}

func (c Client) negotication_check_ip(data negotication.NegoticationData) (err error) {
	ras_size := c.Ras.Prvkey.Size()
	var rsaEncByte []byte
	if len(data.Body) != ras_size*3 {
		return fmt.Errorf("lenth error")
	}
	rsaEncByte, err = c.Ras.RsaDecrypt(data.Body[0:ras_size])
	if err != nil {
		return
	}
	c.TUNDrive.TunInfo.LocalIP = net.ParseIP(string(rsaEncByte))

	rsaEncByte, err = c.Ras.RsaDecrypt(data.Body[ras_size : ras_size*2])
	if err != nil {
		return
	}
	c.TUNDrive.TunInfo.Mask = net.ParseIP(string(rsaEncByte))

	rsaEncByte, err = c.Ras.RsaDecrypt(data.Body[ras_size*2 : ras_size*3])
	if err != nil {
		return
	}
	c.ServerIP = net.ParseIP(string(rsaEncByte))

	return
}
