package getSomethingToPrint

//### 例如：获取erp系统的下的demo工程部署在哪些机器的Tomcat中
//key :/mid/ERP/Tomcat/{projectName}/{hostname/ipaddr}
//sabrectl get mid -t tomcat -a app -n erp

//输出：
//namespace	host		midType	projectName	port	version	monitor running	runningTime
//MNPP		127.0.0.1 	Tomcat 	demo		8099	7.0.78 	True	True	10d
//MNPP		127.0.0.2 	Tomcat 	demo		8099	7.0.78 	True	True	10d
