package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/config"
	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/conversation"
	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/health"
	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/platform/database"
	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/platform/httpserver"
	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/user"

	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/auth"
	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/message"
)

type App struct {
	server   *httpserver.Server
	database *pgxpool.Pool
}

func New(
	ctx context.Context,
	cfg config.Config,
) (*App, error) {
	databasePool, err := database.NewPostgres(
		ctx,
		cfg.DatabaseURL,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"initialize database: %w",
			err,
		)
	}

	mux := http.NewServeMux()

	healthHandler := health.NewHandler(
		databasePool,
	)

	health.RegisterRoutes(
		mux,
		healthHandler,
	)

	userRepository := user.NewRepository(
		databasePool,
	)

	authRepository := auth.NewRepository(
		databasePool,
	)

	authMiddleware := auth.NewMiddleware(authRepository)

	userService := user.NewService(
		userRepository,
		authRepository,
	)

	userHandler := user.NewHandler(userService)

	conversationRepository := conversation.NewRepository(databasePool)

	conversationService := conversation.NewService(
		conversationRepository,
		userRepository,
	)

	conversationHandler := conversation.NewHandler(
		conversationService,
	)

	conversation.RegisterRoutes(
		mux,
		conversationHandler,
		authMiddleware,
	)

	messageRepository := message.NewRepository(databasePool)

	messageService := message.NewService(
		messageRepository,
		conversationRepository,
	)

	messageHandler := message.NewHandler(messageService)

	message.RegisterRoutes(
		mux,
		messageHandler,
		authMiddleware,
	)

	user.RegisterRoutes(
		mux,
		userHandler,
		authMiddleware,
	)

	server := httpserver.New(
		cfg.HTTPAddress,
		mux,
	)

	return &App{
		server:   server,
		database: databasePool,
	}, nil
}

func (app *App) Run(ctx context.Context) error {
	return app.server.Run(ctx)
}

func (app *App) Close() {
	if app.database != nil {
		app.database.Close()
	}
}
