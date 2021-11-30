package repository

import (
	"context"

	"golang.org/x/text/language"

	usr_model "github.com/caos/zitadel/internal/user/model"

	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type IAMRepository interface {
	Languages(ctx context.Context) ([]language.Tag, error)

	SearchIAMMembers(ctx context.Context, request *iam_model.IAMMemberSearchRequest) (*iam_model.IAMMemberSearchResponse, error)

	GetIAMMemberRoles() []string

	SearchIDPConfigs(ctx context.Context, request *iam_model.IDPConfigSearchRequest) (*iam_model.IDPConfigSearchResponse, error)

	SearchDefaultIDPProviders(ctx context.Context, request *iam_model.IDPProviderSearchRequest) (*iam_model.IDPProviderSearchResponse, error)

	IDPProvidersByIDPConfigID(ctx context.Context, idpConfigID string) ([]*iam_model.IDPProviderView, error)
	ExternalIDPsByIDPConfigID(ctx context.Context, idpConfigID string) ([]*usr_model.ExternalIDPView, error)
}
