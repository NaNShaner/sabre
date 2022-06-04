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
	l "sabre/pkg/util/logbase/logscheduled"
	"strings"
	"time"
)

var (
	midRegx  = "mid"
	hostRegx = "hosts"
)

var (
	midTomcat = "tomcat"
	// midJDK    = "jdk"
)

//GetDBCli 获取ETCDCli
func GetDBCli() (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		//Endpoints: []string{"124.71.219.53:2379"},
		Endpoints:   []string{"192.168.3.111:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		l.Log.Errorf("failed to get cli for etcd, %s\n", err)
		return nil, err
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
	_, err = cli.Put(context.TODO(), k, v)

	if err != nil {
		switch err {
		case context.Canceled:
			return fmt.Errorf("ctx is canceled by another routine: %v\n", err)
		case context.DeadlineExceeded:
			return fmt.Errorf("ctx is attached with a deadline is exceeded: %v\n", err)
		case rpctypes.ErrEmptyKey:
			return fmt.Errorf("client-side error: %v\n", err)
		default:
			return fmt.Errorf("bad cluster endpoints, which are not etcd servers: %v\n", err)
		}
	}
	//fmt.Printf("%+v\n", *resp)
	return nil

}

//WatchFromDB 通过API网关对etcd中的资源类型进行watch，进行后续调度
func WatchFromDB(s string) {
	cli, err := GetDBCli()
	if err != nil {
		l.Log.Errorf("connect failed, %s\n", err)
		return
	}
	defer cli.Close()

	for {
		// clientv3.WithPrefix() 监控s作为前缀key值的value变化，默认为精确watch
		rch := cli.Watch(context.Background(), s, clientv3.WithPrefix())
		for wresp := range rch {
			err = wresp.Err()
			if err != nil {
				l.Log.Errorf("etcd watch response err, %s\n", err)
			}
			for _, ev := range wresp.Events {
				//TODO: 判断执行动作，发起调度指令
				l.Log.Info("%s %q %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				keySplit, err := keySplit(ev.Kv.Key)
				if err != nil {
					return
				}
				switch {
				case isMidType(keySplit):
					l.Log.Infof("isMidType: %s\n", ev.Kv.Key)
					switch {
					case isMidTypeOfTomcat(keySplit):
						l.Log.Infof("isMidTypeOfTomcat %s\n", ev.Kv.Key)
					}
				case isMidHost(keySplit):
					l.Log.Infof("isMidHost: %s\n", ev.Kv.Key)
				default:
					return
				}

			}
		}
	}
}

//GetKeyWithPrefix 以k为前缀获取etcd中的key，即模糊查询etcd中所有符合k为前缀的kv
func GetKeyWithPrefix(k string) ([]*mvccpb.KeyValue, error) {
	cli, err := GetDBCli()
	if err != nil {
		return nil, fmt.Errorf("connect failed, %s\n", err)
	}
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, getKeyErr := cli.Get(ctx, k, clientv3.WithPrefix())
	cancel()
	if getKeyErr != nil {

		return nil, fmt.Errorf("get from etcd failed, err:%v\n", getKeyErr)
	}
	return resp.Kvs, nil
}

func keySplit(t []byte) (string, error) {
	s := string(t)
	sSplit := strings.Split(s, "/")
	if len(sSplit) < 1 {
		return "", fmt.Errorf("Etcd key %s is not in normal format\n", s)
	}
	return sSplit[1], nil
}

func isMidType(t string) bool {
	if strings.Contains(t, midRegx) {
		return true
	} else {
		return false
	}
}

func isMidTypeOfTomcat(t string) bool {

	if strings.Contains(t, midTomcat) {
		return true
	} else {
		return false
	}
}

func isMidHost(t string) bool {
	if strings.Contains(t, hostRegx) {
		return true
	} else {
		return false
	}
}
