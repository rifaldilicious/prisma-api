package interfaces

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string
	Email     string
	Password  string
	UserType  string
	Status    string
	Deskripsi string
	Url       string
}

type Account struct {
	gorm.Model
	Type    string
	Name    string
	Balance uint
	UserID  uint
}

type ResponseAccount struct {
	ID      uint
	Name    string
	Balance int
}

type ResponseUser struct {
	ID       uint
	Username string
	Email    string
	UserType string
	Accounts []ResponseAccount
}

type Validation struct {
	Value string
	Valid string
}

type ErrResponse struct {
	Message string
}

type Post struct {
	gorm.Model
	ID       uint
	User_ID  uint
	Name     string
	Skill    string
	Location string
	Position string
	Work     string
	Salary   uint
	Message  string
}

type ResponseCreatePost struct {
	ID uint
}

type ResponseReadPost struct {
	ID       uint
	User_ID  uint
	Name     string
	Skill    string
	Location string
	Position string
	Work     string
	Salary   uint
	Message  string
}

type ResponseDeletePost struct {
	ID uint
}

type ResponseUpdatePost struct {
	gorm.Model
	ID       uint
	User_ID  uint
	Name     string
	Skill    string
	Location string
	Position string
	Work     string
	Salary   uint
	Message  string
}

type ResponseReadJadwal struct {
	ID     uint
	Kuota  uint
	Lokasi string
}

type Jadwal struct {
	gorm.Model
	ID     uint
	Kuota  uint
	Lokasi string
}

type ResponseDeleteJadwal struct {
	ID uint
}

type ResponseUpdateJadwal struct {
	gorm.Model
	ID     uint
	Kuota  uint
	Lokasi string
}
