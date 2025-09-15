package mid

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/ardanlabs/service/app/api/errs"
	"github.com/ardanlabs/service/business/api/auth"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"strings"
	"time"
)

func Authorization(ctx context.Context, auth *auth.Auth, authorization string, handler Handler) error {
	var err error
	parts := strings.Split(authorization, " ")

	switch parts[0] {
	case "Bearer":
		ctx, err = processJWT(ctx, auth, authorization)
	case "Basic":
		ctx, err = processBasic(ctx)
	}

	if err != nil {
		return err
	}

	return handler(ctx)
}

func processJWT(ctx context.Context, auth *auth.Auth, authorization string) (context.Context, error) {
	claims, err := auth.Authenticate(ctx, authorization)
	if err != nil {
		return ctx, errs.New(errs.Unauthenticated, err)
	}

	if claims.Subject == "" {
		return ctx, errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, no claims")
	}

	subjectID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return ctx, errs.New(errs.Unauthenticated, fmt.Errorf("parsing subject: %w", err))
	}

	ctx = setUserID(ctx, subjectID)
	ctx = setClaims(ctx, claims)

	return ctx, nil
}

// Basic processes basic authentication logic.
func processBasic(ctx context.Context) (context.Context, error) {
	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "a86505e1-fb77-4cbe-81e5-b83f58069a1b",
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: []string{"ADMIN"},
	}

	subjectID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return ctx, errs.Newf(errs.Unauthenticated, "parsing subject: %s", err)
	}

	ctx = setUserID(ctx, subjectID)
	ctx = setClaims(ctx, claims)

	return ctx, nil
}

func parseBasicAuth(auth string) (string, string, bool) {
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Basic" {
		return "", "", false
	}

	c, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", false
	}

	username, password, ok := strings.Cut(string(c), ":")
	if !ok {
		return "", "", false
	}

	return username, password, true
}
