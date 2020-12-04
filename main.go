package main

import (
	"fmt"
	"net/http"

	"github.com/neo4j-examples/golang-neo4j-realworld-example/pkg/users"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func main() {
	neo4jUri := "bolt://72.140.181.254:7687/"
	neo4jUsername := "neo4j_dbuser"
	neo4jPassword := "password"

	usersRepository := users.UserNeo4jRepository{
		Driver: driver(neo4jUri, neo4j.BasicAuth(neo4jUsername, neo4jPassword, "")),
	}

	registrationHandler := &users.UserRegistrationHandler{
		Path:           "/users",
		UserRepository: &usersRepository,
	}
	userHandler := &users.UserHandler{
		Path:           "/findUser",
		UserRepository: &usersRepository,
	}

	server := http.NewServeMux()
	server.HandleFunc(registrationHandler.Path, registrationHandler.Register)
	server.HandleFunc(userHandler.Path, userHandler.FindUser)
	server.HandleFunc("/", homePage)
	if err := http.ListenAndServe(":3000", server); err != nil {
		panic(err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Chaos Monkeys Web API!")
}

func driver(target string, token neo4j.AuthToken) neo4j.Driver {
	result, err := neo4j.NewDriver(target, token)
	if err != nil {
		panic(err)
	}
	return result
}
