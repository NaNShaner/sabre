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
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

//GetDBCli 获取ETCDCli
func GetDBCli() (*clientv3.Client, error) {
	// endpoints := []string{"192.168.3.111:2379"}
	endpoints := []string{"124.71.219.53:2379"}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get cli for etcd, %s. The server address is %s\n", err, endpoints)
	}
	return cli, err
}

//SetIntoDB 入库
//TODO： apiserver的地址硬编码
func SetIntoDB(k, v string) error {
	cli, err := GetDBCli()
	if err != nil {
		return err
	}
	defer cli.Close()

	//ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30))
	_, putErr := cli.Put(context.TODO(), k, v)

	if putErr != nil {
		switch err {
		case context.Canceled:
			return fmt.Errorf("ctx is canceled by another routine: %v\n", putErr)
		case context.DeadlineExceeded:
			return fmt.Errorf("ctx is attached with a deadline is exceeded: %v\n", putErr)
		case rpctypes.ErrEmptyKey:
			return fmt.Errorf("client-side error: %v\n", putErr)
		default:
			return fmt.Errorf("bad cluster endpoints, which are not etcd servers: %v\n", putErr)
		}
	}
	//L.Log.Infof("数据库入库成功，%s：%s", k, v)
	return nil

}

//GetKeyFromETCD 以k为前缀获取etcd中的key
//withPrefix 为true时即模糊查询etcd中所有符合k为前缀的kv，false表示精确匹配
func GetKeyFromETCD(k string, withPrefix bool) ([]*mvccpb.KeyValue, error) {
	cli, err := GetDBCli()
	if err != nil {
		return nil, fmt.Errorf("connect failed, %s\n", err)
	}
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	if withPrefix {
		resp, getKeyErr := cli.Get(ctx, k, clientv3.WithPrefix())
		cancel()
		if getKeyErr != nil {

			return nil, fmt.Errorf("get from etcd failed, err:%v\n", getKeyErr)
		}
		return resp.Kvs, nil
	} else {
		resp, getKeyErr := cli.Get(ctx, k)
		cancel()
		if getKeyErr != nil {

			return nil, fmt.Errorf("get from etcd failed, err:%v\n", getKeyErr)
		}
		return resp.Kvs, nil
	}
}
