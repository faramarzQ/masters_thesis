package repository

import (
	"resource_manager/internal/database"
	"resource_manager/internal/database/model"
)

func InsertScalerExecutionLogDetail(scalerExecutionLog model.ScalerExecutionLog, state string, action int8) {
	scalerExecutionLog.ScalerExecutionLogDetails.State = state
	scalerExecutionLog.ScalerExecutionLogDetails.ActionTaken = action

	database.DBConn.Save(&scalerExecutionLog.ScalerExecutionLogDetails)
}
