package cmd

import "testing"

func TestExecute(t *testing.T) {
	err := Execute()
	if err != nil {
		return 
	}
}
