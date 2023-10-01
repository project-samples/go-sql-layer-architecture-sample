package app

import (
	"context"
	"github.com/core-go/core"

	v "github.com/core-go/core/v10"
	"github.com/core-go/health"
	"github.com/core-go/log/zap"
	q "github.com/core-go/sql"

	. "go-service/internal/handler"
	. "go-service/internal/repository"
	. "go-service/internal/service"
)

type ApplicationContext struct {
	Health *health.Handler
	User   UserTransport
}

func NewApp(ctx context.Context, cfg Config) (*ApplicationContext, error) {
	db, err := q.OpenByConfig(cfg.Sql)
	if err != nil {
		return nil, err
	}
	logError := log.LogError
	action := core.InitAction(cfg.Action)
	validator, err := v.NewValidator()
	if err != nil {
		return nil, err
	}

	userRepository, err := NewUserAdapter(db, BuildQuery)
	if err != nil {
		return nil, err
	}
	userService := NewUserService(userRepository)
	userHandler := NewUserHandler(userService, logError, validator.Validate, &action)

	sqlChecker := q.NewHealthChecker(db)
	healthHandler := health.NewHandler(sqlChecker)

	return &ApplicationContext{
		Health: healthHandler,
		User:   userHandler,
	}, nil
}
