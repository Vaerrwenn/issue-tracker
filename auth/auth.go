package auth

import (
	"errors"
	"issue-tracker/models"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JwtWrapper wraps the signing key and the issuer.
type JwtWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

// JwtClaim adds email as a claim to the token.
type JwtClaim struct {
	Email string
	jwt.StandardClaims
}

// GenerateToken generates a JWT token.
func (j *JwtWrapper) GenerateToken(email string) (string, error) {
	claims := &JwtClaim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(j.ExpirationHours)).Unix(),
			Issuer:    j.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// ValidateToken validates the JWT token brought by the Header.
func (j *JwtWrapper) ValidateToken(signedToken string, userID string) (claims *JwtClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)
	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}

	userId, err := strconv.Atoi(userID)
	if err != nil {
		return
	}

	var user models.User
	sourceUser := user.GetUserByID(userId)
	if sourceUser == nil {
		err = errors.New("Could not find User.")
		return
	}

	if claims.Email != sourceUser.Email {
		err = errors.New("INVALID TOKEN. You are not authorized to access this request.")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT is expired")
		return
	}
	return
}
