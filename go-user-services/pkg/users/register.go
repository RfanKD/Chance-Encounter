package users

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

//func(writer http.ResponseWriter, request *http.Request) {
//
//}

type UserRegistration struct {
	User User `json:"user"`
}

type User struct {
	Id           int    `json:"Id,omitempty"`
	Name         string `json:"Name,omitempty"`
	Email        string `json:"Email,omitempty"`
	PhoneNumber  string `json:"PhoneNumber,omitempty"`
	Status       string `json:"Status,omitempty"`
	Availability string `json:"Availability,omitempty"`
}

type UserRegistrationHandler struct {
	Path           string
	UserRepository UserRepository
}

func (u *UserRegistrationHandler) Register(writer http.ResponseWriter, request *http.Request) {
	requestBody, _ := ioutil.ReadAll(request.Body)
	userRegistrationRequest := UserRegistration{}
	_ = json.Unmarshal(requestBody, &userRegistrationRequest)
	requestUser := userRegistrationRequest.User
	if len(requestUser.Name) < 1 {
		log.Println(" 'Name' is missing")
		return
	}
	if len(requestUser.PhoneNumber) < 1 {
		log.Println(" 'PhoneNumber' is missing")
		return
	}

	_ = u.UserRepository.RegisterUser(&requestUser)

	writer.WriteHeader(201)
	writer.Header().Add("Content-Type", "application/json")
	userRegistrationResponse := UserRegistration{
		User: User{
			Name:         requestUser.Name,
			Email:        requestUser.Email,
			PhoneNumber:  requestUser.PhoneNumber,
			Status:       requestUser.Status,
			Availability: requestUser.Availability,
		}}
	bytes, _ := json.Marshal(&userRegistrationResponse)
	_, _ = writer.Write(bytes)
}
