package common

import (
	"context"
	"github.com/youchuangcd/gopkg/common/utils"
)

// GetUserIdByContext
// @Description: 从上下文中获取用户id
// @param ctx
// @return uint
// @return bool
func GetUserIdByContext(ctx context.Context) (uint, bool) {
	return utils.GetUserIdByContext(ctx)
}

// GetSysUserIdByContext
// @Description: 从上下文中获取系统用户id
// @param ctx
// @return uint
// @return bool
func GetSysUserIdByContext(ctx context.Context) (uint, bool) {
	return utils.GetSysUserIdByContext(ctx)
}

// GetTraceIdByContext
//
//	@Description: 从上下文中获取traceId
//	@param ctx
//	@return string
//	@return bool
func GetTraceIdByContext(ctx context.Context) (string, bool) {
	return utils.GetTraceIdByContext(ctx)
}
