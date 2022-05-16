package sabreschedule

import (
	"net/http"
)

type SabreSchedule interface {
	// Watch ETCD中资源变化
	Watch()
	// Cron 资源变更后进行资源调度
	Cron()
	// Add 调度资源创建
	Add()
	// Del 调度资源删除
	Del()
	// Apply 调度资源修改
	Apply()
}

//Schedule 调度器
type Schedule struct {
	Server    string
	ApiVesion string
}

//Watch 接收 sabrectl show 指令，从etcd中获取数据并反馈
func (s Schedule) Watch(wr http.ResponseWriter, r *http.Request) {

}
