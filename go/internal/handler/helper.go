package handler

import (
	"context"
	"encoding/json"
	"fancy-todo/internal/libs"
	"fmt"

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
	
	if body != nil {
		err := json.NewDecoder(c.Request().Body).Decode(&body)
		if err != nil {
			return ctx, libs.DefaultInternalServerError(err)
		}
		
		err = validate.Struct(body)
		if err != nil {
			return ctx, libs.DefaultInternalServerError(err)
		}
	}

	return ctx, nil
}

func GetUserIdFromContext(c echo.Context) (int64, error) {
	userId := c.Get("user_id")

	switch v := userId.(type) {
	case int64:
		return v, nil
	default:
		return 0, libs.DefaultInternalServerError(fmt.Errorf("invalid userID %v", userId))
	}
}