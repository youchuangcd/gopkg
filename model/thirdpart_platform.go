package model

const (
	// 平台id枚举
	ThirdpartPlatformIdWechat       uint = iota + 1 // 个微(已失效)
	ThirdpartPlatformIdWechatWework                 // 企微
	ThirdpartPlatformId788Code                      // 788码(用户还未注册企微，通过手机号搜到了微信)
	ThirdpartPlatformIdKs                           // 快手
	ThirdpartPlatformIdDy                           // 抖音
	ThirdpartPlatformIdQq                           // QQ
	ThirdpartPlatformIdGw                           // 个微
	ThirdpartPlatformIdSph                          // 视频号
	ThirdpartPlatformIdXhs                          // 小红书
	ThirdpartPlatformIdWxOrder                      // 微信小程序订单（云货优选和其它）
	ThirdpartPlatformIdQQAccount                    // QQ号
)

var ThirdpartPlatformIds = map[uint]string{
	ThirdpartPlatformIdWechat:       "个微",
	ThirdpartPlatformIdWechatWework: "企微",
	ThirdpartPlatformId788Code:      "788码",
	ThirdpartPlatformIdKs:           "快手",
	ThirdpartPlatformIdDy:           "抖音",
	ThirdpartPlatformIdQq:           "QQ",
	ThirdpartPlatformIdGw:           "个微",
	ThirdpartPlatformIdSph:          "视频号",
	ThirdpartPlatformIdXhs:          "小红书",
	ThirdpartPlatformIdWxOrder:      "微信小程序订单", //（云货优选和其它）
	ThirdpartPlatformIdQQAccount:    "QQ号",
}

const ThirdpartPlatformIdDwsMobileData = 999         //dws手机号数据
const ThirdpartPlatformIdDwsSuccessMobileData = 1000 //dws success手机号数据
const ThirdpartPlatformIdScanBuildDyMetaData = 10000 // 抖音扫描元数据建立索引表

// 关联类型
type RelationType int

const (
	RelationTypeOrder               RelationType = 1001 // 订单数据
	RelationTypeGood                RelationType = 1002 // 商品数据
	RelationTypeLiveAction          RelationType = 1003 // 直播间行为数据
	RelationTypeTarget              RelationType = 1004 // 关键词目标数据
	RelationTypeVideoList           RelationType = 1005 // 视频列表数据
	RelationTypeVideoComment        RelationType = 1006 // 视频评论数据
	RelationTypeLiveLuckyBagDy      RelationType = 1007 // 福袋数据
	RelationTypeMonitorTaskCallback RelationType = 1008 // 主页链接用户数据
	RelationTypeUserLevel           RelationType = 1009 // 用户等级数据
	RelationTypeGoodsComment        RelationType = 1010 // 商品评论数据
)

// 关联来源
type RelationSource int

const (
	RelationSourceDou1        RelationSource = 1 // 抖一
	RelationSourceBiLin       RelationSource = 2 // 比邻
	RelationSourcePluginOrder RelationSource = 3 // 插件订单
	RelationSourceJinTun      RelationSource = 4 // 鲸吞
	RelationSourceKFPT        RelationSource = 5 // 开放平台
)

// 匹配来源
type MatchSource int

const (
	MatchSourceDefault                         MatchSource = 0 // 默认通讯录匹配
	MatchSourcePluginOrder                     MatchSource = 1 // 插件订单,收货人手机号
	MatchSource788Avatar                       MatchSource = 2 // 788头像昵称匹配（张宇）
	MatchSourceEmptyAccountCheckAvatarNickname MatchSource = 3 // 空号检测头像昵称普配（乐远）
)