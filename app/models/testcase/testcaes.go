package testcase

import "go-oj/app/models"

type TestCase struct {
	models.BaseModel

	Identity        string `json:"identity" db:"identity"`
	ProblemIdentity string `json:"problem_identity" db:"problem_identity"`
	Input           string `json:"input" db:"input"`
	Output          string `json:"output" db:"output"`

	models.CommonTimestampsField
}
