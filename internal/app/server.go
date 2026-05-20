package app

import (
	"database/sql"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"mertani/internal/config"
	"mertani/internal/device"
	"mertani/internal/sensor"
	"mertani/internal/shared/response"
	"mertani/internal/shared/validator"
)

type Server struct {
	echo *echo.Echo
	cfg  config.Config
	db   *sql.DB
}

func NewServer(cfg config.Config, db *sql.DB) *Server {
	e := echo.New()
	e.Validator = validator.New()
	e.HTTPErrorHandler = response.ErrorHandler
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS("*"))

	server := &Server{
		echo: e,
		cfg:  cfg,
		db:   db,
	}
	server.registerRoutes()

	return server
}

func (s *Server) Start() error {
	return s.echo.Start(s.cfg.ServerAddress())
}

func (s *Server) registerRoutes() {
	s.echo.GET("/health", health)
	s.registerOpenAPI()

	api := s.echo.Group("/api/v1")

	deviceRepository := device.NewPostgresRepository(s.db)
	deviceService := device.NewService(deviceRepository)
	deviceHandler := device.NewHandler(deviceService)
	deviceHandler.RegisterRoutes(api)

	sensorRepository := sensor.NewPostgresRepository(s.db)
	sensorService := sensor.NewService(sensorRepository, deviceRepository)
	sensorHandler := sensor.NewHandler(sensorService)
	sensorHandler.RegisterRoutes(api)
}

func health(c *echo.Context) error {
	return response.OK(c, "Service is healthy", map[string]any{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
