package app

import (
	"go-oj/pkg/config"
	"time"
)

func IsLocal() bool {
	return config.Get("app.env") == "local"
}

func TimeNowInTimezone() time.Time {
	Timezone, _ := time.LoadLocation(config.GetString("app.timezone"))
	return time.Now().In(Timezone)
}
