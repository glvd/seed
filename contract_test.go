package seed

import (
	"encoding/json"
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

	var slist []string

	err = json.Unmarshal([]byte(list), &slist)
	if err != nil {
		return
	}
	t.Log(slist)
}
