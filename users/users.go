package users

import (
	"crypto-project/helpers"
	"crypto-project/interfaces"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

func Login(username string, pass string) map[string]interface{} {
	// Add validation to login
	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: pass, Valid: "password"},
		})
	if valid {

		db := helpers.ConnectDB()
		user := &interfaces.User{}

		err := db.Where("username", username).First(&user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return map[string]interface{}{"message": "User not found"}
		}

		passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
		if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
			return map[string]interface{}{"message": "Wrong password"}
		}

		accounts := []interfaces.ResponseAccount{}
		db.Table("accounts").Select("id, name, balance").Where("user_id = ? ", user.ID).Scan(&accounts)

		var response = prepareResponse(user, accounts, true)
		return response
	} else {
		return map[string]interface{}{"message": "not valid values"}
	}

}

func prepareToken(user *interfaces.User) string {
	tokenContent := jwt.MapClaims{
		"user_id": user.ID,
		"expiry":  time.Now().Add(time.Minute * 60).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	helpers.HandleErr(err)

	return token
}

func prepareResponse(user *interfaces.User, accounts []interfaces.ResponseAccount, withToken bool) map[string]interface{} {
	responseUser := &interfaces.ResponseUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		UserType: user.UserType,
		Accounts: accounts,
	}

	var response = map[string]interface{}{"message": "all is fine"}
	if withToken {
		var token = prepareToken(user)
		response["jwt"] = token
	}
	response["data"] = responseUser

	return response
}

func prepareUpdateUserResponse(user *interfaces.User, accounts []interfaces.ResponseAccount, withToken bool) map[string]interface{} {
	responseUser := &interfaces.ResponseUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		UserType: user.UserType,
		Accounts: accounts,
	}

	var response = map[string]interface{}{"message": "all is fine"}
	if withToken {
		var token = prepareToken(user)
		response["jwt"] = token
	}
	response["data"] = responseUser

	return response
}

func Register(username string, email string, pass string, user_type string) map[string]interface{} {

	// Add validation to registration
	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: email, Valid: "email"},
			{Value: pass, Valid: "password"},
			{Value: user_type, Valid: "user_type"},
		})
	if valid {
		db := helpers.ConnectDB()
		checkuser := &interfaces.User{}

		//prevent duplicate username, email
		err := db.Where("username", username).Or("email", email).First(&checkuser).Error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return map[string]interface{}{"message": "User exist"}
		}

		generatedPassword := helpers.HashAndSalt([]byte(pass))
		user := &interfaces.User{Username: username, Email: email, Password: generatedPassword, UserType: user_type}
		db.Create(&user)

		account := &interfaces.Account{Type: "Daily Account", Name: string(username + "'s" + " account"), Balance: 0, UserID: user.ID}
		db.Create(&account)

		accounts := []interfaces.ResponseAccount{}
		respAccount := interfaces.ResponseAccount{ID: account.ID, Name: account.Name, Balance: int(account.Balance)}
		accounts = append(accounts, respAccount)
		var response = prepareResponse(user, accounts, true)

		return response
	} else {
		return map[string]interface{}{"message": "not valid values"}
	}
}

func GetUser(id string, jwt string) map[string]interface{} {
	isValid := helpers.ValidateToken(id, jwt)
	if isValid {
		db := helpers.ConnectDB()

		user := &interfaces.User{}
		err := db.Where("id", id).First(&user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return map[string]interface{}{"message": "User not found"}
		}

		accounts := []interfaces.ResponseAccount{}
		db.Table("accounts").Select("id, name, balance").Where("user_id = ? ", user.ID).Scan(&accounts)

		var response = prepareResponse(user, accounts, false)
		return response
	} else {
		return map[string]interface{}{"Message": "Not valid token"}
	}
}

func UpdateUser(userID string, username string, email string, password string, user_type string, name string, balance uint, jwt string) map[string]interface{} {

	// userID 	= helpers.UserIDStr(jwt)
	isValid := helpers.ValidateToken(userID, jwt)
	
	if isValid {
		db 		:= helpers.ConnectDB()	
		user 	:= &interfaces.User{}
		
		err := db.First(&user, userID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return map[string]interface{} {"message": "record not found"}
		}

		users 	:= &interfaces.User{Username: username, Password: password, Email: email, UserType: user_type}
		db.Where("id = ?", userID).Updates(&users)

		account := &interfaces.Account{Name: name, Balance: balance}
		db.Where("user_id = ?", userID).Updates(&account)

		accounts := []interfaces.ResponseAccount{}
		respAccount := interfaces.ResponseAccount{ID: account.ID, Name: account.Name, Balance: int(account.Balance)}
		accounts = append(accounts, respAccount)

		return prepareUpdateUserResponse(user,accounts, false)
	} else {
		return map[string]interface{}{"Message": "Not valid token"}
	}
}