package repository

import (
	"context"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
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

	GetDefaultLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error)
	GetDefaultPreviewLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error)

	GetDefaultMailTemplate(ctx context.Context) (*iam_model.MailTemplateView, error)

	GetDefaultMessageText(ctx context.Context, textType, language string) (*domain.CustomMessageText, error)
	GetCustomMessageText(ctx context.Context, textType string, language string) (*domain.CustomMessageText, error)
	GetDefaultLoginTexts(ctx context.Context, lang string) (*domain.CustomLoginText, error)
	GetCustomLoginTexts(ctx context.Context, lang string) (*domain.CustomLoginText, error)

	GetDefaultPrivacyPolicy(ctx context.Context) (*iam_model.PrivacyPolicyView, error)

	GetDefaultOrgIAMPolicy(ctx context.Context) (*iam_model.OrgIAMPolicyView, error)
}
