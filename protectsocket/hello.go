package protectsocket

// #cgo CFLAGS:
// #include "ancillary.h"
// #include "ssgo.h"
import "C"
import "errors"
import "strconv"

func ProtectFD(fd int) error {

	var r C.int = C.send_fd(C.int(fd))

	if r != 0 {
		return errors.New("can not send fd." + strconv.Itoa(int(r)))
	}

	return nil
}
