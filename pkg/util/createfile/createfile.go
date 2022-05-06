package createfile

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

// createFile 创建文件
func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}

// YamlToFile 从结构体生成yaml文件
func YamlToFile(s interface{}) error {
	yamlData, err := yaml.Marshal(&s)
	if err != nil {
		return fmt.Errorf("can not create file from yaml, err:%s", err)
	}
	fmt.Println(string(yamlData))
	fileName := "test.yaml"
	err = ioutil.WriteFile(fileName, yamlData, 0644)
	if err != nil {
		return fmt.Errorf("unable to write data into the file:%s", err)
	}
	return nil
}
