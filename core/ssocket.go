package core

import (
	"bytes"
	"errors"
	"io"

	"github.com/Evan2698/chimney/utils"

	"github.com/Evan2698/chimney/security"
)

const (
	bufsize = 4096
	minsize = 1024
)

type sockssocket struct {
	origin io.ReadWriteCloser
	temp   []byte
	I      security.EncryptThings
	pass   string
}

func (s *sockssocket) Read() (buf []byte, err error) {

	if s.I == nil {
		n, err := s.origin.Read(s.temp)
		if err != nil {
			utils.LOG.Println("read raw socket failed: ", err)
			return nil, err
		}

		return s.temp[:n], nil
	}

	cL, err := s.readbytesfromraw(4)
	if err != nil {
		utils.LOG.Println("length of ciphertext read failed: ", err)
		return nil, err
	}

	if len(cL) != 4 {
		return nil, errors.New("read failed from raw reader")
	}

	valen := utils.Bytes2Int(cL)
	raw, err := s.readbytesfromraw(valen)
	if err != nil {
		utils.LOG.Println("read ciphertext failed: ", err)
		return nil, err
	}
	con, err := s.I.Uncompress(raw, security.MakeCompressKey(s.pass))
	return con, err

}

func (s *sockssocket) Write(p []byte) error {
	if p == nil || len(p) == 0 {
		return errors.New("parameter is invalid")
	}

	if s.I == nil {
		_, err := s.origin.Write(p)
		return err
	}

	input, err := s.I.Compress(p, security.MakeCompressKey(s.pass))
	if err != nil {
		utils.LOG.Println("compress failed:", p, " key", security.MakeCompressKey(s.pass))
		return err
	}

	var fo bytes.Buffer
	n := len(input)
	l := utils.Int2Bytes(n)
	fo.Write(l)
	fo.Write(input)
	_, err = s.origin.Write(fo.Bytes())

	return err
}

func (s *sockssocket) Close() error {
	s.origin.Close()
	s.origin = nil
	s.temp = nil
	s.I = nil
	return nil
}

func (s *sockssocket) SetI(i security.EncryptThings) {
	s.I = i
}

func (s *sockssocket) readbytesfromraw(bytes int) ([]byte, error) {

	if bytes <= 0 {
		return nil, errors.New("0 bytes can not read! ")
	}

	buf := make([]byte, bytes)
	index := 0
	var err error
	var n int
	for {
		n, err = s.origin.Read(buf[index:])
		utils.LOG.Println("read from socket size: ", n, err)
		index = index + n
		if err != nil {
			utils.LOG.Println("error on read_bytes_from_socket ", n, err)
			break
		}

		if index >= bytes && index > 0 {
			utils.LOG.Println("read count for output ", index, err)
			break
		}

	}

	if index < bytes && index != 0 {
		utils.LOG.Println("can not run here!!!!!")
	}

	utils.LOG.Println("read result size: ", index, err)
	return buf, err
}

//NewSocksSocket ..
func NewSocksSocket(o io.ReadWriteCloser, passwd string, sec security.EncryptThings) SSocket {
	v := &sockssocket{
		origin: o,
		I:      sec,
		pass:   passwd,
		temp:   make([]byte, minsize),
	}
	return v
}
