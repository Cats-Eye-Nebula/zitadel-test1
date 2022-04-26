package command

import (
	"context"

	"github.com/caos/logging"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddAPIApplication(ctx context.Context, application *domain.APIApp, resourceOwner string) (_ *domain.APIApp, err error) {
	if application == nil || application.AggregateID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "PROJECT-5m9E", "Errors.Application.Invalid")
	}
	project, err := c.getProjectByID(ctx, application.AggregateID, resourceOwner)
	if err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(err, "PROJECT-9fnsf", "Errors.Project.NotFound")
	}
	addedApplication := NewAPIApplicationWriteModel(application.AggregateID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&addedApplication.WriteModel)
	events, stringPw, err := c.addAPIApplication(ctx, projectAgg, project, application, resourceOwner)
	if err != nil {
		return nil, err
	}
	addedApplication.AppID = application.AppID
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedApplication, pushedEvents...)
	if err != nil {
		return nil, err
	}
	result := apiWriteModelToAPIConfig(addedApplication)
	result.ClientSecretString = stringPw
	return result, nil
}

func (c *Commands) addAPIApplication(ctx context.Context, projectAgg *eventstore.Aggregate, proj *domain.Project, apiAppApp *domain.APIApp, resourceOwner string) (events []eventstore.Command, stringPW string, err error) {
	if !apiAppApp.IsValid() {
		return nil, "", caos_errs.ThrowInvalidArgument(nil, "PROJECT-Bff2g", "Errors.Application.Invalid")
	}
	apiAppApp.AppID, err = c.idGenerator.Next()
	if err != nil {
		return nil, "", err
	}

	events = []eventstore.Command{
		project.NewApplicationAddedEvent(ctx, projectAgg, apiAppApp.AppID, apiAppApp.AppName),
	}

	var stringPw string
	err = domain.SetNewClientID(apiAppApp, c.idGenerator, proj)
	if err != nil {
		return nil, "", err
	}
	stringPw, err = domain.SetNewClientSecretIfNeeded(apiAppApp, c.applicationSecretGenerator)
	if err != nil {
		return nil, "", err
	}
	events = append(events, project.NewAPIConfigAddedEvent(ctx,
		projectAgg,
		apiAppApp.AppID,
		apiAppApp.ClientID,
		apiAppApp.ClientSecret,
		apiAppApp.AuthMethodType))

	return events, stringPw, nil
}

func (c *Commands) ChangeAPIApplication(ctx context.Context, apiApp *domain.APIApp, resourceOwner string) (*domain.APIApp, error) {
	if apiApp.AppID == "" || apiApp.AggregateID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-1m900", "Errors.Project.App.APIConfigInvalid")
	}

	existingAPI, err := c.getAPIAppWriteModel(ctx, apiApp.AggregateID, apiApp.AppID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingAPI.State == domain.AppStateUnspecified || existingAPI.State == domain.AppStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-2n8uU", "Errors.Project.App.NotExisting")
	}
	if !existingAPI.IsAPI() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-Gnwt3", "Errors.Project.App.IsNotAPI")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingAPI.WriteModel)
	changedEvent, hasChanged, err := existingAPI.NewChangedEvent(
		ctx,
		projectAgg,
		apiApp.AppID,
		apiApp.AuthMethodType)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-1m88i", "Errors.NoChangesFound")
	}

	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingAPI, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return apiWriteModelToAPIConfig(existingAPI), nil
}

func (c *Commands) ChangeAPIApplicationSecret(ctx context.Context, projectID, appID, resourceOwner string) (*domain.APIApp, error) {
	if projectID == "" || appID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-99i83", "Errors.IDMissing")
	}

	existingAPI, err := c.getAPIAppWriteModel(ctx, projectID, appID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingAPI.State == domain.AppStateUnspecified || existingAPI.State == domain.AppStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-2g66f", "Errors.Project.App.NotExisting")
	}
	if !existingAPI.IsAPI() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-aeH4", "Errors.Project.App.IsNotAPI")
	}
	cryptoSecret, stringPW, err := domain.NewClientSecret(c.applicationSecretGenerator)
	if err != nil {
		return nil, err
	}

	projectAgg := ProjectAggregateFromWriteModel(&existingAPI.WriteModel)

	pushedEvents, err := c.eventstore.Push(ctx, project.NewAPIConfigSecretChangedEvent(ctx, projectAgg, appID, cryptoSecret))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingAPI, pushedEvents...)
	if err != nil {
		return nil, err
	}

	result := apiWriteModelToAPIConfig(existingAPI)
	result.ClientSecretString = stringPW
	return result, err
}

func (c *Commands) VerifyAPIClientSecret(ctx context.Context, projectID, appID, secret string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	app, err := c.getAPIAppWriteModel(ctx, projectID, appID, "")
	if err != nil {
		return err
	}
	if !app.State.Exists() {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-DFnbf", "Errors.Project.App.NoExisting")
	}
	if !app.IsAPI() {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-Bf3fw", "Errors.Project.App.IsNotAPI")
	}
	if app.ClientSecret == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-D3t5g", "Errors.Project.App.APIConfigInvalid")
	}

	projectAgg := ProjectAggregateFromWriteModel(&app.WriteModel)
	ctx, spanPasswordComparison := tracing.NewNamedSpan(ctx, "crypto.CompareHash")
	err = crypto.CompareHash(app.ClientSecret, []byte(secret), c.userPasswordAlg)
	spanPasswordComparison.EndWithError(err)
	if err == nil {
		_, err = c.eventstore.Push(ctx, project.NewAPIConfigSecretCheckSucceededEvent(ctx, projectAgg, app.AppID))
		return err
	}
	_, err = c.eventstore.Push(ctx, project.NewAPIConfigSecretCheckFailedEvent(ctx, projectAgg, app.AppID))
	logging.Log("COMMAND-g3f12").OnError(err).Error("could not push event APIClientSecretCheckFailed")
	return caos_errs.ThrowInvalidArgument(nil, "COMMAND-SADfg", "Errors.Project.App.ClientSecretInvalid")
}

func (c *Commands) getAPIAppWriteModel(ctx context.Context, projectID, appID, resourceOwner string) (*APIApplicationWriteModel, error) {
	appWriteModel := NewAPIApplicationWriteModelWithAppID(projectID, appID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, appWriteModel)
	if err != nil {
		return nil, err
	}
	return appWriteModel, nil
}
