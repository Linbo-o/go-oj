package sms

import (
	"encoding/json"
	aliyun "github.com/KenmyZhang/aliyun-communicate"
	"go-oj/pkg/logger"
)

type AliYun struct {
}

// Send 调用短信服务商接口向phone用户发送短信
func (s *AliYun) Send(phone string, message Message, config map[string]string) bool {
	//获取服务端对象
	smsClient := aliyun.New("http://dysmsapi.aliyuncs.com/")

	//将结构转换为json字符串
	messageTemplate, err := json.Marshal(message.Data)
	if err != nil {
		logger.ErrorString("短信[阿里云]", "解析绑定错误", err.Error())
		return false
	}

	logger.DebugJSON("短信[阿里云]", "配置信息", config)

	//填写获取短信服务必要的信息
	result, err := smsClient.Execute(
		config["access_key_id"],
		config["access_key_secret"],
		phone,
		config["sign_name"],
		message.Template,
		string(messageTemplate),
	)
	logger.DebugJSON("短信[阿里云]", "请求内容", smsClient.Request)
	logger.DebugJSON("短信[阿里云]", "接口响应", result)
	if err != nil {
		logger.ErrorString("短信[阿里云]", "发信失败", err.Error())
		return false
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		logger.ErrorString("短信[阿里云]", "解析响应 JSON 错误", err.Error())
		return false
	}
	if result.IsSuccessful() {
		logger.DebugString("短信[阿里云]", "发信成功", "")
		return true
	} else {
		logger.ErrorString("短信[阿里云]", "服务商返回错误", string(resultJSON))
		return false
	}
}
