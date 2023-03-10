package config

import (
	"github.com/spf13/cast"
	vip "github.com/spf13/viper"
	"go-oj/pkg/helpers"
	"os"
)

var viper *vip.Viper

type ConfigFunc func() map[string]interface{}

var ConfigFuncs map[string]ConfigFunc

func init() {
	viper = vip.New()
	//配置viper
	//1、配置文件类型
	viper.SetConfigType("env")
	//2、配置文件地址
	viper.AddConfigPath(".")
	//3、设置环境变量前缀，用以区分 Go 的系统环境变量
	viper.SetEnvPrefix("appenv")
	//4、读取环境变量（支持 flags）
	viper.AutomaticEnv()

	ConfigFuncs = make(map[string]ConfigFunc)
}

// InitConfig 初始化配置信息，完成对环境变量以及 config 信息的加载
func InitConfig(envSuffix string) {
	// 1. 加载环境变量
	loadEnv(envSuffix)
	// 2. 注册配置信息
	loadConfig()
}

func loadEnv(envSuffix string) {
	path := ".env"
	//默认加载.env文件，如果有传参--env=name，加载.env.name文件
	if len(envSuffix) > 0 {
		filePath := ".env." + envSuffix
		if _, err := os.Stat(filePath); err == nil {
			path = filePath
		}
	}

	viper.SetConfigName(path)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	// 监控 .env 文件，变更时重新加载
	viper.WatchConfig()
}

func loadConfig() {
	for name, f := range ConfigFuncs {
		viper.Set(name, f())
	}
}

// Add 新增配置项,新增配置函数
func Add(name string, configFn ConfigFunc) {
	ConfigFuncs[name] = configFn
}

// 包装viper的方法，采用 包名+方法名 格式调用，即viper.Method()

// Env 读取环境变量，支持默认值
func Env(envName string, defaultValue ...interface{}) interface{} {
	if len(defaultValue) > 0 {
		return internalGet(envName, defaultValue[0])
	}
	return internalGet(envName)
}

// Get 获取配置项
// 第一个参数 path 允许使用点式获取，如：app.name
// 第二个参数允许传参默认值
func Get(path string, defaultValue ...interface{}) string {
	return GetString(path, defaultValue...)
}

func internalGet(path string, defaultValue ...interface{}) interface{} {
	// config 或者环境变量不存在的情况
	if !viper.IsSet(path) || helpers.Empty(viper.Get(path)) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	return viper.Get(path)
}

// GetString 获取 String 类型的配置信息
func GetString(path string, defaultValue ...interface{}) string {
	return cast.ToString(internalGet(path, defaultValue...))
}

// GetInt 获取 Int 类型的配置信息
func GetInt(path string, defaultValue ...interface{}) int {
	return cast.ToInt(internalGet(path, defaultValue...))
}

// GetFloat64 获取 float64 类型的配置信息
func GetFloat64(path string, defaultValue ...interface{}) float64 {
	return cast.ToFloat64(internalGet(path, defaultValue...))
}

// GetInt64 获取 Int64 类型的配置信息
func GetInt64(path string, defaultValue ...interface{}) int64 {
	return cast.ToInt64(internalGet(path, defaultValue...))
}

// GetUint 获取 Uint 类型的配置信息
func GetUint(path string, defaultValue ...interface{}) uint {
	return cast.ToUint(internalGet(path, defaultValue...))
}

// GetBool 获取 Bool 类型的配置信息
func GetBool(path string, defaultValue ...interface{}) bool {
	return cast.ToBool(internalGet(path, defaultValue...))
}

// GetStringMapString 获取结构数据
func GetStringMapString(path string) map[string]string {
	return viper.GetStringMapString(path)
}
