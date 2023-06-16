package model

import "time"

type ScalerExecutionLog struct {
	ID         int    `gorm:"column:id;primaryKey;autoIncrement"`
	ScalerName string `gorm:"size:256"`
	ExecutedAt time.Time
}
