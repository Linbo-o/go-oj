package jwt

import (
	"errors"
	jwtpkg "github.com/golang-jwt/jwt"
	"go-oj/pkg/app"
	"go-oj/pkg/config"
	"go-oj/pkg/logger"
	"time"
)

var (
	ErrTokenExpired           error = errors.New("令牌已过期")
	ErrTokenExpiredMaxRefresh error = errors.New("令牌已过最大刷新时间")
	ErrTokenMalformed         error = errors.New("请求令牌格式有误")
	ErrTokenInvalid           error = errors.New("请求令牌无效")
	ErrHeaderEmpty            error = errors.New("需要认证才能访问！")
	ErrHeaderMalformed        error = errors.New("请求头中 Authorization 格式有误")
)

// JWT 定义一个jwt对象
type JWT struct {
	// 秘钥，用以加密 JWT，读取配置信息 app.key
	SignKey []byte

	// 刷新 Token 的最大过期时间
	MaxRefresh time.Duration
}

func NewJWT() *JWT {
	return &JWT{
		SignKey:    []byte(config.GetString("app.key")),
		MaxRefresh: time.Duration(config.GetInt64("jwt.max_refresh_time")) * time.Minute,
	}
}

type JWTCustomClaims struct {
	UserIdentity string `json:"user_identity"`
	UserName     string `json:"name"`
	ExpireAtTime int64  `json:"expire_time"`

	// StandardClaims 结构体实现了 Claims 接口继承了  Valid() 方法
	jwtpkg.StandardClaims
}

func (jwt *JWT) IssueToken(userIdentity string, userName string) string {
	expireAtTime := getExpireAtTime()
	claim := JWTCustomClaims{
		UserIdentity: userIdentity,
		UserName:     userName,
		ExpireAtTime: expireAtTime,
		StandardClaims: jwtpkg.StandardClaims{
			NotBefore: app.TimeNowInTimezone().Unix(), // 签名生效时间
			IssuedAt:  app.TimeNowInTimezone().Unix(), // 首次签名时间（后续刷新 Token 不会更新）
			ExpiresAt: expireAtTime,                   // 签名过期时间
			Issuer:    config.GetString("app.name"),   // 签名颁发者
		},
	}
	token, err := jwt.createToken(claim)
	if err != nil {
		logger.LogIf(err)
		return ""
	}
	return token
}

func getExpireAtTime() int64 {
	timeNow := app.TimeNowInTimezone()
	expireTime := time.Duration(config.GetInt64("jwt.expire_time")) * time.Minute
	return timeNow.Add(expireTime).Unix()
}

func (jwt *JWT) createToken(claims JWTCustomClaims) (string, error) {
	token := jwtpkg.NewWithClaims(jwtpkg.SigningMethodHS256, claims)
	return token.SignedString(jwt.SignKey)
}
