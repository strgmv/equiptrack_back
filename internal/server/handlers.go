package server

import (
	authHttp "equiptrack/internal/auth/delivery/http"
	authRepository "equiptrack/internal/auth/repository"
	authUseCase "equiptrack/internal/auth/usecase"
	apiMiddlewares "equiptrack/internal/middleware"

	equipHttp "equiptrack/internal/equipment/delivery/http"
	equipRepository "equiptrack/internal/equipment/repository"
	equipUseCase "equiptrack/internal/equipment/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) MapHandlers(e *echo.Echo) error {

	// Init repositories
	aRepo := authRepository.NewAuthRepository(s.db)
	eRepo := equipRepository.NewEquipmentRepository(s.db)

	// Init useCases
	authUC := authUseCase.NewAuthUseCase(s.cfg, aRepo, s.logger)
	equipUC := equipUseCase.NewEquipmentUseCase(s.cfg, eRepo, s.logger)

	// Init handlers
	authHandlers := authHttp.NewAuthHandlers(s.cfg, authUC, s.logger)
	equipmentHandlers := equipHttp.NewEquipmentHandlers(s.cfg, equipUC, s.logger)

	mw := apiMiddlewares.NewMiddlewareManager(authUC, s.cfg, []string{"*"}, s.logger)

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	e.Use(middleware.RequestID())

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Use(middleware.Secure())
	e.Use(mw.RequestLoggerMiddleware)
	// e.Use(middleware.BodyLimit("2M"))

	v1 := e.Group("/api")

	authGroup := v1.Group("/auth")
	equipmentGroup := v1.Group("/equipment")

	authHttp.MapAuthRoutes(authGroup, authHandlers, mw)
	equipHttp.MapEquipmentRoutes(equipmentGroup, equipmentHandlers, mw)

	return nil
}
