package bootstrap

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gohub/pkg/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm/logger"

	"gohub/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// SetupDB 初始化数据库和ORM

func SetupDB{
	var dbConfig gorm.Dialector
	switch config.Get("database.connection") {
	case mysql:
		// 构建DSN信息
		dsn:=fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&multiStatements=true&loc=Local",
			config.Get("database.mysql.username"),
			config.Get("database.mysql.password"),
			config.Get("database.mysql.host"),
			config.Get("database.mysql.port"),
			config.Get("database.mysql.database"),
			config.Get("database.mysql.charset"),
			)
		dbConfig=mysql.New(mysql.Config{
			DSN:dsn,
		})
	case "sqlite":
		// 初始化sqlite
		database:=config.Get("database.sqlite.database")
		dbConfig=sqlite.Open(database)
	default:
		panic(errors.New("database connection not supported"))
	}

	// 连接数据库，并设置GORM的日志模式
	database.Connect(dbConfig,logger.Default.LogMode(logger.Info))

	// 设置最大连接数
	database.SQLDB.SetMaxOpenConns(config.GetInt("database.mysql.max_open_connections"))

	// 设置最大空闲连接
	database.SQLDB.SetMaxIdleConns(config.GetInt("database.mysql.max_idle_connections"))

	// 设置每个连接过期时间
	database.SQLDB.SetConnMaxLifetime(time.Duration(config.GetInt("database.mysql.max_life_seconds"))*time.Seconds)
}