package user

import (
	"go-oj/app/models"
	"go-oj/pkg/database"
	"go-oj/pkg/hash"
	"go-oj/pkg/logger"
	"time"
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

func (u *UserBasic) Create() bool {
	u.CommonTimestampsField.CreatedAt = time.Now()
	u.CommonTimestampsField.UpdatedAt = time.Now()

	sql := "INSERT INTO user_basic (identity,name,password,phone,mail,pass_num,submit_num,is_admin,created_at,updated_at)" +
		"VALUES(:identity,:name,:password,:phone,:mail,:pass_num,:submit_num,:is_admin,:created_at,:updated_at)"
	_, err := database.DB.NamedExec(sql, u)
	if err != nil {
		logger.DebugString("create_user", "create", err.Error())
		return false
	}
	return true
}

func GetByPhone(phone string) UserBasic {
	u := UserBasic{}
	sql := "SELECT identity,name,phone,mail,pass_num,submit_num,is_admin From user_basic Where phone=?"
	err := database.DB.QueryRowx(sql, phone).StructScan(&u)
	if err != nil {
		logger.DebugString("user", "GetByPhone", err.Error())
		return UserBasic{}
	}
	return u
}

func GetByMulti(multiID string) UserBasic {
	u := UserBasic{}
	sql := "SELECT identity,name,phone,mail,pass_num,submit_num,is_admin From user_basic Where phone=? OR name=? OR mail=?"
	err := database.DB.QueryRowx(sql, multiID, multiID, multiID).StructScan(&u)
	if err != nil {
		logger.DebugString("user", "GetByPhone", err.Error())
		return UserBasic{}
	}
	return u
}

func Get(identity string) UserBasic {
	u := UserBasic{}
	sql := "SELECT identity,name,phone,mail,pass_num,submit_num,is_admin From user_basic Where identity=?"
	err := database.DB.QueryRowx(sql, identity).StructScan(&u)
	if err != nil {
		logger.DebugString("user", "GetByPhone", err.Error())
		return UserBasic{}
	}
	return u
}

func (u *UserBasic) GetIdentity() string {
	return u.Identity
}

func (u *UserBasic) ComparePassword(password string) bool {
	return hash.CheckPassword(u.Password, password)
}
