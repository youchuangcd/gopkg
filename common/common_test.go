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

func TestKsTraitCode(t *testing.T) {
	avatars := []string{
		"http://p1.a.yximgs.com/kimg/uhead/AB/2022/10/17/14/CjdCTWpBeU1qRXdNVGN4TkRNeU1URmZOelkzT0RFd05EazFYekZmYUdRMU9UUmZOVE01X3MuanBnEJXM1y8:200x200.heif",
		"http://alimov2.a.yximgs.com/kos/nlav12689/head.jpg",
		"http://p5.a.yximgs.com/kimg/uhead/AB/2020/12/30/18/CjtCTWpBeU1ERXlNekF4T0RFNE5EVmZNakU0T1RNME16Y3dNVjh5WDJoa016STBYek0xTmc9PV9zLmpwZxCVzNcv:200x200.heif",
		"http://p4.a.yximgs.com/kimg/uhead/AB/2017/05/14/14/CjVCTWpBeE56QTFNVFF4TkRFeE16VmZNVGd4TnpnM05UVTJYekpmYUdRMU1qUmZNalE9LmpwZxCVzNcv.heif",
		"http://p4.a.yximgs.com/kimg/uhead/AB/2019/12/20/18/CjtCTWpBeE9URXlNakF4T0RNek5ESmZNVFl6TkRjek9EUTFObDh5WDJoa05UUXlYell3TkE9PV9zLmpwZxCVzNcv:200x200.heif",
	}
	for _, v := range avatars {
		t.Log(GetKsAvatarTraitCode(v))
	}
}
