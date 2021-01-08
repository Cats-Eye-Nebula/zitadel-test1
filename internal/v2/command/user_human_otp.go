package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

func (r *CommandSide) AddHumanOTP(ctx context.Context, userID string) (*domain.OTP, error) {
	if userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-5M0sd", "Errors.User.UserIDMissing")
	}
	human, err := r.getUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	org, err := r.getOrg(ctx, human.ResourceOwner)
	if err != nil {
		return nil, err
	}
	orgPolicy, err := r.getOrgIAMPolicy(ctx, org.AggregateID)
	if err != nil {
		return nil, err
	}
	otpWriteModel, err := r.otpWriteModelByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if otpWriteModel.MFAState == domain.MFAStateReady {
		return nil, caos_errs.ThrowAlreadyExists(nil, "COMMAND-do9se", "Errors.User.MFA.OTP.AlreadyReady")
	}
	userAgg := UserAggregateFromWriteModel(&otpWriteModel.WriteModel)
	accountName := domain.GenerateLoginName(human.UserName, org.PrimaryDomain, orgPolicy.UserLoginMustBeDomain)
	if accountName == "" {
		accountName = human.EmailAddress
	}
	key, secret, err := domain.NewOTPKey(r.multifactors.OTP.Issuer, accountName, r.multifactors.OTP.CryptoMFA)
	if err != nil {
		return nil, err
	}
	userAgg.PushEvents(
		user.NewHumanOTPAddedEvent(ctx, secret),
	)

	err = r.eventstore.PushAggregate(ctx, otpWriteModel, userAgg)
	if err != nil {
		return nil, err
	}
	return &domain.OTP{
		ObjectRoot: models.ObjectRoot{
			AggregateID: human.AggregateID,
		},
		SecretString: key.Secret(),
		Url:          key.URL(),
	}, nil
}

func (r *CommandSide) RemoveHumanOTP(ctx context.Context, userID string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-5M0sd", "Errors.User.UserIDMissing")
	}

	existingOTP, err := r.otpWriteModelByID(ctx, userID)
	if err != nil {
		return err
	}
	if existingOTP.OTPState == domain.OTPStateUnspecified || existingOTP.OTPState == domain.OTPStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-5M0ds", "Errors.User.OTP.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingOTP.WriteModel)
	userAgg.PushEvents(
		user.NewHumanOTPRemovedEvent(ctx),
	)

	return r.eventstore.PushAggregate(ctx, existingOTP, userAgg)
}

func (r *CommandSide) otpWriteModelByID(ctx context.Context, userID string) (writeModel *HumanOTPWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanOTPWriteModel(userID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
