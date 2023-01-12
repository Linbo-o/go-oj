package problem

import (
	"go-oj/app/models"
	pc "go-oj/app/models/problem-category"
	"go-oj/app/models/testcase"
	"go-oj/pkg/database"
	"go-oj/pkg/logger"
	"time"
)

type ProblemBasic struct {
	models.BaseModel

	Identity   string                `json:"identity" db:"identity"`
	Title      string                `json:"title"    db:"title"`
	Content    string                `json:"content"  db:"content"`
	Categories []*pc.ProblemCategory `json:"problem_category"`
	MaxRuntime int                   `json:"max_runtime" db:"max_runtime"`
	MaxMem     int                   `json:"max_mem" db:"max_mem"`
	TestCases  []*testcase.TestCase  `json:"test_cases"`
	SubmitNum  int                   `json:"submit_num" db:"submit_num"`
	PassNum    int                   `json:"pass_num" db:"pass_num"`

	models.CommonTimestampsField
}

func (p *ProblemBasic) Create() bool {
	CreatedAt := time.Now()
	UpdatedAt := time.Now()
	//开启事物写入数据库
	tx, err := database.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		logger.LogIf(err)
		return false
	}
	//处理problem
	_, err = tx.Exec(`INSERT INTO problem_basic (id,identity,title,content,max_runtime,max_mem,pass_num,submit_num,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?)`,
		p.ID, p.Identity, p.Title, p.Content, p.MaxRuntime, p.MaxMem, p.PassNum, p.SubmitNum, CreatedAt, UpdatedAt)
	if err != nil {
		logger.ErrorString("create_porblem", "create", err.Error())
		return false
	}
	//处理problem_category
	for _, cat := range p.Categories {
		_, err = tx.Exec(`INSERT INTO problem_category (problem_identity,category_id,created_at,updated_at) VALUES (?,?,?,?)`,
			cat.ProblemIdentity, cat.CategoryId, CreatedAt, UpdatedAt)
		if err != nil {
			logger.ErrorString("create_porblem", "create", err.Error())
			return false
		}
	}
	//处理testcase
	for _, c := range p.TestCases {
		_, err = tx.Exec(`INSERT INTO test_case (identity,problem_identity,input,output,created_at,updated_at) VALUES (?,?,?,?,?,?)`,
			c.Identity, c.ProblemIdentity, c.Input, c.Output, CreatedAt, UpdatedAt)
		if err != nil {
			logger.ErrorString("create_porblem", "create", err.Error())
			return false
		}
	}
	err = tx.Commit()
	if err != nil {
		logger.LogIf(err)
		return false
	}
	return true
}
