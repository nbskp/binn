package binn

import "github.com/google/uuid"

type Bottle struct {
	Id  string
	Msg string
}

func GenerateID() string {
	uuidObj, _ := uuid.NewUUID()
	return uuidObj.String()
}
