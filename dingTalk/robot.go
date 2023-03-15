package dingTalk

import (
	"context"
	"github.com/youchuangcd/gopkg"

	httpclient "github.com/youchuangcd/gopkg/http-client"
	"github.com/youchuangcd/gopkg/mylog"
	"net/http"
)

// 文本消息
type TextMessage struct {
	Msgtype string             `form:"msgtype" json:"msgtype" `
	Text    TextMessageContent `form:"text" json:"text" `
}

type TextMessageContent struct {
	Content string `form:"content" json:"content" `
}

// 钉钉群：数据关联（机器人：数据关联汇总）
var (
	RobotDataRelationKeyword = "关联数据汇总"
)

type RobotRep struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

// 钉钉群机器人
func PushGroupTextMessage(ctx context.Context, url string, reqParam TextMessage) (err error) {
	var ret RobotRep
	headers := http.Header{}
	if err = httpclient.CallWithJson(ctx, &ret, "POST", url, headers, reqParam); err != nil {
		mylog.WithError(ctx, gopkg.LogCtx, map[string]interface{}{
			"err":   err,
			"param": reqParam,
		}, "发送到钉钉群的机器人消息失败")
		return
	}
	return
}
