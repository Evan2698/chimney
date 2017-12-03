package core

import (
	"errors"
	"net"
	"climbwall/utils"
	"climbwall/sercurity"

)

const (
	BF_SIZE =  8096
)

type SSocketWrapper struct {
	src_socket net.Conn
	cipher string
	iv []byte
}

func NewSSocket(ss net.Conn, c string, i []byte) *SSocketWrapper {
	return & SSocketWrapper {
		src_socket : ss,
		cipher : c, 
		iv : i,		
	}
}


func (ssocket *SSocketWrapper) Copy2RaW(raw net.Conn) (err error) {

	buf := make([]byte, 4)
	n, err := ssocket.src_socket.Read(buf)
	if err != nil || n <= 0 {
		utils.Logger.Print("parse int from socket failed! ", err,  "bytes: ", n)
		return err
	}

	size := utils.Byte2int(buf);
	if (size > 4096  * 4096) {
		utils.Logger.Print("parse int from socket failed! ", err,  "bytes: ", n)
		return errors.New("out of memory size")
	}

	content := make([]byte, size)
	cn, err := ssocket.src_socket.Read(content)
	if err != nil || cn <= 0 {
		utils.Logger.Print("parse content from socket failed! ", err,  "bytes: ", cn)
		return err
	}

	out, err := sercurity.Uncompress(content, ssocket.iv, sercurity.MakeCompressKey(ssocket.cipher))
	if err != nil {
		utils.Logger.Print("content decrypt failed! ", err)
		return err
	}

	on, err := raw.Write(out)
	if err != nil {
		utils.Logger.Print("write content to raw failed! ", err, "write bytes: ", on)
		return err
	}

	return nil
}

func (ssocket *SSocketWrapper) WriteFromRaw(raw net.Conn) (err error){

	buf := make([]byte, BF_SIZE)
	rn, err := raw.Read(buf)
	if err != nil || rn <= 0 {
		utils.Logger.Print("read from raw socket ", err)
		return err
	}

	out, err := sercurity.Compress(buf, ssocket.iv, sercurity.MakeCompressKey(ssocket.cipher))
	if err != nil {
		utils.Logger.Print("content encrypt failed! ", err)
		return err
	}

	on, err := ssocket.src_socket.Write(utils.Int2byte((uint32)(len(out))))
	if err != nil {
		utils.Logger.Print("write content to SSocket failed! ", err, "write bytes: ", on, "bytes.")
		return err
	}

	n, err := ssocket.src_socket.Write(out)
	if err != nil {
		utils.Logger.Print("write content to SSocket failed! ", err, "write bytes: ", n, "bytes.")
		return err
	}

	return nil
}