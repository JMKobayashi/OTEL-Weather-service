package handlers

import (
	"net/http"

	"otel-weather-service/service-a/internal/models"
	"otel-weather-service/service-a/internal/services"

	"github.com/gin-gonic/gin"
)

type CEPHandler struct {
	cepService *services.CEPService
}

func NewCEPHandler(cepService *services.CEPService) *CEPHandler {
	return &CEPHandler{
		cepService: cepService,
	}
}

func (h *CEPHandler) GetWeatherByCEP(c *gin.Context) {
	var request models.CEPRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "invalid request format",
		})
		return
	}

	weather, err := h.cepService.ValidateAndGetWeather(c.Request.Context(), request.CEP)
	if err != nil {
		switch err.Error() {
		case "invalid zipcode":
			c.JSON(http.StatusUnprocessableEntity, models.ErrorResponse{
				Error: "invalid zipcode",
			})
		case "can not find zipcode":
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "can not find zipcode",
			})
		default:
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, weather)
}
