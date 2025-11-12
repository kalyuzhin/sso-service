package model

// JWK – ...
type JWK struct {
	Kty string `json:"kty"`
	KID string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// JWKS – ...
type JWKS struct {
	Keys []JWK `json:"keys"`
}
