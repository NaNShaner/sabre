package hostregister

import "testing"

func TestKeyName(t *testing.T) {
	// e.g. /hosts/erp/machine/app/hostname
	n := "erp"
	a := "app"
	h := "local"
	i := "192.168.3.58"

	t.Log(KeyName(n, h, a, i))
}
