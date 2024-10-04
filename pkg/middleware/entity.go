package middleware

import (
	"fmt"
	"github.com/mikestefanello/pagoda/pkg/context"
	"github.com/mikestefanello/pagoda/pkg/db/sqlc"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// LoadUser loads the user based on the ID provided as a path parameter
func LoadUser(queries *sqlc.Queries) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, err := strconv.Atoi(c.Param("user"))
			if err != nil {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			u, err := queries.Auth_GetUserById(c.Request().Context(), int64(userID))

			switch err.(type) {
			case nil:
				c.Set(context.UserKey, u)
				return next(c)
			/*case *ent.NotFoundError:
			return echo.NewHTTPError(http.StatusNotFound)*/
			default:
				return echo.NewHTTPError(
					http.StatusInternalServerError,
					fmt.Sprintf("error querying user: %v", err),
				)
			}
		}
	}
}
