// go:generate protoc -I ../../../proto --go_out=$GOPATH/src  --go-grpc_out=$GOPATH/src  ../../../proto/account_write.proto
package main

import (
	"log"
	"net"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	_ "github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	pb "github.com/rodrwan/bank/pkg/pb/accounts"
	"github.com/rodrwan/bank/pkg/services/accounts"
	"github.com/rodrwan/bank/pkg/tracer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	defaultPort = ":8011"
)

func main() {
	port := os.Getenv("PORT")
	if port != "" {
		port = defaultPort
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// initialize tracer
	tracer, closer, err := tracer.NewTracer()
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			// add opentracing stream interceptor to chain
			grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(tracer)),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			// add opentracing unary interceptor to chain
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer)),
		)),
	)

	config := accounts.AccountServiceConfig{
		RabbitMQURL: os.Getenv("RABBIT_MQ_URL"),
	}
	pb.RegisterAccountWriteServiceServer(s, accounts.NewWriteServer(config))
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
