package api

import (
	"crypto-project/helpers"
	"crypto-project/interfaces"
	"crypto-project/users"
	"crypto-project/posts"
	"crypto-project/jadwals"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type JSONResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Login struct {
	Username string
	Password string
}

type Register struct {
	Username 	string
	Email    	string
	Password 	string
	User_type 	string
}

type User struct {
	Username 	string
	Password 	string
	Email 		string
	UserType	string
	UserID 		string
	Balance		uint
	Name		string
}

type Post struct {
	User_ID 	uint
	Name    	string
	Skill 		string
	Location 	string
	Position 	string
	Work 		string
	Salary 		uint
	Message 	string
}

type UpdatePost struct {
	ID 			uint
	User_ID 	uint
	Name    	string
	Skill 		string
	Location 	string
	Position 	string
	Work 		string
	Salary 		uint
	Message 	string
}

type Jadwal struct {
	Kuota 	uint
	Lokasi 	string
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

func RespondWithError(w http.ResponseWriter, code int, message string) {
	json := ErrorResponse{Status: "Failed", Message: message}
	RespondWithJSON(w, code, json)
}


func readBody(r *http.Request) []byte {
	body, err := ioutil.ReadAll(r.Body)
	helpers.HandleErr(err)
	return body
}

func apiResponse(call map[string]interface{}, w http.ResponseWriter) {
	if call["message"] == "all is fine" {
		resp := call
		json.NewEncoder(w).Encode(resp)
	} else if call["message"] == "record not found" {
		resp := interfaces.ErrResponse{Message: "Record not found"}
		json.NewEncoder(w).Encode(resp)
	} else {
		resp := interfaces.ErrResponse{Message: "Wrong username or password"}
		json.NewEncoder(w).Encode(resp)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	w.Write([]byte(""))

	body := readBody(r)
	var formattedBody Login
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)
	login := users.Login(formattedBody.Username, formattedBody.Password)

	apiResponse(login, w)
}

func register(w http.ResponseWriter, r *http.Request) {
	body := readBody(r)
	// Handle registration
	var formattedBody Register
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)

	register := users.Register(formattedBody.Username, formattedBody.Email, formattedBody.Password, formattedBody.User_type)
	// Prepare response
	apiResponse(register, w)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	auth := r.Header.Get("Authorization")

	user := users.GetUser(userId, auth)
	apiResponse(user, w)
}

func readAllUser(w http.ResponseWriter, r *http.Request) {

	//TODO: Gunakan apiresponse interfasce
	auth 	:= r.Header.Get("Authorization")	
	userID 	:= helpers.UserIDStr(auth)
	isValid := helpers.ValidateToken(userID, auth)

	if isValid {
		var users []interfaces.User
		db 		:= helpers.ConnectDB()
		db.Find(&users)
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)
	}

}

func updateUser(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")

	var formattedBody User

	vars 	:= mux.Vars(r)
	userID 	:= vars["id"]
	body 	:= readBody(r)
	err 	:= json.Unmarshal(body, &formattedBody)

	helpers.HandleErr(err)

	updateUser := users.UpdateUser(userID, formattedBody.Username, formattedBody.Email, formattedBody.Password, formattedBody.UserType, formattedBody.Name, formattedBody.Balance, auth)
	apiResponse(updateUser, w)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")

	body := readBody(r)
	var formattedBody Post
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)
	createPost := posts.CreatePost(formattedBody.Name, formattedBody.Skill, formattedBody.Location, formattedBody.Position, formattedBody.Work, formattedBody.Salary, formattedBody.Message, auth)
	apiResponse(createPost, w)
}

func readPost(w http.ResponseWriter, r *http.Request) {

	vars 	:= mux.Vars(r)
	postId 	:= vars["id"]
	auth := r.Header.Get("Authorization")
	post 	:= posts.ReadPost(postId, auth)
	apiResponse(post, w)

}

func readAllPost(w http.ResponseWriter, r *http.Request) {

	//TODO: Gunakan apiresponse interfasce
	auth 	:= r.Header.Get("Authorization")	
	userID 	:= helpers.UserIDStr(auth)
	isValid := helpers.ValidateToken(userID, auth)

	if isValid == true {

		var posts []interfaces.Post
		db 		:= helpers.ConnectDB()
		if err := db.Find(&posts).Error; err != nil {
			RespondWithJSON(w, 400, JSONResponse{
				Status:  "failed",
				Message: "Gagal mengambil Data",
				Data:    err,
			})
			return
		}

		RespondWithJSON(w, 200, JSONResponse{
			Status:  "success",
			Message: "Berhasil mengambil data",
			Data:    posts,
		})

	} else {
		fmt.Println("xxx")
	}

}

func deletePost(w http.ResponseWriter, r *http.Request) {
	vars 	:= mux.Vars(r)
	postId 	:= vars["id"]
	auth := r.Header.Get("Authorization")
	post 	:= posts.DeletePost(postId, auth)
	apiResponse(post, w)

}

func updatePost(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")

	var formattedBody Post

	vars 	:= mux.Vars(r)
	postId 	:= vars["id"]
	body 	:= readBody(r)
	err 	:= json.Unmarshal(body, &formattedBody)

	helpers.HandleErr(err)

	updatePost := posts.UpdatePost(postId, formattedBody.User_ID, formattedBody.Name, formattedBody.Skill, formattedBody.Location, formattedBody.Position, formattedBody.Work, formattedBody.Salary, formattedBody.Message, auth)
	apiResponse(updatePost, w)
}

func createJadwal(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	body := readBody(r)

	var formattedBody Jadwal
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)
	createJadwal := jadwals.CreateJadwal(formattedBody.Kuota, formattedBody.Lokasi, auth)
	apiResponse(createJadwal, w)
}

func readJadwal(w http.ResponseWriter, r *http.Request) {

	vars 		:= mux.Vars(r)
	jadwalID 	:= vars["id"]
	auth 		:= r.Header.Get("Authorization")
	jadwal 		:= jadwals.ReadJadwal(jadwalID, auth)
	apiResponse(jadwal, w)

}

func readAllJadwal(w http.ResponseWriter, r *http.Request) {

	//TODO: Gunakan apiresponse interfasce
	auth 	:= r.Header.Get("Authorization")	
	userID 	:= helpers.UserIDStr(auth)
	isValid := helpers.ValidateToken(userID, auth)

	if isValid {

		var jadwals []interfaces.Jadwal
		db 		:= helpers.ConnectDB()
		db.Find(&jadwals)
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(jadwals)
		
	}

}

func deleteJadwal(w http.ResponseWriter, r *http.Request) {
	vars 		:= mux.Vars(r)
	jadwalID 	:= vars["id"]
	auth 		:= r.Header.Get("Authorization")
	jadwal 		:= jadwals.DeleteJadwal(jadwalID, auth)
	apiResponse(jadwal, w)

}

func updateJadwal(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")

	var formattedBody Jadwal

	vars 		:= mux.Vars(r)
	jadwalID	:= vars["id"]
	body 		:= readBody(r)
	err 		:= json.Unmarshal(body, &formattedBody)

	helpers.HandleErr(err)

	updateJadwal := jadwals.UpdateJadwal(jadwalID, formattedBody.Kuota, formattedBody.Lokasi, auth)
	apiResponse(updateJadwal, w)
}

func StartApi() {
	router := mux.NewRouter()
	router.Use(helpers.PanicHandler)

	//USER
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/register", register).Methods("POST")
	router.HandleFunc("/users", readAllUser).Methods("GET")
	router.HandleFunc("/user/{id}", getUser).Methods("GET")
	router.HandleFunc("/user/{id}", updateUser).Methods("PUT")

	//POST
	router.HandleFunc("/post", createPost).Methods("POST")
	router.HandleFunc("/post/{id}", readPost).Methods("GET")
	router.HandleFunc("/posts", readAllPost).Methods("GET")
	router.HandleFunc("/post/{id}", deletePost).Methods("DELETE")
	router.HandleFunc("/post/{id}", updatePost).Methods("PUT")

	//JADWAL
	router.HandleFunc("/jadwal", createJadwal).Methods("POST")
	router.HandleFunc("/jadwal/{id}", readJadwal).Methods("GET")
	router.HandleFunc("/jadwals", readAllJadwal).Methods("GET")
	router.HandleFunc("/jadwal/{id}", deleteJadwal).Methods("DELETE")
	router.HandleFunc("/jadwal/{id}", updateJadwal).Methods("PUT")

	router.Use(mux.CORSMethodMiddleware(router))

	fmt.Println("App is working on port :8888")

	//mux := http.NewServeMux()
	//mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	w.Header().Set("Content-Type", "application/json")
	//	w.Write([]byte("{\"hello\": \"world\"}"))
	//})
	//handler := cors.Default().Handler(mux)

	log.Fatal(http.ListenAndServe(":8888", router))
	//log.Fatal(http.ListenAndServe(":8888", handlers.CORS()(router)))

}
