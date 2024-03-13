package middlewares

import (
	"github.com/SawitProRecruitment/UserService/constants"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func ValidateContentType() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			URI := c.Request().RequestURI
			prefix := strings.Split(URI, "/")[1]

			if prefix == "swagger" {
				return next(c)
			}

			if contentType := c.Request().Header.Get("Content-Type"); contentType != "application/json" {
				return c.JSON(http.StatusBadRequest, generated.ErrorResponse{
					Message: constants.ContentTypeJson,
				})
			}
			return next(c)
		}
	}
}
