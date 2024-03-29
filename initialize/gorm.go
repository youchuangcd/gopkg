package initialize

import (
	"fmt"
	"github.com/youchuangcd/gopkg/mylog"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

type DBConfigItem struct {
	Driver          string `mapstructure:"driver" yaml:"driver"`
	Name            string `mapstructure:"name" yaml:"name"`
	Dsn             string `mapstructure:"dsn" yaml:"dsn"`
	MaxIdle         int    `mapstructure:"maxIdle" yaml:"maxIdle"`                 //空闲连接池中连接的最大数量
	MaxOpen         int    `mapstructure:"maxOpen" yaml:"maxOpen"`                 //打开数据库连接的最大数量
	MaxLifetime     int    `mapstructure:"maxLifetime" yaml:"maxLifetime"`         //连接可复用的最大时间 单位：秒
	TableNamePrefix string `mapstructure:"tableNamePrefix" yaml:"tableNamePrefix"` // 表名前缀
}

func Gorm(dbList []DBConfigItem, gormLevel int, gormDBMap map[string]*gorm.DB) {
	// GORM
	logLevel := logger.Error
	if gormLevel >= int(logger.Silent) && gormLevel <= int(logger.Info) {
		logLevel = logger.LogLevel(gormLevel)
	}
	slowLogger := mylog.NewGormLogger(
		//设置Logger
		mylog.GetLogger(),
		logger.Config{
			//慢SQL阈值
			SlowThreshold: 200 * time.Millisecond,
			//设置日志级别，只有Info以上才会打印sql
			LogLevel: logLevel,
			// 忽略未找到记录的错误日志
			IgnoreRecordNotFoundError: true,
		},
	)
	for _, v := range dbList {
		var (
			dialector           gorm.Dialector
			disableLastInsertId = false
		)
		switch v.Driver {
		case "mysql":
			dialector = mysql.Open(v.Dsn)
		case "postgres":
			disableLastInsertId = true
			dialector = postgres.New(postgres.Config{
				PreferSimpleProtocol: true,
				DSN:                  v.Dsn,
				WithoutReturning:     disableLastInsertId, // 关闭insert返回id
			})
		case "sqlserver":
			dialector = sqlserver.Open(v.Dsn)
		default:
			panic(fmt.Sprintf("数据库: %s; 无效的数据库驱动: %s", v.Name, v.Driver))
		}
		db, err := gorm.Open(dialector, &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   v.TableNamePrefix, // 表名前缀
				SingularTable: true,
			},
			Logger: slowLogger,
			// 为了确保数据一致性，GORM 会在事务里执行写入操作（创建、更新、删除）。
			// 如果没有这方面的要求，您可以在初始化时禁用它，这将获得大约 30%+ 性能提升。
			SkipDefaultTransaction: true,
			DisableLastInsertId:    disableLastInsertId, // TODO 这个魔改的配置，后续官方支持关闭功能，就会撤掉
		})
		if err != nil {
			panic(fmt.Sprintf("%s数据库%s(%s)初始化失败，原因：%s", v.Driver, v.Name, v.Dsn, err.Error()))
		}
		sqlDB, err := db.DB()
		if err != nil {
			panic(fmt.Sprintf("获取%s数据库%s(%s)DB对象失败，原因：%s", v.Driver, v.Name, v.Dsn, err.Error()))
		}
		if v.MaxIdle > 0 {
			// SetMaxIdleConns 设置空闲连接池中连接的最大数量
			sqlDB.SetMaxIdleConns(v.MaxIdle)
		}
		if v.MaxOpen > 0 {
			// SetMaxOpenConns 设置打开数据库连接的最大数量。
			sqlDB.SetMaxOpenConns(v.MaxOpen)
		}
		if v.MaxLifetime > 0 {
			// SetConnMaxLifetime 设置了连接可复用的最大时间。
			sqlDB.SetConnMaxLifetime(time.Duration(v.MaxLifetime) * time.Second)
		}
		//// 设置分布式追踪插件失败
		//if err = db.Use(otelgorm.NewPlugin(otelgorm.WithDBName(v.Name), otelgorm.WithoutMetrics(), otelgorm.WithAttributes(
		//	attribute.Int("MaxOpen", v.MaxOpen),
		//	attribute.Int("MaxIdle", v.MaxIdle),
		//	attribute.Int("MaxLifetime", v.MaxLifetime),
		//))); err != nil {
		//	panic(fmt.Sprintf("设置%s数据库%sDB追踪插件失败，原因：%s", v.Driver, v.Name, err.Error()))
		//}
		gormDBMap[v.Name] = db
	}
}
