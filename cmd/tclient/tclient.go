package main

import (
	"fmt"
	"net"

	"github.com/Evan2698/chimney/security"

	"github.com/Evan2698/chimney/config"

	"github.com/Evan2698/chimney/core"
)

func main() {
	config, err := config.Parse("../../config.json")
	fmt.Print(config)

	con, err := net.Dial("tcp", "127.0.0.1:5000")
	if err != nil {
		fmt.Print("connect failed!")
		return
	}

	temp := make([]byte, 200)
	n, _ := con.Read(temp)

	gcm, _ := security.FromByte(temp[:n])

	v := core.NewSocksSocket(con, "zhangweihua", gcm)

	defer v.Close()
	v.Write([]byte{0x1, 0x2, 0x3})
	buf, err := v.Read()
	fmt.Println(string(buf), err)
}
