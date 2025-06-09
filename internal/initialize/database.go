package initialize

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

// InitDatabase 初始化数据库
func InitDatabase(db *gorm.DB) error {
	// 读取SQL文件
	sqlFile := filepath.Join("scripts", "sql", "init.sql")
	sqlContent, err := os.ReadFile(sqlFile)
	if err != nil {
		return fmt.Errorf("读取SQL文件失败: %v", err)
	}

	// 分割SQL语句
	sqlStatements := strings.Split(string(sqlContent), ";")

	// 执行SQL语句
	for _, stmt := range sqlStatements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		if err := db.Exec(stmt).Error; err != nil {
			return fmt.Errorf("执行SQL语句失败: %v", err)
		}
	}

	return nil
}

// CheckDatabaseExists 检查数据库是否存在
func CheckDatabaseExists(db *sql.DB, dbName string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = '%s')", dbName)
	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// CreateDatabase 创建数据库
func CreateDatabase(db *sql.DB, dbName string) error {
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName)
	_, err := db.Exec(query)
	return err
}
