package config

import "gohub/pkg/config"

func init() {
	config.Add("log", func() map[string]interface{} {
		return map[string]interface{}{
			// 日志级别
			"level": config.Env("LOG_LEVEL", "debug"),
			// 日志类型，single单个文件，daily按日期来
			"type": config.Env("LOG_TYPE", "single"),
			// ---------------------滚动日志配置---------------------------
			// 日志文件路径
			"filename": config.Env("LOG_NAME", "storage/logs/logs.log"),
			// 每个日志保存大小,单位M
			"max_size": config.Env("LOG_MAX_SIZE", 64),
			// 保存日志文件数，0为不限制
			"max_backup": config.Env("LOG_MAX_BACKUP", 5),
			// 最多保存天数，0不限
			"max_age": config.Env("LOG_MAX_AGE", 30),
			// 是否压缩,设置为false不压缩
			"compress": config.Env("LOG_COMPRESS", false),
		}
	})
}
