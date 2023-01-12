package validators

import (
	"fmt"
	"github.com/thedevsaddam/govalidator"
	"go-oj/pkg/database"
	"strings"
)

func init() {
	//"not_exist:users,phone"表示在users表里面的phone不存在
	//跟两个参数，表名和字段名
	govalidator.AddCustomRule("not_exists", func(field string, rule string, message string, value interface{}) error {
		param := strings.Split(strings.TrimPrefix(rule, "not_exists:"), ",")
		tableName := param[0]
		fieldName := param[1]
		v := value.(string)
		if ok := database.IsExist(tableName, fieldName, v); ok {
			return fmt.Errorf("%s 已经存在", field)
		}
		return nil
	})
}
