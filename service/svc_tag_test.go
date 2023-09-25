package service

import (
	"testing"
)

func TestDailySyncConsulSvcs(t *testing.T) {
	err := DailySyncConsulSvcs()
	if err != nil {
		t.Fatal(err)
	}
}
