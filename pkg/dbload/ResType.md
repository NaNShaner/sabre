# 资源类型，定义etcd中的前缀
格式：prefix + "/" + namespace + "/" + 资源类型 `[+ 资源提供方式]` `[+ 工程名称]`

## 资源类型及相关说明
### prefix 资源前缀
| 前缀类型| 简称
--------|------
|系统|default|
|中间件|mid|
|网络资源|net|
|存储资源|storage|
|计算资源|machine|

### namespace，表示当前资源所属的系统简称
例如：ERP

### 资源名称，表示当前资源的种类名称
例如：
- 中间件资源：Tomcat、NGINX
- 缓存资源：Redis
- 数据资源：Mysql、PG
- 网络资源：F5、DNS
- 存储资源：Nas、oss

### 资源提供方式，可选
例如：
- 基于sabre提供的云下资源
- 公司内PaaS平台通过接口的方式给提供
  - 中间件管理平台，提供资源的部署能力
  - 软负载平台，提供的负载能力
- 基于K8s提供资源

### 工程名称，可选
例如：demo



## 事例
### 例如：当前机器的namespace
key :/midRegx/namespace/{hostname/ipaddr}
sabrectl get namespace
输出：
MNPP


### 例如：获取erp系统的下的demo工程部署在哪些机器的Tomcat中
key :/mid/ERP/Tomcat/{projectName}/{hostname/ipaddr}
```shell
sabrectl get tomcat demo
```
输出：
```shell
namespace	host		midType	projectName	version
MNPP		127.0.0.1 	Tomcat 	demo		7.0.78
MNPP		127.0.0.2 	Tomcat 	demo		7.0.78
```
如需获取详细信息可以通过-d获取，应用当前的运行、监控状态以及实例运行的时间
```shell
sabrectl get tomcat demo -d
sabrectl get tomcat demo
```
输出：
```shell
namespace	host		midType	projectName	port	version	monitor running	runningTime
MNPP		127.0.0.1 	Tomcat 	demo		8099	7.0.78 	True	True	10d
MNPP		127.0.0.2 	Tomcat 	demo		8099	7.0.78 	True	True	10d
```