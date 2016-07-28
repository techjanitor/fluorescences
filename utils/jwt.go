package utils

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	// jwt header keys
	jwtHeaderKeyID = "kid"
	// jwt issuer
	jwtIssuer = "fluorescences"
	// jwt expire days
	jwtExpireDays = 90
)

// TokenClaims holds the custom and standard claims for the JWT token
type TokenClaims struct {
	jwt.StandardClaims
}

// MakeToken will create a JWT token
func MakeToken() (newtoken string, err error) {

	// the current timestamp
	now := time.Now()

	claims := TokenClaims{
		jwt.StandardClaims{
			Issuer:    jwtIssuer,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: now.Add(time.Hour * 24 * jwtExpireDays).Unix(),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// set our header info
	token.Header[jwtHeaderKeyID] = 1

	secret, err := GetSecret()
	if err != nil {
		return
	}

	return token.SignedString(secret)

}

// ValidateToken checks all the claims in the provided token
func ValidateToken(token *jwt.Token) ([]byte, error) {

	// check alg to make sure its hmac
	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	// get the claims from the token
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, fmt.Errorf("Couldnt parse claims")
	}

	// get the issuer from claims
	tokenIssuer := claims.Issuer

	// check the issuer
	if tokenIssuer != jwtIssuer {
		return nil, fmt.Errorf("Incorrect issuer")
	}

	// get the stored HMAC secret
	secret, err := GetSecret()
	if err != nil {
		return nil, err
	}

	// compare with secret from settings
	return secret, nil

}
