package session

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
	pb "github.com/rodrwan/bank/pkg/pb/session"
	"github.com/rodrwan/bank/pkg/pb/users"
)

// ServerConfig ...
type ServerConfig struct {
	AccessSecret  []byte
	RefreshSecret []byte

	Store *redis.Client
}

// NewServer initialises an instance of the session server.
func NewServer(sc ServerConfig) pb.SessionServiceServer {
	fmt.Println("Session service")
	return &Server{
		sessionService: &Service{
			accessSecret:  sc.AccessSecret,
			refreshSecret: sc.RefreshSecret,
			store:         sc.Store,
		},
	}
}

// Server ...
type Server struct {
	userReadClient users.UsersReadServiceClient
	sessionService *Service

	pb.UnimplementedSessionServiceServer
}

// GetSessionData ...
func (srv *Server) GetSessionData(ctx context.Context, in *pb.GetSessionDataRequest) (*pb.GetSessionDataResponse, error) {
	token := in.GetToken()
	data, err := srv.sessionService.GetAuthData(ctx, token)
	if err != nil {
		return nil, err
	}

	return &pb.GetSessionDataResponse{
		Data: &pb.Session{
			Data:        data.Data,
			AccessUuid:  data.AccessUUID,
			RefreshUuid: data.RefreshUUID,
			ReferenceId: data.ReferenceID,
		},
	}, nil
}

// CreateSession ...
func (srv *Server) CreateSession(ctx context.Context, in *pb.CreateSessionRequest) (*pb.CreateSessionResponse, error) {
	data := in.GetData()
	referenceID := in.GetReferenceId()

	auth, err := srv.sessionService.CreateAuthData(ctx, referenceID, []byte(data))
	if err != nil {
		return nil, err
	}

	return &pb.CreateSessionResponse{
		Data: &pb.Auth{
			AccessToken:  auth.AccessToken,
			RefreshToken: auth.RefreshToken,
		},
	}, nil
}

// RefreshSession ...
func (srv *Server) RefreshSession(ctx context.Context, in *pb.RefreshSessionRequest) (*pb.RefreshSessionResponse, error) {
	auth := in.GetData()

	newAuth, err := srv.sessionService.RefreshAuthData(ctx, &Auth{
		AccessToken:  auth.GetAccessToken(),
		RefreshToken: auth.GetRefreshToken(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.RefreshSessionResponse{
		Data: &pb.Auth{
			AccessToken:  newAuth.AccessToken,
			RefreshToken: newAuth.RefreshToken,
		},
	}, nil
}

// DeleteSession ...
func (srv *Server) DeleteSession(ctx context.Context, in *pb.DeleteSessionRequest) (*pb.DeleteSessionResponse, error) {
	auth := in.GetData()

	if err := srv.sessionService.BlockAuthData(ctx, &Auth{
		AccessToken:  auth.GetAccessToken(),
		RefreshToken: auth.GetRefreshToken(),
	}); err != nil {
		return nil, err
	}

	return &pb.DeleteSessionResponse{}, nil
}
