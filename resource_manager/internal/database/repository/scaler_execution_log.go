package repository

import (
	"resource_manager/internal/database"
	"resource_manager/internal/database/model"
	"time"
)

// Inserts a scaler execution log record
func InsertScalerExecutionLog(previousScalerExecutionLog *model.ScalerExecutionLog, scalerName string) (model.ScalerExecutionLog, error) {
	var previousState string
	var previousActionTaken int8
	if previousScalerExecutionLog != nil {
		previousState = previousScalerExecutionLog.ScalerExecutionLogDetails.State
		previousActionTaken = previousScalerExecutionLog.ScalerExecutionLogDetails.ActionTaken
	}

	log := model.ScalerExecutionLog{
		ScalerName: scalerName,
		ExecutedAt: time.Now(),
		ScalerExecutionLogDetails: &model.ScalerExecutionLogDetails{
			PreviousState:       previousState,
			PreviousActionTaken: previousActionTaken,
		},
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
		Preload("ScalingLog").
		Last(&log)

	if tx.Error != nil {
		return nil
	}

	return log
}
