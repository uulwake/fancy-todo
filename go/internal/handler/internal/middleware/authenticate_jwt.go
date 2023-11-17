package middleware

import (
	"fancy-todo/internal/config"
	"fancy-todo/internal/libs"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AuthenticateJwt(env *config.Env) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			jwtToken, ok := c.Request().Header["Jwt-Token"]
			if !ok {
				return libs.CustomError{
					HTTPCode: http.StatusUnauthorized,
					Message: "Jwt-Token is not found in request header",
				}
			}
			
			claims := jwt.MapClaims{}
			_, err := jwt.ParseWithClaims(jwtToken[0], claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(env.JwtSecret), nil
			})

			if err != nil {
				return libs.CustomError{
					HTTPCode: http.StatusUnauthorized,
					Message: err.Error(),
				}
			}

			id, ok := claims["id"]
			if !ok {
				return libs.CustomError{
					HTTPCode: http.StatusUnauthorized,
					Message: "Jwt claims does not have id",
				}
			}
			c.Set("user_id", id)

			email, ok := claims["email"]
			if !ok {
				return libs.CustomError{
					HTTPCode: http.StatusUnauthorized,
					Message: "Jwt claims does not have email",
				}
			}
			c.Set("user_email", email)

			return next(c)
		}
	}
	 
}