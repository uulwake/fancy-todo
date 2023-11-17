package handler

import (
	"context"
	"encoding/json"
	"fancy-todo/internal/libs"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func CreateContext(c echo.Context) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, libs.RequestId, c.Response().Header().Get(echo.HeaderXRequestID))
	ctx = context.WithValue(ctx, libs.IpAddress, c.RealIP())
	ctx = context.WithValue(ctx, libs.UserID, c.Get("user_id"))
	ctx = context.WithValue(ctx, libs.UserID, c.Get("user_email"))
	return ctx
}

func PreprocessedRequest(c echo.Context, validate *validator.Validate, body any) (context.Context, error) {
	ctx := CreateContext(c)
	
	err := json.NewDecoder(c.Request().Body).Decode(&body)
	if err != nil {
		return ctx, libs.CustomError{
			HTTPCode: http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	
	err = validate.Struct(body)
	if err != nil {
		return ctx, libs.CustomError{
			HTTPCode: http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	return ctx, nil

}