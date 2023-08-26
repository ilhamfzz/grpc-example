package main

import (
	"context"
	pb "go-grpc-example/user"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedUserServer
}

func (s *server) UnaryGetUser(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	return &pb.UserResponse{
		Id:    "5025201275",
		Name:  "Fakhrii",
		Email: in.Email,
		Age:   22,
	}, nil
}

func isValidAppKey(appKey []string) bool {
	if len(appKey) < 1 {
		return false
	}
	return appKey[0] == "appkey123"
}

func unaryInterceptorImpl(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("incoming request to", info.FullMethod)
	mdKey, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	if isValidAppKey(mdKey["app-key"]) {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid app key")
	}

	md, err := handler(ctx, req)
	if err != nil {
		log.Println("error in unaryInterceptorImpl:", err)
	}
	return md, err
}

func main() {
	list, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptorImpl))
	pb.RegisterUserServer(s, &server{})
	log.Println("Server running on port :9090")
	if err := s.Serve(list); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
