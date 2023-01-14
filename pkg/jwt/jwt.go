package jwt

import (
	"errors"
	"github.com/gin-gonic/gin"
	jwtpkg "github.com/golang-jwt/jwt"
	"go-oj/pkg/app"
	"go-oj/pkg/config"
	"go-oj/pkg/logger"
	"strings"
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

type JWT struct {
	SignKey    []byte
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

func (jwt *JWT) ParesToken(c *gin.Context) (*JWTCustomClaims, error) {
	tokenString, err := jwt.getTokenFromHeader(c)
	if err != nil {
		logger.WarnString("jwt", "parseToken", err.Error())
		return nil, err
	}
	token, er := jwt.parseTokenString(tokenString)
	if er != nil {
		logger.Dump(token)
		return nil, er
	}

	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

func (jwt *JWT) RefreshToken(c *gin.Context) (string, error) {

	// 1. 从 Header 里获取 token
	tokenString, parseErr := jwt.getTokenFromHeader(c)
	if parseErr != nil {
		return "", parseErr
	}

	// 2. 调用 jwt 库解析用户传参的 Token
	token, err := jwt.parseTokenString(tokenString)

	// 3. 解析出错，未报错证明是合法的 Token（甚至未到过期时间）
	if err != nil {
		validationErr, ok := err.(*jwtpkg.ValidationError)
		// 满足 refresh 的条件：只是单一的报错 ValidationErrorExpired
		if !ok || validationErr.Errors != jwtpkg.ValidationErrorExpired {
			return "", err
		}
	}

	// 4. 解析 JWTCustomClaims 的数据
	claims := token.Claims.(*JWTCustomClaims)

	// 5. 检查是否过了『最大允许刷新的时间』
	x := app.TimeNowInTimezone().Add(-jwt.MaxRefresh).Unix()
	if claims.IssuedAt > x {
		// 修改过期时间
		claims.StandardClaims.ExpiresAt = getExpireAtTime()
		return jwt.createToken(*claims)
	}

	return "", ErrTokenExpiredMaxRefresh
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

func (jwt *JWT) getTokenFromHeader(c *gin.Context) (string, error) {
	header := c.GetHeader("Authorization")
	if header == "" {
		return "", errors.New("需要认证才能访问")
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("需要认证才能访问")
	}
	return parts[1], nil
}

func (jwt *JWT) parseTokenString(tokenString string) (*jwtpkg.Token, error) {
	return jwtpkg.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwtpkg.Token) (interface{}, error) {
		return jwt.SignKey, nil
	})
}
