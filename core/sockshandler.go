package core

import (
	"bytes"
	"crypto/hmac"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"

	"chimney/geo"
	"chimney/utils"

	"chimney/security"

	"chimney/config"
)

const (
	//CMDCONNECT ...
	CMDCONNECT = 0x1 // connect
)

type socksreceive struct {
	proxy  SocksProxy
	src    io.ReadWriteCloser
	dst    io.ReadWriteCloser
	appcon *config.AppConfig
}

func (s *socksreceive) createproxy(app *config.AppConfig, p SocketService) (SocksProxy, error) {
	con, err := createclientsocket(p, "tcp", app)
	if err != nil {
		utils.LOG.Print("create socket failed", err)
		return nil, err
	}

	ss := NewSocksSocket(con, app.Password, nil)
	proxysocket := NewSocketProxy(ss, app)

	return proxysocket, nil
}

func (s *socksreceive) Receive(p SocketService) error {

	if s.src == nil {
		return s.handleServerResponse()
	}

	buf := make([]byte, 512)
	n, err := s.src.Read(buf)
	if err != nil {
		return err
	}

	utils.LOG.Println(buf[:n])

	if n < 1 || buf[0] != 5 {
		s.src.Write([]byte{0x05, 0xff})
		utils.LOG.Println("can not support socks flag: ", buf[:n])
		return errors.New("can not support socks version: ")
	}

	// No authentication required
	s.src.Write([]byte{0x05, 0x00})

	n, err = s.src.Read(buf)
	if err != nil {
		utils.LOG.Print("read from client failed!", err)
		return err
	}

	//the type of ip is mini length
	if n < 10 {
		s.src.Write([]byte{0x05, 0x0A, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		utils.LOG.Print("socks format error", err)
		return errors.New("socks format error")
	}

	if CMDCONNECT != buf[1] {
		s.src.Write([]byte{0x05, 0x07, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		utils.LOG.Print("it does not support this method.", err)
		return errors.New("it does not support this method")
	}

	data := buf[:n]
	utils.LOG.Println("CMD:", data)
	if data[3] == 0x1 || data[3] == 0x4 {
		ip := net.IP(data[4 : len(data)-2])
		result := geo.QueryCountryByIP(ip)
		if result == "CN" {
			port := binary.BigEndian.Uint16(data[len(data)-2:])
			host := net.JoinHostPort(ip.String(), strconv.Itoa(int(port)))
			con, err := CreateCommonSocket(host, "tcp", s.appcon.Timeout, p)
			if err != nil {
				s.src.Write([]byte{0x05, 0x03, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
				utils.LOG.Print("it does not support this method.", err)
				return err
			}
			s.proxy = NewDirectProxy(con)
		}
	}

	if s.proxy == nil {
		s.proxy, err = s.createproxy(s.appcon, p)
		if err != nil {
			s.src.Write([]byte{0x05, 0x04, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
			utils.LOG.Print("c", err)
			return errors.New("can not connect proxy sever")
		}
	}

	err = s.proxy.Connect(data)
	if err != nil {
		s.src.Write([]byte{0x05, 0x05, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		utils.LOG.Print("can not connect server.", err)
		return errors.New("can not connect server")
	}

	var ans bytes.Buffer
	ans.Write([]byte{0x05, 0x00, 0x00, 0x1, 0x00, 0x00, 0x00, 0x00})
	ans.WriteByte(byte((s.appcon.LocalPort >> 8) & 0xff))
	ans.WriteByte(byte(s.appcon.LocalPort & 0xff))
	s.src.Write(ans.Bytes())
	utils.LOG.Print("write (ok) to browser connection sucessed!")

	return nil
}

func (s *socksreceive) handleServerResponse() error {

	buf, err := s.proxy.Read()
	if err != nil {
		utils.LOG.Print("can not connect server.", err)
		return err
	}
	if len(buf) < 2 {
		s.proxy.Write([]byte{0x05, 0xff})
		utils.LOG.Print("socks format error.")
		return errors.New("socks format error")
	}
	if buf[0] != 5 {
		s.proxy.Write([]byte{0x05, 0xfA})
		utils.LOG.Print("can not support this version.")
		return errors.New("can not support this version")
	}

	I := security.NewEncryptyMethod("chacha20")

	var out bytes.Buffer
	out.Write([]byte{0x05, 0x02}) // need user name & password
	out.Write(security.ToBytes(I))
	s.proxy.Write(out.Bytes())

	//verify the password
	user, err := s.proxy.Read()
	if err != nil {
		utils.LOG.Print(err)
		return err
	}

	pu := bytes.NewBuffer(user)
	pu.Next(1) // 0x5
	usrlen := pu.Next(1)
	ul := int(usrlen[0])
	name := pu.Next(ul)
	hashlen := pu.Next(1)
	hl := int(hashlen[0])
	hash := pu.Next(hl)

	nametext, err := I.Uncompress(name, security.MakeCompressKey(s.appcon.Password))
	if err != nil {
		utils.LOG.Print("uncompress user name failed")
		return err
	}

	nhash := security.BuildMacHash(I.GetIV(), string(nametext))
	if !hmac.Equal(nhash, hash) {
		s.proxy.Write([]byte{0x05, 0x2})
		utils.LOG.Print("user name verify failed!")
		return errors.New("user name verify failed")
	}

	s.proxy.Write([]byte{0x05, 0x0})

	// handle connect
	oush, err := s.proxy.Read()
	if err != nil {
		s.proxy.Write([]byte{0x05, 0x05, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		utils.LOG.Print(err)
		return err
	}
	if len(oush) < 10 {
		s.proxy.Write([]byte{0x05, 0x0A, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		utils.LOG.Print("socks format error", err)
		return errors.New("socks format error")
	}

	domain, err := I.Uncompress(oush[4:len(oush)-2], security.MakeCompressKey(s.appcon.Password))
	if err != nil {
		s.proxy.Write([]byte{0x05, 0x0B, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		utils.LOG.Print("socks format error", err)
		return err
	}

	ip := "127.0.0.1"
	if oush[3] == 0x1 || oush[3] == 0x4 {
		ip = net.IP(domain).String()
	} else if oush[3] == 0x3 {
		ip = string(domain[1:])
	}

	port := binary.BigEndian.Uint16(oush[len(oush)-2:])

	host := net.JoinHostPort(ip, strconv.Itoa(int(port)))
	host = strings.Trim(host, " \n")

	utils.LOG.Println("server host: ", "|"+host)

	remote, err := CreateCommonSocket(host, "tcp", s.appcon.Timeout, nil)
	if err != nil {
		s.proxy.Write([]byte{0x05, 0x03, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		utils.LOG.Print("socks format error", err)
		return err
	}
	s.dst = remote
	if 443 == port && s.appcon.SSLRaw {
		s.proxy.Write([]byte{0x05, 0xEF, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		s.proxy.SetEncrypt(security.NewEncryptyMethod("raw"))
		utils.LOG.Println("use 443 protocol for encryption.")
	} else {
		s.proxy.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		s.proxy.SetEncrypt(I)
		utils.LOG.Println("use normal protocol for encryption.")
	}
	return nil
}

func (s *socksreceive) Close() error {
	if s.src != nil {
		s.src.Close()
		s.src = nil
	}

	if s.dst != nil {
		s.dst.Close()
		s.dst = nil
	}

	if s.proxy != nil {
		s.proxy.Close()
		s.proxy = nil
	}
	s.appcon = nil

	return nil

}

func (s *socksreceive) Run(f DataFlow) {

	ch := make(chan int, 1)
	if s.src != nil {
		go readrawfirst(s.src, s.proxy, ch)
		readproxyfirst(s.src, s.proxy, nil)
	} else {
		go readproxyfirst(s.dst, s.proxy, ch)
		readrawfirst(s.dst, s.proxy, nil)
	}

	<-ch
	close(ch)
	utils.LOG.Print("one time is over!!!!!!")
}

func readrawfirst(raw io.ReadWriteCloser, proxy SocksProxy, ch chan int) {

	buf := make([]byte, bufsize)
	var err error
	for {
		n, err := raw.Read(buf)
		if err != nil {
			utils.LOG.Print("read raw socket failed!", err)
			break
		}
		//utils.LOG.Println("READ FROM RAW: ", buf[:n])
		err = proxy.Write(buf[:n])
		if err != nil {
			utils.LOG.Print("proxy write failed!", err)
			break
		}
	}
	buf = nil
	if ch != nil {
		ch <- 1
	}
	utils.LOG.Print("END of readrawfirst: ", err)
}

func readproxyfirst(raw io.ReadWriteCloser, proxy SocksProxy, ch chan int) {
	var err error
	for {
		pout, err := proxy.Read()
		if err != nil {
			utils.LOG.Print("proxy read failed!", err)
			break
		}

		//utils.LOG.Println("READ PROXY: ", pout)

		_, err = raw.Write(pout)
		if err != nil {
			utils.LOG.Print("raw socket write failed", err)
			break
		}
	}

	if ch != nil {
		ch <- 1
	}
	utils.LOG.Print("END of readproxyfirst: ", err)
}

// NewSocksHandler ...
func NewSocksHandler(s net.Conn, p SocksProxy, app *config.AppConfig) SocksHandler {

	v := &socksreceive{
		src:    s,
		dst:    nil,
		proxy:  p,
		appcon: app,
	}
	return v
}
