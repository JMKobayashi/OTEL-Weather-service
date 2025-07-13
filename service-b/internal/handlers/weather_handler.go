package handlers

import (
	"net/http"

	"otel-weather-service/service-b/internal/services"

	"github.com/gin-gonic/gin"
)

type WeatherHandler struct {
	weatherService *services.WeatherService
}

func NewWeatherHandler(weatherService *services.WeatherService) *WeatherHandler {
	return &WeatherHandler{
		weatherService: weatherService,
	}
}

func (h *WeatherHandler) GetWeather(c *gin.Context) {
	zipcode := c.Param("zipcode")

	weather, err := h.weatherService.GetWeatherByZipcode(c.Request.Context(), zipcode)
	if err != nil {
		switch err.Error() {
		case "invalid zipcode":
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid zipcode"})
		case "can not find zipcode":
			c.JSON(http.StatusNotFound, gin.H{"error": "can not find zipcode"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, weather)
}
