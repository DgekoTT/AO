package auth

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
	"yourTeamAuth/internal/domain/models"
	"yourTeamAuth/storage"
)

type Auth struct {
	log                *slog.Logger
	usrSaver           UserSaver
	usrProvider        UserProvider
	JWTRefreshTTL      time.Duration
	JWTAccessTTL       time.Duration
	JWTVerificationTTL time.Duration
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		HashedPassword []byte,
	) (uid string, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.Users, error)
	IsAdmin(ctx context.Context, userID string) (bool, error)
}

func NewAuth(
	log *slog.Logger,
	usrSaver UserSaver,
	usrProvider UserProvider,
	JWTRefreshTTL time.Duration,
	JWTAccessTTL time.Duration,
	JWTVerificationTTL time.Duration,
) *Auth {
	return &Auth{
		log:                log,
		usrSaver:           usrSaver,
		usrProvider:        usrProvider,
		JWTRefreshTTL:      JWTRefreshTTL,
		JWTAccessTTL:       JWTAccessTTL,
		JWTVerificationTTL: JWTVerificationTTL,
	}
}

func (a *Auth) RegisterNewUser(ctx context.Context, email, password string) (string, error) {
	const op = "auth.registerNewUser"

	log := a.log.With(slog.String("op", op))

	log.Info("Register new user")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Error("failed to hashed password", err)
		return "0", fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, hashedPassword)

	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("modelUser already exists", err)
			return "0", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		log.Error("failed to save user in db", err)

		return "0", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil

}

func (a *Auth) Login(ctx context.Context, email, password string) (string, string, string, error) {
	const op = "auth.Login"

	log := a.log.With(slog.String("op", op))

	log.Info("User is trying to login")

	user, err := a.usrProvider.User(ctx, email)

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", err)

			return "0", "0", "0", err
		}
		a.log.Error("failed to get modelUser", err)
		return "0", "0", "0", fmt.Errorf("%s : %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password)); err != nil {
		a.log.Info("invalid credentials", err)
		return "", "", "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("User logged in successfully")

	accessToken, err := jwt.NewAccessToken(user, a.JWTAccessTTL)
	if err != nil {
		log.Warn("failed to create access token", err)
		return "0", "0", "0", fmt.Errorf("%s : %w", op, ErrInvalidCredentials)
	}
	refreshToken, err := jwt.NewRefreshToken(user, a.JWTRefreshTTL)
	if err != nil {
		log.Warn("failed to create refresh token", err)
		return "0", "0", "0", fmt.Errorf("%s : %w", op, ErrInvalidCredentials)
	}

	return accessToken, refreshToken, user.ID.String(), nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID string) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(slog.String("op", op))

	log.Info("start check if user ia admin")

	res, err := a.usrProvider.IsAdmin(ctx, userID)

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("User not found", err)

			return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
	}

	return res, nil
}
