package jwt

import (
	errorpkg "github.com/kalyuzhin/sso-service/internal/error"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/kalyuzhin/sso-service/internal/model"
)

// GenerateToken – ...
func GenerateToken(app model.App, user model.User, ttl time.Duration) (token string, err error) {
	tokenObj := jwt.New(jwt.SigningMethodHS256)

	claims := tokenObj.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(ttl).Unix()
	claims["app_id"] = app.ID

	token, err = tokenObj.SignedString([]byte("ключ"))
	if err != nil {
		return token, errorpkg.WrapErr(err, "can't convert token object into string")
	}

	return token, nil
}
