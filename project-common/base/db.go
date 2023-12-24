package base

import (
	"com.levi/project-common/config"
	"github.com/natefinch/lumberjack"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

var Mysql *gorm.DB

func InitMysql() {
	dbConfig := config.GlobalConf.MySql
	dsn := dbConfig.Username + ":" + dbConfig.Password + "@tcp(" + dbConfig.Host + ":" + strconv.Itoa(dbConfig.Port) + ")/" + dbConfig.Database + "?charset=" + dbConfig.Charset + "&parseTime=True&loc=Local"
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束
		Logger:                                   getGormLogger(),
	}); err != nil {
		log.Fatalf("mysql connect error %v", err)
		return
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
		sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
		Mysql = db
	}
}

func getGormLogger() logger.Interface {
	var lm logger.LogLevel
	switch config.GlobalConf.MySql.LogMode {
	case "silent":
		lm = logger.Silent
	case "error":
		lm = logger.Error
	case "warn":
		lm = logger.Warn
	case "info":
		lm = logger.Info
	default:
		lm = logger.Info
	}
	return logger.New(getGormLogWriter(), logger.Config{
		SlowThreshold:             200 * time.Millisecond,                      // 慢 SQL 阈值
		LogLevel:                  lm,                                          // 日志级别
		IgnoreRecordNotFoundError: false,                                       // 忽略ErrRecordNotFound（记录未找到）错误
		Colorful:                  config.GlobalConf.MySql.EnableFileLogWriter, // 禁用彩色打印
	})
}

func getGormLogWriter() logger.Writer {
	var writer io.Writer

	// 是否启用日志文件
	if config.GlobalConf.MySql.EnableFileLogWriter {
		// 自定义 Writer
		writer = &lumberjack.Logger{
			Filename:   "/Users/levi/go/src/project-gin-grpc/project-common/log/db/db.log",
			MaxSize:    500,
			MaxBackups: 60,
			MaxAge:     1,
			Compress:   true,
		}
	} else {
		writer = os.Stdout
	}
	return log.New(writer, "\r\n", log.LstdFlags)
}
