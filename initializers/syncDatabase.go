package initializers

import "JWT-AUTH-GIN/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})

}
