package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/telemetry/tracing"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
)

func (c *Commands) AddDefaultIDPConfig(ctx context.Context, config *domain.IDPConfig) (*domain.IDPConfig, error) {
	if config.OIDCConfig == nil && config.JWTConfig == nil {
		return nil, caos_errs.ThrowInvalidArgument(nil, "IDP-s8nn3", "Errors.IDPConfig.Invalid")
	}
	idpConfigID, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	addedConfig := NewIAMIDPConfigWriteModel(idpConfigID)

	iamAgg := IAMAggregateFromWriteModel(&addedConfig.WriteModel)
	events := []eventstore.Command{
		iam_repo.NewIDPConfigAddedEvent(
			ctx,
			iamAgg,
			idpConfigID,
			config.Name,
			config.Type,
			config.StylingType,
			config.AutoRegister,
		),
	}
	if config.OIDCConfig != nil {
		clientSecret, err := crypto.Encrypt([]byte(config.OIDCConfig.ClientSecretString), c.idpConfigSecretCrypto)
		if err != nil {
			return nil, err
		}

		events = append(events, iam_repo.NewIDPOIDCConfigAddedEvent(
			ctx,
			iamAgg,
			config.OIDCConfig.ClientID,
			idpConfigID,
			config.OIDCConfig.Issuer,
			config.OIDCConfig.AuthorizationEndpoint,
			config.OIDCConfig.TokenEndpoint,
			clientSecret,
			config.OIDCConfig.IDPDisplayNameMapping,
			config.OIDCConfig.UsernameMapping,
			config.OIDCConfig.Scopes...,
		))
	} else if config.JWTConfig != nil {
		events = append(events, iam_repo.NewIDPJWTConfigAddedEvent(
			ctx,
			iamAgg,
			idpConfigID,
			config.JWTConfig.JWTEndpoint,
			config.JWTConfig.Issuer,
			config.JWTConfig.KeysEndpoint,
			config.JWTConfig.HeaderName,
		))
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedConfig, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToIDPConfig(&addedConfig.IDPConfigWriteModel), nil
}

func (c *Commands) ChangeDefaultIDPConfig(ctx context.Context, config *domain.IDPConfig) (*domain.IDPConfig, error) {
	if config.IDPConfigID == "" {
		return nil, errors.ThrowInvalidArgument(nil, "IAM-4m9gs", "Errors.IDMissing")
	}
	existingIDP, err := c.iamIDPConfigWriteModelByID(ctx, config.IDPConfigID)
	if err != nil {
		return nil, err
	}
	if existingIDP.State == domain.IDPConfigStateRemoved || existingIDP.State == domain.IDPConfigStateUnspecified {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-m0e3r", "Errors.IDPConfig.NotExisting")
	}

	iamAgg := IAMAggregateFromWriteModel(&existingIDP.WriteModel)
	changedEvent, hasChanged := existingIDP.NewChangedEvent(ctx, iamAgg, config.IDPConfigID, config.Name, config.StylingType, config.AutoRegister)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-4M9vs", "Errors.IAM.LabelPolicy.NotChanged")
	}
	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingIDP, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToIDPConfig(&existingIDP.IDPConfigWriteModel), nil
}

func (c *Commands) DeactivateDefaultIDPConfig(ctx context.Context, idpID string) (*domain.ObjectDetails, error) {
	existingIDP, err := c.iamIDPConfigWriteModelByID(ctx, idpID)
	if err != nil {
		return nil, err
	}
	if existingIDP.State != domain.IDPConfigStateActive {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-2n0fs", "Errors.IAM.IDPConfig.NotActive")
	}
	iamAgg := IAMAggregateFromWriteModel(&existingIDP.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, iam_repo.NewIDPConfigDeactivatedEvent(ctx, iamAgg, idpID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingIDP, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingIDP.IDPConfigWriteModel.WriteModel), nil
}

func (c *Commands) ReactivateDefaultIDPConfig(ctx context.Context, idpID string) (*domain.ObjectDetails, error) {
	existingIDP, err := c.iamIDPConfigWriteModelByID(ctx, idpID)
	if err != nil {
		return nil, err
	}
	if existingIDP.State != domain.IDPConfigStateInactive {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-5Mo0d", "Errors.IAM.IDPConfig.NotInactive")
	}
	iamAgg := IAMAggregateFromWriteModel(&existingIDP.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, iam_repo.NewIDPConfigReactivatedEvent(ctx, iamAgg, idpID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingIDP, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingIDP.IDPConfigWriteModel.WriteModel), nil
}

func (c *Commands) RemoveDefaultIDPConfig(ctx context.Context, idpID string, idpProviders []*domain.IDPProvider, externalIDPs ...*domain.UserIDPLink) (*domain.ObjectDetails, error) {
	existingIDP, err := c.iamIDPConfigWriteModelByID(ctx, idpID)
	if err != nil {
		return nil, err
	}
	if existingIDP.State == domain.IDPConfigStateRemoved || existingIDP.State == domain.IDPConfigStateUnspecified {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-4M0xy", "Errors.IDPConfig.NotExisting")
	}

	iamAgg := IAMAggregateFromWriteModel(&existingIDP.WriteModel)
	events := []eventstore.Command{
		iam_repo.NewIDPConfigRemovedEvent(ctx, iamAgg, idpID, existingIDP.Name),
	}

	for _, idpProvider := range idpProviders {
		if idpProvider.AggregateID == domain.IAMID {
			userEvents := c.removeIDPProviderFromDefaultLoginPolicy(ctx, iamAgg, idpProvider, true, externalIDPs...)
			events = append(events, userEvents...)
		}
		orgAgg := OrgAggregateFromWriteModel(&NewOrgIdentityProviderWriteModel(idpProvider.AggregateID, idpID).WriteModel)
		orgEvents := c.removeIDPProviderFromLoginPolicy(ctx, orgAgg, idpID, true)
		events = append(events, orgEvents...)
	}

	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingIDP, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingIDP.IDPConfigWriteModel.WriteModel), nil
}

func (c *Commands) getIAMIDPConfigByID(ctx context.Context, idpID string) (*domain.IDPConfig, error) {
	config, err := c.iamIDPConfigWriteModelByID(ctx, idpID)
	if err != nil {
		return nil, err
	}
	if !config.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-p0pFF", "Errors.IDPConfig.NotExisting")
	}
	return writeModelToIDPConfig(&config.IDPConfigWriteModel), nil
}

func (c *Commands) iamIDPConfigWriteModelByID(ctx context.Context, idpID string) (policy *IAMIDPConfigWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMIDPConfigWriteModel(idpID)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
