package core

import (
	"errors"
	"net"
	"climbwall/utils"
	"climbwall/sercurity"

)

const (
	BF_SIZE =  5120
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


func (ssocket *SSocketWrapper) CopyFromC2Raw(raw net.Conn) (err error) {

	buf := make([]byte, 4)
	n, err := ssocket.src_socket.Read(buf)
	if err != nil || n <= 0 {
		utils.Logger.Print("read length from C failed!", err,  "bytes: ", n)
		return err
	}
	utils.Logger.Println("read buffer size",  buf)

	size := utils.Byte2int(buf);
	if (size > BF_SIZE  * BF_SIZE * 100) {
		utils.Logger.Print("out of memory: ", size)
		return errors.New("out of memory size")
	}

	utils.Logger.Println("read content size", size)

	content := make([]byte, size)
	cn, err := ssocket.src_socket.Read(content)
	if err != nil {
		utils.Logger.Print("Read the content from C failed ", err, " bytes: ", cn)
		return err
	}

	utils.Logger.Println("read C content", size)

	out, err := sercurity.DecompressWithChaCha20(content, ssocket.iv[:8], sercurity.MakeCompressKey(ssocket.cipher))
	if err != nil {
		utils.Logger.Print("uncompressed content failed! ", err)
		return err
	}

	on, err := raw.Write(out)
	if err != nil {
		utils.Logger.Print("send to remote failed ", err, "write bytes: ", on)
		return err
	}

	return nil
}

func (ssocket *SSocketWrapper) CopyFromRaw2C(raw net.Conn) (err error){

	buf := make([]byte, BF_SIZE)
	n, err := raw.Read(buf)
	if err != nil{
		utils.Logger.Print("read content from raw socket failed", err)
		return err
	}

	utils.Logger.Print("bowser content: ", buf[0:n])

	out, err := sercurity.CompressWithChaCha20(buf[0:n], ssocket.iv[:8], sercurity.MakeCompressKey(ssocket.cipher))
	if err != nil {
		utils.Logger.Print("compress content failed! ", err)
		return err
	}

	utils.Logger.Print("length of buf ", len(out))
	start := utils.Int2byte((uint32)(len(out)))
	ll := append(start, out...)
	utils.Logger.Print("bowser content:(ALL): ", ll)

	on, err := ssocket.src_socket.Write(ll)
	if err != nil {
		utils.Logger.Print("write content to SSocket failed! ", err, "write bytes: ", on, "bytes.")
		return err
	}

	return nil
}

func Copy_C2RAW(ssl *SSocketWrapper, raw net.Conn){

	for {
		neterr := ssl.CopyFromC2Raw(raw)
		if neterr != nil {
			utils.Logger.Println("failed or compeleted (C -->Remote)", neterr)
			break
		}
	}
}

func Copy_RAW2C(ssl *SSocketWrapper, raw net.Conn){
	
	for {
		neterr := ssl.CopyFromRaw2C(raw)
		if neterr != nil {
			utils.Logger.Println("failed or completed (remote--->C) ", neterr)
			break
		}
	}
}