package db

import (
	"Anubis/pkg/config"
	"Anubis/pkg/logger"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	MysqlType  = "mysql"
	SQLiteType = "sqlite"
)

var DB *gorm.DB

func Initial() {
	cfg := config.GlobalConfig

	switch cfg.DBType {
	case MysqlType:
		dsn := ObtainDSN(cfg)
		DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			logger.Errorf("Open mysql connection err: %v", err)
		}

		sqlDB, err := DB.DB()
		if err != nil {
			logger.Errorf("Get sql.db err: %v", err)
		}
		// 设置连接池
		sqlDB.SetConnMaxIdleTime(10)
		sqlDB.SetMaxOpenConns(20)
	case SQLiteType:
		// TODO implement SQLite configuration
	}
}

func ObtainDSN(cfg *config.Config) string {
	connUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.DBUser, cfg.DBPass, cfg.DBAddress, cfg.DBPort, cfg.DBTable)
	if cfg.DBParams == nil {
		return connUrl
	}

	connUrl += "?"
	for key, value := range cfg.DBParams {
		param := fmt.Sprintf("%s=%s&", key, value)
		connUrl += param
	}
	return connUrl[:len(cfg.DBParams)-2]
}
