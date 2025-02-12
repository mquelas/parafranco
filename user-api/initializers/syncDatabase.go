// initializers/syncDatabase.go
package initializers

import "user-reservation-api/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
