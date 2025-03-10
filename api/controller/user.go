package controller

import (
	"context"
	"net/http"
	"time"

	"fit-tracker/api/service"

	"github.com/labstack/echo/v4"
)

type (
	UserService interface {
		GetUserData(ctx context.Context, input *service.GetUserDataInput) (*service.GetUserDataResult, error)
	}

	UserController struct {
		userService UserService
	}

	GetUserDataRequest struct {
		UserID string  `param:"userID"`
		Date   string  `query:"date"`
		Weight float64 `query:"weight"`
	}

	GetUserDataResponse struct {
		Steps            int64   `json:"steps"`
		Distance         float64 `json:"distance"`
		AverageHeartBeat float64 `json:"averageHeartBeat"`
		KcalBurned       float64 `json:"kcalBurned"`
	}
)

func New(userService UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (u *UserController) GetUserData(c echo.Context) error {
	var req = new(GetUserDataRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	reqDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid date value"})
	}

	if req.Weight == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid weight value"})
	}

	res, err := u.userService.GetUserData(c.Request().Context(), &service.GetUserDataInput{
		UserID: req.UserID,
		Date:   reqDate,
		Weight: req.Weight,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, GetUserDataResponse{
		Steps:            res.Steps,
		Distance:         res.Distance,
		AverageHeartBeat: res.AverageHeartBeat,
		KcalBurned:       res.KcalBurned,
	})
}
