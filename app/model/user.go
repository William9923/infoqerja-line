package model

import (
	"log"

	"github.com/Kamva/mgm/v3"
)

// UserData : A model to represent the user data in the database
type UserData struct {
	mgm.DefaultModel `bson:",inline"`
	State            string `json:"state" bson:"state"`
	SourceID         string `json:"sourceID" bson:"sourceID"`
}

// NewUserData : default constructor for UserData struct
func NewUserData(source string, state string) *UserData {
	return &UserData{
		SourceID: source,
		State:    state,
	}
}

func (user *UserData) Create() error {
	if err := mgm.Coll(user).Create(user); err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (user *UserData) Update() error {
	if err := mgm.Coll(user).Update(user); err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (user *UserData) Delete() error {
	if err := mgm.Coll(&UserData{}).Delete(user); err != nil {
		log.Print(err)
		return err
	}

	return nil
}
