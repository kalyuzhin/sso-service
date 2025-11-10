package refreshtoken

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateRefreshToken() ([]byte, error) {
	refreshToken, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(refreshToken.String()), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}
