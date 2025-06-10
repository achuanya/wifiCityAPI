package database

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin/wifiCityAPI/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

// DB 是一个全局的 gorm.DB 实例
var DB *gorm.DB

// Init 函数用于初始化数据库连接
func Init() {
	var err error

	// 基础 GORM 配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // GORM 日志配置
	}

	// 连接主库
	db, err := gorm.Open(mysql.Open(config.Cfg.Database.Master.DSN), gormConfig)
	if err != nil {
		log.Fatalf("无法连接到主数据库: %v", err)
	}

	// 配置从库
	replicas := make([]gorm.Dialector, len(config.Cfg.Database.Slaves))
	for i, slave := range config.Cfg.Database.Slaves {
		replicas[i] = mysql.Open(slave.DSN)
	}

	// 配置读写分离插件
	err = db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(config.Cfg.Database.Master.DSN)}, // 主库（写）
		Replicas: replicas,                                                     // 从库（读）
		Policy:   dbresolver.RandomPolicy{},                                    // 负载均衡策略：随机
	}).SetConnMaxIdleTime(config.Cfg.Database.Settings.ConnMaxIdleTime).
		SetConnMaxLifetime(config.Cfg.Database.Settings.ConnMaxLifetime).
		SetMaxIdleConns(config.Cfg.Database.Settings.MaxIdleConns).
		SetMaxOpenConns(config.Cfg.Database.Settings.MaxOpenConns))

	if err != nil {
		log.Fatalf("无法配置数据库读写分离: %v", err)
	}

	// 获取底层的 *sql.DB 对象以进行 ping 测试
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("无法获取数据库实例: %v", err)
	}

	// Ping 数据库以确保连接成功
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("无法 Ping 数据库: %v", err)
	}

	// 将初始化好的实例赋值给全局变量
	DB = db
	fmt.Println("数据库连接成功并且读写分离已配置。")
}
