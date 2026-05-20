package id

import "github.com/google/uuid"

type ID = uuid.UUID

func New() ID {
	return uuid.New()
}

func Parse(value string) (ID, error) {
	return uuid.Parse(value)
}
