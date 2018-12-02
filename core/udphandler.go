package core

import (
	"errors"
	"net"

	"github.com/Evan2698/chimney/sercurity"
	"github.com/Evan2698/chimney/utils"
)

// UDPSocket ..
type UDPSocket struct {
	srcsocket net.Conn     // ..
	cipher    string       // ..
	iv        []byte       // ..
	done      chan string  // ..
	info      *ConnectInfo // ..
}

// NewUDPSocket ...
//
func NewUDPSocket(ss net.Conn, c string, i []byte, al *ConnectInfo, ch chan string) *UDPSocket {
	return &UDPSocket{
		srcsocket: ss,
		cipher:    c,
		iv:        i,
		info:      al,
		done:      ch,
	}
}

func sos(ssocket *UDPSocket) ([]byte, []byte, error) {

	buf, err := read_bytes_from_socket(ssocket.srcsocket, 4)
	if err != nil {
		utils.Logger.Print("read length from C failed!", err)
		return nil, nil, err
	}
	utils.Logger.Println("read buffer size", len(buf))

	size := utils.Byte2int(buf)
	if size > BF_SIZE*BF_SIZE*100 || size == 0 {
		utils.Logger.Print("out of memory: ", size)
		return nil, nil, errors.New("out of memory size")
	}

	utils.Logger.Println("read UDP package size", size)

	content, err := read_bytes_from_socket(ssocket.srcsocket, (int)(size))
	if err != nil {
		utils.Logger.Print("Read the UDP from C failed ", err)
		return nil, nil, err
	}

	ori, err := sercurity.DecompressWithChaCha20(content[4:], ssocket.iv[:8], sercurity.MakeCompressKey(ssocket.cipher))
	if err != nil {
		utils.Logger.Print("Decompress the UDP from C failed ", err)
		return nil, nil, err
	}

	return content[:4], ori, nil
}

//DUDP2RAW just client use it
//
func (ssocket *UDPSocket) DUDP2RAW(raw net.Conn) error {

	header, ori, err := sos(ssocket)
	if err != nil {
		ssocket.done <- "done"
		return errors.New("read packet failed")
	}
	addressLen := caladdress(header[3], ori[0])
	if checkUDPTerminates(ori[:addressLen]) {
		ssocket.done <- "done"
		return errors.New("done")
	}

	// client need to keep UDP header format.
	full := append(header, ori...)

	n, err := raw.Write(full)

	utils.Logger.Print("output: ", n, " bytes. ", err)
	if err != nil {
		ssocket.done <- "done"
	}

	return err

}

func caladdress(c, b byte) int {

	var length int
	if c == 0x1 {
		length = 4
	}

	if c == 0x4 {
		length = 16
	}

	if c == 0x3 {
		length = (int)(b) + 1 // 1 is c[4]
	}

	length = length + 2

	return length
}

// Raw2UDP .. just client use it
func (ssocket *UDPSocket) Raw2UDP(raw net.Conn) error {

	buf := make([]byte, BF_SIZE)
	n, err := raw.Read(buf)
	if err != nil {
		utils.Logger.Print("read UDP packet from raw socket failed", err, "bytes: ", n)
		ssocket.done <- "done"
		return err
	}

	utils.Logger.Print("bowser udp: ", n, err)
	if n < 4+2 {
		utils.Logger.Print("UDP format incorrect from raw socket!!!")
		ssocket.done <- "done"
		return errors.New("UDP format incorectly")
	}

	if checkUDPTerminates(buf[4 : 4+caladdress(buf[3], buf[4])]) {
		ssocket.done <- "done"
		return errors.New("client done")
	}

	out, err := sercurity.CompressWithChaCha20(buf[4:], ssocket.iv[:8], sercurity.MakeCompressKey(ssocket.cipher))
	if err != nil {
		utils.Logger.Print("compress UDP failed! ", err)
		ssocket.done <- "done"
		return err
	}

	start := utils.Int2byte((uint32)(len(out) + 4))
	start = append(start, buf[:4]...)
	full := append(start, out...)

	on, err := raw.Write(full)

	utils.Logger.Print("write UDP to SSocket failed! ", err, "write bytes: ", on, "bytes.")

	if err != nil {
		ssocket.done <- "done"
	}

	return err

}

func udpsend(ss *UDPSocket, raw net.Conn) error {

	go func() {
		for {

			err := ss.Raw2UDP(raw)
			if err != nil {
				break
			}

		}
	}()

	go func() {
		for {
			err := ss.DUDP2RAW(raw)
			if err != nil {
				break
			}
		}
	}()

	<-ss.done

	return nil
}

func checkUDPTerminates(buf []byte) bool {

	zero := true

	if buf == nil {
		return false
	}

	for v := range buf {
		if v != 0x0 {
			zero = false
			break
		}
	}

	return zero
}

func (ssocket *UDPSocket) readUDPD() ([]byte, net.Conn, error) {

	header, ori, err := sos(ssocket)
	if err != nil {
		ssocket.done <- "done"
		return nil, nil, errors.New("read packet failed")
	}
	addrLen := caladdress(header[3], ori[0])
	if checkUDPTerminates(ori[:addrLen]) {
		ssocket.done <- "done"
		return nil, nil, errors.New("done")
	}

	// must write raw data without header.
	host := parsehost(header[3], ori[:addrLen])
	udpConnect, err := createDUPConnect(host)
	if err != nil {
		ssocket.done <- "done"
		return nil, nil, errors.New("address invalid for UDP")
	}

	n, err := udpConnect.Write(ori[addrLen:])
	if err != nil {
		ssocket.done <- "done"
		utils.Logger.Print("write UDP failed!", err)
		udpConnect.Close()
		return nil, nil, err
	}

	utils.Logger.Print("write to destination UDP", n, " bytes.")

	return append(header, ori[:addrLen]...), udpConnect, nil
}

func (ssocket *UDPSocket) rawToUPD(raw net.Conn, header []byte) error {

	defer func() {
		if raw != nil {
			raw.Close()
		}
	}()

	buf := make([]byte, BF_SIZE)
	n, err := raw.Read(buf)
	if err != nil {
		ssocket.done <- "done"
		utils.Logger.Print("read udp host failed: ", err, "bytes: ", n)
		return err
	}

	// encapsulate the data with header and compress
	lop := append(header[4:], buf[:n]...)
	ori, err := sercurity.CompressWithChaCha20(lop, ssocket.iv[:8], sercurity.MakeCompressKey(ssocket.cipher))
	if err != nil {
		ssocket.done <- "done"
		utils.Logger.Print("compress UDP failed", err, "bytes: ")
		return err
	}

	full := append(header[:4], ori...)
	ll := len(full)
	out := append(utils.Int2byte((uint32)(ll)), full...)
	n, err = ssocket.srcsocket.Write(out)
	if err != nil {
		ssocket.done <- "done"
		utils.Logger.Print("remote write UDP failed!", err)
	}
	utils.Logger.Print("write ", n, "bytes.", err)

	return err
}

func udpserverRoutine(ss *UDPSocket, h net.Conn) error {

	go func() {

		for {
			header, raw, err := ss.readUDPD()
			if err != nil {
				break
			}

			err = ss.rawToUPD(raw, header)
			if err != nil {
				break
			}

		}

	}()

	<-ss.done

	return nil
}
