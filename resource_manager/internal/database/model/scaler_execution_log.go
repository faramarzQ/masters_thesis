package model

import (
	"time"

	"gorm.io/gorm"
)

type ScalerExecutionLog struct {
	ID         int    `gorm:"column:id;primaryKey;autoIncrement"`
	ScalerName string `gorm:"size:256"`
	Step       int
	ExecutedAt time.Time
	gorm.Model
	ScalingLog                []*ScalingLog
	ScalerExecutionLogDetails *ScalerExecutionLogDetails
}
