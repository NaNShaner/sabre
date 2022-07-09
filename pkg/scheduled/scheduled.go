package scheduled

import (
	"context"
	"fmt"
	"github.com/sevlyar/go-daemon"
	clientv3 "go.etcd.io/etcd/client/v3"
	"os"
	"os/signal"
	"sabre/pkg/dbload"
	"sabre/pkg/sabstruct"
	"sabre/pkg/sabstruct/res"
	"sabre/pkg/util/callsabrelet"
	"sabre/pkg/util/commontools"
	L "sabre/pkg/util/logbase/logscheduled"
	"sabre/pkg/yamlfmt"
	"strings"
	"syscall"
)

var (
	midRegx   = "mid"
	hostRegx  = "hosts"
	midTomcat = "tomcat"
	// midJDK    = "jdk"
)

type SabreSchedule interface {
	// Cron 资源变更后进行资源调度
	Cron(m *commontools.Basest) (string, error)
	// Deploy 调度资源创建
}

//Schedule 调度器
type Schedule struct {
	Server    string
	ApiVesion string
}

func Cron(m *commontools.Basest) (string, error) {
	return "", nil
}

//Watch 监控ETCD中资源状态变化，发起调度逻辑
func Watch() {
	cntxt := &daemon.Context{
		PidFileName: "/var/run/sabreschedule.pid",
		PidFilePerm: 0644,
		LogFileName: "/var/log/sabreschedule.log",
		LogFilePerm: 0640,
		Umask:       027,
		Args:        []string{"[sabreschedule]"},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		L.Log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	rs := res.Register()
	for _, regx := range rs.ResRegx {
		go WatchFromDB(regx)
		L.Log.Infof("Watch etcd key的名称为%s\n", regx)
	}
	// 主goroutine堵塞
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
}

//WatchFromDB 通过API网关对etcd中的资源类型进行watch，进行后续调度
func WatchFromDB(s string) {
	L.Logfile = "dbload.log"
	cli, err := dbload.GetDBCli()
	if err != nil {
		L.Log.Infof("connect failed, %s\n", err)
		return
	}
	defer cli.Close()

	for {
		// clientv3.WithPrefix() 监控s作为前缀key值的value变化，默认为精确watch
		rch := cli.Watch(context.Background(), s, clientv3.WithPrefix())
		for wresp := range rch {
			err = wresp.Err()
			if err != nil {
				L.Log.Infof("etcd watch response err, %s\n", err)
			}
			for _, ev := range wresp.Events {
				L.Log.Info("%s %q %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				keySplit, err := keySplit(ev.Kv.Key)
				if err != nil {
					return
				}
				switch {
				// 判断etcd中的变化的key的类型，枚举值
				case isMidType(keySplit):
					L.Log.Infof("isMidType: %s\n", ev.Kv.Key)
					switch {
					case isMidTypeOfTomcat(keySplit):
						DoActionOfMidType(string(ev.Kv.Value))
					}
				case isMidHost(keySplit):
					L.Log.Infof("isMidHost: %s\n", ev.Kv.Key)
				default:
					return
				}

			}
		}
	}
}

//DoActionOfMidType 针对不同中间件类型执行安装动作
func DoActionOfMidType(s string) {
	yamlFmt, err := yamlfmt.YamlFmt(s, sabstruct.Config{})
	if err != nil {
		return
	}
	b := (*callsabrelet.Basest)(yamlFmt)
	callsabrelet.CallFaceOfSabrelet(b, b.DeployHost)
}

func keySplit(t []byte) (string, error) {
	s := string(t)
	sSplit := strings.Split(s, "/")
	if len(sSplit) < 1 {
		return "", fmt.Errorf("Etcd key %s is not in normal format\n", s)
	}
	return sSplit[1], nil
}

func isMidType(t string) bool {
	if strings.Contains(t, midRegx) {
		return true
	} else {
		return false
	}
}

func isMidTypeOfTomcat(t string) bool {

	if strings.Contains(t, midTomcat) {
		return true
	} else {
		return false
	}
}

func isMidHost(t string) bool {
	if strings.Contains(t, hostRegx) {
		return true
	} else {
		return false
	}
}
