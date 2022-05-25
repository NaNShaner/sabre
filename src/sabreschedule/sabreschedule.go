package main

import "sabre/pkg/scheduled"

//import _ "sabre/pkg/sabstruct/res"

//Watch 接收 sabrectl show 指令，从etcd中获取数据并反馈
//func (s Schedule) Watch(wr http.ResponseWriter, r *http.Request) error {
//
//}

func main() {
	scheduled.Watch()
}
