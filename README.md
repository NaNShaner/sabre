# Sabre (佩剑)
以K8s的资源管理模式，全面纳管非云上的计算、存储以及网络资源。

## 项目简介
### 主要功能（当前功能）
- 基于yaml文件完成Tomcat、JDK的资源以及启动脚本的部署

### 组件说明
#### 控制节点
- sabreapi

平台资源操作的唯一入口，接受用户输入的命令
- sabreschedule 

负责集群资源的调度，按照预定的调度策略将资源调度到相应的计算节点上
#### 工作节点
- sabrectl

命令行程序
- sabrelet

部署在计算节点上，负责维护集群的状态，比如程序部署安排，故障检测，信息上送api
### 执行流程
- 注册主机
```shell
$ sabrectl hosted 192.168.3.58 -a APP -n ERP
Server /xxx/192.168.3.58 registration succeeded
```

- 部署Tomcat
```shell
$sabrectl create tomcat /opt/sabre/pkg/util/tomcat/install/deployTomcat.yaml
```
deployTomcat.yaml 文件如下
```yaml
apiversion: beta
kind: Config
metadata:
    namespace: ERP
    netarea: APP
    appname: erp
spec:
    midtype: Tomcat
    version: 7.0.75
    installpath: /u01/app
    pkgdownloadpath: http://124.71.219.53:8001/uploads/uploads/2022/05/07/apache-tomcat-7.0.75.tar.gz
    midruntype:
        - cluster
    user:
        name: miduser
        group: miduser
    defaultconfig:
        tomcat:
            javaopts: -server -Xms1024M -Xmx1024M -Xss512k
            listeningport: "8099"
            ajpport: "8009"
            shutdownport: "8005"
    deployaction:
        action: Install
        deployhost:
            - 192.168.3.182
            - 192.168.3.58

```
## Quick start
### 环境依赖
- 操作系统
```
CentOS Linux release 8.2
```
- go版本
```shell
go version go1.17.3 
```
- 安装步骤
# 拓扑架构
![img.png](docs/imgs/img.png)

# Todo list
- [x] 中间单资源部署，已完成Tomcat、Jdk
- [ ] 完成本机器资源获取
- [ ] 中间单资源部署，完成Nginx、Redis
- [ ] 缓存资源部署，完成Redis集群部署
- [ ] 完成Etcd存储交付资源
- [ ] 完成资源集群模式的运维操作，包括启动、停止、重启
- [ ] 完成监控接口
- [ ] 完成CMDB接口

# 致谢
开发工具由[Jetbrains](https://www.jetbrains.com/)赞助的Pycharm