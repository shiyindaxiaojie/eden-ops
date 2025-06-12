package main

import (
	"context"
	"eden-ops/internal/handler"
	"eden-ops/internal/pkg/database"
	"eden-ops/internal/repository"
	"eden-ops/internal/router"
	"eden-ops/internal/service"
	"eden-ops/internal/task"
	"eden-ops/pkg/auth"
	"eden-ops/pkg/config"
	"eden-ops/pkg/logger"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 获取本机IP地址
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	// 优先查找内网地址
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				// 检查是否为内网地址
				ip := ipnet.IP.To4()
				// 排除链路本地地址 169.254.x.x
				if ip[0] != 169 || ip[1] != 254 {
					// 优先使用内网地址
					// 10.x.x.x, 172.16.x.x - 172.31.x.x, 192.168.x.x
					if ip[0] == 10 ||
						(ip[0] == 172 && ip[1] >= 16 && ip[1] <= 31) ||
						(ip[0] == 192 && ip[1] == 168) {
						return ipnet.IP.String()
					}
				}
			}
		}
	}

	// 如果都没找到，使用回环地址
	return "127.0.0.1"
}

// 统计API接口数量
func countAPIEndpoints(r *gin.Engine) int {
	count := 0
	for _, routeInfo := range r.Routes() {
		if strings.HasPrefix(routeInfo.Path, "/api/") {
			count++
		}
	}
	return count
}

// 获取API接口列表
func getAPIEndpoints(r *gin.Engine) []string {
	var endpoints []string
	for _, routeInfo := range r.Routes() {
		if strings.HasPrefix(routeInfo.Path, "/api/") {
			endpoints = append(endpoints, fmt.Sprintf("%s %s", routeInfo.Method, routeInfo.Path))
		}
	}
	return endpoints
}

// 获取初始化脚本信息
func getInitScripts() []string {
	var scripts []string

	// 检查SQL脚本
	sqlDir := "scripts/sql"
	files, err := os.ReadDir(sqlDir)
	if err == nil {
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
				scripts = append(scripts, filepath.Join(sqlDir, file.Name()))
			}
		}
	}

	// 检查其他初始化脚本
	initDir := "init"
	files, err = os.ReadDir(initDir)
	if err == nil {
		for _, file := range files {
			if !file.IsDir() {
				scripts = append(scripts, filepath.Join(initDir, file.Name()))
			}
		}
	}

	return scripts
}

func main() {
	// 记录开始时间
	startTime := time.Now()

	// 加载配置
	logger.Info("正在加载配置...")

	// 列出所有配置文件
	configFiles := []string{"configs/config.yaml"}
	envFile := ".env"
	if _, err := os.Stat(envFile); err == nil {
		configFiles = append(configFiles, envFile)
	}

	logger.Info("加载配置文件: %s", strings.Join(configFiles, ", "))

	cfg, err := config.LoadConfigFromYAML("configs/config.yaml")
	if err != nil {
		logger.Error("加载配置文件失败: %v", err)
		os.Exit(1)
	}

	logger.Info("配置加载成功, 服务端口: %d, 数据库: %s@%s:%d/%s",
		cfg.Server.Port, cfg.Database.Username, cfg.Database.Host,
		cfg.Database.Port, cfg.Database.DBName)

	// 初始化日志
	logger.Info("正在初始化日志...")
	if err := logger.Init(cfg.Log); err != nil {
		logger.Error("初始化日志失败: %v", err)
		os.Exit(1)
	}
	logger.Info("日志初始化成功")

	// 获取日志实例
	logInstance := logger.GetLogger()
	if logInstance == nil {
		logger.Error("获取日志实例失败，日志未正确初始化")
		os.Exit(1)
	}

	// 初始化脚本信息
	initScripts := getInitScripts()
	logger.Info("加载脚本:\n%s", strings.Join(initScripts, "\n"))

	// 初始化数据库
	logger.Info("初始化数据库...")
	dbInstance, err := database.InitDB(cfg)
	if err != nil {
		logger.Error("初始化数据库失败: %v", err)
		os.Exit(1)
	}
	logger.Info("数据库初始化成功")

	// 获取GORM DB实例
	db := dbInstance.DB

	// 初始化JWT
	jwtAuth := auth.NewJWTAuth(cfg.JWT.Secret, cfg.JWT.Expire)

	// 初始化仓库
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	cloudAccountRepo := repository.NewCloudAccountRepository(db)
	cloudProviderRepo := repository.NewCloudProviderRepository(db)
	databaseConfigRepo := repository.NewDatabaseConfigRepository(db)
	serverConfigRepo := repository.NewServerConfigRepository(db)
	k8sConfigRepo := repository.NewK8sConfigRepository(db)
	k8sWorkloadRepo := repository.NewK8sWorkloadRepository(db)

	// 初始化服务
	userService := service.NewUserService(userRepo, jwtAuth)
	roleService := service.NewRoleService(roleRepo)
	menuService := service.NewMenuService(menuRepo)
	cloudAccountService := service.NewCloudAccountService(cloudAccountRepo)
	cloudProviderService := service.NewCloudProviderService(cloudProviderRepo)
	databaseConfigService := service.NewDatabaseConfigService(databaseConfigRepo)
	serverConfigService := service.NewServerConfigService(serverConfigRepo)
	k8sWorkloadService := service.NewK8sWorkloadService(k8sWorkloadRepo)
	k8sConfigService := service.NewK8sConfigService(k8sConfigRepo, k8sWorkloadService)

	// 创建日志记录器
	logrusLogger := logrus.New()
	logrusLogger.SetFormatter(&logger.CustomFormatter{})

	// 启动K8s同步任务
	k8sSyncTask := task.NewK8sSyncTask(db, k8sConfigService, logrusLogger)
	syncCtx, syncCancel := context.WithCancel(context.Background())
	defer syncCancel()

	go func() {
		if err := k8sSyncTask.Start(syncCtx); err != nil {
			logger.Error("启动K8s同步任务失败: %v", err)
		}
	}()

	// 初始化处理器
	userHandler := handler.NewUserHandler(userService, jwtAuth)
	roleHandler := handler.NewRoleHandler(roleService)
	menuHandler := handler.NewMenuHandler(menuService)
	authHandler := handler.NewAuthHandler(userService, jwtAuth)
	cloudAccountHandler := handler.NewCloudAccountHandler(cloudAccountService)
	cloudProviderHandler := handler.NewCloudProviderHandler(cloudProviderService)
	databaseConfigHandler := handler.NewDatabaseConfigHandler(databaseConfigService)
	serverConfigHandler := handler.NewServerConfigHandler(serverConfigService, logrusLogger)
	k8sConfigHandler := handler.NewK8sConfigHandler(k8sConfigService)

	// 初始化路由
	logger.Info("初始化路由...")
	r := router.NewRouter(
		jwtAuth,
		cloudAccountHandler,
		cloudProviderHandler,
		databaseConfigHandler,
		serverConfigHandler,
		k8sConfigHandler,
		userHandler,
		roleHandler,
		menuHandler,
		authHandler,
	)

	// 列出所有API接口
	apiEndpoints := getAPIEndpoints(r)
	for _, endpoint := range apiEndpoints {
		logger.Debug("注册路由: %s", endpoint)
	}

	// 统计API接口
	apiCount := countAPIEndpoints(r)
	logger.Info("路由初始化成功，共有 %d 个 API 接口", apiCount)

	// 启动服务器
	localIP := getLocalIP()
	// 监听所有接口，以便前端代理可以连接
	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Server.Port)
	logger.Info("正在启动服务器，监听地址: %s (本机IP: %s)", addr, localIP)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 在一个单独的goroutine中启动服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("服务器启动失败: %v", err)
			os.Exit(1)
		}
	}()

	// 计算启动耗时
	elapsedTime := time.Since(startTime)

	logger.Info("服务启动成功，耗时 %.3f 毫秒", float64(elapsedTime.Microseconds())/1000.0)

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	// kill (无参数) 默认发送 syscall.SIGTERM
	// kill -2 发送 syscall.SIGINT
	// kill -9 发送 syscall.SIGKILL，但无法被捕获，所以不需要添加
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("正在关闭服务器...")

	// 创建一个5秒的上下文用于超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("服务器强制关闭: %v", err)
		os.Exit(1)
	}

	logger.Info("服务器已优雅关闭")
}
