package service

import (
	"context"
	"fancy-todo/internal/config"
	"fancy-todo/internal/libs"
	"fancy-todo/internal/repository"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func NewUserService(env *config.Env, userRepo UserRepo) *UserService {
	return &UserService{
		env: env,
		userRepo: userRepo,
	}
}

type UserService struct {
	env *config.Env
	userRepo UserRepo
}

func (us *UserService) CreateJwtToken(ctx context.Context, userId int, email string) (string, error) {
	fmt.Println("User Service: CreateJwtToken")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": userId,
		"email": email,
		"exp": time.Now().Add(time.Duration(us.env.JwtExpired * int(time.Hour))).Unix(),
	})

	stringToken, err := token.SignedString([]byte(us.env.JwtSecret))
	if err != nil {
		return "", err
	}

	return stringToken, nil
}

func (us *UserService) Register(ctx context.Context, data UserServiceRegisterInput) (int, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(data.Password), us.env.Salt)
	if err != nil {
		return 0, libs.CustomError{
			HTTPCode: http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	us.userRepo.Create(ctx, repository.CreateUserInput{
		Name: data.Name,
		Email: data.Email,
		Password: string(hashedPwd),
	})
	return 1, nil
}