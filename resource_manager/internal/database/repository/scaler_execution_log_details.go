package repository

import (
	"resource_manager/internal/database"
	"resource_manager/internal/database/model"
)

func InsertScalerExecutionLogDetail(scalerExecutionLog model.ScalerExecutionLog,
	state string,
	action int8,
	epsilon float64,
	energyConsumption float64,
	successRequestRate float64) {

	scalerExecutionLog.ScalerExecutionLogDetails.State = state
	scalerExecutionLog.ScalerExecutionLogDetails.ActionTaken = action
	scalerExecutionLog.ScalerExecutionLogDetails.Epsilon = epsilon
	scalerExecutionLog.ScalerExecutionLogDetails.EnergyConsumption = energyConsumption
	scalerExecutionLog.ScalerExecutionLogDetails.SuccessRequestRate = successRequestRate

	database.DBConn.Save(&scalerExecutionLog.ScalerExecutionLogDetails)
}

func UpdateScalerExecutionLogDetail(scalerExecutionLog model.ScalerExecutionLog, reward float64) {
	scalerExecutionLog.ScalerExecutionLogDetails.Reward = reward

	database.DBConn.Save(&scalerExecutionLog.ScalerExecutionLogDetails)
}
