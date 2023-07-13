package model

import (
	"gorm.io/gorm"
)

type ScalerExecutionLogDetails struct {
	ID                   int `gorm:"column:id;primaryKey;autoIncrement"`
	ScalerExecutionLogID int `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	State                string
	ActionTaken          int8
	EpsilonValue         uint8
	ScalerExecutionLog   ScalerExecutionLog
	gorm.Model
}