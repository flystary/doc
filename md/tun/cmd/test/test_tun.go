package main

import (
	"context"
	abyssclinet "tunnel/abyss_clinet"
)

func main() {
	ctx, _ := context.WithCancel(context.Background())
	c, err := abyssclinet.NewClinet(ctx, "/home/code/tunnel_go/config/test.yaml")
	if err != nil {
		println(err.Error())
		return
	}

	// r.Ras.ServerEncInfo.Key =
	if err = c.StartNegotication(); err != nil {
		println(err.Error())
		return
	}
	c.StartTransport()
	println("fd :")
	// key := "vQ981znJ1fhW0e4D2OQ6FKXiBb5A015J"
	// // key := "89c54b0d3bc3c397d5039058c220685f"
	// c.Ras.ServerEncInfo.Key = []byte(key)
	// if err != nil {
	// 	println(err.Error())
	// 	return
	// }

	// iv := "4pLvbnDyL165"
	// c.Ras.ServerEncInfo.IV = []byte(iv)

	// aad := "XbCxkqwq3YcHb183"
	// c.Ras.ServerEncInfo.ADD = []byte(aad)

	// if err = c.Ras.GenAESCipher(); err != nil {
	// 	println(err.Error())
	// 	return
	// }
	// str := "1234567899123456780"
	// ret := c.Ras.Aes256GcmEnc([]byte(str))

	// fmt.Printf("% x\n", ret)
	// ret2, err := c.Ras.Aes256GcmDec(ret)
	// if err != nil {
	// 	println(err.Error())
	// 	return
	// }

	// fmt.Printf("%s \n", ret2)
	// println("1111")

	// clientId := "1234567890"
	// data := negotication.NegoticationData{}
	// data.Type = 1
	// data.Length = uint32(len(clientId))
	// data.Body = []byte(clientId)
	// fmt.Printf("% x \n", data.GetSendByte())

	// Config, err := config.NewConfigPath("/home/code/tunnel_go/config/test.yaml")
	// if err != nil {
	// 	println(err.Error())
	// 	return
	// }
	// RSA, err := wfrsa.NewRSA(Config)
	// if err != nil {
	// 	println(err.Error())
	// 	return
	// }
	// println(RSA.Prvkey.Size())

	// tuninfo := tundrive.TUNInfo{}
	// tuninfo.Name = "tun_1"
	// tuninfo.LocalIP = net.ParseIP("192.168.100.22")
	// tuninfo.Mask = net.ParseIP("255.255.0.0")
	// tuninfo.Mtu = 1111
	// tuninfo.Type = tundrive.TUN_TYPE_TUN

	// err := tuninfo.CreateTunDrive()
	// if err != nil {
	// 	println(err.Error())
	// }

	// tuninfo.RunTunDrive()

	// println("fd :")
}
