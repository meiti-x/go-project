package handlers

import (
	"net/http"
	"os"

	"agentic/commerce/internal/utils"

	"github.com/go-json-experiment/json/v1"
	"github.com/labstack/echo/v4"
)

type ICommonResource interface {
	GetCities() echo.HandlerFunc
}

type commonResource struct{}

func NewCommonResource() ICommonResource {
	return &commonResource{}
}

type City struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	StateID   int    `json:"state_id"`
	StateName string `json:"state_name"`
}
type CitiesData struct {
	Cities []City `json:"cities"`
}

func (v *commonResource) GetCities() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cwd, _ := os.Getwd()
		citiesData, err := os.ReadFile(cwd + "/cities.json")
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to read JSON file: " + err.Error(),
			})
		}

		var data CitiesData
		if err := json.Unmarshal(citiesData, &data); err != nil {
			return utils.ErrorResponse(ctx, err, err.Error())
		}

		return utils.SuccessResponse(ctx, data)
	}
}
