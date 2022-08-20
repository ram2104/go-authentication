package initializer

import model "github.com/ram2104/go-authentication/models"

func SyncDatabase() {
	DB.AutoMigrate(&model.User{})
}
