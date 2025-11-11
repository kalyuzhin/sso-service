package refreshtoken

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// GenerateRefreshToken – ...
func GenerateRefreshToken() (refreshTokenStr string, hash []byte, err error) {
	refreshToken, err := uuid.NewRandom()
	if err != nil {
		return refreshTokenStr, hash, err
	}

	hash, err = bcrypt.GenerateFromPassword([]byte(refreshToken.String()), bcrypt.DefaultCost)
	if err != nil {
		return refreshTokenStr, hash, err
	}

	return refreshToken.String(), hash, err
}
