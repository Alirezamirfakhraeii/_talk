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
	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/realtime"
)

type App struct {
	server      *httpserver.Server
	database    *pgxpool.Pool
	realtimeHub *realtime.Hub
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

	mux.Handle(
		"GET /uploads/",
		http.StripPrefix(
			"/uploads/",
			http.FileServer(
				http.Dir("./uploads"),
			),
		),
	)

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

	realtimeHub := realtime.NewHub()

	realtimeHandler := realtime.NewHandler(realtimeHub)

	realtime.RegisterRoutes(
		mux,
		realtimeHandler,
		authMiddleware,
	)

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

	messageHandler := message.NewHandler(
		messageService,
		realtimeHub,
	)

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
		server:      server,
		database:    databasePool,
		realtimeHub: realtimeHub,
	}, nil
}

func (app *App) Run(ctx context.Context) error {
	go app.realtimeHub.Run(ctx)
	return app.server.Run(ctx)
}

func (app *App) Close() {
	if app.database != nil {
		app.database.Close()
	}
}
