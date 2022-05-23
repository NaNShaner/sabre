package sabreschedule

import (
	"sabre/pkg/util/commontools"
)

type Bt commontools.Basest

type SabreSchedule interface {
	// Watch ETCD中资源变化
	Watch(*Bt) error
	// Cron 资源变更后进行资源调度
	Cron(*Bt) error
	// Add 调度资源创建
	deploy(*Bt) error
	// Del 调度资源删除
	Del(*Bt) error
	// Apply 调度资源修改
	Apply(*Bt) error
}

//Schedule 调度器
type Schedule struct {
	Server    string
	ApiVesion string
}

//Watch 接收 sabrectl show 指令，从etcd中获取数据并反馈
//func (s Schedule) Watch(wr http.ResponseWriter, r *http.Request) error {
//
//}
