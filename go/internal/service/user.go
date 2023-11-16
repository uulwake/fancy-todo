package service

import (
	"context"
	"fancy-todo/internal/config"
	"fancy-todo/internal/repository"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewUserService(env *config.Env, userRepo UserRepo) *UserService {
	return &UserService{
		Env: env,
		UserRepo: userRepo,
	}
}

type UserService struct {
	Env *config.Env
	UserRepo UserRepo
}

func (us *UserService) Register(ctx context.Context, data UserServiceRegisterInput) (int, error) {
	fmt.Println("User Service: Register")
	us.UserRepo.Create(ctx, repository.CreateUserInput{})
	return 1, nil
}

func (us *UserService) CreateJwtToken(ctx context.Context, userId int, email string) (string, error) {
	fmt.Println("User Service: CreateJwtToken")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": userId,
		"email": email,
		"exp": time.Now().Add(time.Duration(us.Env.JwtExpired * int(time.Hour))).Unix(),
	})

	stringToken, err := token.SignedString([]byte(us.Env.JwtSecret))
	if err != nil {
		return "", err
	}

	return stringToken, nil
}