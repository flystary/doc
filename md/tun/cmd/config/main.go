package main

import "tunnel/internal/config"

func main() {
	if err := config.NewConfigPath("/home/code/tunnel_go/config/test.config"); err != nil {
		println(err.Error())
		return
	}
	println(config.C.GetBool("enable_kcp", true))
}
