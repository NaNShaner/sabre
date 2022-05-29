package changefile

import "testing"

func TestChangefile(t *testing.T) {
	//f := "/Users/bijingrui/sabre/pkg/util/changefile/s.xml"
	//replace := make(map[string]string)
	//replace["JAVAOPTS"] = "-server -Xms1024M -Xmx1024M -Xss512k -XX:PermSize=1024M -XX:MaxPermSize=1024M -XX:+DisableExplicitGC -XX:MaxTenuringThreshold=31 -XX:+UseConcMarkSweepGC -XX:+UseParNewGC -XX:+CMSParallelRemarkEnabled -XX:+UseCMSCompactAtFullCollection -XX:LargePageSizeInBytes=128m -XX:+UseFastAccessorMethods -XX:+UseCMSInitiatingOccupancyOnly -XX:+PrintGCDetails -XX:+PrintGCTimeStamps -XX:+PrintHeapAtGC -XX:+HeapDumpOnOutOfMemoryError -XX:+HeapDumpOnCtrlBreak -Djava.awt.headless=true"
	//err := Changefile(f, replace)
	//if err != nil {
	//	t.Errorf("error ==> %s", err)
	//}

	// 修改server.xml
	serverXmlReplace := make(map[string]string)
	serverXmlReplace["shutdownport"] = "8099"
	serverXmlReplace["listeningport"] = "8009"
	serverXmlReplace["ajpport"] = "8005"
	//serverXmlReplace["ajprirectport"] = "8443"
	fc := "/Users/bijingrui/sabre/pkg/util/changefile/s.xml"
	err := Changefile(fc, serverXmlReplace)
	if err != nil {
		t.Errorf("error ==> %s", err)
	}

}
