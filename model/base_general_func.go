package model

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//
// Add
// @Description:
// @param ctx
// @param tx
// @param data 必须是指针
// @param args ...interface{} 支持传递model和批量插入大小参数
// @return err
//
func Add(ctx context.Context, tx *gorm.DB, data interface{}, args ...interface{}) (err error) {
	db := GetDB(ctx, tx)
	argLen := len(args)
	if argLen > 0 {
		if argLen == 1 {
			db = db.Model(args[0])
		} else if argLen == 2 {
			if size, ok := args[1].(int); ok {
				db = db.Model(args[0])
				db.CreateBatchSize = size
			}
		}
	}
	err = db.Create(data).Error
	return
}

//
// AddIgnore
//  @Description: 插入忽略错误
//  @param ctx
//  @param tx
//  @param data
//  @param args
//  @return err
//
func AddIgnore(ctx context.Context, tx *gorm.DB, data interface{}, args ...interface{}) (err error) {
	db := GetDB(ctx, tx)
	argLen := len(args)
	if argLen > 0 {
		if argLen == 1 {
			db = db.Model(args[0])
		} else if argLen == 2 {
			if size, ok := args[1].(int); ok {
				db = db.Model(args[0])
				db.CreateBatchSize = size
			}
		}
	}
	err = db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(data).Error
	return
}

//
// Update
// @Description:
// @param ctx
// @param tx
// @param fields
// @param data 必须是指针
// @param args
// @return affectedRow
// @return err
//
func Update(ctx context.Context, tx *gorm.DB, fields []string, data interface{}, args ...interface{}) (affectedRow int64, err error) {
	db := GetDB(ctx, tx).Select(fields)
	// 可变参数传入model
	if len(args) == 1 {
		db = db.Model(args[0])
	}
	res := db.Updates(data)
	affectedRow = res.RowsAffected
	err = res.Error
	return
}

func UpdateByIds(ctx context.Context, tx *gorm.DB, ids []uint, fields []string, data interface{}, args ...interface{}) (affectedRow int64, err error) {
	db := GetDB(ctx, tx).Select(fields).Where("id IN ?", ids)
	// 可变参数传入model
	if len(args) == 1 {
		db = db.Model(args[0])
	}
	res := db.Updates(data)
	affectedRow = res.RowsAffected
	err = res.Error
	return
}

func UpdateSetMapById(ctx context.Context, tx *gorm.DB, model Interface, fields []string, data map[string]interface{}) (affectedRow int64, err error) {
	res := GetDB(ctx, tx).Model(model).Select(fields).Updates(data)
	affectedRow = res.RowsAffected
	err = res.Error
	return
}

//
// GetDetailById
// @Description: 根据id获取详情
// @param ctx
// @param model
// @param id
// @param res 必须是指针，要不然无法写值
// @param args
// @return error
//
func GetDetailById(ctx context.Context, model Interface, id uint, res interface{}, args ...interface{}) error {
	db := GetDB(ctx, nil).Model(&model).Where("id = ?", id)
	if len(args) == 1 {
		db = db.Select(args[0])
	}
	err := db.Take(&res).Error
	return err
}

//
// GetDetailByPrimary
// @Description: 根据主键id获取详情
// @param ctx
// @param model 必须是指针，要不然无法写值
// @param args
// @return error
//
func GetDetailByPrimary(ctx context.Context, model Interface, args ...interface{}) error {
	db := GetDB(ctx, nil)
	if len(args) == 1 {
		db = db.Select(args[0])
	}
	err := db.Take(model).Error
	return err
}

//
// GetDetailByPrimaryInTransaction
// @Description: 根据主键id获取详情; 在事务中使用
// @param ctx
// @param tx
// @param model
// @param args
// @return error
//
func GetDetailByPrimaryInTransaction(ctx context.Context, tx *gorm.DB, model Interface, args ...interface{}) error {
	db := GetDB(ctx, tx)
	if len(args) == 1 {
		db = db.Select(args[0])
	}
	err := db.Take(model).Error
	return err
}

//
// GetPreloadDetailById
// @Description: 关联获取详情
// @param ctx
// @param model
// @param id
// @param res
// @param args
// @return error
//
func GetPreloadDetailById(ctx context.Context, model Interface, id uint, res interface{}, args ...interface{}) error {
	db := GetDB(ctx, nil).Model(&model).Where("id = ?", id)
	if len(args) > 0 {
		for _, v := range args {
			if column, ok := v.(string); ok {
				db = db.Preload(column)
			}
		}
	}
	err := db.Take(&res).Error
	return err
}

//
// ChunkGetDataByCondition
// @Description: 根据自定义条件分批获取数据
// @param ctx
// @param model
// @param fields
// @param res
// @param limit
// @param whereStr
// @param param
// @return err
//
func ChunkGetDataByCondition(ctx context.Context, model Interface, fields []string, res interface{}, limit int, whereStr string, param ...interface{}) (err error) {
	err = GetDB(ctx, nil).Model(model).Select(fields).Where(whereStr, param...).Limit(limit).Find(res).Error
	return
}
