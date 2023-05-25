package initialize

import (
	"fmt"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"github.com/youchuangcd/gopkg/mylog"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type DBConfigItem struct {
	Driver      string `mapstructure:"driver"`
	Name        string `mapstructure:"name"`
	Dsn         string `mapstructure:"dsn"`
	MaxIdle     int    `mapstructure:"maxIdle"`     //空闲连接池中连接的最大数量
	MaxOpen     int    `mapstructure:"maxOpen"`     //打开数据库连接的最大数量
	MaxLifetime int    `mapstructure:"maxLifetime"` //连接可复用的最大时间 单位：秒
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
		var dialector gorm.Dialector
		switch v.Driver {
		case "mysql":
			dialector = mysql.Open(v.Dsn)
		case "postgres":
			dialector = postgres.New(postgres.Config{
				PreferSimpleProtocol: true,
				DSN:                  v.Dsn,
				WithoutReturning:     false, // 关闭insert返回id
			})
		case "sqlserver":
			dialector = sqlserver.Open(v.Dsn)
		default:
			panic(fmt.Sprintf("数据库: %s; 无效的数据库驱动: %s", v.Name, v.Driver))
		}
		db, err := gorm.Open(dialector, &gorm.Config{
			Logger: slowLogger,
			// 为了确保数据一致性，GORM 会在事务里执行写入操作（创建、更新、删除）。
			// 如果没有这方面的要求，您可以在初始化时禁用它，这将获得大约 30%+ 性能提升。
			SkipDefaultTransaction: true,
		})
		if err != nil {
			panic(fmt.Sprintf("%s数据库%s初始化失败，原因：%s", v.Driver, v.Name, err.Error()))
		}
		sqlDB, err := db.DB()
		if err != nil {
			panic(fmt.Sprintf("获取%s数据库%sDB对象失败，原因：%s", v.Driver, v.Name, err.Error()))
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
		// 设置分布式追踪插件失败
		if err = db.Use(otelgorm.NewPlugin(otelgorm.WithDBName(v.Name))); err != nil {
			panic(fmt.Sprintf("设置%s数据库%sDB追踪插件失败，原因：%s", v.Driver, v.Name, err.Error()))
		}
		gormDBMap[v.Name] = db
	}
}
