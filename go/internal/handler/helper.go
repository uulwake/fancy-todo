package handler

import (
	"context"
	"fancy-todo/internal/libs"

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