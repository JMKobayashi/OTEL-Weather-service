package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"otel-weather-service/service-b/internal/handlers"
	"otel-weather-service/service-b/internal/services"
)

func main() {
	// Configurar Gin para modo release em produção
	gin.SetMode(gin.ReleaseMode)

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Configurar OTEL (opcional)
	if err := setupOTEL(); err != nil {
		log.Printf("Warning: Failed to setup OTEL: %v", err)
		log.Println("Continuing without OTEL tracing...")
	} else {
		log.Println("OTEL tracing configured successfully")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Starting Service B on port %s", port)

	weatherAPIKey := os.Getenv("WEATHER_API_KEY")
	if weatherAPIKey == "" {
		log.Fatal("WEATHER_API_KEY environment variable is required")
	}

	log.Println("Weather API Key configured successfully")

	weatherService := services.NewWeatherService(weatherAPIKey)
	weatherHandler := handlers.NewWeatherHandler(weatherService)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "service-b",
		})
	})

	r.GET("/weather/:zipcode", weatherHandler.GetWeather)

	// Configurar graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		if err := r.Run("0.0.0.0:" + port); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	select {
	case <-sigCh:
		log.Println("Shutting down gracefully, CTRL+C pressed...")
	case <-time.After(30 * time.Second):
		log.Println("Shutting down due to timeout...")
	}
}

func setupOTEL() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("service-b"),
		),
	)
	if err != nil {
		return err
	}

	// Aguardar um pouco para o OTEL Collector inicializar
	time.Sleep(5 * time.Second)

	conn, err := grpc.DialContext(ctx, "otel-collector:4317",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(20*time.Second),
	)
	if err != nil {
		return err
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return err
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return nil
}
