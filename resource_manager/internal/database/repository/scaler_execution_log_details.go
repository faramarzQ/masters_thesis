package repository

import (
	"resource_manager/internal/database"
	"resource_manager/internal/database/model"
)

func InsertScalerExecutionLogDetail(scalerExecutionLog model.ScalerExecutionLog, state string, action int8, epsilon uint8) (model.ScalerExecutionLogDetails, error) {
	log := model.ScalerExecutionLogDetails{
		ScalerExecutionLog: scalerExecutionLog,
		State:              state,
		ActionTaken:        action,
		EpsilonValue:       epsilon,
	}

	database.DBConn.Create(&log)

	return log, nil
}
