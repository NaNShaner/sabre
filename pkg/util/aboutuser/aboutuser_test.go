package aboutuser

import (
	"testing"
)

func TestIsUserExist(t *testing.T) {
	exist, err := IsUserExist("miduser")
	if err != nil {
		return
	}
	t.Log(exist)
}
