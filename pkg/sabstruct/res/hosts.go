package res

import "net"

type HostStatus struct {
	Status bool
}

type MonitorStaus struct {
	Ping bool
	Mem  map[float64]float64
}

type Hosts struct {
	HostName     string
	IPAddr       net.IP
	Online       bool
	OnlineStatus []string
	BelongTo     string
	Area         string
	annotation   string
	HostStatus
	MonitorStaus
}

type HostRegister interface {
	ServerRegister() (string, error)
	Ping(ip net.IP) (bool, error)
	IpFmt(ip string) (net.IP, error)
}
