package users

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type UserRepository interface {
	RegisterUser(user *User) error
	FindByPhoneNumber(PhoneNumber string) (*User, error)
}

type UserNeo4jRepository struct {
	Driver neo4j.Driver
}

func (u *UserNeo4jRepository) RegisterUser(user *User) (err error) {
	session := u.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer func() {
		err = session.Close()
	}()
	if _, err := session.
		WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			return u.persistUser(tx, user)
		}); err != nil {
		return err
	}
	return nil
}

func (u *UserNeo4jRepository) FindByPhoneNumber(PhoneNumber string) (user *User, err error) {
	session := u.Driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		err = session.Close()
	}()
	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return u.FindByPhone(tx, PhoneNumber)
	})
	if result == nil {
		return nil, err
	}
	user = result.(*User)
	return user, err
}

func (u *UserNeo4jRepository) persistUser(tx neo4j.Transaction, user *User) (interface{}, error) {
	query := "CREATE (:User {Email: $Email, Name: $Name, PhoneNumber: $PhoneNumber, Status: $Status, Availability: $Availability})"

	parameters := map[string]interface{}{
		"Email":        user.Email,
		"Name":         user.Name,
		"PhoneNumber":  user.PhoneNumber,
		"Status":       user.Status,
		"Availability": user.Availability,
	}
	_, err := tx.Run(query, parameters)
	return nil, err
}

func (u *UserNeo4jRepository) FindByPhone(tx neo4j.Transaction, PhoneNumber string) (*User, error) {
	query := "MATCH (u:User {PhoneNumber: $PhoneNumber}) RETURN  u.Name AS Name, u.Email AS Email, u.PhoneNumber AS PhoneNumber, u.Status AS Status, u.Availability AS Availability"

	parameters := map[string]interface{}{
		"PhoneNumber": PhoneNumber,
	}
	result, err := tx.Run(
		query,
		parameters,
	)
	if err != nil {
		return nil, err
	}
	record, err := result.Single()
	if err != nil {
		return nil, err
	}

	//returnedUser := record.GetByIndex(0)
	//json.Marshal(&returnedUser)
	//fmt.Println(record)

	Name, _ := record.Get("Name")
	Status, _ := record.Get("Status")
	Email, _ := record.Get("Email")
	Availability, _ := record.Get("Availability")
	return &User{
		Name:         Name.(string),
		PhoneNumber:  PhoneNumber,
		Status:       Status.(string),
		Email:        Email.(string),
		Availability: Availability.(string),
	}, nil

}
