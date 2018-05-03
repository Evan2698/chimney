//  +build DRCLO

package core

import (
	"bytes"
	"github.com/Evan2698/climbwall/utils"
	"encoding/binary"
	"math/rand"
	"net"
	"time"
)

var last time.Time

func init() {
	t := int64(time.Now().Nanosecond())
	rand.Seed(t)
}

func StatPackage(s uint64, r uint64) {

	t := time.Now()

	elapsed := t.Sub(last)

	if elapsed.Seconds() > 1.0 {

		r := rand.Int()
		v := uint64(r)
		if v < 0 {
			v = -v
		}

		v = v % 100

		go sendimp(v, v)
		last = t
	}
}

func sendimp(s uint64, r uint64) {
	conn, err := net.Dial("unix", "stat_path")
	if err != nil {
		utils.Logger.Println("create statistics socket failed！", err)
		return
	}

	defer conn.Close()

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, s)
	bufh := new(bytes.Buffer)
	binary.Write(bufh, binary.LittleEndian, r)
	out := append(buf.Bytes(), bufh.Bytes()...)

	conn.Write(out)
	utils.Logger.Println("statistics send ok!!")
	n:=conn.Read(out)
	utils.Logger.Println("statistics result:", out[0:n])
}
