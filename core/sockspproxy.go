package core

import (
	"bytes"
	"encoding/hex"
	"errors"
	"net"

	"github.com/Evan2698/chimney/utils"

	"github.com/Evan2698/chimney/security"

	"github.com/Evan2698/chimney/config"
)

type psocks struct {
	path SSocket
	app  *config.AppConfig
}

func (s *psocks) Read() (buf []byte, err error) {
	return s.path.Read()
}

func (s *psocks) Write(p []byte) error {
	return s.path.Write(p)
}
func (s *psocks) SetEncrypt(I security.EncryptThings) {
	s.path.SetI(I)
}

func (s *psocks) Close() error {
	if s.path != nil {
		s.path.Close()
	}
	s.app = nil
	return nil
}

func (s *psocks) Connect(remoteaddr []byte) error {

	if remoteaddr == nil {
		return errors.New("remote address is invalid")
	}

	err := s.path.Write([]byte{0x05, 0x01, 0x00})
	if err != nil {
		return errors.New("can not write socks flag to server")
	}

	buf, err := s.path.Read()
	if err != nil {
		utils.LOG.Print("read failed from proxy server")
		return err
	}

	if len(buf) < 2 || buf[1] != 0x2 || buf[0] != 0x5 {
		utils.LOG.Print("failed: server return value: ", hex.EncodeToString(buf))
		return errors.New("server can not handle it")
	}

	I, err := security.FromByte(buf[2:])
	if err != nil {
		utils.LOG.Print("negotiate encrypting algorithm	failed", hex.EncodeToString(buf[2:]))
		return errors.New("negotiate encrypting algorithm failed")
	}

	pw, err := I.Compress([]byte(s.app.Password), security.MakeCompressKey(s.app.Password))
	if err != nil {
		utils.LOG.Print("passwd encrypt failed", err)
		return err
	}

	hmac := security.BuildMacHash(I.GetIV(), s.app.Password)

	var non bytes.Buffer
	non.WriteByte(0x5)
	non.WriteByte(byte(len(pw)))
	non.Write(pw)
	non.WriteByte(byte(len(hmac)))
	non.Write(hmac)

	err = s.path.Write(non.Bytes())
	if err != nil {
		utils.LOG.Print("write user and passwd to server failed", err)
		return err
	}

	buf, err = s.path.Read()
	if err != nil || len(buf) < 2 {
		utils.LOG.Print("server verify failed", err)
		return err
	}

	if buf[0] != 5 || buf[1] != 0 {
		utils.LOG.Print("user and passwd incorrect.  server code: ", hex.EncodeToString(buf))
		return errors.New("user and passwd incorrect")
	}

	// connect
	m := remoteaddr[4 : len(remoteaddr)-2]
	mn, err := I.Compress(m, security.MakeCompressKey(s.app.Password))
	if err != nil {
		utils.LOG.Print("compress the remote address failed", err)
		return errors.New("compress the remote address failed")
	}

	var connect bytes.Buffer
	connect.Write(remoteaddr[:4])
	connect.Write(mn)
	connect.Write(remoteaddr[len(remoteaddr)-2:])
	s.path.Write(connect.Bytes())

	ans, err := s.path.Read()
	if err != nil {
		utils.LOG.Print("read connnect response failed", err)
		return err
	}

	if len(ans) < 2 {
		utils.LOG.Print("connect response format incorrect", ans)
		return errors.New("connect response format incorrect")
	}

	utils.LOG.Print("connect result: ", ans)

	if ans[1] != 0 {
		if ans[1] != 0xEF {
			utils.LOG.Print("connect failed code: ", hex.EncodeToString(ans))
			return errors.New("connect failed")
		}
		I = security.NewEncryptyMethod("raw")
	}
	utils.LOG.Print("encryption name : " + I.GetName())
	s.path.SetI(I)
	return nil
}

//NewSocketProxy ...
func NewSocketProxy(c SSocket, a *config.AppConfig) SocksProxy {
	return &psocks{
		path: c,
		app:  a,
	}
}

type directcn struct {
	rawCon net.Conn
}

func (s *directcn) Read() (buf []byte, err error) {

	buf = make([]byte, bufsize)
	n, err := s.rawCon.Read(buf)
	return buf[:n], err
}

func (s *directcn) Write(p []byte) error {
	_, err := s.rawCon.Write(p)
	return err
}

func (s *directcn) SetEncrypt(I security.EncryptThings) {

}

func (s *directcn) Close() error {
	if s.rawCon != nil {
		s.rawCon.Close()
	}
	s.rawCon = nil
	return nil
}

func (s *directcn) Connect(remoteaddr []byte) error {
	return nil
}

//NewDirectProxy ..
func NewDirectProxy(con net.Conn) SocksProxy {
	return &directcn{
		rawCon: con,
	}
}
