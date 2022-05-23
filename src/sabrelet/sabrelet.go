package sabrelet

import (
	"net"
)

type NetInfo struct {
}

type ServerInfo struct {
	NameSpace string
	AreaZone  string
	HostName  string
	IpAddr    net.IPAddr
	//MemInfo
	//CpuInfo
	NetInfo
}

type ListenAndAction interface {
}
