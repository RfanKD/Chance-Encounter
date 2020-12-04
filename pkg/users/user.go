package users

import (
	"encoding/json"
	"log"
	"net/http"
)

type UserHandler struct {
	Path           string
	UserRepository UserRepository
}

func (u *UserHandler) FindUser(writer http.ResponseWriter, request *http.Request) {
	keys, ok := request.URL.Query()["PhoneNumber"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'PhoneNumber' is missing")
		return
	}
	PhoneNumber := keys[0]
	//log.Println("Url Param 'PhoneNumber' is " + PhoneNumber)
	user, _ := u.UserRepository.FindByPhoneNumber(PhoneNumber)

	writer.WriteHeader(200)
	writer.Header().Add("Content-Type", "application/json")
	userResponse := User{
		Id:           user.Id,
		Name:         user.Name,
		Email:        user.Email,
		PhoneNumber:  user.PhoneNumber,
		Status:       user.Status,
		Availability: user.Availability,
	}

	bytes, _ := json.Marshal(&userResponse)
	_, _ = writer.Write(bytes)
}
