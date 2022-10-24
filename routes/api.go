package routes

import (
	"github.com/gin-gonic/gin"
	"gohub/app/http/controllers/api/v1/auth"
)

func RegisterAPIRoutes(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		authGroup := v1.Group("/auth")
		{
			suc := new(auth.SignupController)
			// 判断手机是否已注册
			authGroup.POST("/signup/phone/exist", suc.IsPhoneExist)
			// 判断Email是否已注册
			authGroup.POST("/signup/email/exist", suc.IsEmailExist)
			// 发送验证码
			vcc := new(auth.VerifyCodeController)
			authGroup.POST("/verify_codes/captcha", vcc.ShowCaptcha)
		}
	}
}
