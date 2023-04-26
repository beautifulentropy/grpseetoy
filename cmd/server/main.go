package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"

	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello " + in.GetName()}, nil
}

func main() {
	grpclog.SetLoggerV2(grpclog.NewLoggerV2WithVerbosity(os.Stdout, os.Stdout, os.Stdout, 99))

	listener, err := net.Listen("tcp", ":1337")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(
		// grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		// 	MinTime:             1 * time.Second,
		// 	PermitWithoutStream: true,
		// }),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: 30 * time.Second,
			// MaxConnectionAgeGrace: 5 * time.Second,
		}),
	)

	pb.RegisterGreeterServer(grpcServer, &server{})
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
