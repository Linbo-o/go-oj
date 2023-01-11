package sms

import (
	"go-oj/pkg/config"
	"sync"
)

type Message struct {
	Template string
	Data     map[string]string
	Content  string
}

// SMS 抽象的短信服务对象
type SMS struct {
	Driver Driver
}

var internalSMS *SMS
var once sync.Once

func NewSMS() *SMS {
	once.Do(func() {
		internalSMS = &SMS{
			Driver: &AliYun{},
		}
	})
	return internalSMS
}

func (s SMS) Send(phone string, message Message) bool {
	return s.Driver.Send(phone, message, config.GetStringMapString("sms.aliyun"))
}
