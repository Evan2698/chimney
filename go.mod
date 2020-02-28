module github.com/Evan2698/chimney

go 1.13

require (
	github.com/eycorsican/go-tun2socks v1.16.8
	github.com/lucas-clemente/quic-go v0.14.4
	github.com/miekg/dns v1.1.27
	golang.org/x/crypto v0.0.0-20200221231518-2aa609cf4a9d
	golang.org/x/net v0.0.0-20191021144547-ec77196f6094
	golang.org/x/text v0.3.0

)

replace golang.org/x/crypto v0.0.0-20200221231518-2aa609cf4a9d => github.com/golang/crypto v0.0.0-20200221231518-2aa609cf4a9d
