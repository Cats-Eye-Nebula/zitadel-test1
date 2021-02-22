package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	usr_repo "github.com/caos/zitadel/internal/v2/repository/user"
)

func (r *CommandSide) getHumanU2FTokens(ctx context.Context, userID, resourceowner string) ([]*domain.WebAuthNToken, error) {
	tokenReadModel := NewHumanU2FTokensReadModel(userID, resourceowner)
	err := r.eventstore.FilterToQueryReducer(ctx, tokenReadModel)
	if err != nil {
		return nil, err
	}
	if tokenReadModel.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-4M0ds", "Errors.User.NotFound")
	}
	return readModelToU2FTokens(tokenReadModel), nil
}

func (r *CommandSide) getHumanPasswordlessTokens(ctx context.Context, userID, resourceowner string) ([]*domain.WebAuthNToken, error) {
	tokenReadModel := NewHumanPasswordlessTokensReadModel(userID, resourceowner)
	err := r.eventstore.FilterToQueryReducer(ctx, tokenReadModel)
	if err != nil {
		return nil, err
	}
	if tokenReadModel.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-Mv9sd", "Errors.User.NotFound")
	}
	return readModelToPasswordlessTokens(tokenReadModel), nil
}

func (r *CommandSide) getHumanU2FLogin(ctx context.Context, userID, authReqID, resourceowner string) (*domain.WebAuthNLogin, error) {
	tokenReadModel := NewHumanU2FLoginReadModel(userID, authReqID, resourceowner)
	err := r.eventstore.FilterToQueryReducer(ctx, tokenReadModel)
	if err != nil {
		return nil, err
	}
	if tokenReadModel.State == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-5m88U", "Errors.User.NotFound")
	}
	return &domain.WebAuthNLogin{
		Challenge: tokenReadModel.Challenge,
	}, nil
}

func (r *CommandSide) getHumanPasswordlessLogin(ctx context.Context, userID, authReqID, resourceowner string) (*domain.WebAuthNLogin, error) {
	tokenReadModel := NewHumanPasswordlessLoginReadModel(userID, authReqID, resourceowner)
	err := r.eventstore.FilterToQueryReducer(ctx, tokenReadModel)
	if err != nil {
		return nil, err
	}
	if tokenReadModel.State == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-fm84R", "Errors.User.NotFound")
	}
	return &domain.WebAuthNLogin{
		Challenge: tokenReadModel.Challenge,
	}, nil
}

func (r *CommandSide) HumanAddU2FSetup(ctx context.Context, userID, resourceowner string, isLoginUI bool) (*domain.WebAuthNToken, error) {
	u2fTokens, err := r.getHumanU2FTokens(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}
	addWebAuthN, userAgg, webAuthN, err := r.addHumanWebAuthN(ctx, userID, resourceowner, isLoginUI, u2fTokens)
	if err != nil {
		return nil, err
	}

	events, err := r.eventstore.PushEvents(ctx, usr_repo.NewHumanU2FAddedEvent(ctx, userAgg, addWebAuthN.WebauthNTokenID, webAuthN.Challenge))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addWebAuthN, events...)
	if err != nil {
		return nil, err
	}

	createdWebAuthN := writeModelToWebAuthN(addWebAuthN)
	createdWebAuthN.CredentialCreationData = webAuthN.CredentialCreationData
	createdWebAuthN.AllowedCredentialIDs = webAuthN.AllowedCredentialIDs
	createdWebAuthN.UserVerification = webAuthN.UserVerification
	return createdWebAuthN, nil
}

func (r *CommandSide) HumanAddPasswordlessSetup(ctx context.Context, userID, resourceowner string, isLoginUI bool) (*domain.WebAuthNToken, error) {
	passwordlessTokens, err := r.getHumanPasswordlessTokens(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}
	addWebAuthN, userAgg, webAuthN, err := r.addHumanWebAuthN(ctx, userID, resourceowner, isLoginUI, passwordlessTokens)
	if err != nil {
		return nil, err
	}

	events, err := r.eventstore.PushEvents(ctx, usr_repo.NewHumanPasswordlessAddedEvent(ctx, userAgg, addWebAuthN.WebauthNTokenID, webAuthN.Challenge))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addWebAuthN, events...)
	if err != nil {
		return nil, err
	}

	createdWebAuthN := writeModelToWebAuthN(addWebAuthN)
	createdWebAuthN.CredentialCreationData = webAuthN.CredentialCreationData
	createdWebAuthN.AllowedCredentialIDs = webAuthN.AllowedCredentialIDs
	createdWebAuthN.UserVerification = webAuthN.UserVerification
	return createdWebAuthN, nil
}

func (r *CommandSide) addHumanWebAuthN(ctx context.Context, userID, resourceowner string, isLoginUI bool, tokens []*domain.WebAuthNToken) (*HumanWebAuthNWriteModel, *eventstore.Aggregate, *domain.WebAuthNToken, error) {
	if userID == "" || resourceowner == "" {
		return nil, nil, nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3M0od", "Errors.IDMissing")
	}
	user, err := r.getHuman(ctx, userID, resourceowner)
	if err != nil {
		return nil, nil, nil, err
	}
	org, err := r.getOrg(ctx, user.ResourceOwner)
	if err != nil {
		return nil, nil, nil, err
	}
	orgPolicy, err := r.getOrgIAMPolicy(ctx, org.AggregateID)
	if err != nil {
		return nil, nil, nil, err
	}
	accountName := domain.GenerateLoginName(user.GetUsername(), org.PrimaryDomain, orgPolicy.UserLoginMustBeDomain)
	if accountName == "" {
		accountName = user.EmailAddress
	}
	webAuthN, err := r.webauthn.BeginRegistration(user, accountName, domain.AuthenticatorAttachmentUnspecified, domain.UserVerificationRequirementDiscouraged, isLoginUI, tokens...)
	if err != nil {
		return nil, nil, nil, err
	}
	tokenID, err := r.idGenerator.Next()
	if err != nil {
		return nil, nil, nil, err
	}
	addWebAuthN, err := r.webauthNWriteModelByID(ctx, userID, tokenID, resourceowner)
	if err != nil {
		return nil, nil, nil, err
	}

	userAgg := UserAggregateFromWriteModel(&addWebAuthN.WriteModel)
	return addWebAuthN, userAgg, webAuthN, nil
}

func (r *CommandSide) HumanVerifyU2FSetup(ctx context.Context, userID, resourceowner, tokenName, userAgentID string, credentialData []byte) error {
	u2fTokens, err := r.getHumanU2FTokens(ctx, userID, resourceowner)
	if err != nil {
		return err
	}
	userAgg, webAuthN, verifyWebAuthN, err := r.verifyHumanWebAuthN(ctx, userID, resourceowner, tokenName, userAgentID, credentialData, u2fTokens)
	if err != nil {
		return err
	}

	_, err = r.eventstore.PushEvents(ctx,
		usr_repo.NewHumanU2FVerifiedEvent(
			ctx,
			userAgg,
			verifyWebAuthN.WebauthNTokenID, //TODO: webAuthN andverifyWebAuthN same TokenID?
			webAuthN.WebAuthNTokenName,
			webAuthN.AttestationType,
			webAuthN.KeyID,
			webAuthN.PublicKey,
			webAuthN.AAGUID,
			webAuthN.SignCount,
		),
	)
	return err
}

func (r *CommandSide) HumanHumanPasswordlessSetup(ctx context.Context, userID, resourceowner, tokenName, userAgentID string, credentialData []byte) error {
	u2fTokens, err := r.getHumanPasswordlessTokens(ctx, userID, resourceowner)
	if err != nil {
		return err
	}
	userAgg, webAuthN, verifyWebAuthN, err := r.verifyHumanWebAuthN(ctx, userID, resourceowner, tokenName, userAgentID, credentialData, u2fTokens)
	if err != nil {
		return err
	}

	_, err = r.eventstore.PushEvents(ctx,
		usr_repo.NewHumanPasswordlessVerifiedEvent(
			ctx,
			userAgg,
			verifyWebAuthN.WebauthNTokenID, //TODO: webAuthN andverifyWebAuthN same TokenID?
			webAuthN.WebAuthNTokenName,
			webAuthN.AttestationType,
			webAuthN.KeyID,
			webAuthN.PublicKey,
			webAuthN.AAGUID,
			webAuthN.SignCount,
		),
	)
	return err
}

func (r *CommandSide) verifyHumanWebAuthN(ctx context.Context, userID, resourceowner, tokenName, userAgentID string, credentialData []byte, tokens []*domain.WebAuthNToken) (*eventstore.Aggregate, *domain.WebAuthNToken, *HumanWebAuthNWriteModel, error) {
	if userID == "" || resourceowner == "" {
		return nil, nil, nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3M0od", "Errors.IDMissing")
	}
	user, err := r.getHuman(ctx, userID, resourceowner)
	if err != nil {
		return nil, nil, nil, err
	}
	_, token := domain.GetTokenToVerify(tokens)
	webAuthN, err := r.webauthn.FinishRegistration(user, token, tokenName, credentialData, userAgentID != "")
	if err != nil {
		return nil, nil, nil, err
	}

	verifyWebAuthN, err := r.webauthNWriteModelByID(ctx, userID, token.WebAuthNTokenID, resourceowner)
	if err != nil {
		return nil, nil, nil, err
	}

	userAgg := UserAggregateFromWriteModel(&verifyWebAuthN.WriteModel)
	return userAgg, webAuthN, verifyWebAuthN, nil
}

func (r *CommandSide) HumanBeginU2FLogin(ctx context.Context, userID, resourceOwner string, authRequest *domain.AuthRequest, isLoginUI bool) (*domain.WebAuthNLogin, error) {
	u2fTokens, err := r.getHumanU2FTokens(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}

	userAgg, webAuthNLogin, err := r.beginWebAuthNLogin(ctx, userID, resourceOwner, u2fTokens, isLoginUI)
	if err != nil {
		return nil, err
	}

	_, err = r.eventstore.PushEvents(ctx,
		usr_repo.NewHumanU2FBeginLoginEvent(
			ctx,
			userAgg,
			webAuthNLogin.Challenge,
			authRequestDomainToAuthRequestInfo(authRequest),
		),
	)

	return webAuthNLogin, err
}

func (r *CommandSide) HumanBeginPasswordlessLogin(ctx context.Context, userID, resourceOwner string, authRequest *domain.AuthRequest, isLoginUI bool) (*domain.WebAuthNLogin, error) {
	u2fTokens, err := r.getHumanPasswordlessTokens(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}

	userAgg, webAuthNLogin, err := r.beginWebAuthNLogin(ctx, userID, resourceOwner, u2fTokens, isLoginUI)
	if err != nil {
		return nil, err
	}
	_, err = r.eventstore.PushEvents(ctx,
		usr_repo.NewHumanPasswordlessBeginLoginEvent(
			ctx,
			userAgg,
			webAuthNLogin.Challenge,
			authRequestDomainToAuthRequestInfo(authRequest),
		),
	)
	return webAuthNLogin, err
}

func (r *CommandSide) beginWebAuthNLogin(ctx context.Context, userID, resourceOwner string, tokens []*domain.WebAuthNToken, isLoginUI bool) (*eventstore.Aggregate, *domain.WebAuthNLogin, error) {
	if userID == "" {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-hh8K9", "Errors.IDMissing")
	}

	human, err := r.getHuman(ctx, userID, resourceOwner)
	if err != nil {
		return nil, nil, err
	}
	webAuthNLogin, err := r.webauthn.BeginLogin(human, domain.UserVerificationRequirementDiscouraged, isLoginUI, tokens...)
	if err != nil {
		return nil, nil, err
	}

	writeModel, err := r.webauthNWriteModelByID(ctx, userID, "", resourceOwner)
	if err != nil {
		return nil, nil, err
	}
	userAgg := UserAggregateFromWriteModel(&writeModel.WriteModel)

	return userAgg, webAuthNLogin, nil
}

func (r *CommandSide) HumanFinishU2FLogin(ctx context.Context, userID, resourceOwner string, credentialData []byte, authRequest *domain.AuthRequest, isLoginUI bool) error {
	webAuthNLogin, err := r.getHumanU2FLogin(ctx, userID, authRequest.ID, resourceOwner)
	if err != nil {
		return err
	}
	u2fTokens, err := r.getHumanU2FTokens(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}

	userAgg, token, signCount, err := r.finishWebAuthNLogin(ctx, userID, resourceOwner, credentialData, webAuthNLogin, u2fTokens, isLoginUI)
	if err != nil {
		return err
	}

	_, err = r.eventstore.PushEvents(ctx,
		usr_repo.NewHumanU2FSignCountChangedEvent(
			ctx,
			userAgg,
			token.WebAuthNTokenID,
			signCount,
		),
	)

	return err
}

func (r *CommandSide) HumanFinishPasswordlessLogin(ctx context.Context, userID, resourceOwner string, credentialData []byte, authRequest *domain.AuthRequest, isLoginUI bool) error {
	webAuthNLogin, err := r.getHumanPasswordlessLogin(ctx, userID, authRequest.ID, resourceOwner)
	if err != nil {
		return err
	}

	passwordlessTokens, err := r.getHumanPasswordlessTokens(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}

	userAgg, token, signCount, err := r.finishWebAuthNLogin(ctx, userID, resourceOwner, credentialData, webAuthNLogin, passwordlessTokens, isLoginUI)
	if err != nil {
		return err
	}

	_, err = r.eventstore.PushEvents(ctx,
		usr_repo.NewHumanPasswordlessSignCountChangedEvent(
			ctx,
			userAgg,
			token.WebAuthNTokenID,
			signCount,
		),
	)
	return err
}

func (r *CommandSide) finishWebAuthNLogin(ctx context.Context, userID, resourceOwner string, credentialData []byte, webAuthN *domain.WebAuthNLogin, tokens []*domain.WebAuthNToken, isLoginUI bool) (*eventstore.Aggregate, *domain.WebAuthNToken, uint32, error) {
	if userID == "" {
		return nil, nil, 0, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-hh8K9", "Errors.IDMissing")
	}

	human, err := r.getHuman(ctx, userID, resourceOwner)
	if err != nil {
		return nil, nil, 0, err
	}
	keyID, signCount, err := r.webauthn.FinishLogin(human, webAuthN, credentialData, isLoginUI, tokens...)
	if err != nil && keyID == nil {
		return nil, nil, 0, err
	}

	_, token := domain.GetTokenByKeyID(tokens, keyID)
	if token == nil {
		return nil, nil, 0, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3b7zs", "Errors.User.WebAuthN.NotFound")
	}

	writeModel, err := r.webauthNWriteModelByID(ctx, userID, "", resourceOwner)
	if err != nil {
		return nil, nil, 0, err
	}
	userAgg := UserAggregateFromWriteModel(&writeModel.WriteModel)

	return userAgg, token, signCount, nil
}

func (r *CommandSide) HumanRemoveU2F(ctx context.Context, userID, webAuthNID, resourceOwner string) error {
	event := usr_repo.PrepareHumanU2FRemovedEvent(ctx, webAuthNID)
	return r.removeHumanWebAuthN(ctx, userID, webAuthNID, resourceOwner, event)
}

func (r *CommandSide) HumanRemovePasswordless(ctx context.Context, userID, webAuthNID, resourceOwner string) error {
	event := usr_repo.PrepareHumanPasswordlessRemovedEvent(ctx, webAuthNID)
	return r.removeHumanWebAuthN(ctx, userID, webAuthNID, resourceOwner, event)
}

func (r *CommandSide) removeHumanWebAuthN(ctx context.Context, userID, webAuthNID, resourceOwner string, preparedEvent func(*eventstore.Aggregate) eventstore.EventPusher) error {
	if userID == "" || webAuthNID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-6M9de", "Errors.IDMissing")
	}

	existingWebAuthN, err := r.webauthNWriteModelByID(ctx, userID, webAuthNID, resourceOwner)
	if err != nil {
		return err
	}
	if existingWebAuthN.State == domain.MFAStateUnspecified || existingWebAuthN.State == domain.MFAStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-2M9ds", "Errors.User.ExternalIDP.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingWebAuthN.WriteModel)
	_, err = r.eventstore.PushEvents(ctx, preparedEvent(userAgg))
	return err
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
