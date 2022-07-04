package api

import (
	"crypto-project/helpers"
	"crypto-project/interfaces"
	"crypto-project/users"
	"crypto-project/posts"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type Login struct {
	Username string
	Password string
}

type Register struct {
	Username string
	Email    string
	Password string
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
	register := users.Register(formattedBody.Username, formattedBody.Email, formattedBody.Password)
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

func createPost(w http.ResponseWriter, r *http.Request) {
	body := readBody(r)
	var formattedBody Post
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)
	createPost := posts.CreatePost(formattedBody.User_ID, formattedBody.Name, formattedBody.Skill, formattedBody.Location, formattedBody.Position, formattedBody.Work, formattedBody.Salary, formattedBody.Message)
	apiResponse(createPost, w)
}

func readPost(w http.ResponseWriter, r *http.Request) {

	vars 	:= mux.Vars(r)
	postId 	:= vars["id"]
	post 	:= posts.ReadPost(postId)
	apiResponse(post, w)

}

func deletePost(w http.ResponseWriter, r *http.Request) {
	vars 	:= mux.Vars(r)
	postId 	:= vars["id"]
	post 	:= posts.DeletePost(postId)
	apiResponse(post, w)

}

func updatePost(w http.ResponseWriter, r *http.Request) {

	var formattedBody Post

	vars 	:= mux.Vars(r)
	postId 	:= vars["id"]
	body 	:= readBody(r)
	err 	:= json.Unmarshal(body, &formattedBody)

	helpers.HandleErr(err)

	updatePost := posts.UpdatePost(postId, formattedBody.User_ID, formattedBody.Name, formattedBody.Skill, formattedBody.Location, formattedBody.Position, formattedBody.Work, formattedBody.Salary, formattedBody.Message)
	apiResponse(updatePost, w)
}

func StartApi() {
	router := mux.NewRouter()
	router.Use(helpers.PanicHandler)
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/register", register).Methods("POST")
	router.HandleFunc("/user/{id}", getUser).Methods("GET")

	//POST
	router.HandleFunc("/post", createPost).Methods("POST")
	router.HandleFunc("/post/{id}", readPost).Methods("GET")
	router.HandleFunc("/post/{id}", deletePost).Methods("DELETE")
	router.HandleFunc("/post/{id}", updatePost).Methods("PUT")

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
