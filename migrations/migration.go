package migrations

import (
	"gin-boilerplate/infra/database"
	"gin-boilerplate/models"
)

// Migrate Add list of model add for migrations
// TODO later separate migration each models
func Migrate() {
	var migrationModels = []interface{}{
		&models.User{},
		&models.UserProfile{},
		&models.SalesDepartment{},
		&models.Zone{},
		&models.Saler{},
		&models.WorkLog{},
		&models.FinanceSpecialist{},
		&models.Accountant{},
		&models.Customer{},
		&models.Contract{},
		&models.SystemLog{},
	}
	err := database.DB.AutoMigrate(migrationModels...)
	if err != nil {
		return
	}
}
