package server

import (
	"context"
	"database/sql"
	"equiptrack/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const (
	ctxTimeout = 5
)

type Server struct {
	echo   *echo.Echo
	cfg    *config.Config
	db     *sql.DB
	logger *logrus.Logger
}

func NewServer(cfg *config.Config, db *sql.DB, logger *logrus.Logger) *Server {
	return &Server{echo: echo.New(), cfg: cfg, db: db, logger: logger}
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:         s.cfg.Server.Port,
		ReadTimeout:  time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout: time.Second * s.cfg.Server.WriteTimeout,
	}

	go func() {
		s.logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
		if err := s.echo.StartServer(server); err != nil {
			s.logger.Fatalf("Error starting Server: %s", err)
		}
	}()

	// go func() {
	// 	s.logger.Infof("Starting Debug Server on PORT: %s", s.cfg.Server.PprofPort)
	// 	if err := http.ListenAndServe(s.cfg.Server.PprofPort, http.DefaultServeMux); err != nil {
	// 		s.logger.Errorf("Error PPROF ListenAndServe: %s", err)
	// 	}
	// }()

	if err := s.MapHandlers(s.echo); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	s.logger.Info("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}
