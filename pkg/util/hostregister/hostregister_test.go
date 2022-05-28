package hostregister

import "testing"

func TestKeyName(t *testing.T) {
	// e.g. /hosts/erp/machine/app/hostname
	n := "erp"
	a := "app"
	h := "local"

	t.Log(KeyName(n, h, a))
}
