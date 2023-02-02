package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gohub/pkg/logger"
	"gorm.io/gorm"
)

// 响应处理工具

// JSON响应200和JSON数据

func JSON(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// JSON响应200和预设操作成功的JSON数据
// 执行某个 没有具体返回数据的变更操作成功后调用，如删除、修改秘密、修改手机号

func Success(c *gin.Context) {
	JSON(c, gin.H{
		"success": true,
		"message": "操作成功！",
	})
}

// data响应200和带data键的JSON数据
// 执行更新成功后调用，比如更新话题，成功后返回已更新的话题

func Data(c *gin.Context, data interface{}) {
	JSON(c, gin.H{
		"success": true,
		"data":    data,
	})
}

// Created响应201和data键的JSON数据
// 执行更新成功后调用，比如更新话题，成功后返回已更新的话题

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    data,
	})
}

// CreatedJSON 响应201和JSON数据
func CreatedJSON(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

// Abort404响应404，未传参msg使用默认消息

func Abort404(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
		"message": defaultMessage("数据不存在，请确认请求正确", msg...),
	})
}

// Abort403响应403，未传参msg使用默认消息

func Abort403(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"message": defaultMessage("权限不足，请确认是否有对应权限", msg...),
	})
}

// Abort500响应500，未传参msg使用默认消息

func Abort500(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"message": defaultMessage("服务内部错误，请稍后再试", msg...),
	})
}

// BadRequest 响应400，传参err对象，未传参msg使用默认消息
// 解析用户请求，请求的格式或者方法不符合预期时调用

func BadRequest(c *gin.Context, err error, msg ...string) {
	logger.LogIf(err)

	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"message": defaultMessage("请求解析错误，请确认请求格式是否正确。上传文件使用multipart头，参数使用json格式", msg...),
		"error":   err.Error(),
	})
}

// Error 响应404或422，未传参msg使用默认消息
// 处理请求时出现错误err,会返回error信息，如登录错误，找不到ID对应的Model
func Error(c *gin.Context, err error, msg ...string) {
	logger.LogIf(err)
	// error 类型未数据库未找到内容
	if err == gorm.ErrRecordNotFound {
		Abort404(c)
		return
	}
	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
		"message": defaultMessage("请求处理失败，请查看error信息", msg...),
		"error":   err.Error(),
	})
}

// ValidationError 处理表单验证不通过的错误，返回JSON

func ValidationError(c *gin.Context, errors map[string][]string) {
	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
		"message": "请求验证不通过，具体查看errors",
		"errors":  errors,
	})
}

// Unauthorized响应401 未传参msg使用默认消息，登录失败，jwt解析失败

func Unauthorized(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"message": defaultMessage("请求解析错误，请确认请求格式是否正确。", msg...),
	})
}

func defaultMessage(defaultMsg string, msg ...string) (message string) {
	if len(msg) > 0 {
		message = msg[0]
	} else {
		message = defaultMsg
	}
	return
}
