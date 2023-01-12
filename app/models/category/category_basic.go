package category

import "go-oj/app/models"

type CategoryBasic struct {
	models.BaseModel

	Identity string `json:"identity"`
	Name     string `json:"name"`
	ParentId string `json:"parentid"`

	models.CommonTimestampsField
}
