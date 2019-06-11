package seed

import (
	"testing"
)

// TestETH_CheckExist ...
func TestGetHostList(t *testing.T) {
	err := InitGlobalETH("", "")
	if err != nil {
		return
	}
	list := GetHostList()
	t.Log(list)
}

// TestUpdateHotList ...
func TestUpdateHotList(t *testing.T) {
	err := InitGlobalETH("", "")
	if err != nil {
		return
	}
	list := UpdateHotList()
	t.Log(list)

}
