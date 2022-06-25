package scheduled

import (
	"github.com/sevlyar/go-daemon"
	"sabre/pkg/dbload"
	"sabre/pkg/sabstruct/res"
	"sabre/pkg/util/commontools"
	l "sabre/pkg/util/logbase/logscheduled"
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
		l.Log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	rs := res.Register()
	for _, regx := range rs.ResRegx {
		go dbload.WatchFromDB(regx)
		l.Log.Infof("Watch etcd key的名称为%s\n", regx)
	}
}
