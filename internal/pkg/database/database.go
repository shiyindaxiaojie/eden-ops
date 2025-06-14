package database

import (
	"fmt"
	"path/filepath"
	"runtime"
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
	// 使用项目日志系统
	var pkgLogger = struct {
		Info  func(string, ...interface{})
		Error func(string, ...interface{})
	}{
		Info: func(format string, args ...interface{}) {
			timestamp := time.Now().Format("2006/01/02 15:04:05.000")
			_, file, line, _ := runtime.Caller(1)
			fileName := filepath.Base(file)
			fmt.Printf("%s %s:%d %s\n", timestamp, fileName, line, fmt.Sprintf(format, args...))
		},
		Error: func(format string, args ...interface{}) {
			timestamp := time.Now().Format("2006/01/02 15:04:05.000")
			_, file, line, _ := runtime.Caller(1)
			fileName := filepath.Base(file)
			fmt.Printf("%s %s:%d ERROR: %s\n", timestamp, fileName, line, fmt.Sprintf(format, args...))
		},
	}

	pkgLogger.Info("开始初始化数据库连接")

	logObj := logrus.New()
	logObj.SetFormatter(&logger.CustomFormatter{
		TimestampFormat: "2006/01/02 15:04:05.000",
	})

	// 首先创建数据库
	pkgLogger.Info("尝试创建数据库: %s", cfg.Database.DBName)
	if err := createDatabase(cfg); err != nil {
		pkgLogger.Error("创建数据库失败: %v", err)
		return nil, fmt.Errorf("创建数据库失败: %v", err)
	}
	pkgLogger.Info("数据库创建成功或已存在")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)
	// 创建安全的DSN用于日志显示（隐藏密码）
	safeDSN := fmt.Sprintf("%s:***@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true",
		cfg.Database.Username,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)
	pkgLogger.Info("数据库连接DSN: %s", safeDSN)

	// 配置 GORM
	gormConfig := &gorm.Config{
		Logger:                                   logger.NewGormLogger(logObj),
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
	}

	// 连接数据库
	pkgLogger.Info("尝试连接数据库")
	gormDB, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		pkgLogger.Error("连接数据库失败: %v", err)
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}
	pkgLogger.Info("数据库连接成功")

	db := &DB{DB: gormDB}

	// 获取底层的 *sql.DB 对象
	sqlDB, err := gormDB.DB()
	if err != nil {
		pkgLogger.Error("获取 *sql.DB 失败: %v", err)
		return nil, fmt.Errorf("获取 *sql.DB 失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)
	pkgLogger.Info("数据库连接池参数设置完成")

	// 执行初始化脚本
	pkgLogger.Info("开始执行数据库初始化")
	if err := initializeDatabase(db.DB, logObj); err != nil {
		pkgLogger.Error("执行初始化脚本失败: %v", err)
		return nil, fmt.Errorf("执行初始化脚本失败: %v", err)
	}
	pkgLogger.Info("数据库初始化完成")

	// 初始化数据库迁移
	pkgLogger.Info("开始执行数据库迁移")
	migrationService := NewMigrationService(db)
	if err := migrationService.Migrate("scripts/sql"); err != nil {
		pkgLogger.Error("数据库迁移失败: %v", err)
		return nil, fmt.Errorf("数据库迁移失败: %v", err)
	}
	pkgLogger.Info("数据库迁移完成")

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
