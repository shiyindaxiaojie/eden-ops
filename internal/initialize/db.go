package initialize

import (
	"eden-ops/pkg/config"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB 初始化数据库连接
func InitDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	return db, nil
}
