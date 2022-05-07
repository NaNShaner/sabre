package changefile

import "testing"

func TestChangefile(t *testing.T) {
	f := "/tmp/apache-tomcat-7.0.75/bin/catalina.sh"
	replace := make(map[string]string)
	replace["JAVAOPTS"] = "-server -Xms1024M -Xmx1024M -Xss512k -XX:PermSize=1024M -XX:MaxPermSize=1024M -XX:+DisableExplicitGC -XX:MaxTenuringThreshold=31 -XX:+UseConcMarkSweepGC -XX:+UseParNewGC -XX:+CMSParallelRemarkEnabled -XX:+UseCMSCompactAtFullCollection -XX:LargePageSizeInBytes=128m -XX:+UseFastAccessorMethods -XX:+UseCMSInitiatingOccupancyOnly -XX:+PrintGCDetails -XX:+PrintGCTimeStamps -XX:+PrintHeapAtGC -XX:+HeapDumpOnOutOfMemoryError -XX:+HeapDumpOnCtrlBreak -Djava.awt.headless=true"
	err := Changefile(f, replace)
	if err != nil {
		t.Errorf("error ==> %s", err)
	}
}
