module github.com/adlotsof/filetun

go 1.22.4

require (
	github.com/alecthomas/kong v0.9.0
	github.com/google/gopacket v1.1.19
	github.com/songgao/water v0.0.0-20200317203138-2b4b6d7c09d8
	github.com/vishvananda/netlink v1.1.0
)

require (
	github.com/vishvananda/netns v0.0.0-20191106174202-0a2b9b5464df // indirect
	golang.org/x/sys v0.0.0-20190606203320-7fc4e5ec1444 // indirect
)

replace github.com/adlotsof/filetun => ./

replace github.com/adlotsof/filetun/config => ./config
