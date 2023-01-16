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
	logger.Dump(p.TestCases)
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
			logger.ErrorString("create_problem", "create", err.Error())
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

func (p *ProblemBasic) Modify() bool {
	//1、开启事务
	tx, err := database.DB.Beginx()
	if err != nil {
		logger.ErrorString("problem", "modify", err.Error())
		return false
	}
	defer tx.Rollback()

	//2、修改problem_basic数据
	p.UpdatedAt = time.Now()
	sql := "UPDATE problem_basic SET title=?,content=?,max_runtime=?,max_mem=?,updated_at=?  WHERE identity=?"
	_, err = tx.Exec(sql, p.Title, p.Content, p.MaxRuntime, p.MaxMem, p.UpdatedAt, p.Identity)
	if err != nil {
		logger.ErrorString("problem", "modify", err.Error())
		return false
	}

	//commit
	err = tx.Commit()
	if err != nil {
		logger.ErrorString("problem", "modify", err.Error())
		return false
	}
	return true
}

// GetProblemList 返回问题的 title-identity 列表
func GetProblemList(size, page int) map[string]string {
	//1、计算偏移量
	offset := (page - 1) * size

	//2、访问数据库
	sql := "SELECT identity,title FROM problem_basic WHERE id>? and id<=?"
	rows, err := database.DB.Queryx(sql, offset, offset+size)
	logger.Dump(offset, "offset")
	logger.Dump(offset+size, "offset+size")
	if err != nil {
		logger.LogIf(err)
		return nil
	}
	defer rows.Close()
	list := make(map[string]string)
	var identity, title string
	for rows.Next() {
		err := rows.Scan(&identity, &title)
		if err != nil {
			logger.LogIf(err)
			return nil
		}
		list[title] = identity
	}
	return list
}

// GetProblemDetail 获取问题详细信息
func GetProblemDetail(identity string) ProblemBasic {
	p := ProblemBasic{}
	//1、开启事务
	tx, err := database.DB.Beginx()
	if err != nil {
		logger.LogIf(err)
		return ProblemBasic{}
	}
	defer tx.Rollback()

	//2、获取问题基础信息
	err = tx.Get(&p, "SELECT * FROM problem_basic WHERE identity=?", identity)
	if err != nil {
		logger.Dump(p)
		logger.LogIf(err)
	}
	//3、获取问题分类信息
	err = tx.Select(&p.Categories, "SELECT * FROM problem_category WHERE problem_identity=?", identity)
	if err != nil {
		logger.Dump(p.Categories)
		logger.LogIf(err)
	}
	//4、获取问题测试用例
	err = tx.Select(&p.TestCases, "SELECT * FROM test_case WHERE problem_identity=?", identity)
	if err != nil {
		logger.Dump(p.TestCases)
		logger.LogIf(err)
	}
	tx.Commit()
	return p
}
