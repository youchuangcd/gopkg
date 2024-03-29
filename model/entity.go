package model

const (
	// 实体类型
	EntityTypeMobile uint16 = iota + 1 // 手机号
)

// entity数据来源
type EntityChannel uint16

const (
	EntityChannelInner    EntityChannel = iota + 1 //内部提供
	EntityChannelMerchant                          //同步788数据
	EntityChannelPlugin
	EntityChannelSms          // 短信导入
	EntityChannelWechatOrder  // 小程序订单
	EntityChannelWechat       // 个微关联数据
	EntityChannelJintunWechat // 鲸吞个微成为好友数据
	EntityChannel8            // DATE4.19_NUM_2.8Y
	EntityChannel9            // 特征匹配来源
	EntityChannel10           // ga
	EntityChannel11           // f_xxx_gb_4_23 来源
)

var EntityChannelMap = map[EntityChannel]string{
	EntityChannelInner:        "内部提供",
	EntityChannelMerchant:     "同步788数据",
	EntityChannelPlugin:       "插件来源",
	EntityChannelSms:          "短信导入",
	EntityChannelWechatOrder:  "小程序订单",
	EntityChannelWechat:       "个微关联数据",
	EntityChannelJintunWechat: "鲸吞个微成为好友数据",
	EntityChannel8:            "DATE4.19_NUM_2.8Y",
	EntityChannel9:            "特征匹配",
	EntityChannel10:           "ga",
	EntityChannel11:           "f_xxx_gb_4_23",
}
