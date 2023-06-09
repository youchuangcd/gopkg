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
	DeviceTypeIOS                   = 10 //IOS协议设备
	DeviceTypeScanCode              = 12 //扫码号设备
	DeviceTypeAndroid               = 13 //安卓协议设备
	DeviceTypeAndroidSimulator      = 14 //安卓模拟器
	DeviceTypeWindows               = 0  //Windows协议设备
	DeviceTypeAccountTrusteeship    = 0  //安卓托管设备
	DeviceTypeChatRoomGuard         = 17 //群卫士协议设备类型
	DeviceTypeWindowsSuperAddFrined = 18 //Windows超级加人设备
	DeviceTypeWindowsProtocol       = 19 //windows协议版设备类型
	DeviceTypeAZProtocol            = 20 //安卓协议_新
	DeviceTypeWindowsPdd            = 21 //拼多多分身windows协议设备
	DeviceTypeMacProtocol           = 22 //Mac协议设备类型
	DeviceTypeAndroidPadProtocol    = 23 //安卓平板协议
	DeviceTypeCloudProtocol         = 24 //云机设备
	DeviceTypeWindowsLocal          = 25 //Windows本地部署
)

const (
	ProductTypeQQ         = 10 //QQ
	ProductTypeWeChat     = 20 //个人微信
	ProductTypeWorkWechat = 30 //企业微信
)
