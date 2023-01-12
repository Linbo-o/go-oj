package app

import (
	"go-oj/pkg/config"
	"time"
)

func IsLocal() bool {
	return config.Get("app.env") == "local"
}

func IsProduction() bool {
	return config.Get("app.env") == "production"
}

func IsTesting() bool {
	return config.Get("app.env") == "testing"
}

func TimeNowInTimezone() time.Time {
	Timezone, _ := time.LoadLocation(config.GetString("app.timezone"))
	return time.Now().In(Timezone)
}
