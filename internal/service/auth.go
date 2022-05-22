package service

import (
	"crypto/sha1"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/model"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/store"
	"time"
)

const (
	salt       = "asfafihfas8ga7sg7ashgah9s"
	signingKey = "alsgjasg8gas6k5whwgjaso9"
	tokenTTL   = 12 * time.Hour
)

type Authorization interface {
	CreateUser(user model.User) (uint64, error)
	GenerateToken(username, password string) (string, error)
}

type AuthService struct {
	repo store.AuthRepository
}

func NewAuthService(repo store.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user model.User) (uint64, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId uint64 `json:"user_id"`
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
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

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
