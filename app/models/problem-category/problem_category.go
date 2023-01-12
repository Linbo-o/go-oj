package problem_category

import "go-oj/app/models"

type ProblemCategory struct {
	models.BaseModel

	ProblemIdentity string `json:"problem_identity" db:"problem_identity"`
	CategoryId      uint   `json:"category_id" db:"category_id"`

	models.CommonTimestampsField
}
