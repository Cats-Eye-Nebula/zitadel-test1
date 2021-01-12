package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

func (r *CommandSide) RemoveHumanU2F(ctx context.Context, userID, webAuthNID, resourceOwner string) error {
	event := user.NewHumanU2FRemovedEvent(ctx, webAuthNID)
	return r.removeHumanWebAuthN(ctx, userID, webAuthNID, resourceOwner, event)
}

func (r *CommandSide) RemoveHumanPasswordless(ctx context.Context, userID, webAuthNID, resourceOwner string) error {
	event := user.NewHumanPasswordlessRemovedEvent(ctx, webAuthNID)
	return r.removeHumanWebAuthN(ctx, userID, webAuthNID, resourceOwner, event)
}

func (r *CommandSide) removeHumanWebAuthN(ctx context.Context, userID, webAuthNID, resourceOwner string, event eventstore.EventPusher) error {
	if userID == "" || webAuthNID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-6M9de", "Errors.IDMissing")
	}

	existingWebAuthN, err := r.webauthNWriteModelByID(ctx, userID, webAuthNID, resourceOwner)
	if err != nil {
		return err
	}
	if existingWebAuthN.State == domain.WebAuthNStateUnspecified || existingWebAuthN.State == domain.WebAuthNStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-5M0ds", "Errors.User.ExternalIDP.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingWebAuthN.WriteModel)
	userAgg.PushEvents(event)

	return r.eventstore.PushAggregate(ctx, existingWebAuthN, userAgg)
}

func (r *CommandSide) webauthNWriteModelByID(ctx context.Context, userID, webAuthNID, resourceOwner string) (writeModel *HumanWebAuthNWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanWebAuthNWriteModel(userID, webAuthNID, resourceOwner)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
