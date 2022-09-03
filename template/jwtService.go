package service

import (
	"[project]/entity"
	"errors"

	"github.com/dgrijalva/jwt-go"
)

type JwtService interface {
	GenerateToken(userId int, userName string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
}

func NewJwtService() *jwtService {
	return &jwtService{}
}

var secret = []byte("your_token")

func (s *jwtService) GenerateToken(userId int, userName string) (string, error) {
	claim := jwt.MapClaims{}

	claim["user_id"] = userId
	claim["user_name"]= userName

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString(secret)

	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

func (s *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	encodeToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return encodeToken, err
	}

	return encodeToken, nil
}