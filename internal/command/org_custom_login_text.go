package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
)

func (c *Commands) SetOrgLoginText(ctx context.Context, resourceOwner string, loginText *domain.CustomLoginText) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-m29rF", "Errors.ResourceOwnerMissing")
	}
	iamAgg := org.NewAggregate(resourceOwner, resourceOwner)
	events, existingLoginText, err := c.setOrgLoginText(ctx, &iamAgg.Aggregate, loginText)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingLoginText, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingLoginText.WriteModel), nil
}

func (c *Commands) setOrgLoginText(ctx context.Context, orgAgg *eventstore.Aggregate, loginText *domain.CustomLoginText) ([]eventstore.EventPusher, *OrgCustomLoginTextReadModel, error) {
	if !loginText.IsValid() {
		return nil, nil, caos_errs.ThrowInvalidArgument(nil, "ORG-PPo2w", "Errors.CustomText.Invalid")
	}

	existingLoginText, err := c.orgCustomLoginTextWriteModelByID(ctx, orgAgg.ID, loginText.Language)
	if err != nil {
		return nil, nil, err
	}
	events := make([]eventstore.EventPusher, 0)
	events = append(events, c.getSelectLoginTextEvents(ctx, orgAgg, &existingLoginText.CustomLoginTextReadModel, loginText, false)...)
	events = append(events, c.getLoginTextEvents(ctx, orgAgg, &existingLoginText.CustomLoginTextReadModel, loginText, false)...)
	events = append(events, c.getPasswordTextEvents(ctx, orgAgg, &existingLoginText.CustomLoginTextReadModel, loginText, false)...)
	events = append(events, c.getPasswordResetTextEvents(ctx, orgAgg, &existingLoginText.CustomLoginTextReadModel, loginText, false)...)
	events = append(events, c.getInitUserEvents(ctx, orgAgg, &existingLoginText.CustomLoginTextReadModel, loginText, false)...)
	events = append(events, c.getInitDoneEvents(ctx, orgAgg, &existingLoginText.CustomLoginTextReadModel, loginText, false)...)
	events = append(events, c.getInitMFAPromptEvents(ctx, orgAgg, &existingLoginText.CustomLoginTextReadModel, loginText, false)...)
	events = append(events, c.getInitMFAOTPEvents(ctx, orgAgg, &existingLoginText.CustomLoginTextReadModel, loginText, false)...)
	events = append(events, c.getInitMFAU2FEvents(ctx, orgAgg, &existingLoginText.CustomLoginTextReadModel, loginText, false)...)
	events = append(events, c.getInitMFADoneEvents(ctx, orgAgg, &existingLoginText.CustomLoginTextReadModel, loginText, false)...)
	events = append(events, c.getVerifyMFAOTPEvents(ctx, orgAgg, &existingLoginText.CustomLoginTextReadModel, loginText, false)...)
	events = append(events, c.getVerifyMFAU2FEvents(ctx, orgAgg, &existingLoginText.CustomLoginTextReadModel, loginText, false)...)
	events = append(events, c.getRegistrationOptionEvents(ctx, orgAgg, &existingLoginText.CustomLoginTextReadModel, loginText, false)...)
	events = append(events, c.getRegistrationUserEvents(ctx, orgAgg, &existingLoginText.CustomLoginTextReadModel, loginText, false)...)
	events = append(events, c.getRegistrationOrgEvents(ctx, orgAgg, &existingLoginText.CustomLoginTextReadModel, loginText, false)...)
	events = append(events, c.getPasswordlessEvents(ctx, orgAgg, &existingLoginText.CustomLoginTextReadModel, loginText, false)...)
	events = append(events, c.getSuccessLoginEvents(ctx, orgAgg, &existingLoginText.CustomLoginTextReadModel, loginText, false)...)
	return events, existingLoginText, nil
}

func (c *Commands) RemoveOrgLoginTexts(ctx context.Context, resourceOwner string, lang language.Tag) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-1B8dw", "Errors.ResourceOwnerMissing")
	}
	if lang == language.Und {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-5ZZmo", "Errors.CustomMailText.Invalid")
	}
	customText, err := c.orgCustomLoginTextWriteModelByID(ctx, resourceOwner, lang)
	if err != nil {
		return nil, err
	}
	if customText.State == domain.PolicyStateUnspecified || customText.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "Org-9ru44", "Errors.CustomMailText.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&customText.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, org.NewCustomTextTemplateRemovedEvent(ctx, orgAgg, domain.LoginCustomText, lang))
	err = AppendAndReduce(customText, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&customText.WriteModel), nil
}

func (c *Commands) orgCustomLoginTextWriteModelByID(ctx context.Context, orgID string, lang language.Tag) (*OrgCustomLoginTextReadModel, error) {
	writeModel := NewOrgCustomLoginTextReadModel(orgID, lang)
	err := c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
