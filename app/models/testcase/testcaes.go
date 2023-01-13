package testcase

import (
	"go-oj/app/models"
	"go-oj/pkg/database"
	"go-oj/pkg/logger"
)

type TestCase struct {
	models.BaseModel

	Identity        string `json:"identity" db:"identity"`
	ProblemIdentity string `json:"problem_identity" db:"problem_identity"`
	Input           string `json:"input" db:"input"`
	Output          string `json:"output" db:"output"`

	models.CommonTimestampsField
}

func GetTestCases(problemIdentity string) []*TestCase {
	testCases := make([]*TestCase, 0)
	err := database.DB.Select(&testCases, "SELECT * FROM test_case WHERE problem_identity=?", problemIdentity)
	if err != nil {
		logger.ErrorString("submit", "GetTestCase", err.Error())
		return nil
	}
	return testCases
}
