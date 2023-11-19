package initializers

import "final-assignment/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.Admin{})
}
