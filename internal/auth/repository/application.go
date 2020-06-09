package repository

import (
	"context"

	"github.com/caos/zitadel/internal/project/model"
)

type ApplicationRepository interface {
	ApplicationByClientID(ctx context.Context, clientID string) (*model.ApplicationView, error)
	AuthorizeOIDCApplication(ctx context.Context, clientID, secret string) error
}
