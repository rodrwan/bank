package graph

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/opentracing/opentracing-go"
	"github.com/rodrwan/bank/pkg/auth"
	accounts "github.com/rodrwan/bank/pkg/pb/accounts"
	users "github.com/rodrwan/bank/pkg/pb/users"
	"github.com/rodrwan/bank/pkg/services/graph/graph"
	"github.com/rodrwan/bank/pkg/services/graph/graph/generated"
	tracerPkg "github.com/rodrwan/bank/pkg/tracer"
)

// ServerConfig ...
type ServerConfig struct {
	Port string

	AccountsReadClient  accounts.AccountReadServiceClient
	AccountsWriteClient accounts.AccountWriteServiceClient

	UsersReadClient  users.UsersReadServiceClient
	UsersWriteClient users.UsersWriteServiceClient
}

const defaultPort = "8080"

// NewServer initialises an instance of the graph server.
func NewServer(sc ServerConfig, tracer opentracing.Tracer) {
	if sc.Port == "" {
		sc.Port = defaultPort
	}

	gc := generated.Config{
		Resolvers: &graph.Resolver{
			QueryService: &graph.QueryService{
				AccountsReadClient: sc.AccountsReadClient,
				UsersReadClient:    sc.UsersReadClient,
			},
			CommandHandler: &graph.CommandHandler{
				AccountsWriteClient: sc.AccountsWriteClient,
				UsersWriteClient:    sc.UsersWriteClient,
			},
		},
	}

	gc.Directives.IsAuthenticated = auth.MiddlewareDirective()
	gc.Directives.Trace = tracerPkg.MiddlewareDirective(tracer)
	g := generated.NewExecutableSchema(gc)
	srv := handler.NewDefaultServer(g)

	withTracer := tracerPkg.Middleware(playground.Handler("GraphQL playground", "/query"), tracer)
	withAuth := auth.Middleware(withTracer, "localhost", true)

	http.Handle("/", withAuth)
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", sc.Port)
	log.Fatal(http.ListenAndServe(":"+sc.Port, nil))
}
