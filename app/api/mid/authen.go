package mid

import (
	"context"
	"fmt"
	"github.com/ardanlabs/service/app/api/errs"
	"github.com/ardanlabs/service/business/api/auth"
	"github.com/google/uuid"
	"strings"
)

func Authorization(ctx context.Context, auth *auth.Auth, authorization string, handler Handler) error {
	var err error
	parts := strings.Split(authorization, " ")

	switch parts[0] {
	case "Bearer":
		ctx, err = processJWT(ctx, auth, authorization)
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
