package submit

import (
	"bytes"
	"fmt"
	"go-oj/app/models"
	tc "go-oj/app/models/testcase"
	"go-oj/pkg/config"
	"go-oj/pkg/database"
	"go-oj/pkg/helpers"
	"go-oj/pkg/logger"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

type SubmitBasic struct {
	models.BaseModel

	Identity        string `json:"identity" db:"identity"`
	ProblemIdentity string `json:"problem_identity" db:"problem_identity"`
	UserIdentity    string `json:"user_identity" db:"user_identity"`
	Path            string `json:"path" db:"path"`
	Status          int    `json:"status" db:"status"`
	IsPass          int

	models.CommonTimestampsField
}

func (su *SubmitBasic) Save(code []byte) bool {
	//1、在指定位置创建文件夹
	dirName := config.GetString("code.dir_path") + "/" + helpers.GetUUID()
	fmt.Println("dirname", dirName)
	if err := os.Mkdir(dirName, 0777); err != nil {
		logger.LogIf(err)
		return false
	}

	//2、在文件夹中创建文件，写入代码
	f, err := os.Create(dirName + "/main.go")
	if err != nil {
		logger.LogIf(err)
		return false
	}
	defer f.Close()
	_, err = f.Write(code)
	if err != nil {
		logger.LogIf(err)
		return false
	}
	//3、保存文件名，返回true
	su.Path = dirName + "/main.go"
	return true
}

const (
	Accepted = iota + 1
	WrongAnswer
	RunTimeERR
	MemorizeErr
	CompilationERR
)

func (su *SubmitBasic) Judge(testCases []*tc.TestCase) error {
	//1、获取问题的最大运行时间和内存
	var maxRuntime, maxMem int
	err := database.DB.QueryRow("SELECT max_runtime,max_mem FROM problem_basic WHERE identity=?", su.ProblemIdentity).Scan(&maxRuntime, &maxMem)
	if err != nil {
		logger.ErrorString("submit", "judge", err.Error())
		return err
	}

	//2、判断代码相关逻辑
	var msg string
	var passCnt int
	// 答案错误的channel
	WA := make(chan int)
	// 超内存的channel
	OOM := make(chan int)
	// 编译错误的channel
	CE := make(chan int)
	// 答案正确的channel
	AC := make(chan int)

	//判断代码,暂时只支持golang
	var lock sync.Mutex
	for _, cas := range testCases {
		//对于每一个样例（goroutine）
		go func() {
			cmd := exec.Command("go", "run", su.Path)
			var stderr, stdout bytes.Buffer
			cmd.Stderr = &stderr
			cmd.Stdout = &stdout
			stdin, err := cmd.StdinPipe()
			if err != nil {
				logger.LogIf(err)
				return
			}
			_, err = io.WriteString(stdin, cas.Input+"\n")
			if err != nil {
				logger.LogIf(err)
				return
			}

			//获取运行的内存信息
			var begin, end runtime.MemStats
			runtime.ReadMemStats(&begin)
			//运行代码，捕获编译状态
			if err = cmd.Run(); err != nil {
				logger.LogIf(err)
				msg = err.Error()
				CE <- 1
				return
			}
			runtime.ReadMemStats(&end)

			//获取输出，与标准输出相比较
			if stdout.String() != cas.Output {
				WA <- 1
				return
			}

			//判断是否超内存
			if end.Alloc/1024-(begin.Alloc/1024) > uint64(maxMem) {
				OOM <- 1
				return
			}

			//答案正确
			lock.Lock()
			passCnt += 1
			if passCnt == len(testCases) {
				AC <- 1
			}
			lock.Unlock()
		}()
	}
	//等待代码检测返回结果
	select {
	case <-CE:
		msg = "编译出错 err:" + msg
		su.Status = CompilationERR
	case <-WA:
		msg = "答案错误"
		su.Status = WrongAnswer
	case <-AC:
		msg = "答案正确"
		su.Status = Accepted
		su.IsPass = 1
	case <-OOM:
		msg = "运行超内存"
		su.Status = MemorizeErr
	case <-time.After(time.Millisecond * time.Duration(maxRuntime)):
		if passCnt == len(testCases) {
			msg = "答案正确"
			su.Status = Accepted
		} else {
			msg = "运行超时"
			su.Status = RunTimeERR
		}
	}

	//3、更新评测状态
	if err := su.UpdateStatus(); err != nil {
		logger.LogIf(err)
		return err
	}
	return nil
}

func (su *SubmitBasic) UpdateStatus() error {
	_, err := database.DB.Exec("UPDATE submit_basic SET status=? WHERE identity=?", su.Status, su.Identity)
	return err
}

func (su *SubmitBasic) Commit() bool {
	//1、开启事务
	tx, err := database.DB.Beginx()
	if err != nil {
		logger.ErrorString("submit", "commit", err.Error())
		return false
	}
	defer tx.Rollback()
	//2、存储submit信息
	su.CreatedAt = time.Now()
	su.UpdatedAt = time.Now()
	sql1 := "INSERT INTO submit_basic (identity,problem_identity,user_identity,path,status,created_at,updated_at)VALUES" +
		"(:identity,:problem_identity,:user_identity,:path,:status,:created_at,:updated_at)"
	_, err = tx.NamedExec(sql1, su)
	if err != nil {
		logger.WarnString("submit", "commit_submit", err.Error())
		return false
	}

	//3、更新用户信息
	_, err = tx.Exec("UPDATE problem_basic SET submit_num=submit_num+1,pass_num=pass_num+? WHERE identity=?", su.IsPass, su.ProblemIdentity)
	if err != nil {
		logger.WarnString("submit", "commit_user", err.Error())
		return false
	}

	//4、更新题目信息
	_, err = tx.Exec("UPDATE user_basic SET submit_num=submit_num+1,pass_num=pass_num+? WHERE identity=?", su.IsPass, su.UserIdentity)
	if err != nil {
		logger.WarnString("submit", "commit_problem", err.Error())
		return false
	}
	err = tx.Commit()
	if err != nil {
		logger.ErrorString("submit", "commit", err.Error())
		return false
	}
	return true
}

func GetStatus(identity string) string {
	var status int
	err := database.DB.QueryRow("SELECT status FROM submit_basic WHERE identity=?", identity).Scan(&status)
	if err != nil {
		logger.LogIf(err)
		return ""
	}
	switch status {
	case -1:
		return "正在等待评测"
	case 1:
		return "Accepted"
	case 2:
		return "WrongAnswer"
	case 3:
		return "RuntimeError"
	case 4:
		return "MemError"
	case 5:
		return "CompilationError"
	}
	return ""
}
