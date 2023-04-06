package common

import "testing"

func TestIsMobile(t *testing.T) {
	mobile := "+8615456657895"
	t.Log(mobile, IsMobile(mobile))
	mobile = "8615456657895"
	t.Log(mobile, IsMobile(mobile))
	mobile = "15456657895"
	t.Log(mobile, IsMobile(mobile))
	mobile = "5456657895"
	t.Log(mobile, IsMobile(mobile))
	mobile = "123e213yu21y1378"
	t.Log(mobile, IsMobile(mobile))
}
