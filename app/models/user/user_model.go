package user

import "gohub/app/models"

// Package user 存放用户Model相关逻辑

// User用户模型

type User struct {
	models.BaseModel

	Name     string `json:"name,omitempty"`
	Email    string `json:"-"`
	Phone    string `json:"-"`
	Password string `json:"-"`

	models.CommonTimestampsField
}
