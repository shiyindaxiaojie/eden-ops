package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPage 获取页码
func GetPage(c *gin.Context) int {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page <= 0 {
		page = 1
	}
	return page
}

// GetPageSize 获取每页数量
func GetPageSize(c *gin.Context) int {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if pageSize <= 0 {
		pageSize = 10
	}
	return pageSize
}
