package hostregister

import (
	"fmt"
	"github.com/tatsushid/go-fastping"
	"net"
	"sabre/pkg/sabstruct/res"
	"time"
)

type Hosts res.Hosts

func (h *Hosts) ServerRegister() (string, error) {
	return "", nil
}

func (h *Hosts) IpFmt(ip string) (net.IP, error) {
	i := net.ParseIP(ip)
	if i == nil {
		return nil, fmt.Errorf("ip 地址%s格式错误\n", ip)
	}
	return i, nil
}

func (h *Hosts) Ping(ip net.IP) (bool, error) {
	p := fastping.NewPinger()
	fmt.Println(ip.String())
	ra, err := net.ResolveIPAddr("ip4:icmp", ip.String())
	if err != nil {
		return false, err
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {}
	err = p.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}
