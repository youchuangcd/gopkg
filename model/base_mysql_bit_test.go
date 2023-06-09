package model

import (
	"encoding/json"
	"testing"
)

func TestBitBool(t *testing.T) {
	var b BitBool = false
	res, err := json.Marshal(b)
	t.Log(string(res), err)
	var c BitBool
	err = json.Unmarshal(res, &c)
	t.Log(c, err)
}
