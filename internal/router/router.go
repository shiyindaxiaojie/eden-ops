package router

import (
	"eden-ops/internal/handler"
	middleware "eden-ops/internal/pkg/middleware"
	"eden-ops/pkg/auth"

	"github.com/gin-gonic/gin"
)

// NewRouter 创建路由
func NewRouter(
	jwtAuth *auth.JWTAuth,
	cloudAccountHandler *handler.CloudAccountHandler,
	cloudProviderHandler *handler.CloudProviderHandler,
	databaseConfigHandler *handler.DatabaseConfigHandler,
	serverConfigHandler *handler.ServerConfigHandler,
	k8sConfigHandler *handler.K8sConfigHandler,
	userHandler *handler.UserHandler,
	roleHandler *handler.RoleHandler,
	menuHandler *handler.MenuHandler,
	authHandler *handler.AuthHandler,
) *gin.Engine {
	// 使用自定义Logger替代默认Logger
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// 注册中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Cors())

	// API路由组
	api := r.Group("/api/v1")

	// 公共路由
	{
		// 用户认证
		api.POST("/login", authHandler.Login)
		api.POST("/logout", authHandler.Logout)
	}

	// 需要认证的路由
	auth := api.Group("/", middleware.JWT(jwtAuth))
	{
		// 获取当前用户信息
		auth.GET("/users/info", userHandler.GetUserInfo)
		// 用户管理
		auth.GET("/users", userHandler.List)
		auth.GET("/users/:id", userHandler.Get)
		auth.POST("/users", userHandler.Create)
		auth.PUT("/users/:id", userHandler.Update)
		auth.DELETE("/users/:id", userHandler.Delete)
		auth.GET("/users/:id/roles", userHandler.GetRoles)
		auth.PUT("/users/:id/roles", userHandler.AssignRoles)

		// 角色管理
		auth.GET("/roles", roleHandler.List)
		auth.GET("/roles/:id", roleHandler.Get)
		auth.POST("/roles", roleHandler.Create)
		auth.PUT("/roles/:id", roleHandler.Update)
		auth.DELETE("/roles/:id", roleHandler.Delete)
		auth.PUT("/roles/:id/menus", roleHandler.AssignMenus)

		// 菜单管理
		auth.GET("/menus", menuHandler.List)
		auth.GET("/menus/:id", menuHandler.Get)
		auth.POST("/menus", menuHandler.Create)
		auth.PUT("/menus/:id", menuHandler.Update)
		auth.DELETE("/menus/:id", menuHandler.Delete)

		// 云账号管理
		auth.GET("/cloud-accounts", cloudAccountHandler.List)
		auth.GET("/cloud-accounts/:id", cloudAccountHandler.Get)
		auth.POST("/cloud-accounts", cloudAccountHandler.Create)
		auth.PUT("/cloud-accounts/:id", cloudAccountHandler.Update)
		auth.DELETE("/cloud-accounts/:id", cloudAccountHandler.Delete)
		auth.POST("/cloud-accounts/test", cloudAccountHandler.TestConnection)

		// 数据库配置管理
		auth.GET("/database-configs", databaseConfigHandler.List)
		auth.GET("/database-configs/:id", databaseConfigHandler.Get)
		auth.POST("/database-configs", databaseConfigHandler.Create)
		auth.PUT("/database-configs/:id", databaseConfigHandler.Update)
		auth.DELETE("/database-configs/:id", databaseConfigHandler.Delete)
		auth.POST("/database-configs/test", databaseConfigHandler.TestConnection)

		// 服务器配置管理
		auth.GET("/server-configs", serverConfigHandler.List)
		auth.GET("/server-configs/:id", serverConfigHandler.Get)
		auth.POST("/server-configs", serverConfigHandler.Create)
		auth.PUT("/server-configs/:id", serverConfigHandler.Update)
		auth.DELETE("/server-configs/:id", serverConfigHandler.Delete)
		auth.POST("/server-configs/test", serverConfigHandler.TestConnection)

		// Kubernetes配置管理
		auth.GET("/k8s-configs", k8sConfigHandler.List)
		auth.GET("/k8s-configs/:id", k8sConfigHandler.Get)
		auth.POST("/k8s-configs", k8sConfigHandler.Create)
		auth.PUT("/k8s-configs/:id", k8sConfigHandler.Update)
		auth.DELETE("/k8s-configs/:id", k8sConfigHandler.Delete)
		auth.POST("/k8s-configs/test", k8sConfigHandler.TestConnection)

		// 基础设施路由组
		infrastructure := auth.Group("/infrastructure")
		{
			// 云厂商
			infrastructure.GET("/cloud-providers", cloudProviderHandler.List)
			infrastructure.GET("/cloud-providers/:id", cloudProviderHandler.Get)
			infrastructure.POST("/cloud-providers", cloudProviderHandler.Create)
			infrastructure.PUT("/cloud-providers/:id", cloudProviderHandler.Update)
			infrastructure.DELETE("/cloud-providers/:id", cloudProviderHandler.Delete)

			// 云账号
			infrastructure.GET("/cloud-accounts", cloudAccountHandler.List)
			infrastructure.GET("/cloud-accounts/:id", cloudAccountHandler.Get)
			infrastructure.POST("/cloud-accounts", cloudAccountHandler.Create)
			infrastructure.PUT("/cloud-accounts/:id", cloudAccountHandler.Update)
			infrastructure.DELETE("/cloud-accounts/:id", cloudAccountHandler.Delete)

			// 数据库
			infrastructure.GET("/database", databaseConfigHandler.List)
			infrastructure.GET("/database/:id", databaseConfigHandler.Get)
			infrastructure.POST("/database", databaseConfigHandler.Create)
			infrastructure.PUT("/database/:id", databaseConfigHandler.Update)
			infrastructure.DELETE("/database/:id", databaseConfigHandler.Delete)

			// 服务器
			infrastructure.GET("/server", serverConfigHandler.List)
			infrastructure.GET("/server/:id", serverConfigHandler.Get)
			infrastructure.POST("/server", serverConfigHandler.Create)
			infrastructure.PUT("/server/:id", serverConfigHandler.Update)
			infrastructure.DELETE("/server/:id", serverConfigHandler.Delete)

			// Kubernetes
			infrastructure.GET("/kubernetes", k8sConfigHandler.List)
			infrastructure.GET("/kubernetes/:id", k8sConfigHandler.Get)
			infrastructure.POST("/kubernetes", k8sConfigHandler.Create)
			infrastructure.PUT("/kubernetes/:id", k8sConfigHandler.Update)
			infrastructure.DELETE("/kubernetes/:id", k8sConfigHandler.Delete)
		}
	}

	return r
}
