// Package dbload
/*
格式：prefix + "/" + 资源类型 + "/" + namespace + "/" + 具体资源名

例如：当前机器的namespace
key :/midRegx/namespace/{hostname/ipaddr}
sabrectl get namespace
输出：
MNPP


例如：获取mnpp下的Tomcat
key :/registry/deployments/MNPP/Tomcat/{projectName}/{hostname/ipaddr}
sabrectl get tomcat demo
输出：
namespace	host		midType	projectName	version
MNPP		127.0.0.1 	Tomcat 	demo		7.0.78
MNPP		127.0.0.2 	Tomcat 	demo		7.0.78

sabrectl get tomcat demo -d
输出：
namespace	host		midType	projectName	port	version	monitor running	runningTime
MNPP		127.0.0.1 	Tomcat 	demo		8099	7.0.78 	True	True	10d
MNPP		127.0.0.2 	Tomcat 	demo		8099	7.0.78 	True	True	10d

*/
package dbload

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func SetIntoDB(k string, v interface{}) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"124.71.219.53:3380"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	defer cli.Close()

	//ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30))
	resp, err := cli.Put(context.TODO(), "name", "sabre")

	if err != nil {
		switch err {
		case context.Canceled:
			return fmt.Errorf("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			return fmt.Errorf("ctx is attached with a deadline is exceeded: %v", err)
		case rpctypes.ErrEmptyKey:
			return fmt.Errorf("client-side error: %v", err)
		default:
			return fmt.Errorf("bad cluster endpoints, which are not etcd servers: %v", err)
		}
	}
	fmt.Printf("%+v\n", *resp)
	return nil

}
