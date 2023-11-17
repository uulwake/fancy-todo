package service

import (
	"context"
	"fancy-todo/internal/config"
	"fancy-todo/internal/libs"
	"fancy-todo/internal/model"
	"fancy-todo/internal/repository"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func NewUserService(env *config.Env, userRepo IUserRepo) *UserService {
	return &UserService{
		env: env,
		userRepo: userRepo,
	}
}

type UserService struct {
	env *config.Env
	userRepo IUserRepo
}

func (us *UserService) CreateJwtToken(ctx context.Context, userId int64, email string) (string, error) {
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

func (us *UserService) Register(ctx context.Context, data UserServiceRegisterInput) (int64, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(data.Password), us.env.Salt)
	if err != nil {
		return 0, libs.CustomError{
			HTTPCode: http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return us.userRepo.Create(ctx, repository.CreateUserInput{
		Name: data.Name,
		Email: data.Email,
		Password: string(hashedPwd),
	})
}

func (us *UserService) Login(ctx context.Context, data UserServiceLoginInput) (int64, error) {
	var user model.User
	err := us.userRepo.GetDetail(ctx, repository.GetDetailUserInput{
		Email: data.Email,
		Cols: []string{"id", "email", "password"},
		Values: []any{&user.ID, &user.Email, &user.Password},
	})

	if err != nil {
		return 0, libs.CustomError{
			HTTPCode: http.StatusBadRequest,
			Message: "Invalid email/password",
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return 0, libs.CustomError{
			HTTPCode: http.StatusBadRequest,
			Message: "Invalid email/password",
		}
	}

	return user.ID, nil
}