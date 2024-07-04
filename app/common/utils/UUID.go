package utils

import (
	"github.com/google/uuid"
)

func GenerateUUID() uuid.UUID {
	return uuid.New()
}

func ParseUUID(stringUUID *string) uuid.UUID {
	if *stringUUID == "" || *stringUUID == "null" {
		return uuid.Nil
	}
	parsedUUID, err := uuid.Parse(*stringUUID)
	if err != nil {
		return uuid.Nil
	}
	return parsedUUID
}

func IsUUIDNil(u uuid.UUID) bool {
	return u == uuid.Nil
}

func IsValidUUID(rawUUID string) bool {
	_, err := uuid.Parse(rawUUID)
	return err == nil
}
