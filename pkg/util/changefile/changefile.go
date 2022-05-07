package changefile

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// ServerXml Tomcat 的server.xml
type ServerXml struct {
	ShutdownPort    int `json:"SHUTDOWN_PORT"`
	ListeningPort   int `json:"LISTENING_PORT"`
	AjpPort         int `json:"AJP_PORT"`
	AjpRedirectPort int `json:"AJP_Redirect_PORT"`
}

// CatalinaSh Tomcat 的Catalina.sh
type CatalinaSh struct {
	JavaOpts string `json:"SET_JAVA_OPTS"`
}

type TomcatConfigFile struct {
	ServerXml
	CatalinaSh
}

type NginxConf struct {
}

type NginxConfigFile struct {
	NginxConf
}

type StartUp struct {
}

type JDKConfigFile struct {
	StartUp
}

// Changefile 修改配置文件
// TODO: 替换关键字内容
func Changefile(f string, r ...map[string]string) error {
	var filePerLine []string
	in, err := os.Open(f)
	if err != nil {
		return fmt.Errorf("打开%s文件失败: %s", f, err)
	}
	defer in.Close()
	out, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return fmt.Errorf("打开%s文件失败: %s", f, err)
	}
	defer out.Close()

	br := bufio.NewReader(in)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("读取%s文件失败: %s", f, err)
		}
		for _, m := range r {
			// k 是在配置文件中已经标记可以替换的文件
			// v 是被替换的内容
			for k, v := range m {
				if strings.Contains(string(line), k) {
					newLine := strings.Replace(string(line), k, v, -1)
					filePerLine = append(filePerLine, newLine)
				} else {
					filePerLine = append(filePerLine, string(line))
				}
			}
		}
	}

	for _, s := range filePerLine {
		fmt.Println(s + "\n")
		_, err = out.WriteString(s + "\n")
		if err != nil {
			return fmt.Errorf("写%s文件失败: %s", f, err)
		}
	}
	return nil

}
