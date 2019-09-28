# Chimney

Chimney is a component for vpn. on PCs(windows & linux), server & client communicate with socks5 protocol,It is a separate configurable executable program. It can communicate through configuration and quic protocol, which is fast and simple to configure.

On Android, it is just a component(AAR package) for vpn.


# Summary

- [HowToBuild](#HowToBuild "HowToBuild")
- [AboutUs](#AboutUs "AboutUs")
- [License](#License "License")

# HowToBuild
 
   ### Build PC's executable
-   go get github.com/Evan2698/chimney/cmd/client **or**  go get github.com/Evan2698/chimney/cmd/server
-   go build github.com/Evan2698/chimney/cmd/client
-   go build github.com/Evan2698/chimney/cmd/server
-  modify the config.json and put it in the same directory with executable.

  ### Build android component
  build command:
  - prepare gomobile and quic.
  - gomobile bind -target=android  -ldflags="-s -w" github.com/Evan2698/chimney/android 

  BTW:  If quic fails to compile, you can modify it to compile. Quic protocol is not supported on Android.

# About
 ☺ ☺ ☺ ☺ ☺ 

# License
```
And of course:

MIT: https://opensource.org/licenses/MIT