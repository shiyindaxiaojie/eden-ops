package database

import (
	"fmt"
	"time"

	"eden-ops/internal/model"
	"eden-ops/internal/pkg/logger"
	"eden-ops/pkg/config"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB 数据库实例
type DB struct {
	*gorm.DB
}

// DBVersion 数据库版本记录
type DBVersion struct {
	ID            int64  `gorm:"primaryKey"`
	Version       string `gorm:"uniqueIndex"`
	Description   string
	Script        string
	Checksum      string
	InstalledBy   string
	InstalledOn   time.Time
	ExecutionTime int
	Success       bool
}

func (DBVersion) TableName() string {
	return "sys_db_version"
}

// InitDB 初始化数据库连接
func InitDB(cfg *config.Config) (*DB, error) {
	log := logrus.New()
	log.SetFormatter(&logger.CustomFormatter{
		TimestampFormat: "2006/01/02 15:04:05.000",
	})

	// 首先创建数据库
	if err := createDatabase(cfg); err != nil {
		return nil, fmt.Errorf("创建数据库失败: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	// 配置 GORM
	gormConfig := &gorm.Config{
		Logger:                                   logger.NewGormLogger(log),
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
	}

	// 连接数据库
	gormDB, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	db := &DB{DB: gormDB}

	// 获取底层的 *sql.DB 对象
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("获取 *sql.DB 失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// 执行初始化脚本
	if err := initializeDatabase(db.DB, log); err != nil {
		return nil, fmt.Errorf("执行初始化脚本失败: %v", err)
	}

	// 初始化数据库迁移
	migrationService := NewMigrationService(db)
	if err := migrationService.Migrate("scripts/sql"); err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %v", err)
	}

	return db, nil
}

// createDatabase 创建数据库
func createDatabase(cfg *config.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
	})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	// 创建数据库
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.Database.DBName)
	return db.Exec(sql).Error
}

// initializeDatabase 初始化数据库
func initializeDatabase(db *gorm.DB, log *logrus.Logger) error {
	// 确保版本表存在
	if err := db.AutoMigrate(&model.Migration{}); err != nil {
		return fmt.Errorf("创建版本表失败: %v", err)
	}

	return nil
}
