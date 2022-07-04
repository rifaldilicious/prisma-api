package posts

import (
	"crypto-project/helpers"
	"crypto-project/interfaces"
	"gorm.io/gorm"
	"errors"
	"math/rand"
	"time"
)

func prepareCreatePostResponse(post *interfaces.Post) map[string]interface{} {
	responsePost := &interfaces.ResponseCreatePost{
		ID:   		post.ID,
	}

	var response = map[string]interface{}{"message": "all is fine"}
	response["data"] = responsePost

	return response
}

func prepareReadPostResponse(post *interfaces.Post) map[string]interface{} {
	responseReadPost := &interfaces.ResponseReadPost{
		ID:   		post.ID,
		User_ID:   	post.User_ID,
		Name:   	post.Name,
		Skill:   	post.Skill,
		Location:   post.Location,
		Position:   post.Position,
		Work:   	post.Work,
		Salary:   	post.Salary,
		Message:   	post.Message,
	}

	var response = map[string]interface{}{"message": "all is fine"}
	response["data"] = responseReadPost

	return response
}

func prepareDeletePostResponse(post *interfaces.Post) map[string]interface{} {
	responsePost := &interfaces.ResponseDeletePost{
		ID:   		post.ID,
	}

	var response = map[string]interface{}{"message": "all is fine"}
	response["data"] = responsePost

	return response
}

func prepareUpdatePostResponse(post *interfaces.Post) map[string]interface{} {
	responsePost := &interfaces.ResponseUpdatePost{
		User_ID:   	post.User_ID,
		Name:   	post.Name,
		Skill:   	post.Skill,
		Location:   post.Location,
		Position:   post.Position,
		Work:   	post.Work,
		Salary:   	post.Salary,
		Message:   	post.Message,
	}

	var response = map[string]interface{}{"message": "all is fine"}
	response["data"] = responsePost

	return response
}

func CreatePost(user_id uint, name string, skill string, location string, position string, work string, salary uint, message string) map[string]interface{} {

	//generate random number for post ID
	rand.Seed(time.Now().UnixNano())
	id := uint(rand.Intn(99999999 - 10000000) + 10000000)

	db 		:= helpers.ConnectDB()	
	post 	:= &interfaces.Post{ID: id, User_ID: user_id, Name: name, Skill: skill, Location: location, Position: position, Work: work, Salary: salary, Message: message}
	db.Create(&post)

	return prepareCreatePostResponse(post)
}

func ReadPost(id string) map[string]interface{} {

	db 		:= helpers.ConnectDB()
	post 	:= &interfaces.Post{}
	
	err 	:= db.First(&post, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return map[string]interface{}{"message": "record not found"}
	}

	return prepareReadPostResponse(post)
}

func UpdatePost(id string, user_id uint, name string, skill string, location string, position string, work string, salary uint, message string) map[string]interface{} {

	db 		:= helpers.ConnectDB()	
	post 	:= &interfaces.Post{}
	
	err := db.First(&post, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return map[string]interface{} {"message": "record not found"}
	}

	post 	= &interfaces.Post{User_ID: user_id, Name: name, Skill: skill, Location: location, Position: position, Work: work, Salary: salary, Message: message}
	db.Where("id = ?", id).Updates(&post)
	return prepareUpdatePostResponse(post)
}

func DeletePost(id string) map[string]interface{} {

	db 		:= helpers.ConnectDB()
	post 	:= &interfaces.Post{}
	
	err := db.First(&post, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return map[string]interface{} {"message": "record not found"}
	}

	db.Where("id = ?", id).Delete(&post)
	return prepareDeletePostResponse(post)
}