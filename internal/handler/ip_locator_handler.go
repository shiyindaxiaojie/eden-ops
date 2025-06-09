package handler

import (
	"eden-ops/pkg/iplocator"
	"eden-ops/pkg/response"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// IPLocatorHandler IP定位处理器
type IPLocatorHandler struct {
	ipLocator *iplocator.IPLocator
	logger    *logrus.Logger
}

// NewIPLocatorHandler 创建IP定位处理器
func NewIPLocatorHandler(accessKey, secretKey, baseURL string, logger *logrus.Logger) *IPLocatorHandler {
	return &IPLocatorHandler{
		ipLocator: iplocator.NewIPLocator(accessKey, secretKey, baseURL),
		logger:    logger,
	}
}

// Locate 定位IP地址
func (h *IPLocatorHandler) Locate(c *gin.Context) {
	ip := c.Param("ip")
	if ip == "" {
		response.Failed(c, fmt.Errorf("IP地址不能为空"))
		return
	}

	result, err := h.ipLocator.Locate(ip)
	if err != nil {
		h.logger.WithError(err).WithField("ip", ip).Error("IP定位失败")
		response.Failed(c, err)
		return
	}

	response.Success(c, result)
}
