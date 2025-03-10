package route

import (
	"fit-tracker/api/controller"

	"github.com/labstack/echo/v4"
)

func RegisterUserRoutes(g *echo.Group, uc *controller.UserController) {
  g.GET("/:userID", uc.GetUserData)
}
