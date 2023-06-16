package repository

import (
	"resource_manager/internal/database"
	"resource_manager/internal/database/model"
	"time"
)

// Inserts a scaler execution log record
func InsertScalerExecutionLog(scalerName string) (model.ScalerExecutionLog, error) {
	log := model.ScalerExecutionLog{
		ScalerName: scalerName,
		ExecutedAt: time.Now(),
	}

	database.DBConn.Create(&log)

	return log, nil
}

// Selects execution logs of a given scaler name
func SelectScalingLogByScalerName(scalerName string) []model.ScalerExecutionLog {
	logs := []model.ScalerExecutionLog{}

	database.DBConn.Where("scaler_name >= ?", scalerName).Find(&logs)

	return logs
}
