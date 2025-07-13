package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type CEPService struct {
	weatherServiceURL string
	httpClient        *http.Client
}

func NewCEPService(weatherServiceURL string) *CEPService {
	return &CEPService{
		weatherServiceURL: weatherServiceURL,
		httpClient:        &http.Client{},
	}
}

func (s *CEPService) ValidateAndGetWeather(ctx context.Context, cep string) (*WeatherResponse, error) {
	tracer := otel.Tracer("service-a")
	if tracer == nil {
		// Se não há tracer configurado, usar um no-op tracer
		tracer = trace.NewNoopTracerProvider().Tracer("service-a")
	}

	// Span para validação do CEP
	ctx, span := tracer.Start(ctx, "validate_cep")
	defer span.End()

	// Validar formato do CEP
	if !isValidCEP(cep) {
		span.SetAttributes(attribute.String("cep.validation", "invalid"))
		return nil, fmt.Errorf("invalid zipcode")
	}

	span.SetAttributes(attribute.String("cep.validation", "valid"))
	span.SetAttributes(attribute.String("cep.value", cep))

	// Span para chamada ao Serviço B
	ctx, span2 := tracer.Start(ctx, "call_weather_service")
	defer span2.End()

	// Fazer chamada HTTP para o Serviço B
	url := fmt.Sprintf("%s/weather/%s", s.weatherServiceURL, cep)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		span2.SetAttributes(attribute.String("weather_service.error", err.Error()))
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Propagar contexto de tracing (se disponível)
	if propagator := otel.GetTextMapPropagator(); propagator != nil {
		propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		span2.SetAttributes(attribute.String("weather_service.error", err.Error()))
		return nil, fmt.Errorf("failed to call weather service: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		span2.SetAttributes(attribute.String("weather_service.error", "failed to read response body"))
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	span2.SetAttributes(attribute.Int64("weather_service.status_code", int64(resp.StatusCode)))

	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(bodyBytes, &errorResp); err != nil {
			span2.SetAttributes(attribute.String("weather_service.error", "failed to parse error response"))
			return nil, fmt.Errorf("weather service error: %d", resp.StatusCode)
		}
		span2.SetAttributes(attribute.String("weather_service.error", errorResp.Error))
		return nil, fmt.Errorf(errorResp.Error)
	}

	var weatherResp WeatherResponse
	if err := json.Unmarshal(bodyBytes, &weatherResp); err != nil {
		span2.SetAttributes(attribute.String("weather_service.error", "failed to parse weather response"))
		return nil, fmt.Errorf("failed to parse weather response: %w", err)
	}

	span2.SetAttributes(attribute.String("weather_service.city", weatherResp.City))
	span2.SetAttributes(attribute.Float64("weather_service.temp_c", weatherResp.TempC))

	return &weatherResp, nil
}

func isValidCEP(cep string) bool {
	// Remove hífens e espaços
	cleanCEP := strings.ReplaceAll(strings.ReplaceAll(cep, "-", ""), " ", "")

	// Verifica se tem exatamente 8 dígitos
	pattern := `^\d{8}$`
	match, _ := regexp.MatchString(pattern, cleanCEP)
	return match
}

type WeatherResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}
