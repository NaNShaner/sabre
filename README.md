# Sabre (佩剑)
以IaC的方式，管理你想管理的资源。

## 项目简介
### 主要功能（当前功能）
- 基于yaml文件完成Tomcat、JDK的资源以及启动脚本的部署

### 组件说明
控制面 Master 节点主要包含以下组件：
- kube-apiserver，负责对外提供集群各类资源的增删改查及 Watch 接口，它是 Kubernetes 集群中各组件数据交互和通信的枢纽。kube-apiserver 在设计上可水平扩展，高可用 Kubernetes 集群中一般多副本部署。当收到一个创建 Pod 写请求时，它的基本流程是对请求进行认证、限速、授权、准入机制等检查后，写入到 etcd 即可。
- kube-scheduler 是调度器组件，负责集群 Pod 的调度。基本原理是通过监听 kube-apiserver 获取待调度的 Pod，然后基于一系列筛选和评优算法，为 Pod 分配最佳的 Node 节点。
- kube-controller-manager 包含一系列的控制器组件，比如 Deployment、StatefulSet 等控制器。控制器的核心思想是监听、比较资源实际状态与期望状态是否一致，若不一致则进行协调工作使其最终一致。
- etcd 组件，Kubernetes 的元数据存储。
Node 节点主要包含以下组件：
- kubelet，部署在每个节点上的 Agent 的组件，负责 Pod 的创建运行。基本原理是通过监听 APIServer 获取分配到其节点上的 Pod，然后根据 Pod 的规格详情，调用运行时组件创建 pause 和业务容器等。
- kube-proxy，部署在每个节点上的网络代理组件。基本原理是通过监听 APIServer 获取 Service、Endpoint 等资源，基于 Iptables、IPVS 等技术实现数据包转发等功能。
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