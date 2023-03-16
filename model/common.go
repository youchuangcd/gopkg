package model

/**
 * 公共的model
 */

// PageInfo Paging common input parameter structure
type PageInfo struct {
	Page   int `form:"page" binding:"omitempty,gt=0"`      // 页码
	Limit  int `form:"page_size" binding:"omitempty,gt=0"` // 每页大小
	Offset int
}

// GetById Find by id structure
type GetById struct {
	ID int `json:"id" form:"id"` // 主键ID
}

func (r *GetById) Uint() uint {
	return uint(r.ID)
}

type Detail struct {
	ID uint `form:"id" json:"id" binding:"required"`
}

//
// CalculateOffsetLimit
// @Description: 计算分页偏移量
// @receiver p
//
func (p *PageInfo) CalculateOffsetLimit() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 10
	}
	if p.Limit >= 1000 {
		p.Limit = 1000
	}
	p.Offset = p.Limit * (p.Page - 1)
}

type ResponseList struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

//
// ResponseNoPageList
// @Description: 不分页列表
//
type ResponseNoPageList struct {
	List interface{} `json:"list"`
}
