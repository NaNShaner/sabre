/*
@spec: 用于解析yaml文件
*/

package yamlfmt

import (
	"awesomeProject/pkg/sabstruct"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

// Config 创建一个结构体
//apiversion: apps/v1
//kind: Deployment
//metadata:
//	namespace: MNPP
//	netarea: app
//	appname: entry
//spec:
//	name: tomcat
//	version: 7.0.78
//	port: 8099
//	installpath: /u01
//	user:
//		name: miduser
//		group: miduser

// YamlFmt 解析yaml文件为json。f为yaml文件的绝对路径，s为解析的结构体
func YamlFmt(f string, s sabstruct.Config) (*sabstruct.Config, error) {
	//从外部的conf.yaml文件读取数据
	// TODO: 改为命令行进行输入
	data, _ := ioutil.ReadFile(f)
	//使用yaml包，把读取到的data格式化后解析到config实例中
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		panic("decode error")
	}

	//resultJson, err := PrintResultJson(s)
	//if err != nil {
	//	return nil, err
	//}
	return &s, err
}

// PrintResultJson 解析结果，并输出json
func PrintResultJson(s interface{}) ([]byte, error) {
	// 字典格式化为json
	//data, err := json.Marshal(s)
	//if err != nil {
	//	fmt.Printf("JSON marshaling failed: %s", err)
	//	return nil
	//}

	// 针对json增加人类的可读性
	data, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	return data, err
}
