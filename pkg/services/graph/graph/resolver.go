package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
import (
	"github.com/rodrwan/bank/pkg/pb/session"
)

type Resolver struct {
	QueryService IQueryService

	CommandHandler ICommandHandler

	SessionClient session.SessionServiceClient
}
