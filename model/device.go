package model

// 设备状态
const (
	DeviceStatusClosed           = 0  //停用
	DeviceStatusEnabled          = 10 //启用
	DeviceStatusDisabled         = 11 //断网
	DeviceStatusException        = 14 //异常
	DeviceStatusRegionalAnomaly  = 15 //地区异常
	DeviceStatusDeviceProxyExp   = 16 //代理异常
	DeviceStatusServerException  = 17 //服务器异常
	DeviceStatusIpChangeAbnormal = 18 // IP变动异常
)

// 设备类型
const (
	DeviceTypeWechatIOS                = 10 //个微·IOS协议
	DeviceTypeWechatIPad               = 12 //个微·iPad协议
	DeviceTypeWechatAndroid            = 13 //个微·安卓协议
	DeviceTypeWechatAndroidSimulator   = 14 //个微·安卓模拟器
	DeviceTypeWechatWindowsProtocol    = 19 //个微·Windows协议
	DeviceTypeWechatMacProtocol        = 22 //个微·Mac协议
	DeviceTypeWechatAndroidPadProtocol = 23 //个微·安卓Pad协议
	DeviceTypeWorkWechatXposed         = 30 //企微·Xposed
	DeviceTypeWorkWechatWindows        = 31 //企微·Windows
	DeviceTypeWorkWechatIPad           = 32 //企微·iPad协议
	DeviceTypeWorkWechatIPhone         = 33 //企微·iPhone协议
	DeviceTypeWorkWechatAndroidCloud   = 34 //企微·安卓云真机
	DeviceTypeWorkWechatAndroid        = 35 //企微·安卓协议
	DeviceTypeQQAndroid                = 50 //QQ·安卓协议
	DeviceTypeQQWindowsProtocol        = 51 //QQ·Windows协议
)

const (
	DeviceProductTypeQQ         = 10 //QQ
	DeviceProductTypeWeChat     = 20 //个人微信
	DeviceProductTypeWorkWechat = 30 //企业微信
)

var (
	// 设备类型映射产品类型
	DeviceProductTypeMap = map[int]int{
		// 个微
		DeviceTypeWechatIOS:                DeviceProductTypeWeChat,
		DeviceTypeWechatIPad:               DeviceProductTypeWeChat,
		DeviceTypeWechatAndroid:            DeviceProductTypeWeChat,
		DeviceTypeWechatAndroidSimulator:   DeviceProductTypeWeChat,
		DeviceTypeWechatWindowsProtocol:    DeviceProductTypeWeChat,
		DeviceTypeWechatMacProtocol:        DeviceProductTypeWeChat,
		DeviceTypeWechatAndroidPadProtocol: DeviceProductTypeWeChat,
		// 企微
		DeviceTypeWorkWechatXposed:       DeviceProductTypeWorkWechat,
		DeviceTypeWorkWechatWindows:      DeviceProductTypeWorkWechat,
		DeviceTypeWorkWechatIPad:         DeviceProductTypeWorkWechat,
		DeviceTypeWorkWechatIPhone:       DeviceProductTypeWorkWechat,
		DeviceTypeWorkWechatAndroidCloud: DeviceProductTypeWorkWechat,
		DeviceTypeWorkWechatAndroid:      DeviceProductTypeWorkWechat,
		// QQ
		DeviceTypeQQAndroid:         DeviceProductTypeQQ,
		DeviceTypeQQWindowsProtocol: DeviceProductTypeQQ,
	}
)
