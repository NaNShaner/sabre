package cmdline

import "testing"

func TestDeployTomcat(t *testing.T) {
	_, err := DeployTomcat()
	if err != nil {
		return
	}
}
