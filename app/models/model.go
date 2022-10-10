package models

import "time"

// package modles 模型通用属性和方法

// BaseModel 模型基类
type BaseModel struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement;"json:"id,omitempty"`
}

// CommonTimestampsField时间戳
type CommonTimestampsField struct {
	CreatedAt time.Time `gorm:"column:created_at;index;"json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;index;"json:"Updated_at,o,omitempty"`
}
