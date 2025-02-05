package middleware

import (
	"equiptrack/config"
	"equiptrack/internal/auth"

	"github.com/sirupsen/logrus"
)

type MiddlewareManager struct {
	authUC  auth.UseCase
	cfg     *config.Config
	origins []string
	logger  *logrus.Logger
}

// Middleware manager constructor
func NewMiddlewareManager(authUC auth.UseCase, cfg *config.Config, origins []string, logger *logrus.Logger) *MiddlewareManager {
	return &MiddlewareManager{authUC: authUC, cfg: cfg, origins: origins, logger: logger}
}
