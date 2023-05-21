package utils

import (
	"github.com/google/uuid"
)

func BinUuidToString(uuidBin []byte) (string, error) {
	id, err := uuid.FromBytes(uuidBin)

	if err != nil {
		return "", err
	}

	return id.String(), nil
}
