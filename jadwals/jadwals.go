package jadwals

import (
	"crypto-project/helpers"
	"crypto-project/interfaces"
	"errors"
	// "github.com/dgrijalva/jwt-go"
	// "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	// "time"
)

func prepareCreateJadwalResponse(jadwal *interfaces.Jadwal) map[string]interface{} {
	responsePost := &interfaces.ResponseCreatePost{
		ID:   		jadwal.ID,
	}

	var response = map[string]interface{}{"message": "all is fine"}
	response["data"] = responsePost
	return response
}

func prepareReadJadwalResponse(jadwal *interfaces.Jadwal) map[string]interface{} {
	responseReadJadwal := &interfaces.ResponseReadJadwal{
		ID:   		jadwal.ID,
		Kuota:   	jadwal.Kuota,
		Lokasi:   	jadwal.Lokasi,
	}

	var response = map[string]interface{}{"message": "all is fine"}
	response["data"] = responseReadJadwal

	return response
}

func prepareDeleteJadwalResponse(jadwal *interfaces.Jadwal) map[string]interface{} {
	responseJadwal := &interfaces.ResponseDeleteJadwal{
		ID:   		jadwal.ID,
	}

	var response = map[string]interface{}{"message": "all is fine"}
	response["data"] = responseJadwal

	return response
}

func prepareUpdateJadwalResponse(jadwal *interfaces.Jadwal) map[string]interface{} {
	responsePost := &interfaces.ResponseUpdateJadwal{
		ID:   		jadwal.ID,
		Kuota:   	jadwal.Kuota,
		Lokasi:   	jadwal.Lokasi,
	}

	var response = map[string]interface{}{"message": "all is fine"}
	response["data"] = responsePost
	return response
}

func CreateJadwal(kuota uint, lokasi string, jwt string) map[string]interface{} {

	userID 	:= helpers.UserIDStr(jwt)
	isValid := helpers.ValidateToken(userID, jwt)

	if isValid {
		// userID := helpers.UserID(jwt)
		db 		:= helpers.ConnectDB()	
		jadwal 	:= &interfaces.Jadwal{Kuota: kuota, Lokasi: lokasi}
		db.Create(&jadwal)
		return prepareCreateJadwalResponse(jadwal)
	} else {
		return map[string]interface{}{"Message": "Not valid token"}
	}
}

func ReadJadwal(jadwalID string, jwt string) map[string]interface{} {

	userID 	:= helpers.UserIDStr(jwt)
	isValid := helpers.ValidateToken(userID, jwt)
	if isValid {
		db 		:= helpers.ConnectDB()
		jadwal 	:= &interfaces.Jadwal{}		
		err 	:= db.First(&jadwal, jadwalID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return map[string]interface{}{"message": "record not found"}
		}
		return prepareReadJadwalResponse(jadwal)
	} else {
		return map[string]interface{}{"Message": "Not valid token"}
	}
}

func DeleteJadwal(jadwalID string, jwt string) map[string]interface{} {

	userID 	:= helpers.UserIDStr(jwt)
	isValid := helpers.ValidateToken(userID, jwt)
	if isValid {
		db 		:= helpers.ConnectDB()
		jadwal 	:= &interfaces.Jadwal{}
		
		err := db.First(&jadwal, jadwalID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return map[string]interface{} {"message": "record not found"}
		}

		db.Where("id = ?", jadwalID).Delete(&jadwal)
		return prepareDeleteJadwalResponse(jadwal)
	} else {
		return map[string]interface{}{"Message": "Not valid token"}
	}
}

func UpdateJadwal(jadwalID string, kuota uint, lokasi string, jwt string) map[string]interface{} {

	userID 	:= helpers.UserIDStr(jwt)
	isValid := helpers.ValidateToken(userID, jwt)
	if isValid {
		db 		:= helpers.ConnectDB()	
		jadwal 	:= &interfaces.Jadwal{}
		
		err := db.First(&jadwal, jadwalID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return map[string]interface{} {"message": "record not found"}
		}

		jadwal 	= &interfaces.Jadwal{Kuota: kuota, Lokasi: lokasi}
		db.Where("id = ?", jadwalID).Updates(&jadwal)
		return prepareUpdateJadwalResponse(jadwal)
	} else {
		return map[string]interface{}{"Message": "Not valid token"}
	}
}