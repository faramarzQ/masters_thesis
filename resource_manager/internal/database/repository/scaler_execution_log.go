package repository

import (
	"resource_manager/internal/database"
	"resource_manager/internal/database/model"
	"time"

	"k8s.io/klog/v2"
)

// Inserts a scaler execution log record
func InsertScalerExecutionLog(previousScalerExecutionLog *model.ScalerExecutionLog, scalerName string) (model.ScalerExecutionLog, error) {
	step := 1
	if previousScalerExecutionLog != nil {
		step = previousScalerExecutionLog.Step + 1
	}

	log := model.ScalerExecutionLog{
		ScalerName:                scalerName,
		Step:                      step,
		ExecutedAt:                time.Now(),
		ScalerExecutionLogDetails: &model.ScalerExecutionLogDetails{},
	}

	database.DBConn.Create(&log)
	return log, nil
}

// Selects execution logs of a given scaler name
func GetPreviousScalerExecutionLog(scalerName string) *model.ScalerExecutionLog {
	log := &model.ScalerExecutionLog{}
	tx := database.DBConn.Where("scaler_name = ?", scalerName).
		// This join only retrieves logs which have details
		Joins("join scaler_execution_log_details on scaler_execution_logs.id = scaler_execution_log_details.scaler_execution_log_id").
		Preload("ScalerExecutionLogDetails").
		Last(&log)

	if tx.Error != nil {
		klog.Error(tx.Error)
		return nil
	}

	return log
}

// Deletes scaling log record with it's dependencies
func DeleteScalingExecutionLog(scalerExecutionLog *model.ScalerExecutionLog) {
	for _, transition := range (*scalerExecutionLog).ScalingLog {
		database.DBConn.Unscoped().Delete(transition)
	}
	database.DBConn.Unscoped().Delete((*scalerExecutionLog).ScalerExecutionLogDetails)
	database.DBConn.Unscoped().Delete(&scalerExecutionLog)
}
