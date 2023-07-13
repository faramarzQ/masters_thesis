package model

import (
	"resource_manager/internal/consts"

	"gorm.io/gorm"
)

type ScalingLog struct {
	ID                   int `gorm:"column:id;primaryKey;autoIncrement"`
	ScalerExecutionLogID int
	NodeName             string             `gorm:"size:256"`
	Class                consts.NODE_CLASS  `gorm:"size:32"`
	ScalerExecutionLog   ScalerExecutionLog `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	gorm.Model
}
