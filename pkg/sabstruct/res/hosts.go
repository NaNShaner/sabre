package res

import (
	"bytes"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/tatsushid/go-fastping"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type MonitorStatus struct {
	//Online 根据ping确认主机在线状态
	Online     bool
	Mem        map[uint64]uint64
	OnlineDay  int
	OnlineTime string
}

//HostInfo 主机信息
type HostInfo struct {
	Platform string
	OS       string
	Core     string
	CPUs     int
}

//Hosts 注册主机信息
type Hosts struct {
	//HostName 主机名
	HostName string
	//IPAddr 主机IP地址
	IPAddr net.IP
	//BelongTo 主机所属业务系统简称
	BelongTo string
	//Area 主机所属网络区域
	Area string
	//Annotation 其他声明，例如冷备节点
	Annotation []string
	MonitorStatus
	HostInfo
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

type HostRegister interface {
	ServerRegister(ip, beloogto, area string) (Hosts, error)
	Ping(ip net.IP) (bool, error)
	IpFmt(ip string) (net.IP, error)
	GetMem() (map[uint64]uint64, error)
	GetOsInfo() (HostInfo, error)
	ReturnHost() Hosts
	AddAnnotation(a ...string) []string
}

//ServerRegister 主机信息注册入库
func (h *Hosts) ServerRegister(ip, belongto, area string) (Hosts, error) {
	h.HostInfo, _ = h.GetOsInfo()
	hostName, getHostNameErr := os.Hostname()
	if getHostNameErr != nil {
		h.HostName = "unknown hostname"
	}
	h.HostName = hostName

	ipAddr, IpFmtErr := h.IpFmt(ip)
	if IpFmtErr != nil {
		return Hosts{}, IpFmtErr
	}
	h.IPAddr = ipAddr
	h.BelongTo = belongto
	h.Area = area

	h.Online = true
	memInfo, getMemErr := h.GetMem()
	if getMemErr != nil {
		return Hosts{}, getMemErr
	}
	h.Mem = memInfo
	return *h, nil
}

func (h *Hosts) ReturnHost() Hosts {
	return *h
}

func (h *Hosts) AddAnnotation(a ...string) []string {
	for _, s := range a {
		h.Annotation = append(h.Annotation, s)
	}
	return h.Annotation
}

//IpFmt 判断及格式化IP地址，返回net.IP
func (h *Hosts) IpFmt(ip string) (net.IP, error) {
	i := net.ParseIP(ip)
	if i == nil {
		return nil, fmt.Errorf("IP address %s format error\n", ip)
	}
	return i, nil
}

func (h *Hosts) Ping(ip net.IP) (bool, error) {
	p := fastping.NewPinger()
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
	//TODO 方便调试，后续删除
	fmt.Printf("主机地址%s ping正常\n", ip.String())
	return true, nil
}

func (h *Hosts) GetMem() (map[uint64]uint64, error) {
	ramCheck, err := RAMCheck()
	if err != nil {
		return nil, err
	}
	return ramCheck, nil
}

//GetOsInfo 获取当前服务器Os的相关信息
func (h *Hosts) GetOsInfo() (HostInfo, error) {
	cmd := exec.Command("uname", "-srm")
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	osStr := strings.Replace(out.String(), "\n", "", -1)
	osStr = strings.Replace(osStr, "\r\n", "", -1)
	osInfo := strings.Split(osStr, " ")

	h.CPUs = runtime.NumCPU()

	if len(osInfo) != 3 {
		h.OS = ""
		h.Core = ""
		h.Platform = ""
		return h.HostInfo, err
	}
	h.OS = osInfo[0]
	h.Core = osInfo[1]
	h.Platform = osInfo[2]

	return h.HostInfo, err
}

//DiskCheck 服务器硬盘使用量
func DiskCheck() {
	u, _ := disk.Usage("/")
	usedMB := int(u.Used) / MB
	usedGB := int(u.Used) / GB
	totalMB := int(u.Total) / MB
	totalGB := int(u.Total) / GB
	usedPercent := int(u.UsedPercent)
	fmt.Printf("Free space: %dMB (%dGB) / %dMB (%dGB) | Used: %d%%\n", usedMB, usedGB, totalMB, totalGB, usedPercent)
}

//CPUCheck CPU核数
func CPUCheck() (int, error) {
	cores, err := cpu.Counts(false)
	if err != nil {
		return 0, err
	}

	//cpus, err := cpu.Percent(time.Duration(200)*time.Millisecond, true)
	//if err == nil {
	//	for i, c := range cpus {
	//		fmt.Printf("cpu%d : %f%%\n", i, c)
	//	}
	//}
	// CPU的平均负载
	//a, _ := load.Avg()
	//l1 := a.Load1
	//l5 := a.Load5
	//l15 := a.Load15

	return cores, nil
}

//RAMCheck 内存使用量
func RAMCheck() (map[uint64]uint64, error) {
	u, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	usedMB := u.Used / MB
	totalMB := u.Total / MB
	// usedPercent := int(u.UsedPercent)
	m := make(map[uint64]uint64)
	m[usedMB] = totalMB
	return m, nil
}
