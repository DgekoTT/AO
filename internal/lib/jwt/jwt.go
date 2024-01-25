package jwtM

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
	"yourTeamAuth/internal/domain/models"
)

func NewAccessToken(userN models.Users, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = userN.ID.String()
	claims["email"] = userN.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["roles"] = userN.Roles

	tokenString, err := token.SignedString([]byte("sfsfsffffsss--3452!!"))
	if err != nil {
		return "0", err
	}

	return tokenString, nil
}

func NewRefreshToken(userN models.Users, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = userN.ID.String()
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte("sfsfsffffsss--3452!!"))
	if err != nil {
		return "0", err
	}

	return tokenString, nil
}
