package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"otel-weather-service/service-b/internal/models"
)

type WeatherService struct {
	weatherAPIKey string
}

func NewWeatherService(weatherAPIKey string) *WeatherService {
	return &WeatherService{
		weatherAPIKey: weatherAPIKey,
	}
}

func (s *WeatherService) GetWeatherByZipcode(ctx context.Context, zipcode string) (*models.WeatherResponse, error) {
	tracer := otel.Tracer("service-b")
	if tracer == nil {
		// Se não há tracer configurado, usar um no-op tracer
		tracer = trace.NewNoopTracerProvider().Tracer("service-b")
	}

	// Span para validação do CEP
	ctx, span := tracer.Start(ctx, "validate_zipcode")
	defer span.End()

	// Validar formato do CEP
	if !isValidZipcode(zipcode) {
		span.SetAttributes(attribute.String("zipcode.validation", "invalid"))
		log.Printf("Invalid zipcode: %s", zipcode)
		return nil, fmt.Errorf("invalid zipcode")
	}

	span.SetAttributes(attribute.String("zipcode.validation", "valid"))
	span.SetAttributes(attribute.String("zipcode.value", zipcode))

	// Span para busca de localização
	ctx, span2 := tracer.Start(ctx, "get_location_by_zipcode")
	defer span2.End()

	// Buscar localização pelo CEP
	location, err := s.getLocationByZipcode(ctx, zipcode)
	if err != nil {
		span2.SetAttributes(attribute.String("location.error", err.Error()))
		log.Printf("Erro ao buscar localização para o CEP %s: %v", zipcode, err)
		return nil, fmt.Errorf("can not find zipcode")
	}

	// Verificar se a localização foi encontrada
	if location.Localidade == "" {
		span2.SetAttributes(attribute.String("location.error", "localidade not found"))
		log.Printf("Localidade não encontrada para o CEP: %s", zipcode)
		return nil, fmt.Errorf("can not find zipcode")
	}

	span2.SetAttributes(attribute.String("location.city", location.Localidade))

	// Span para busca de temperatura
	ctx, span3 := tracer.Start(ctx, "get_temperature_by_location")
	defer span3.End()

	// Buscar temperatura pela localização
	tempC, err := s.getTemperatureByLocation(ctx, location.Localidade)
	if err != nil {
		span3.SetAttributes(attribute.String("temperature.error", err.Error()))
		log.Printf("Erro ao buscar temperatura para a localidade %s (CEP %s): %v", location.Localidade, zipcode, err)
		return nil, err
	}

	// Converter temperaturas
	tempF := tempC*1.8 + 32
	tempK := tempC + 273

	span3.SetAttributes(attribute.Float64("temperature.celsius", tempC))
	span3.SetAttributes(attribute.Float64("temperature.fahrenheit", tempF))
	span3.SetAttributes(attribute.Float64("temperature.kelvin", tempK))

	return &models.WeatherResponse{
		City:  location.Localidade,
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}, nil
}

func (s *WeatherService) getLocationByZipcode(ctx context.Context, zipcode string) (*models.ViaCEPResponse, error) {
	tracer := otel.Tracer("service-b")
	if tracer == nil {
		tracer = trace.NewNoopTracerProvider().Tracer("service-b")
	}

	ctx, span := tracer.Start(ctx, "viacep_request")
	defer span.End()

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", zipcode)
	log.Printf("Consultando ViaCEP: %s", url)

	span.SetAttributes(attribute.String("viacep.url", url))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		span.SetAttributes(attribute.String("viacep.error", err.Error()))
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		span.SetAttributes(attribute.String("viacep.error", err.Error()))
		log.Printf("Erro de requisição HTTP para ViaCEP: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	span.SetAttributes(attribute.Int64("viacep.status_code", int64(resp.StatusCode)))

	if resp.StatusCode != http.StatusOK {
		log.Printf("ViaCEP retornou status %d para o CEP %s", resp.StatusCode, zipcode)
		// Log do corpo da resposta para debug
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Resposta do ViaCEP: %s", string(bodyBytes))

		if resp.StatusCode == http.StatusBadGateway || resp.StatusCode == http.StatusServiceUnavailable {
			span.SetAttributes(attribute.String("viacep.error", "service temporarily unavailable"))
			return nil, fmt.Errorf("viacep service temporarily unavailable")
		}
		span.SetAttributes(attribute.String("viacep.error", "zipcode not found"))
		return nil, fmt.Errorf("zipcode not found")
	}

	var location models.ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		span.SetAttributes(attribute.String("viacep.error", "failed to decode response"))
		log.Printf("Erro ao decodificar resposta do ViaCEP: %v", err)
		return nil, err
	}

	// Verificar se o CEP foi encontrado (ViaCEP retorna erro quando não encontra)
	if location.Cep == "" {
		span.SetAttributes(attribute.String("viacep.error", "zipcode not found"))
		log.Printf("ViaCEP não encontrou o CEP: %s", zipcode)
		return nil, fmt.Errorf("zipcode not found")
	}

	span.SetAttributes(attribute.String("viacep.city", location.Localidade))
	return &location, nil
}

func (s *WeatherService) getTemperatureByLocation(ctx context.Context, location string) (float64, error) {
	tracer := otel.Tracer("service-b")
	if tracer == nil {
		tracer = trace.NewNoopTracerProvider().Tracer("service-b")
	}

	ctx, span := tracer.Start(ctx, "weatherapi_request")
	defer span.End()

	// Usar encoding URL correto para caracteres especiais
	encodedLocation := url.QueryEscape(location)
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", s.weatherAPIKey, encodedLocation)
	log.Printf("Consultando WeatherAPI: %s", url)

	span.SetAttributes(attribute.String("weatherapi.url", url))
	span.SetAttributes(attribute.String("weatherapi.location", location))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		span.SetAttributes(attribute.String("weatherapi.error", err.Error()))
		return 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		span.SetAttributes(attribute.String("weatherapi.error", err.Error()))
		log.Printf("Erro de requisição HTTP para WeatherAPI: %v", err)
		return 0, err
	}
	defer resp.Body.Close()

	span.SetAttributes(attribute.Int64("weatherapi.status_code", int64(resp.StatusCode)))

	if resp.StatusCode != http.StatusOK {
		log.Printf("WeatherAPI retornou status %d para localidade %s", resp.StatusCode, location)
		// Log do corpo da resposta para debug
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Resposta da WeatherAPI: %s", string(bodyBytes))
		span.SetAttributes(attribute.String("weatherapi.error", fmt.Sprintf("status %d", resp.StatusCode)))
		return 0, fmt.Errorf("weather API error: %d", resp.StatusCode)
	}

	var weatherResp models.WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		span.SetAttributes(attribute.String("weatherapi.error", "failed to decode response"))
		log.Printf("Erro ao decodificar resposta da WeatherAPI: %v", err)
		return 0, err
	}

	span.SetAttributes(attribute.Float64("weatherapi.temp_c", weatherResp.Current.TempC))
	return weatherResp.Current.TempC, nil
}

func isValidZipcode(zipcode string) bool {
	// Remove hífens e espaços
	cleanZipcode := strings.ReplaceAll(strings.ReplaceAll(zipcode, "-", ""), " ", "")

	// Verifica se tem exatamente 8 dígitos
	pattern := `^\d{8}$`
	match, _ := regexp.MatchString(pattern, cleanZipcode)
	return match
}
