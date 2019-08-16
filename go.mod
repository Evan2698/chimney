module github.com/Evan2698/chimney

go 1.12

require (
	github.com/Evan2698/android-netstack v0.0.0-20190101110143-3a7a9258406d
	github.com/lucas-clemente/quic-go v0.11.2
	github.com/miekg/dns v1.1.14 // indirect
)

replace golang.org/x/net => github.com/golang/net v0.0.0-20190813141303-74dc4d7220e7

replace golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190701094942-4def268fd1a4

replace golang.org/x/sys => github.com/golang/sys v0.0.0-20190813064441-fde4db37ae7a

replace golang.org/x/text => github.com/golang/text v0.3.2

replace golang.org/x/sync => github.com/golang/sync v0.0.0-20190423024810-112230192c58
