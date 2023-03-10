package helpers

import (
	"crypto-project/interfaces"
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// "strconv"
	"fmt"
	"strings"
)

func HandleErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func HashAndSalt(pass []byte) string {
	hashed, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	HandleErr(err)

	return string(hashed)
}

func ConnectDB() *gorm.DB {
	// dsn := "host=labourpool-db-dev.mareca.vc user=postgres password=labpool dbname=labourpool port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	//dsn := "host=localhost user=dikdik password=kurnia dbname=restapi port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	dsn := "host=localhost user=postgres password=rifaldilicious5 dbname=labourpool port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	HandleErr(err)
	return db
}

func Validation(values []interfaces.Validation) bool {
	username := regexp.MustCompile(`^([A-Za-z0-9]{5,})+$`)
	email := regexp.MustCompile(`^[A-Za-z0-9]+[@]+[A-Za-z0-9]+[.]+[A-Za-z]+$`)

	fmt.Println(values)
	for i := 0; i < len(values); i++ {
		switch values[i].Valid {
		case "username":
			if username.MatchString(values[i].Value) {
				fmt.Println("gagal user", username.MatchString(values[i].Value))
				return false
			}
		case "email":
			if !email.MatchString(values[i].Value) {
				fmt.Println("gagal email")
				return false
			}
		case "password":
			if len(values[i].Value) < 5 {
				fmt.Println("gagal password")
				return false
			}
		case "user_type":
			s := []string{"peserta", "perusahaan", "admin"}
			if contains(s, values[i].Value) {
				return false
			}
		}
	}
	return true
}

func PanicHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			error := recover()
			if error != nil {
				log.Println(error)
				resp := interfaces.ErrResponse{Message: "Internal Server Error"}
				json.NewEncoder(w).Encode(resp)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func ValidateToken(id string, jwtToken string) bool {
	cleanJWT := strings.Replace(jwtToken, "Bearer ", "", -1)
	tokenData := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cleanJWT, tokenData, func(token *jwt.Token) (interface{}, error) {
		return []byte("TokenPassword"), nil
	})
	HandleErr(err)
	// var userId, _ = strconv.ParseFloat(id, 8)
	// if token.Valid && tokenData["user_id"] == userId {
	if token.Valid {
		return true
	} else {
		return false
	}
}

func UserID(jwtToken string) uint {
	cleanJWT := strings.Replace(jwtToken, "Bearer ", "", -1)
	tokenData := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cleanJWT, tokenData, func(token *jwt.Token) (interface{}, error) {
		return []byte("TokenPassword"), nil
	})
	HandleErr(err)
	userID := tokenData["user_id"]
	id := uint(userID.(float64))
	if token.Valid {
		return id
	} else {
		return 0
	}

}

func UserIDStr(jwtToken string) string {
	cleanJWT := strings.Replace(jwtToken, "Bearer ", "", -1)
	tokenData := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cleanJWT, tokenData, func(token *jwt.Token) (interface{}, error) {
		return []byte("TokenPassword"), nil
	})
	HandleErr(err)
	userID := tokenData["user_id"]
	id := fmt.Sprintf("%v", userID)
	if token.Valid {
		return id
	} else {
		return "0"
	}

}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Print(w, err.Error())
	}
}
