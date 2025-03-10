package api

import (
	"context"
	"net/http"
	"time"

	user_api_controller "fit-tracker/api/controller"
	api_routes "fit-tracker/api/route"
	user_api_service "fit-tracker/api/service"
	"fit-tracker/database"

	"github.com/labstack/echo/v4"
)

const API_VERSION = "v1"

type (
	api struct {
		DB *database.Database
	}
)

func New(db *database.Database) *api {
	return &api{
		DB: db,
	}
}

func (a *api) Run(ctx context.Context) {
	e := echo.New()

	userAPIGroup := e.Group("/api/" + API_VERSION + "/user")

	userAPIService := user_api_service.New(user_api_service.WithIngestorRepository(a.DB.IngestorRepository))
	userAPIController := user_api_controller.New(userAPIService)
	api_routes.RegisterUserRoutes(userAPIGroup, userAPIController)

	go func() {
		if err := e.Start(":8081"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
