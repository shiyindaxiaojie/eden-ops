package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 响应结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// PageData represents paginated data
type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Failed 失败响应
func Failed(c *gin.Context, err error) {
	c.JSON(http.StatusOK, Response{
		Code:    -1,
		Message: err.Error(),
		Data:    nil,
	})
}

// FailedWithMessage 失败响应（只有消息）
func FailedWithMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    -1,
		Message: message,
		Data:    nil,
	})
}

// FailedWithCode 失败响应（带错误码）
func FailedWithCode(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// PageSuccess 分页成功响应
func PageSuccess(c *gin.Context, items interface{}, total int64) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: PageData{
			List:     items,
			Total:    total,
			Page:     1,
			PageSize: 10,
		},
	})
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, message string) {
	FailedWithCode(c, 401, message)
}

// BadRequest 请求参数错误响应
func BadRequest(c *gin.Context, message string) {
	FailedWithCode(c, 400, message)
}

// Forbidden 禁止访问响应
func Forbidden(c *gin.Context, message string) {
	FailedWithCode(c, 403, message)
}

// NotFound 资源不存在响应
func NotFound(c *gin.Context, message string) {
	FailedWithCode(c, 404, message)
}

// InternalServerError 服务器内部错误响应
func InternalServerError(c *gin.Context, message string) {
	FailedWithCode(c, 500, message)
}
