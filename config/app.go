package config

import "gohub/pkg/config"

func init() {
	config.Add("app", func() map[string]interface{} {
		return map[string]interface{}{
			// 应用名称
			"name": config.Env("APP_NAME", "Gohub"),
			// 当前环境
			"debug": config.Env("APP_DEBUG", false),

			// 应用服务端口
			"port": config.Env("APP_PORT", 3000),
			// 加密会话
			"key": config.Env("APP_KEY", "33446a9dcf9ea060a0a6532b166da32f304af0de"),
			// 生成链接
			"url": config.Env("APP_URL", "http://localhost:3000"),
			// 时区
			"timezone": config.Env("APP_TIMEZONE", "Asia/Shanghai"),
		}
	})
}
