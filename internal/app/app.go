package app

import (
	"yourTeamAuth/internal/config"
	"yourTeamAuth/internal/storage/postgres"

	"log/slog"
	"time"
	grpcapp "yourTeamAuth/internal/app/grpc"
	"yourTeamAuth/internal/services/auth"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	confDB *config.DatabaseConfig,
	JWTRefreshTTL time.Duration,
	JWTAccessTTL time.Duration,
	JWTVerificationTTL time.Duration,

) *App {
	db, err := postgres.NewPostgres(confDB)
	if err != nil {
		panic(err)
	}
	authService := auth.NewAuth(log, db, db, JWTRefreshTTL, JWTAccessTTL, JWTVerificationTTL)
	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{GRPCSrv: grpcApp}
}
