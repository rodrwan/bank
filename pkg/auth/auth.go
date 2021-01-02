package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql"
)

const (
	accessTokenCookieName = "access_token"
	authTokenCookieName   = "auth_token"
	tokenTypePrefix       = "Bearer "
	tokenHeaderKey        = "authorization"
)

type contextKey struct {
	name string
}

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var accessTokenCtxKey = &contextKey{accessTokenCookieName}
var authTokenCtxKey = &contextKey{authTokenCookieName}
var userIDCtxKey = &contextKey{"user-id"}

// Middleware decodes the share session cookie and packs the session into context
func Middleware(next http.HandlerFunc, domain string, secure bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer, _ := parseAuthToken(r)

		// Check if bearer is not blank
		if bearer != "" { // && accessToken != nil {
			rCtx := r.Context()
			// put it in context.
			ctx := context.WithValue(rCtx, authTokenCtxKey, bearer)
			ctx = context.WithValue(ctx, userIDCtxKey, bearer)
			// overwrite request
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

// MiddlewareDirective ...
func MiddlewareDirective() func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		ctxUserID := ctx.Value(userIDCtxKey)
		fmt.Println("ctxUserID", ctxUserID)
		if ctxUserID != nil {
			return next(ctx)
		}

		return nil, errors.New("session expired")
	}
}

// GetAuthTokenFromContext finds the auth token from the context.
func GetAuthTokenFromContext(ctx context.Context) string {
	raw, _ := ctx.Value(authTokenCtxKey).(string)
	return raw
}

// GetAccessTokenFromContext finds the access token from the context.
func GetAccessTokenFromContext(ctx context.Context) string {
	raw, _ := ctx.Value(accessTokenCtxKey).(string)
	return raw
}

// GetUserIDFromContext finds the user id from the context.
func GetUserIDFromContext(ctx context.Context) string {
	raw, _ := ctx.Value(userIDCtxKey).(string)
	return raw
}

func parseAuthToken(r *http.Request) (string, error) {
	header := r.Header.Get(tokenHeaderKey)

	if !strings.HasPrefix(header, tokenTypePrefix) {
		return "", errors.New("auth: no token authorization header present")
	}

	return header[len(tokenTypePrefix):], nil
}
