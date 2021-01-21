package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

func (r *CommandSide) getHuman(ctx context.Context, userID, resourceowner string) (*domain.Human, error) {
	writeModel, err := r.getHumanWriteModelByID(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}
	if writeModel.UserState == domain.UserStateUnspecified || writeModel.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-M9dsd", "Errors.User.NotFound")
	}
	return writeModelToHuman(writeModel), nil
}

func (r *CommandSide) AddHuman(ctx context.Context, orgID string, human *domain.Human) (*domain.Human, error) {
	userAgg, addedHuman, err := r.addHuman(ctx, orgID, human)
	if err != nil {
		return nil, err
	}
	err = r.eventstore.PushAggregate(ctx, addedHuman, userAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToHuman(addedHuman), nil
}

func (r *CommandSide) addHuman(ctx context.Context, orgID string, human *domain.Human) (*user.Aggregate, *HumanWriteModel, error) {
	if !human.IsValid() {
		return nil, nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-4M90d", "Errors.User.Invalid")
	}
	return r.createHuman(ctx, orgID, human, nil, nil, false)
}

func (r *CommandSide) RegisterHuman(ctx context.Context, orgID string, human *domain.Human, externalIDP *domain.ExternalIDP, orgMember *domain.Member) (*domain.Human, error) {
	userAgg, addedHuman, err := r.registerHuman(ctx, orgID, human, externalIDP, orgMember)
	if err != nil {
		return nil, err
	}
	err = r.eventstore.PushAggregate(ctx, addedHuman, userAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToHuman(addedHuman), nil
}

func (r *CommandSide) registerHuman(ctx context.Context, orgID string, human *domain.Human, externalIDP *domain.ExternalIDP, orgMember *domain.Member) (*user.Aggregate, *HumanWriteModel, error) {
	if !human.IsValid() || externalIDP == nil && (human.Password == nil || human.SecretString == "") {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-9dk45", "Errors.User.Invalid")
	}
	return r.createHuman(ctx, orgID, human, externalIDP, orgMember, true)
}

func (r *CommandSide) createHuman(ctx context.Context, orgID string, human *domain.Human, externalIDP *domain.ExternalIDP, orgMember *domain.Member, selfregister bool) (*user.Aggregate, *HumanWriteModel, error) {
	userID, err := r.idGenerator.Next()
	if err != nil {
		return nil, nil, err
	}
	human.AggregateID = userID
	orgMember.UserID = userID
	orgIAMPolicy, err := r.getOrgIAMPolicy(ctx, orgID)
	if err != nil {
		return nil, nil, err
	}
	pwPolicy, err := r.getOrgPasswordComplexityPolicy(ctx, orgID)
	if err != nil {
		return nil, nil, err
	}

	addedHuman := NewHumanWriteModel(human.AggregateID, orgID)
	//TODO: Check Unique Username or unique external idp
	if err := human.CheckOrgIAMPolicy(human.Username, orgIAMPolicy); err != nil {
		return nil, nil, err
	}
	human.SetNamesAsDisplayname()
	if err := human.HashPasswordIfExisting(pwPolicy, r.userPasswordAlg, true); err != nil {
		return nil, nil, err
	}

	userAgg := UserAggregateFromWriteModel(&addedHuman.WriteModel)
	var createEvent eventstore.EventPusher
	if selfregister {
		createEvent = createRegisterHumanEvent(ctx, human.Username, human)
	} else {
		createEvent = createAddHumanEvent(ctx, human.Username, human)
	}
	userAgg.PushEvents(createEvent)

	if externalIDP != nil {
		if !externalIDP.IsValid() {
			return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4Dj9s", "Errors.User.ExternalIDP.Invalid")
		}
		//TODO: check if idpconfig exists
		userAgg.PushEvents(user.NewHumanExternalIDPAddedEvent(ctx, externalIDP.IDPConfigID, externalIDP.DisplayName))
	}
	if human.IsInitialState() {
		initCode, err := domain.NewInitUserCode(r.initializeUserCode)
		if err != nil {
			return nil, nil, err
		}
		userAgg.PushEvents(user.NewHumanInitialCodeAddedEvent(ctx, initCode.Code, initCode.Expiry))
	}
	if human.Email != nil && human.EmailAddress != "" && human.IsEmailVerified {
		userAgg.PushEvents(user.NewHumanEmailVerifiedEvent(ctx))
	}
	if human.Phone != nil && human.PhoneNumber != "" && !human.IsPhoneVerified {
		phoneCode, err := domain.NewPhoneCode(r.phoneVerificationCode)
		if err != nil {
			return nil, nil, err
		}
		userAgg.PushEvents(user.NewHumanPhoneCodeAddedEvent(ctx, phoneCode.Code, phoneCode.Expiry))
	} else if human.Phone != nil && human.PhoneNumber != "" && human.IsPhoneVerified {
		userAgg.PushEvents(user.NewHumanPhoneVerifiedEvent(ctx))
	}

	if orgMember != nil {
		//TODO: Add org Member
	}

	return userAgg, addedHuman, nil
}

//ResendInitialMail resend inital mail and changes email if provided
func (r *CommandSide) ResendInitialMail(ctx context.Context, userID, email, resourceowner string) (err error) {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M9fs", "Errors.User.UserIDMissing")
	}

	existingEmail, err := r.emailWriteModel(ctx, userID, resourceowner)
	if err != nil {
		return err
	}
	if existingEmail.UserState == domain.UserStateUnspecified || existingEmail.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-2M9df", "Errors.User.NotFound")
	}
	if existingEmail.UserState != domain.UserStateInitial {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M9sd", "Errors.User.AlreadyInitialised")
	}
	userAgg := UserAggregateFromWriteModel(&existingEmail.WriteModel)
	if email != "" && existingEmail.Email != email {
		changedEvent, _ := existingEmail.NewChangedEvent(ctx, email)
		userAgg.PushEvents(changedEvent)
	}
	initCode, err := domain.NewInitUserCode(r.initializeUserCode)
	if err != nil {
		return err
	}
	userAgg.PushEvents(user.NewHumanInitialCodeAddedEvent(ctx, initCode.Code, initCode.Expiry))
	return r.eventstore.PushAggregate(ctx, existingEmail, userAgg)
}

func (r *CommandSide) HumanSkipMFAInit(ctx context.Context, userID, resourceowner string) (err error) {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2xpX9", "Errors.User.UserIDMissing")
	}

	existingHuman, err := r.getHumanWriteModelByID(ctx, userID, resourceowner)
	if err != nil {
		return err
	}
	if existingHuman.UserState == domain.UserStateUnspecified || existingHuman.UserState == domain.UserStateDeleted {
		return caos_errs.ThrowNotFound(nil, "COMMAND-m9cV8", "Errors.User.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingHuman.WriteModel)
	userAgg.PushEvents(user.NewHumanMFAInitSkippedEvent(ctx))
	return r.eventstore.PushAggregate(ctx, existingHuman, userAgg)
}

func createAddHumanEvent(ctx context.Context, username string, human *domain.Human) *user.HumanAddedEvent {
	addEvent := user.NewHumanAddedEvent(
		ctx,
		username,
		human.FirstName,
		human.LastName,
		human.NickName,
		human.DisplayName,
		human.PreferredLanguage,
		human.Gender,
		human.EmailAddress,
	)
	if human.Phone != nil {
		addEvent.AddPhoneData(human.PhoneNumber)
	}
	if human.Address != nil {
		addEvent.AddAddressData(
			human.Country,
			human.Locality,
			human.PostalCode,
			human.Region,
			human.StreetAddress)
	}
	if human.Password != nil {
		addEvent.AddPasswordData(human.SecretCrypto, human.ChangeRequired)
	}
	return addEvent
}

func createRegisterHumanEvent(ctx context.Context, username string, human *domain.Human) *user.HumanRegisteredEvent {
	addEvent := user.NewHumanRegisteredEvent(
		ctx,
		username,
		human.FirstName,
		human.LastName,
		human.NickName,
		human.DisplayName,
		human.PreferredLanguage,
		human.Gender,
		human.EmailAddress,
	)
	if human.Phone != nil {
		addEvent.AddPhoneData(human.PhoneNumber)
	}
	if human.Address != nil {
		addEvent.AddAddressData(
			human.Country,
			human.Locality,
			human.PostalCode,
			human.Region,
			human.StreetAddress)
	}
	if human.Password != nil {
		addEvent.AddPasswordData(human.SecretCrypto, human.ChangeRequired)
	}
	return addEvent
}

func (r *CommandSide) getHumanWriteModelByID(ctx context.Context, userID, resourceowner string) (*HumanWriteModel, error) {
	humanWriteModel := NewHumanWriteModel(userID, resourceowner)
	err := r.eventstore.FilterToQueryReducer(ctx, humanWriteModel)
	if err != nil {
		return nil, err
	}
	return humanWriteModel, nil
}
