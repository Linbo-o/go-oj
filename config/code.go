package config

import "go-oj/pkg/config"

func init() {
	config.Add("code", func() map[string]interface{} {
		return map[string]interface{}{
			"dir_path": config.Env("DIR_PATH", "storage/code"),
		}
	})
}
