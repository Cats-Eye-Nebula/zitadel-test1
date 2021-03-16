package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/user"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) ChangeHumanPhone(ctx context.Context, phone *domain.Phone) (*domain.Phone, error) {
	if !phone.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-6M0ds", "Errors.Phone.Invalid")
	}

	existingPhone, err := c.phoneWriteModelByID(ctx, phone.AggregateID, phone.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingPhone.UserState.Exists() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3M0fs", "Errors.User.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingPhone.WriteModel)
	changedEvent, hasChanged := existingPhone.NewChangedEvent(ctx, userAgg, phone.PhoneNumber)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-wF94r", "Errors.User.Phone.NotChanged")
	}

	events := []eventstore.EventPusher{changedEvent}
	if phone.IsPhoneVerified {
		events = append(events, user.NewHumanPhoneVerifiedEvent(ctx, userAgg))
	} else {
		phoneCode, err := domain.NewPhoneCode(c.phoneVerificationCode)
		if err != nil {
			return nil, err
		}
		events = append(events, user.NewHumanPhoneCodeAddedEvent(ctx, userAgg, phoneCode.Code, phoneCode.Expiry))
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPhone, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToPhone(existingPhone), nil
}

func (c *Commands) VerifyHumanPhone(ctx context.Context, userID, code, resourceowner string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-Km9ds", "Errors.User.UserIDMissing")
	}
	if code == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-wMe9f", "Errors.User.Code.Empty")
	}

	existingCode, err := c.phoneWriteModelByID(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}
	if !existingCode.UserState.Exists() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Rsj8c", "Errors.User.NotFound")
	}
	if !existingCode.State.Exists() || existingCode.Code == nil {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-Rsj8c", "Errors.User.Code.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingCode.WriteModel)
	err = crypto.VerifyCode(existingCode.CodeCreationDate, existingCode.CodeExpiry, existingCode.Code, code, c.phoneVerificationCode)
	if err == nil {
		pushedEvents, err := c.eventstore.PushEvents(ctx, user.NewHumanPhoneVerifiedEvent(ctx, userAgg))
		if err != nil {
			return nil, err
		}
		err = AppendAndReduce(existingCode, pushedEvents...)
		if err != nil {
			return nil, err
		}
		return writeModelToObjectDetails(&existingCode.WriteModel), nil
	}
	_, err = c.eventstore.PushEvents(ctx, user.NewHumanPhoneVerificationFailedEvent(ctx, userAgg))
	logging.LogWithFields("COMMAND-5M9ds", "userID", userAgg.ID).OnError(err).Error("NewHumanPhoneVerificationFailedEvent push failed")
	return nil, caos_errs.ThrowInvalidArgument(err, "COMMAND-sM0cs", "Errors.User.Code.Invalid")
}

func (c *Commands) CreateHumanPhoneVerificationCode(ctx context.Context, userID, resourceowner string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-4M0ds", "Errors.User.UserIDMissing")
	}

	existingPhone, err := c.phoneWriteModelByID(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}

	if !existingPhone.UserState.Exists() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M0fs", "Errors.User.NotFound")
	}
	if !existingPhone.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-2b7Hf", "Errors.User.Phone.NotFound")
	}
	if existingPhone.IsPhoneVerified {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M9sf", "Errors.User.Phone.AlreadyVerified")
	}

	phoneCode, err := domain.NewPhoneCode(c.phoneVerificationCode)
	if err != nil {
		return nil, err
	}

	userAgg := UserAggregateFromWriteModel(&existingPhone.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, user.NewHumanPhoneCodeAddedEvent(ctx, userAgg, phoneCode.Code, phoneCode.Expiry))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPhone, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPhone.WriteModel), nil
}

func (c *Commands) HumanPhoneVerificationCodeSent(ctx context.Context, orgID, userID string) (err error) {
	if userID == "" {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-3m9Fs", "Errors.User.UserIDMissing")
	}

	existingPhone, err := c.phoneWriteModelByID(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if !existingPhone.UserState.Exists() {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3M9fs", "Errors.User.NotFound")
	}
	if !existingPhone.State.Exists() {
		return caos_errs.ThrowNotFound(nil, "COMMAND-66n8J", "Errors.User.Phone.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingPhone.WriteModel)
	_, err = c.eventstore.PushEvents(ctx, user.NewHumanPhoneCodeSentEvent(ctx, userAgg))
	return err
}

func (c *Commands) RemoveHumanPhone(ctx context.Context, userID, resourceOwner string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-6M0ds", "Errors.User.UserIDMissing")
	}

	existingPhone, err := c.phoneWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !existingPhone.UserState.Exists() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3M9fs", "Errors.User.NotFound")
	}
	if !existingPhone.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-p6rsc", "Errors.User.Phone.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingPhone.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, user.NewHumanPhoneRemovedEvent(ctx, userAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPhone, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPhone.WriteModel), nil
}

func (c *Commands) phoneWriteModelByID(ctx context.Context, userID, resourceOwner string) (writeModel *HumanPhoneWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanPhoneWriteModel(userID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
