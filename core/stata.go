//  +build DRCLO

package core

import (
	"math/rand"
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
	if GFlow != nil {
		GFlow.Update((int64)(s), (int64)(r))
	}
}
