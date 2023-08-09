package repository

import (
	"resource_manager/internal/database"
	"resource_manager/internal/database/model"
)

func InsertScalerExecutionLogDetail(scalerExecutionLog model.ScalerExecutionLog, state string, action int8, epsilon uint8) {
	scalerExecutionLog.ScalerExecutionLogDetails.State = state
	scalerExecutionLog.ScalerExecutionLogDetails.ActionTaken = action
	scalerExecutionLog.ScalerExecutionLogDetails.EpsilonValue = epsilon

	database.DBConn.Save(&scalerExecutionLog.ScalerExecutionLogDetails)
}
