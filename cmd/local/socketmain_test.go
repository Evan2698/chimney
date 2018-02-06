package main

import (
	"climbwall/core"
	"climbwall/utils"
)

func main() {
	utils.Logger.Println("start-------------------------------")

	for j := 0; j < 20; j = j + 1 {

		//176.122.157.41
		//61.135.169.125
		//61.135.169.125
		conn, fd, err := core.Build_low_socket("176.122.157.41", 443)

		if err != nil {
			utils.Logger.Println("[golang] socket error", err)
		}

		utils.Logger.Println("[golang] socket fd", fd)

		utils.Logger.Println("[golang] socket ok~~")

		//fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
		//status, err := bufio.NewReader(conn).ReadString('\n')

		//utils.Logger.Println(status)
		utils.Logger.Println(err)

		conn.Close()
	}

	ch := make(chan int, 1)

	i := <-ch
	utils.Logger.Println("[golang] wait i", i)
	close(ch)

}
