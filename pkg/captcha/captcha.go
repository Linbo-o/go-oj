package captcha

import (
	"github.com/mojocn/base64Captcha"
	"go-oj/pkg/config"
	"go-oj/pkg/redis"
	"sync"
)

type Captcha struct {
	Base64Captcha *base64Captcha.Captcha
}

var internalCaptcha *Captcha
var once sync.Once

// NewCaptcha 创建captcha单例，只会初始化一次
func NewCaptcha() *Captcha {
	once.Do(func() {
		internalCaptcha = &Captcha{}

		store := &StoreRedis{
			Client:    redis.Redis,
			KeyPrefix: config.GetString("app.name") + ":captcha:",
		}

		//数字图片验证码
		driver := &base64Captcha.DriverDigit{
			Height:   config.GetInt("captcha.height"),
			Width:    config.GetInt("captcha.width"),
			Length:   config.GetInt("captcha.length"),
			MaxSkew:  config.GetFloat64("captcha.maxskew"),
			DotCount: config.GetInt("captcha.dotcount"),
		}

		internalCaptcha.Base64Captcha = base64Captcha.NewCaptcha(driver, store)
	})

	return internalCaptcha
}

func (cp *Captcha) GenCaptcha() (id string, b64s string, err error) {
	return cp.Base64Captcha.Generate()
}

func (cp *Captcha) Verify(id, answer string, clear bool) bool {
	return cp.Base64Captcha.Verify(id, answer, clear)
}
