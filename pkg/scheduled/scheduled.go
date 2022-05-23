package scheduled

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
