package user

import (
	"go-oj/app/models"
)

type UserBasic struct {
	models.BaseModel

	Identity  string `db:"identity" json:"identity"`     // 用户的唯一标识
	Name      string `db:"name" json:"name"`             // 用户名
	Password  string `db:"password" json:"password"`     // 密码
	Phone     string `db:"phone" json:"phone"`           // 手机号
	Mail      string `db:"mail" json:"mail"`             // 邮箱
	PassNum   int64  `db:"pass_num" json:"pass_num"`     // 通过的次数
	SubmitNum int64  `db:"submit_num" json:"submit_num"` // 提交次数
	IsAdmin   int    `db:"is_admin" json:"is_admin"`     // 是否是管理员【0-否，1-是】

	models.CommonTimestampsField
}
