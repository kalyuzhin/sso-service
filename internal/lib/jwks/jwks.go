package jwks

import (
	"crypto/rsa"
	"encoding/base64"
	"math/big"

	"github.com/kalyuzhin/sso-service/internal/model"
)

const (
	kty = "RSA"
	kid = "ab1c1"
	use = "sig"
	alg = "rs256"
)

// MakeJWKS – ...
func MakeJWKS(key *rsa.PublicKey) *model.JWKS {
	return &model.JWKS{
		Keys: []model.JWK{
			{
				Kty: kty,
				Use: use,
				Alg: alg,
				KID: kid,
				N:   convertToBase64URL(key.N.Bytes()),
				E:   convertToBase64URL(big.NewInt(int64(key.E)).Bytes()),
			},
		},
	}
}

func convertToBase64URL(bytes []byte) string {
	return base64.RawURLEncoding.EncodeToString(bytes)
}
