package migrations

import (
	"crypto-project/helpers"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
	UserType string
}

type Account struct {
	gorm.Model
	Type    string
	Name    string
	Balance uint
	UserID  uint
}

func createAccounts() {
	db := helpers.ConnectDB()

	users := [2]User{
		{Username: "Martin", Email: "martin@martin.com"},
		{Username: "Michael", Email: "michael@michael.com"},
	}

	for i := 0; i < len(users); i++ {
		generatedPassword := helpers.HashAndSalt([]byte(users[i].Username))
		user := User{Username: users[i].Username, Email: users[i].Email, Password: generatedPassword}
		db.Create(&user)

		account := Account{Type: "Daily Account", Name: string(users[i].Username + "'s" + " account"), Balance: uint(10000 * int(i+1)), UserID: user.ID}
		db.Create(&account)
	}
}

func Migrate() {
	db := helpers.ConnectDB()
	if er := db.AutoMigrate(&User{}, &Account{}); er != nil {
		println('2')
	}

	// createAccounts()
}
