package core

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io"
	"math/big"
	"net"
	"strconv"
	"time"

	"chimney/utils"

	quic "github.com/lucas-clemente/quic-go"

	"chimney/config"
)

const protocolName = "quic-chimney-what-why"

// Runclientsservice ...
func Runclientsservice(host string, app *config.AppConfig, p SocketService, f DataFlow, quit <-chan int32) {
	all, err := net.Listen("tcp", host)
	if err != nil {
		utils.LOG.Print("local listen on   ip =", host, err)
		return
	}

	defer func() {
		utils.LOG.Println("listener will be close.^_^")
		all.Close()
		utils.LOG.Println("Runclientsservice is over!!!!^_^")
	}()

	for {
		someone, err := all.Accept()
		if err != nil {
			utils.LOG.Print("Accept failed: ", err)
			break
		}
		go handclientonesocket(someone, app, p, f)

		if quit != nil {
			select {
			case <-quit:
				utils.LOG.Println("will be exit!!")
				return
			default:
			}
		}
	}

	utils.LOG.Print("exit exit exit exit", err)
}

func handclientonesocket(o net.Conn, app *config.AppConfig, p SocketService, f DataFlow) {

	utils.SetSocketTimeout(o, uint32(app.Timeout))
	h := NewSocksHandler(o, nil, app)
	defer h.Close()

	err := h.Receive(p)
	if err != nil {
		utils.LOG.Print("client recv failed: ", err)
		return
	}
	h.Run(f)
}

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		utils.LOG.Print(err)
		return nil
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		utils.LOG.Print(err)
		return nil
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		utils.LOG.Print(err)
		return nil
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{protocolName},
	}
}

func makeQuicSocket(sess quic.Session, stream quic.Stream) *CStream {
	ct := &CStream{
		MainStream: stream,
		Hold:       sess,
	}
	return ct
}

func (socket *CStream) setQuickTimeout(tm uint32) CReadWriteCloser {
	readTimeout := time.Duration(tm) * time.Second
	v := time.Now().Add(readTimeout)
	stream, ok := socket.MainStream.(quic.Stream)
	if ok {
		stream.SetReadDeadline(v)
		stream.SetWriteDeadline(v)
		stream.SetDeadline(v)
	}
	return socket
}

// RunServerservice ..
func RunServerservice(host string, app *config.AppConfig, p SocketService, f DataFlow) {

	if app.UseQuic {
		quickHost := net.JoinHostPort(app.Server, strconv.Itoa(int(app.QuicPort)))
		listener, err := quic.ListenAddr(quickHost, generateTLSConfig(), nil)
		if err != nil {
			utils.LOG.Print("Create quick socket failed", host, err)
			return
		}

		for {
			sess, err := listener.Accept(context.Background())
			if err != nil {
				utils.LOG.Print("quic accept session failed ", host, err)
				break
			}
			go handleQuicSession(sess, app, p, f)
		}

	} else {
		all, err := net.Listen("tcp", host)
		if err != nil {
			utils.LOG.Print("local listen on   ip =", host, err)
			return
		}
		defer all.Close()
		for {
			someone, err := all.Accept()
			if err != nil {
				utils.LOG.Print("remote socket failed to open", err)
				break
			}
			utils.SetSocketTimeout(someone, uint32(app.Timeout))
			go handServeronesocket(someone, app, p, f)
		}
	}
}

func handleQuicSession(section quic.Session, app *config.AppConfig, p SocketService, f DataFlow) {

	defer func(s quic.Session) {
		if s != nil {
			s.Close()
		}
		utils.LOG.Print("quic session closed!!")
	}(section)

	for {
		stream, err := section.AcceptStream(context.Background())
		if err != nil {
			utils.LOG.Print("quic accept stream failed", section.LocalAddr().String())
			break
		}
		quick := makeQuicSocket(nil, stream).setQuickTimeout(app.Timeout)
		go handServeronesocket(quick, app, p, f)
	}
}

func handServeronesocket(o io.ReadWriteCloser, app *config.AppConfig, p SocketService, f DataFlow) {
	ss := NewSocksSocket(o, app.Password, nil)
	proxysocket := NewSocketProxy(ss, app)
	h := NewSocksHandler(nil, proxysocket, app)
	defer h.Close()

	err := h.Receive(p)
	if err != nil {
		utils.LOG.Print("client recv failed", err)
		return
	}
	h.Run(f)
}
