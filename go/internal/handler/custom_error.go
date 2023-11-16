package handler

import (
	"fancy-todo/internal/libs"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorMsg struct {
	Message string `json:"message"`
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	customErr, ok := err.(libs.CustomError)
	if !ok {
		err := c.JSON(http.StatusInternalServerError, "Internal server error")
		if err != nil {
			c.Logger().Error(err)
		}

		return 
	}

	err = c.JSON(customErr.HTTPCode, ErrorMsg{
		Message: customErr.Message,
	})

	if err != nil {
		c.Logger().Error(err)
	}
}