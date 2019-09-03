module chimney

go 1.12

require (
	github.com/Evan2698/android-netstack v0.0.0-20190628092857-bf7224ae5dde
	github.com/lucas-clemente/quic-go v0.12.0
)

replace golang.org/x/net => github.com/golang/net v0.0.0-20190827160401-ba9fcec4b297

replace golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190829043050-9756ffdc2472

replace golang.org/x/sys => github.com/golang/sys v0.0.0-20190902133755-9109b7679e13

replace golang.org/x/text => github.com/golang/text v0.3.2

replace golang.org/x/sync => github.com/golang/sync v0.0.0-20190423024810-112230192c58
