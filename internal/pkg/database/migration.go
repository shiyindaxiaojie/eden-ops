package database

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"eden-ops/internal/model"

	"gorm.io/gorm"
)

// Migration 迁移记录
type Migration = model.Migration

// MigrationService 迁移服务
type MigrationService struct {
	db *gorm.DB
}

// NewMigrationService 创建迁移服务
func NewMigrationService(db *DB) *MigrationService {
	return &MigrationService{db: db.DB}
}

// splitSQLStatements 分割SQL语句
func splitSQLStatements(content string) []string {
	// 移除注释
	content = regexp.MustCompile(`--.*$`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`/\*[\s\S]*?\*/`).ReplaceAllString(content, "")

	// 分割SQL语句
	statements := make([]string, 0)
	var currentStmt strings.Builder
	var inString bool
	var stringChar rune
	var escaped bool

	for _, char := range content {
		switch {
		case escaped:
			currentStmt.WriteRune(char)
			escaped = false
		case char == '\\':
			currentStmt.WriteRune(char)
			escaped = true
		case inString && char == stringChar:
			currentStmt.WriteRune(char)
			inString = false
		case !inString && (char == '\'' || char == '"'):
			currentStmt.WriteRune(char)
			inString = true
			stringChar = char
		case !inString && char == ';':
			stmt := strings.TrimSpace(currentStmt.String())
			if stmt != "" {
				statements = append(statements, stmt)
			}
			currentStmt.Reset()
		default:
			currentStmt.WriteRune(char)
		}
	}

	// 添加最后一个语句（如果有）
	lastStmt := strings.TrimSpace(currentStmt.String())
	if lastStmt != "" {
		statements = append(statements, lastStmt)
	}

	return statements
}

// parseScriptVersion 解析脚本版本信息
func parseScriptVersion(filename string) (version, description string, err error) {
	// Flyway命名格式: V1.0.0__Description.sql
	pattern := regexp.MustCompile(`^V(\d+\.\d+\.\d+)__(.+)\.sql$`)
	matches := pattern.FindStringSubmatch(filename)
	if len(matches) != 3 {
		return "", "", fmt.Errorf("无效的脚本文件名格式: %s", filename)
	}
	return matches[1], matches[2], nil
}

// isVersionTableExists 检查版本表是否存在
func (s *MigrationService) isVersionTableExists() (bool, error) {
	return s.db.Migrator().HasTable(&Migration{}), nil
}

// createVersionTable 创建版本表
func (s *MigrationService) createVersionTable() error {
	return s.db.AutoMigrate(&Migration{})
}

// getExecutedVersions 获取已执行的版本记录
func (s *MigrationService) getExecutedVersions() (map[string]*Migration, error) {
	var migrations []*Migration
	if err := s.db.Where("success = ?", true).Find(&migrations).Error; err != nil {
		return nil, fmt.Errorf("获取已执行版本记录失败: %v", err)
	}

	versionMap := make(map[string]*Migration)
	for _, m := range migrations {
		versionMap[m.Version] = m
	}
	return versionMap, nil
}

// executeSQLStatements 执行SQL语句
func (s *MigrationService) executeSQLStatements(statements []string) error {
	for _, stmt := range statements {
		if strings.TrimSpace(stmt) == "" {
			continue
		}
		if err := s.db.Exec(stmt).Error; err != nil {
			return fmt.Errorf("执行SQL语句失败: %v\nSQL: %s", err, stmt)
		}
	}
	return nil
}

// Migrate 执行数据库迁移
func (s *MigrationService) Migrate(scriptDir string) error {
	// 检查版本表是否存在
	exists, err := s.isVersionTableExists()
	if err != nil {
		return fmt.Errorf("检查版本表是否存在失败: %v", err)
	}

	// 如果版本表不存在，创建它
	if !exists {
		if err := s.createVersionTable(); err != nil {
			return fmt.Errorf("创建版本表失败: %v", err)
		}
	}

	// 获取已执行的版本记录
	executedVersions, err := s.getExecutedVersions()
	if err != nil {
		return err
	}

	// 获取所有SQL文件
	files, err := filepath.Glob(filepath.Join(scriptDir, "V*.sql"))
	if err != nil {
		return fmt.Errorf("读取SQL文件失败: %v", err)
	}

	// 按文件名排序
	sort.Strings(files)

	// 遍历所有SQL文件
	for _, file := range files {
		filename := filepath.Base(file)

		// 解析版本信息
		version, description, err := parseScriptVersion(filename)
		if err != nil {
			return err
		}

		// 检查是否已经执行过
		if executed, ok := executedVersions[version]; ok {
			// 读取SQL文件内容并验证校验和
			content, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("读取SQL文件失败: %v", err)
			}

			checksum := fmt.Sprintf("%x", md5.Sum(content))
			if checksum != executed.Checksum {
				return fmt.Errorf("脚本文件 %s 已被修改", filename)
			}
			continue
		}

		// 读取SQL文件内容
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("读取SQL文件失败: %v", err)
		}

		// 开始事务
		tx := s.db.Begin()
		if tx.Error != nil {
			return fmt.Errorf("开始事务失败: %v", tx.Error)
		}

		// 分割并执行SQL语句
		startTime := time.Now()
		statements := splitSQLStatements(string(content))
		if err := s.executeSQLStatements(statements); err != nil {
			tx.Rollback()
			return fmt.Errorf("执行SQL失败: %v", err)
		}

		// 记录执行结果
		migration := &Migration{
			Version:       version,
			Description:   description,
			Script:        filename,
			Checksum:      fmt.Sprintf("%x", md5.Sum(content)),
			InstalledBy:   "system",
			InstalledOn:   time.Now(),
			ExecutionTime: int(time.Since(startTime).Milliseconds()),
			Success:       true,
		}

		if err := tx.Create(migration).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("记录版本信息失败: %v", err)
		}

		// 提交事务
		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("提交事务失败: %v", err)
		}
	}

	return nil
}
