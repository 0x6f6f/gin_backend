package migrations

import (
	"gin-boilerplate/infra/database"
	"gin-boilerplate/models"
)

// Migrate Add list of model add for migrations
// TODO later separate migration each models
func Migrate() {
	var migrationModels = []interface{}{
		&models.Zone{},       // 创建没有外键依赖的表
		&models.Department{}, // 依赖Zone
		&models.User{},       // 可能依赖其他表，如Department
		&models.UserProfile{},
		&models.WorkLog{},
		&models.Customer{},
		&models.Contract{},
		&models.SystemLog{},
	}
	err := database.DB.AutoMigrate(migrationModels...)
	if err != nil {
		return
	}
}
