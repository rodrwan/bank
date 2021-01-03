package main

import (
	"log"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	accounts "github.com/rodrwan/bank/pkg/pb/accounts"
	"github.com/rodrwan/bank/pkg/pb/session"
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

	accountsReadConn := connect(os.Getenv("ACCOUNTS_READ_URL"), tracer)
	defer accountsReadConn.Close()
	accountsReadClient := accounts.NewAccountReadServiceClient(accountsReadConn)

	accountsWriteConn := connect(os.Getenv("ACCOUNTS_WRITE_URL"), tracer)
	defer accountsWriteConn.Close()
	accountsWriteClient := accounts.NewAccountWriteServiceClient(accountsWriteConn)

	usersReadConn := connect(os.Getenv("USERS_READ_URL"), tracer)
	defer usersReadConn.Close()
	usersReadClient := users.NewUsersReadServiceClient(usersReadConn)

	usersWriteConn := connect(os.Getenv("USERS_WRITE_URL"), tracer)
	defer usersWriteConn.Close()
	usersWriteClient := users.NewUsersWriteServiceClient(usersWriteConn)

	sessionConn := connect(os.Getenv("SESSION_URL"), tracer)
	defer sessionConn.Close()
	sessionClient := session.NewSessionServiceClient(sessionConn)

	graph.NewServer(graph.ServerConfig{
		Port:                port,
		AccountsReadClient:  accountsReadClient,
		AccountsWriteClient: accountsWriteClient,
		UsersReadClient:     usersReadClient,
		UsersWriteClient:    usersWriteClient,
		SessionClient:       sessionClient,
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
