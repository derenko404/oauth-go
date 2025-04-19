package controllers

import (
	"oauth-go/internal/app"

	"github.com/gin-gonic/gin"
)

type healthController struct {
	app *app.App
}

func NewHelathController(app *app.App) *healthController {
	return &healthController{
		app: app,
	}
}

type healthCheckResponse struct {
	Status string `json:"status" example:"ok"`
}

// @Summary		Health
// @Description	App health check
// @Tags			  health
// @Accept			json
// @Produce		  json
// @Success     200 {object} healthCheckResponse "All is ok"
// @Failure 500 {object} healthCheckResponse "Database error"
// @Router			/health [get]
func (controller *healthController) HealthCheck(c *gin.Context) {
	err := controller.app.DB.Ping(c.Request.Context())
	if err != nil {
		controller.app.Logger.Error("db ping error", "error", err)

		c.JSON(500, &healthCheckResponse{
			Status: "db error",
		})
		return
	}

	_, err = controller.app.RDB.Ping(c.Request.Context()).Result()

	if err != nil {
		controller.app.Logger.Error("redis ping error", "error", err)

		c.JSON(500, &healthCheckResponse{
			Status: "redis error",
		})
		return
	}

	c.JSON(200, &healthCheckResponse{
		Status: "ok",
	})
}
