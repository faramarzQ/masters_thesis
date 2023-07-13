package repository

import (
	"resource_manager/internal/consts"
	"resource_manager/internal/database"
	"resource_manager/internal/database/model"
	"time"
)

// Inserts a scaling log record
func InsertScalingLog(nodeName string, class consts.NODE_CLASS) (model.ScalingLog, error) {
	log := model.ScalingLog{
		NodeName: nodeName,
		Class:    class,
	}

	database.DBConn.Create(&log)

	return log, nil
}

// Selects scaling_log records inserted after a given time
func SelectScalingLogByScaledAt(scaledAt time.Time) []model.ScalingLog {
	logs := []model.ScalingLog{}

	database.DBConn.Where("scaled_at >= ?", scaledAt).Find(&logs)

	return logs
}
