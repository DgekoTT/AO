package grpc

import (
	"context"
	"errors"
	authv1 "github.com/DgekoTT/protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"yourTeamAuth/internal/services/auth"
	"yourTeamAuth/storage"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
	) (accessToken, refreshToken, userID string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID string, err error)
	IsAdmin(
		ctx context.Context,
		userID string) (isAdmin bool, err error)
	Logout(
		ctx context.Context,
		accessToken, refreshToken string,
	) (success bool, message string, err error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth Auth
}

// регистрируем обработчик
func Register(gRPC *grpc.Server, auth Auth) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	accessToken, refreshToken, userID, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())

	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "unable login")
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &authv1.LoginResponse{
		UserId:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *serverAPI) LogOut(ctx context.Context, req *authv1.LogOutRequest) (*authv1.LogOutResponse, error) {
	success, messageS, err := s.auth.Logout(ctx, req.GetAccessToken(), req.GetRefreshToken())
	if err != nil {
		// TODO:
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &authv1.LogOutResponse{
		Success: success,
		Message: messageS,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *authv1.IsAdminRequest) (*authv1.IsAdminResponse, error) {
	res, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		return &authv1.IsAdminResponse{
			IsAdmin: false,
		}, status.Errorf(codes.Internal, "internal error")
	}

	return &authv1.IsAdminResponse{
		IsAdmin: res,
	}, nil
}
