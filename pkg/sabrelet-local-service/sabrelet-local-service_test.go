package sabrelet_local_service

import (
	"github.com/golang-module/carbon/v2"
	"sabre/pkg/dbload"
	"sabre/pkg/yamlfmt"
	"testing"
	"time"
)

func TestGetInfoList(t *testing.T) {
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	t1 := time.Now().In(cstSh).Format("2006-01-02 15:04:05.1234")

	t.Log(carbon.Parse("2022-06-02 15:04:05.1234").DiffInDays(carbon.Parse(t1)))

}

func TestQueryDB(t *testing.T) {
	serverName, serverNameErr := LocalServerName()
	if serverNameErr != nil {
		t.Log(serverNameErr)
		return
	}
	queryKey, err := GetInfoList(serverName)
	if err != nil {
		t.Log(err)
		return
	}

	db, err := QueryDB(serverName, queryKey)
	if err != nil {
		t.Log(err)
		return
	}

	t.Log(db)

	reportHostSabreletStatus, err := ReportHostSabreletStatus(db)
	if err != nil {
		t.Log(err)
		return
	}
	printResultJson, err := yamlfmt.PrintResultJson(reportHostSabreletStatus)
	if err != nil {
		t.Log(err)
		return
	}
	t.Logf("==> %s\n", string(printResultJson))
	setIntoDBErr := dbload.SetIntoDB(db, string(printResultJson))
	if setIntoDBErr != nil {
		t.Log(setIntoDBErr)
		return
	}
	t.Log(setIntoDBErr)
}
