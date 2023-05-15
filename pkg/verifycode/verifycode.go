package verifycode

import (
	"strings"
	"sync"

	"gohub/pkg/app"
	"gohub/pkg/config"
	"gohub/pkg/helpers"
	"gohub/pkg/logger"
	"gohub/pkg/redis"
	"gohub/pkg/sms"
)

type VerifyCode struct {
	Store Store
}

var once sync.Once
var internalVerifyCode *VerifyCode

// NewVerifyCode 单例模式获取
func NewVerifyCode() *VerifyCode {
	once.Do(func() {
		internalVerifyCode = &VerifyCode{
			Store: &RedisStore{
				RedisClient: redis.Redis,
				// 增加前缀保持数据库整洁
				KeyPrefix: config.GetString("app.name") + ":verifycode:",
			},
		}
	})
	return internalVerifyCode
}

// SendSMS 发送短信验证码
func (vc *VerifyCode) SendSMS(phone string) bool {
	// 生成验证码
	code := vc.generateVerifyCode(phone)
	// 本地、API 自动测试
	if !app.IsProduction() && strings.HasPrefix(phone, config.GetString("verifycode.debug_phone_prefix")) {
		return true
	}
	// 发送短信
	return sms.NewSMS().Send(phone, sms.Message{
		Template: config.GetString("sms.aliyun.template_code"),
		Data:     map[string]string{"code": code},
	})
}

// CheckAnswer 检查用户提交的验证码是否正确，key 可以是手机号或者 email
func (vc *VerifyCode) CheckAnswer(key string, answer string) bool {
	logger.DebugJSON("验证码", "检查验证码", map[string]string{key: answer})
	// 方便开发，在非生产环境下，具备特殊前缀的手机号和 Email后缀，会直接验证成功
	if app.IsProduction() &&
		(strings.HasSuffix(key, config.GetString("verifycode.debug_phone_suffix")) ||
			strings.HasSuffix(key, config.GetString("verifycode.debug_email_suffix"))) {
		return true
	}
	return vc.Store.Verify(key, answer, false)
}

// generateVerifyCode 生成验证码，放置在 redis 中
func (vc *VerifyCode) generateVerifyCode(key string) string {
	// 生成随机码
	code := helpers.RandomNumber(config.GetInt("verifycode.code_length"))

	// 方便调试，本地使用固定的验证码
	if app.IsLocal() {
		code = config.GetString("verifycode_debug_code")
	}
	logger.DebugJSON("验证码", "生成验证码", map[string]string{key: code})
	// 将验证码 及 KEY存储到 redis 中，并设置过期时间
	vc.Store.Set(key, code)
	return code
}
