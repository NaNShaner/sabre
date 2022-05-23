package scheduled

import "sabre/pkg/util/commontools"

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
