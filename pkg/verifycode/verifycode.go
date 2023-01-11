package verifycode

import (
	"go-oj/pkg/config"
	"go-oj/pkg/helpers"
	"go-oj/pkg/logger"
	"go-oj/pkg/redis"
	"go-oj/pkg/sms"
	"sync"
)

type VerifyCode struct {
	Store *StoreRedis
}

var internalVerifyCode *VerifyCode
var once sync.Once

func NewVerifyCode() *VerifyCode {
	once.Do(func() {
		internalVerifyCode = &VerifyCode{
			Store: &StoreRedis{
				Client:    redis.Redis,
				KeyPrefix: config.GetString("app.name") + ":verifycode:",
			},
		}
	})
	return internalVerifyCode
}

func (vc *VerifyCode) SendSMS(phone string) bool {
	//1、生成验证码
	code := vc.generateVerifyCode(phone)
	//2、利用sms发送短信给手机用户
	ok := sms.NewSMS().Send(phone, sms.Message{
		Template: config.GetString("sms.aliyun.template_code"),
		Data:     map[string]string{"code": code},
	})
	return ok
}

func (vc *VerifyCode) CheckAnswer(key string, answer string) bool {
	logger.DebugJSON("验证码", "检查验证码", map[string]string{key: answer})
	return vc.Store.Verify(key, answer, false)
}

func (vc *VerifyCode) generateVerifyCode(id string) string {
	//1、生成随机验证码
	code := helpers.RandomNumber(config.GetInt("verifycode.code_length"))
	//2、将 id(这里为手机号)-验证码 键值对存入缓存
	if ok := vc.Store.Set(id, code); !ok {
		return ""
	}
	return code
}
