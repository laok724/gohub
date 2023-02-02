package middlewares

import (
	"net"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gohub/pkg/logger"
	"gohub/pkg/response"
)

// Recovery使用zap.Error()来记录Panic和call stack
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取用户请求信息
				httpRequest, _ := httputil.DumpRequest(c.Request, true)
				// 链接中断，客户端中断连接是正常行为，不需要记录堆栈信息
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						errStr := strings.ToLower(se.Error())
						if strings.Contains(errStr, "broken pipe") || strings.Contains(errStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				// 连接中断情况
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error))
					c.Abort()
					// 连接断开，无法写状态码
					return
				}
				// 不是连接中断，记录堆栈信息
				logger.Error("recovery from paninc",
					zap.Time("time", time.Now()),
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
					zap.Stack("stackrace"),
				)
				// 返回500状态码

				response.Abort500(c)
			}
		}()
		c.Next()
	}
}
