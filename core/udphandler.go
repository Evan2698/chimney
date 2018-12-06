package core

import (
	"errors"
	"net"
	"strconv"

	"github.com/Evan2698/chimney/sercurity"
	"github.com/Evan2698/chimney/utils"
)

// UDPSocket ..
type UDPSocket struct {
	srcsocket net.Conn     // ..
	config    *AppConfig   // ..
	iv        []byte       // ..
	done      chan string  // ..
	info      *ConnectInfo // ..
}

// NewUDPSocket ...
//
func NewUDPSocket(ss net.Conn, c *AppConfig, i []byte, al *ConnectInfo, ch chan string) *UDPSocket {
	return &UDPSocket{
		srcsocket: ss,
		config:    c,
		iv:        i,
		info:      al,
		done:      ch,
	}
}

func sos(ssocket *UDPSocket) ([]byte, []byte, error) {

	buf, err := read_bytes_from_socket(ssocket.srcsocket, 8)
	if err != nil {
		utils.Logger.Print("read length from C failed!", err)
		return nil, nil, err
	}
	utils.Logger.Println("read buffer size", len(buf))
	if len(buf) != 8 {
		utils.Logger.Print("compress format is incorrect!")
		return nil, nil, errors.New("compress format is incorrectly")
	}

	size := utils.Byte2int(buf[4:])
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

	ori, err := sercurity.DecompressWithChaCha20(content, ssocket.iv[:8], sercurity.MakeCompressKey(ssocket.config.Password))
	if err != nil {
		utils.Logger.Print("Decompress the UDP from C failed ", err)
		return nil, nil, err
	}

	return buf[:4], ori, nil
}

//DUDP2RAW just client use it
//
func (ssocket *UDPSocket) DUDP2RAW(raw *net.UDPConn, udpaddr *net.UDPAddr) error {

	header, ori, err := sos(ssocket)
	if err != nil {
		ssocket.done <- "done"
		return errors.New("read packet failed")
	}
	addressLen := caladdress(header[3], ori[0])
	if checkUDPTerminates(ori[:addressLen]) {
		ssocket.done <- "done"
		utils.Logger.Print("check ip address is finish marker.")
		return errors.New("done")
	}

	// client need to keep UDP header format.
	full := append(header, ori...)
	utils.Logger.Print("UDP response: ", full)

	n, err := raw.WriteToUDP(full, udpaddr)

	utils.Logger.Print("output: ", n, " bytes. ", err)
	if err != nil {
		ssocket.done <- "done"
		return errors.New("UDP write ERROR")
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
func (ssocket *UDPSocket) Raw2UDP(raw *net.UDPConn) (*net.UDPAddr, error) {

	buf := make([]byte, BF_SIZE)
	n, udpaddr, err := raw.ReadFromUDP(buf)
	if err != nil {
		utils.Logger.Print("read UDP packet from raw socket failed", err, "bytes: ", n)
		ssocket.done <- "done"
		return nil, err
	}

	utils.Logger.Print("bowser udp: ", buf[:n], err, "read bytes: ", n)
	if n < 4+2 {
		utils.Logger.Print("UDP format incorrect from raw socket!!!")
		ssocket.done <- "done"
		return nil, errors.New("UDP format incorectly")
	}

	if checkUDPTerminates(buf[4 : 4+caladdress(buf[3], buf[4])]) {
		ssocket.done <- "done"
		return nil, errors.New("client done")
	}

	out, err := sercurity.CompressWithChaCha20(buf[4:n], ssocket.iv[:8], sercurity.MakeCompressKey(ssocket.config.Password))
	if err != nil {
		utils.Logger.Print("compress UDP failed! ", err)
		ssocket.done <- "done"
		return nil, err
	}

	start := utils.Int2byte((uint32)(len(out)))
	start = append(buf[:4], start...)
	full := append(start, out...)
	on, err := ssocket.srcsocket.Write(full)

	utils.Logger.Print("write compress UDP packet result: ", err, "write bytes: ", on, "bytes.", full)

	if err != nil {
		ssocket.done <- "done"
		return nil, err
	}

	return udpaddr, err

}

func createUDPListen(ip string, port uint16) (*net.UDPConn, error) {
	address := ip + ":" + strconv.Itoa((int)(port))
	utils.Logger.Print("UDP listen on :", address)
	addr, err := net.ResolveUDPAddr("udp", address)

	udpConn, err := net.ListenUDP("udp", addr)

	return udpConn, err
}

func udpsend(ss *UDPSocket) error {

	defer func() {
		ss.info.udpConnect.Close()
		GPortQueue.Enqueue((Item)(ss.info.udpport))
		utils.Logger.Print("will relese port: ", ss.info.udpport, "rest port size: ", GPortQueue.Size())
	}()

	udpaddress, err := ss.Raw2UDP(ss.info.udpConnect)
	if err != nil {
		utils.Logger.Print("write to remote proxy chanel failed!", err)
		return err
	}

	err = ss.DUDP2RAW(ss.info.udpConnect, udpaddress)
	if err != nil {
		utils.Logger.Print("write to UDP chanel failed!", err)
		return err
	}
	utils.Logger.Print("tcp will shutdown for UDP.")
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

func (ssocket *UDPSocket) readUDPD() error {

	header, ori, err := sos(ssocket)
	if err != nil {

		return errors.New("read packet failed")
	}
	addrLen := caladdress(header[3], ori[0])
	if checkUDPTerminates(ori[:addrLen]) {
		utils.Logger.Print("check finish flag in ip address!")
		return errors.New("done")
	}

	// must write raw data without header.
	host := parsehost(header[3], ori[:addrLen])
	udpConnect, err := createDUPConnect(host)
	if err != nil {
		return errors.New("address invalid for UDP")
	}

	defer func() {
		udpConnect.Close()

	}()

	utils.Logger.Print("send udp data:", ori[addrLen:])

	n, err := udpConnect.Write(ori[addrLen:])
	if err != nil {
		utils.Logger.Print("write UDP failed!", err)
		return err
	}

	buf := make([]byte, BF_SIZE)
	utils.Logger.Print("----------------------------------------")
	n, err = udpConnect.Read(buf)
	utils.Logger.Print("----------------------------------------")
	if err != nil {
		ssocket.done <- "done"
		utils.Logger.Print("read udp host failed: ", err, "bytes: ", n)
		return err
	}
	utils.Logger.Print("real udp Response: ", buf[:n])

	// encapsulate the data with header and compress
	lop := append(ori[:addrLen], buf[:n]...)
	udppacket, err := sercurity.CompressWithChaCha20(lop, ssocket.iv[:8], sercurity.MakeCompressKey(ssocket.config.Password))
	if err != nil {
		ssocket.done <- "done"
		utils.Logger.Print("compress UDP failed", err, "bytes: ")
		return err
	}

	ll := len(udppacket)
	full := append(header[:4], utils.Int2byte((uint32)(ll))...)
	out := append(full, udppacket...)
	n, err = ssocket.srcsocket.Write(out)
	if err != nil {
		ssocket.done <- "done"
		utils.Logger.Print("remote write UDP failed!", err)
	}
	utils.Logger.Print("write ", n, "bytes.", err)

	utils.Logger.Print("write to destination UDP", n, " bytes.")

	return nil
}

func (ssocket *UDPSocket) rawToUPD(raw *net.Conn, header []byte) error {

	defer func() {
		if raw != nil {
			(*raw).Close()
		}
	}()
	utils.Logger.Print("start read form UDP server: ")

	buf := make([]byte, BF_SIZE)
	utils.Logger.Print("----------------------------------------")
	n, err := (*raw).Read(buf)
	utils.Logger.Print("----------------------------------------")
	if err != nil {
		ssocket.done <- "done"
		utils.Logger.Print("read udp host failed: ", err, "bytes: ", n)
		return err
	}
	utils.Logger.Print("real udp Response: ", buf[:n])

	// encapsulate the data with header and compress
	lop := append(header[4:], buf[:n]...)
	ori, err := sercurity.CompressWithChaCha20(lop, ssocket.iv[:8], sercurity.MakeCompressKey(ssocket.config.Password))
	if err != nil {
		ssocket.done <- "done"
		utils.Logger.Print("compress UDP failed", err, "bytes: ")
		return err
	}

	ll := len(ori)
	full := append(header[:4], utils.Int2byte((uint32)(ll))...)
	out := append(full, ori...)
	n, err = ssocket.srcsocket.Write(out)
	if err != nil {
		ssocket.done <- "done"
		utils.Logger.Print("remote write UDP failed!", err)
	}
	utils.Logger.Print("write ", n, "bytes.", err)

	return err
}

func udpserverRoutine(ss *UDPSocket, h net.Conn) error {

	err := ss.readUDPD()

	utils.Logger.Print("server dup done!!!!")

	return err
}
