package model

import (
	"resource_manager/internal/consts"
	"time"
)

type ScalingLog struct {
	ID       int               `gorm:"column:id;primaryKey;autoIncrement"`
	NodeName string            `gorm:"size:256"`
	Class    consts.NODE_CLASS `gorm:"size:32"`
	ScaledAt time.Time
}
