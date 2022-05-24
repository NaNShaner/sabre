package scheduled

import (
	"fmt"
	"os"
	"os/signal"
	"sabre/pkg/dbload"
	"sabre/pkg/sabstruct/res"
	"sabre/pkg/util/commontools"
	"syscall"
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
	rs := res.Register()
	for _, regx := range rs.ResRegx {
		go dbload.WatchFromDB(regx)
		fmt.Printf("Watch etcd key的名称为%s\n", regx)
	}
	// 主goroutine堵塞
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
}
