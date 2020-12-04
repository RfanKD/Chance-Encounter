package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type User struct {
	Id           int    `json:"Id,omitempty"`
	Name         string `json:"Name,omitempty"`
	Email        string `json:"Email,omitempty"`
	PhoneNumber  string `json:"PhoneNumber,omitempty"`
	Status       string `json:"Status,omitempty"`
	Availability string `json:"Availability,omitempty"`
}

type response struct {
	Numbers []string `json:"numbers"`
}

//	"github.com/aws/aws-lambda-go/lambda"
func main() {
	//lambda.Start(timeDisplayHandler)

	configForNeo4j40 := func(conf *neo4j.Config) { conf.Encrypted = false }

	//user := new([]User)
	driver, err := neo4j.NewDriver("bolt://72.140.181.254:7687", neo4j.BasicAuth("neo4j_dbuser", "password", ""), configForNeo4j40)
	if err != nil {
		fmt.Println(err)
	}
	// handle driver lifetime based on your application lifetime requirements
	// driver's lifetime is usually bound by the application lifetime, which usually implies one driver instance per application
	defer driver.Close()

	// For multidatabase support, set sessionConfig.DatabaseName to requested database
	sessionConfig := neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead}
	session, err := driver.NewSession(sessionConfig)
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()

	callList, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		callList := [][]string{}

		//"WHERE exists((u)-[:Available_at]->(:Time {Start_time1:'"+ time.Now() +"'}))"+

		result, err := tx.Run("Match (u:User)-[r:Encounters_with]->(encounter:User)"+
			"WHERE exists((u)-[:Available_at]->(:Time {Start_time1:'2:00pm'}))"+
			"AND exists((encounter)-[:Available_at]->(:Time {Start_time1:'2:00pm'}))"+
			"Return u.name, u.PhoneNumber, encounter.name, encounter.PhoneNumber, r.Available_at", nil)
		if err != nil {
			log.Fatalln(err)
		}

		for result.Next() {
			numlist := []string{result.Record().GetByIndex(1).(string), result.Record().GetByIndex(3).(string)}
			callList = append(callList, numlist)
		}

		if err = result.Err(); err != nil {
			log.Fatalln(err)
		}

		return callList, nil
	})

	if err != nil {
		log.Fatalln(err)
	}

	interateInterfaceVariable(callList)

}

func interateInterfaceVariable(t interface{}) {
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)

		for i := 0; i < s.Len(); i++ {
			//fmt.Println(s.Index(i))
			y := s.Index(i).Interface().([]string)
			sendNumbers(y)
		}
	}
}

func timeDisplayHandler() {
	//fmt.Println(strings.ToLower(time.Now().Format(time.Kitchen)))
	// configForNeo4j35 := func(conf *neo4j.Config) {}
}

func sendNumbers(numList []string) {

	res1D := &response{
		Numbers: numList}
	requestBody, err := json.Marshal(res1D)
	fmt.Println(string(requestBody))

	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post("https://jsonplaceholder.typicode.com/posts", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

}
