module github.com/Evan2698/chimney

go 1.12

require (
	github.com/eycorsican/go-tun2socks v1.16.2
	github.com/lucas-clemente/quic-go v0.12.0
	github.com/miekg/dns v1.1.15
	golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3
	golang.org/x/text v0.3.0
)

replace golang.org/x/net => github.com/golang/net v0.0.0-20190827160401-ba9fcec4b297

replace golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190829043050-9756ffdc2472

replace golang.org/x/text => github.com/golang/text v0.3.2

replace golang.org/x/sync => github.com/golang/sync v0.0.0-20190423024810-112230192c58

replace golang.org/x/sys => github.com/golang/sys v0.0.0-20190813064441-fde4db37ae7a
