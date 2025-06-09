package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageData 分页数据结构
type PageData struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

// Error 自定义错误
type Error struct {
	Message string `json:"message"`
}

// Error 实现error接口
func (e *Error) Error() string {
	return e.Message
}

// NewError 创建新的错误
func NewError(message string) *Error {
	return &Error{Message: message}
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// Failed 失败响应
func Failed(c *gin.Context, err error) {
	c.JSON(http.StatusOK, Response{
		Code:    500,
		Message: err.Error(),
	})
}

// FailedWithCode 带状态码的失败响应
func FailedWithCode(c *gin.Context, code int, err error) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: err.Error(),
	})
}

// PageSuccess 分页成功响应
func PageSuccess(c *gin.Context, list interface{}, total int64) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data: PageData{
			List:  list,
			Total: total,
		},
	})
}

// BadRequest 请求参数错误响应
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    400,
		Message: message,
	})
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    401,
		Message: message,
	})
}

// Forbidden 禁止访问响应
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Code:    403,
		Message: message,
	})
}

// NotFound 资源不存在响应
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    404,
		Message: message,
	})
}

// InternalServerError 服务器内部错误响应
func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    500,
		Message: message,
	})
}
