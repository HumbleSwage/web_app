package dao

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"strings"
	"time"
	"web_app/config"
)

var _db *gorm.DB

func InitMySql() {
	mConfig := config.Config.MySql["default"]
	conn := strings.Join([]string{mConfig.UserName, ":", mConfig.Password,
		"@tcp(", mConfig.DbHost, ":", mConfig.Password, ")/", mConfig.DbName, "?character=", mConfig.Charset, "&parseTime=True"}, "")

	var ormLogger = logger.Default
	if gin.Mode() == "" {
		ormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                           conn,
		DefaultStringSize:             256,   // string类型的默认长度
		DisableDatetimePrecision:      true,  // 禁用datetime精度，MySQL5.6之前的数据库以前不支持
		DontSupportRenameColumn:       true,  // 重命名索引时采用删除并新建的方式，MySQL5.7之前的数据库和MariaDB重命名索引
		DontSupportRenameColumnUnique: true,  // 用`change`重命名列，MySQL8之前的数据库和MariaDB不支持
		SkipInitializeWithVersion:     false, // 根据版本自动配置

	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 单数形式
		},
	})
	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(20)  // 设置最大的连接池
	sqlDB.SetMaxOpenConns(100) // 设置最大打开
	sqlDB.SetConnMaxLifetime(time.Second * 20)

	_db = db
}

func NewSQLClient(ctx context.Context) *gorm.DB {
	return _db.WithContext(ctx)
}
