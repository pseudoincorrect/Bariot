package server

import (
	"context"
	"log"
	"net"

	auth "github.com/pseudoincorrect/bariot/pkg/auth/grpc"
	"google.golang.org/grpc"

	"github.com/pseudoincorrect/bariot/auth/service"
)

type server struct {
	auth.UnimplementedAuthServer
	AuthService service.Auth
}

func (s *server) GetAdminToken(ctx context.Context, in *auth.GetAdminTokenRequest) (*auth.GetAdminTokenResponse, error) {
	token, err := s.AuthService.GetAdminToken()
	if err != nil {
		return nil, err
	}
	return &auth.GetAdminTokenResponse{Jwt: token}, nil
}

func (s *server) GetUserToken(ctx context.Context, in *auth.GetUserTokenRequest) (*auth.GetUserTokenResponse, error) {
	token, err := s.AuthService.GetUserToken(in.UserId)
	if err != nil {
		return nil, err
	}
	return &auth.GetUserTokenResponse{Jwt: token}, nil
}

func (s *server) GetThingToken(ctx context.Context, in *auth.GetThingTokenRequest) (*auth.GetThingTokenResponse, error) {
	token, err := s.AuthService.GetThingToken(in.ThingId, in.UserId)
	if err != nil {
		return nil, err
	}
	return &auth.GetThingTokenResponse{Jwt: token}, nil
}

// func (s *server) ValidateToken(ctx context.Context, in *auth.ValidateTokenRequest) (*auth.ValidateTokenResponse, error) {
// 	valid, err := s.AuthService.ValidateToken(in.Jwt)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &auth.ValidateTokenResponse{Valid: valid}, nil
// }

func (s *server) GetClaimsUserToken(ctx context.Context, in *auth.GetClaimsUserTokenRequest) (*auth.GetClaimsUserTokenResponse, error) {
	claims, err := s.AuthService.GetClaimsUserToken(in.Jwt)
	if err != nil {
		return nil, err
	}
	return &auth.GetClaimsUserTokenResponse{
		Role:      claims.Role,
		Subject:   claims.Subject,
		IssuedAt:  claims.IssuedAt,
		ExpiresAt: claims.ExpiresAt,
		Issuer:    claims.Issuer,
	}, nil
}

func (s *server) GetClaimsThingToken(ctx context.Context, in *auth.GetClaimsThingTokenRequest) (*auth.GetClaimsThingTokenResponse, error) {
	claims, err := s.AuthService.GetClaimsThingToken(in.Jwt)
	if err != nil {
		return nil, err
	}
	return &auth.GetClaimsThingTokenResponse{
		UserId:    claims.UserId,
		Subject:   claims.Subject,
		IssuedAt:  claims.IssuedAt,
		ExpiresAt: claims.ExpiresAt,
		Issuer:    claims.Issuer,
	}, nil
}

type ServerConf struct {
	AuthService service.Auth
	Port        string
}

func Start(c ServerConf) error {
	addr := ":" + c.Port
	log.Println("Starting Auth GRPC on", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}
	s := grpc.NewServer()
	auth.RegisterAuthServer(s, &server{auth.UnimplementedAuthServer{}, c.AuthService})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}
	return nil
}
