package apiserver

import (
	"sabre/pkg/dbload"
	"sabre/pkg/yamlfmt"
	"testing"
)

func TestHttpReq(t *testing.T) {
	var dbInfo ToDBServer

	dbInfo.Kname = "/mid/MNPP/tomcat"

	resultJson, err := yamlfmt.PrintResultJson(dbInfo)
	if err != nil {
		return
	}

	dbloadErr := dbload.SetIntoDB("/mid/MNPP/tomcat", string(resultJson))
	if dbloadErr != nil {
		t.Errorf("入库失败, %s\n", dbloadErr)
	}
	t.Log("入库成功\n")

}
