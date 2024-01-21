package grpc

import (
	"google.golang.org/grpc"
	authv1 "yourTeamAuth/pkg/proto/gen/go/auth"
)

type Auth interface {
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
}

// регистрируем обработчик
func Register(gRPC *grpc.Server) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{})
}
