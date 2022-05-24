//go:generate mockgen -source=auth.go -destination=mocks/auth_mock.go
package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/model"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/store"
	"time"
)

const (
	salt       = "asfafihfas8ga7sg7ashgah9s"
	signingKey = "alsgjasg8gas6k5whwgjaso9"
	tokenTTL   = 2 * time.Hour
)

var (
	ErrSigningMethod        = errors.New("invalid signing method")
	ErrWrongAccessTokenType = errors.New("token claims are not of type *tokenClaims")
)

type Authorization interface {
	CreateUser(user model.User) (uint64, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(accessToken string) (uint64, error)
}

type AuthService struct {
	userRepository store.UserRepository
}

func NewAuthService(repo store.UserRepository) *AuthService {
	return &AuthService{userRepository: repo}
}

func (s *AuthService) CreateUser(user model.User) (uint64, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.userRepository.Create(user)
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId uint64 `json:"user_id"`
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.userRepository.Get(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (uint64, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrSigningMethod
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, ErrWrongAccessTokenType
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
