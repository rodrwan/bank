package main

import (
	"log"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	accounts "github.com/rodrwan/bank/pkg/pb/accounts"
	users "github.com/rodrwan/bank/pkg/pb/users"
	"github.com/rodrwan/bank/pkg/services/graph"
	"github.com/rodrwan/bank/pkg/tracer"
	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("PORT")

	// initialize tracer
	tracer, closer, err := tracer.NewTracer()
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	accountsReadConn := connect("accounts-read:8010", tracer)
	defer accountsReadConn.Close()
	accountsReadClient := accounts.NewAccountReadServiceClient(accountsReadConn)

	accountsWriteConn := connect("accounts-write:8011", tracer)
	defer accountsWriteConn.Close()
	accountsWriteClient := accounts.NewAccountWriteServiceClient(accountsWriteConn)

	usersReadConn := connect("users-read:8020", tracer)
	defer usersReadConn.Close()
	usersReadClient := users.NewUsersReadServiceClient(usersReadConn)

	usersWriteConn := connect("users-write:8021", tracer)
	defer usersWriteConn.Close()
	usersWriteClient := users.NewUsersWriteServiceClient(usersWriteConn)

	graph.NewServer(graph.ServerConfig{
		Port:                port,
		AccountsReadClient:  accountsReadClient,
		AccountsWriteClient: accountsWriteClient,
		UsersReadClient:     usersReadClient,
		UsersWriteClient:    usersWriteClient,
	}, tracer)
}

func connect(addr string, tracer opentracing.Tracer) *grpc.ClientConn {
	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(tracer)),
		)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(tracer)),
		)))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	return conn
}
