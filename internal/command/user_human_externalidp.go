package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/user"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) BulkAddedHumanExternalIDP(ctx context.Context, userID, resourceOwner string, externalIDPs []*domain.ExternalIDP) (err error) {
	if userID == "" {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-03j8f", "Errors.IDMissing")
	}
	if len(externalIDPs) == 0 {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-Ek9s", "Errors.User.ExternalIDP.MinimumExternalIDPNeeded")
	}

	events := make([]eventstore.EventPusher, len(externalIDPs))
	for i, externalIDP := range externalIDPs {
		externalIDPWriteModel := NewHumanExternalIDPWriteModel(userID, externalIDP.IDPConfigID, externalIDP.ExternalUserID, resourceOwner)
		userAgg := UserAggregateFromWriteModel(&externalIDPWriteModel.WriteModel)

		events[i], err = c.addHumanExternalIDP(ctx, userAgg, externalIDP)
		if err != nil {
			return err
		}
	}

	_, err = c.eventstore.PushEvents(ctx, events...)
	return err
}

func (c *Commands) addHumanExternalIDP(ctx context.Context, humanAgg *eventstore.Aggregate, externalIDP *domain.ExternalIDP) (eventstore.EventPusher, error) {
	if externalIDP.AggregateID != "" && humanAgg.ID != externalIDP.AggregateID {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-33M0g", "Errors.IDMissing")
	}
	if !externalIDP.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-6m9Kd", "Errors.User.ExternalIDP.Invalid")
	}
	_, err := c.getOrgIDPConfigByID(ctx, externalIDP.IDPConfigID, humanAgg.ResourceOwner)
	if caos_errs.IsNotFound(err) {
		_, err = c.getIAMIDPConfigByID(ctx, externalIDP.IDPConfigID)
	}
	if err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(err, "COMMAND-39nfs", "Errors.IDPConfig.NotExisting")
	}
	return user.NewHumanExternalIDPAddedEvent(ctx, humanAgg, externalIDP.IDPConfigID, externalIDP.DisplayName, externalIDP.ExternalUserID), nil
}

func (c *Commands) RemoveHumanExternalIDP(ctx context.Context, externalIDP *domain.ExternalIDP) (*domain.ObjectDetails, error) {
	event, externalIDPWriteModel, err := c.removeHumanExternalIDP(ctx, externalIDP, false)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(externalIDPWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&externalIDPWriteModel.WriteModel), nil
}

func (c *Commands) removeHumanExternalIDP(ctx context.Context, externalIDP *domain.ExternalIDP, cascade bool) (eventstore.EventPusher, *HumanExternalIDPWriteModel, error) {
	if !externalIDP.IsValid() || externalIDP.AggregateID == "" {
		return nil, nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-3M9ds", "Errors.IDMissing")
	}

	existingExternalIDP, err := c.externalIDPWriteModelByID(ctx, externalIDP.AggregateID, externalIDP.IDPConfigID, externalIDP.ExternalUserID, externalIDP.ResourceOwner)
	if err != nil {
		return nil, nil, err
	}
	if existingExternalIDP.State == domain.ExternalIDPStateUnspecified || existingExternalIDP.State == domain.ExternalIDPStateRemoved {
		return nil, nil, caos_errs.ThrowNotFound(nil, "COMMAND-1M9xR", "Errors.User.ExternalIDP.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingExternalIDP.WriteModel)
	if cascade {
		return user.NewHumanExternalIDPCascadeRemovedEvent(ctx, userAgg, externalIDP.IDPConfigID, externalIDP.ExternalUserID), existingExternalIDP, nil
	}
	return user.NewHumanExternalIDPRemovedEvent(ctx, userAgg, externalIDP.IDPConfigID, externalIDP.ExternalUserID), existingExternalIDP, nil
}

func (c *Commands) HumanExternalLoginChecked(ctx context.Context, orgID, userID string, authRequest *domain.AuthRequest) (err error) {
	if userID == "" {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-5n8sM", "Errors.IDMissing")
	}

	existingHuman, err := c.getHumanWriteModelByID(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if existingHuman.UserState == domain.UserStateUnspecified || existingHuman.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-dn88J", "Errors.User.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingHuman.WriteModel)
	_, err = c.eventstore.PushEvents(ctx, user.NewHumanExternalIDPCheckSucceededEvent(ctx, userAgg, authRequestDomainToAuthRequestInfo(authRequest)))
	return err
}

func (c *Commands) externalIDPWriteModelByID(ctx context.Context, userID, idpConfigID, externalUserID, resourceOwner string) (writeModel *HumanExternalIDPWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanExternalIDPWriteModel(userID, idpConfigID, externalUserID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
