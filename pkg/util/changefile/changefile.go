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
func Changefile(f string, r ...map[string]string) {
	in, err := os.Open(f)
	if err != nil {
		fmt.Println("open file fail:", err)
		os.Exit(-1)
	}
	defer in.Close()
	newF := f + "now"
	out, err := os.OpenFile(newF, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		fmt.Println("Open write file fail:", err)
		os.Exit(-1)
	}
	defer out.Close()

	br := bufio.NewReader(in)
	index := 1
	for _, m := range r {
		// k 是在配置文件中已经标记可以替换的文件
		// v 是被替换的内容
		for k, v := range m {

			for {
				line, _, err := br.ReadLine()
				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Println("read err:", err)
					os.Exit(-1)
				}
				newLine := strings.Replace(string(line), k, v, -1)
				_, err = out.WriteString(newLine + "\n")
				if err != nil {
					fmt.Println("write to file fail:", err)
					os.Exit(-1)
				}
				fmt.Println("done ", index)
				index++
			}
		}

	}
}
