// initializers/syncDatabase.go
package initializers

import "reservation-api/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.Reservation{})
}
