package config

import "go-oj/pkg/config"

func init() {
	config.Add("jwt", func() map[string]interface{} {
		return map[string]interface{}{

			// 过期时间，单位是分钟
			"expire_time": config.Env("JWT_EXPIRE_TIME", 12000),

			// 允许刷新时间，单位分钟，86400 为两个月，从 Token 的签名时间算起
			"max_refresh_time": config.Env("JWT_MAX_REFRESH_TIME", 86400),
		}
	})
}
