package sabrelet_local_service

import (
	"github.com/golang-module/carbon/v2"
	"testing"
	"time"
)

func TestGetInfoList(t *testing.T) {
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	t1 := time.Now().In(cstSh).Format("2006-01-02 15:04:05.1234")

	t.Log(carbon.Parse("2022-06-02 15:04:05.1234").DiffInDays(carbon.Parse(t1)))

}
