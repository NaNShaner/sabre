package res

import "fmt"

//SabRes 平台纳管的资源类型
type SabRes struct {
	//ResRegx 纳管的资源类型的前缀， e.g. mid
	ResRegx []string
	//ResApply 资源提供方，e.g. k8s、PaaS平台
	//ResApply string
}

//SetResRegx 格式化资源类型，方便后续写入etcd
func SetResRegx(regx []string) SabRes {
	var r SabRes
	for _, rType := range regx {
		r.ResRegx = append(r.ResRegx, "/"+rType)
	}
	return r
}

//Register 资源注册
func Register() SabRes {
	res := []string{"mid", "net", "hosts"}
	rs := SetResRegx(res)
	fmt.Printf("纳管资源注册信息列表%q\n", rs)
	return rs
}
