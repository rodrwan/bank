package tracer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/opentracing/opentracing-go"
)

// Middleware decodes the share session cookie and packs the session into context
func Middleware(next http.HandlerFunc, tracer opentracing.Tracer) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("extract trace from context")
		var parentCtx opentracing.SpanContext
		parentSpan := opentracing.SpanFromContext(r.Context())
		if parentSpan != nil {
			parentCtx = parentSpan.Context()
		}
		// start a new Span to wrap HTTP request
		span := tracer.StartSpan(
			"graphql",
			opentracing.ChildOf(parentCtx),
		)

		defer span.Finish()
		next.ServeHTTP(w, r)
	})
}

// MiddlewareDirective ...
func MiddlewareDirective(tracer opentracing.Tracer) func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		fmt.Println("extract trace from context directive")
		var parentCtx opentracing.SpanContext
		parentSpan := opentracing.SpanFromContext(ctx)
		if parentSpan != nil {
			parentCtx = parentSpan.Context()
		}
		// start a new Span to wrap HTTP request
		span := tracer.StartSpan(
			"graphql",
			opentracing.ChildOf(parentCtx),
		)

		defer span.Finish()
		return next(ctx)
	}
}
