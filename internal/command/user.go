package command

import (
	"context"
	"fmt"
	"github.com/caos/zitadel/internal/eventstore"
	"time"

	"github.com/caos/zitadel/internal/eventstore/v1/models"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/user"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) ChangeUsername(ctx context.Context, orgID, userID, userName string) error {
	if orgID == "" || userID == "" || userName == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2N9fs", "Errors.IDMissing")
	}

	existingUser, err := c.userWriteModelByID(ctx, userID, orgID)
	if err != nil {
		return err
	}

	if !isUserStateExists(existingUser.UserState) {
		return caos_errs.ThrowNotFound(nil, "COMMAND-5N9ds", "Errors.User.NotFound")
	}

	if existingUser.UserName == userName {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-6m9gs", "Errors.User.UsernameNotChanged")
	}

	orgIAMPolicy, err := c.getOrgIAMPolicy(ctx, orgID)
	if err != nil {
		return err
	}

	if err := CheckOrgIAMPolicyForUserName(userName, orgIAMPolicy); err != nil {
		return err
	}
	userAgg := UserAggregateFromWriteModel(&existingUser.WriteModel)

	_, err = c.eventstore.PushEvents(ctx,
		user.NewUsernameChangedEvent(ctx, userAgg, existingUser.UserName, userName, orgIAMPolicy.UserLoginMustBeDomain))

	return err
}

func (c *Commands) DeactivateUser(ctx context.Context, userID, resourceOwner string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-m0gDf", "Errors.User.UserIDMissing")
	}

	existingUser, err := c.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if !isUserStateExists(existingUser.UserState) {
		return caos_errs.ThrowNotFound(nil, "COMMAND-3M9ds", "Errors.User.NotFound")
	}
	if isUserStateInactive(existingUser.UserState) {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-5M0sf", "Errors.User.AlreadyInactive")
	}

	_, err = c.eventstore.PushEvents(ctx,
		user.NewUserDeactivatedEvent(ctx, UserAggregateFromWriteModel(&existingUser.WriteModel)))
	return err
}

func (c *Commands) ReactivateUser(ctx context.Context, userID, resourceOwner string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4M9ds", "Errors.User.UserIDMissing")
	}

	existingUser, err := c.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if !isUserStateExists(existingUser.UserState) {
		return caos_errs.ThrowNotFound(nil, "COMMAND-4M0sd", "Errors.User.NotFound")
	}
	if !isUserStateInactive(existingUser.UserState) {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-6M0sf", "Errors.User.NotInactive")
	}

	_, err = c.eventstore.PushEvents(ctx,
		user.NewUserReactivatedEvent(ctx, UserAggregateFromWriteModel(&existingUser.WriteModel)))
	return err
}

func (c *Commands) LockUser(ctx context.Context, userID, resourceOwner string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M0sd", "Errors.User.UserIDMissing")
	}

	existingUser, err := c.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if !isUserStateExists(existingUser.UserState) {
		return caos_errs.ThrowNotFound(nil, "COMMAND-5M9fs", "Errors.User.NotFound")
	}
	if !hasUserState(existingUser.UserState, domain.UserStateActive, domain.UserStateInitial) {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3NN8v", "Errors.User.ShouldBeActiveOrInitial")
	}

	_, err = c.eventstore.PushEvents(ctx,
		user.NewUserLockedEvent(ctx, UserAggregateFromWriteModel(&existingUser.WriteModel)))
	return err
}

func (c *Commands) UnlockUser(ctx context.Context, userID, resourceOwner string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-M0dse", "Errors.User.UserIDMissing")
	}

	existingUser, err := c.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if !isUserStateExists(existingUser.UserState) {
		return caos_errs.ThrowNotFound(nil, "COMMAND-M0dos", "Errors.User.NotFound")
	}
	if !hasUserState(existingUser.UserState, domain.UserStateLocked) {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4M0ds", "Errors.User.NotLocked")
	}

	_, err = c.eventstore.PushEvents(ctx,
		user.NewUserUnlockedEvent(ctx, UserAggregateFromWriteModel(&existingUser.WriteModel)))
	return err
}

func (c *Commands) RemoveUser(ctx context.Context, userID, resourceOwner string, cascadingGrantIDs ...string) error {
	if userID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M0ds", "Errors.User.UserIDMissing")
	}

	existingUser, err := c.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if !isUserStateExists(existingUser.UserState) {
		return caos_errs.ThrowNotFound(nil, "COMMAND-5M0od", "Errors.User.NotFound")
	}

	orgIAMPolicy, err := c.getOrgIAMPolicy(ctx, existingUser.ResourceOwner)
	if err != nil {
		return err
	}
	var events []eventstore.EventPusher
	userAgg := UserAggregateFromWriteModel(&existingUser.WriteModel)
	events = append(events, user.NewUserRemovedEvent(ctx, userAgg, existingUser.UserName, orgIAMPolicy.UserLoginMustBeDomain))

	for _, grantID := range cascadingGrantIDs {
		removeEvent, err := c.removeUserGrant(ctx, grantID, "", true)
		if err != nil {
			logging.LogWithFields("COMMAND-5m9oL", "usergrantid", grantID).WithError(err).Warn("could not cascade remove role on user grant")
			continue
		}
		events = append(events, removeEvent)
	}

	_, err = c.eventstore.PushEvents(ctx, events...)
	return err
}

func (c *Commands) AddUserToken(ctx context.Context, orgID, agentID, clientID, userID string, audience, scopes []string, lifetime time.Duration) (*domain.Token, error) {
	if orgID == "" || userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-55n8M", "Errors.IDMissing")
	}

	existingUser, err := c.userWriteModelByID(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingUser.UserState) {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-1d6Gg", "Errors.User.NotFound")
	}

	audience = domain.AddAudScopeToAudience(audience, scopes)

	preferredLanguage := ""
	existingHuman, err := c.getHumanWriteModelByID(ctx, userID, orgID)
	if existingHuman != nil {
		preferredLanguage = existingHuman.PreferredLanguage.String()
	}
	expiration := time.Now().UTC().Add(lifetime)
	tokenID, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}

	userAgg := UserAggregateFromWriteModel(&existingUser.WriteModel)
	_, err = c.eventstore.PushEvents(ctx,
		user.NewUserTokenAddedEvent(ctx, userAgg, tokenID, clientID, agentID, preferredLanguage, audience, scopes, expiration))
	if err != nil {
		return nil, err
	}

	return &domain.Token{
		ObjectRoot: models.ObjectRoot{
			AggregateID: userID,
		},
		TokenID:           tokenID,
		UserAgentID:       agentID,
		ApplicationID:     clientID,
		Audience:          audience,
		Scopes:            scopes,
		Expiration:        expiration,
		PreferredLanguage: preferredLanguage,
	}, nil
}

func (c *Commands) userDomainClaimed(ctx context.Context, userID string) (events []eventstore.EventPusher, _ *UserWriteModel, err error) {
	existingUser, err := c.userWriteModelByID(ctx, userID, "")
	if err != nil {
		return nil, nil, err
	}
	if existingUser.UserState == domain.UserStateUnspecified || existingUser.UserState == domain.UserStateDeleted {
		return nil, nil, caos_errs.ThrowNotFound(nil, "COMMAND-ii9K0", "Errors.User.NotFound")
	}
	changedUserGrant := NewUserWriteModel(userID, existingUser.ResourceOwner)
	userAgg := UserAggregateFromWriteModel(&changedUserGrant.WriteModel)

	orgIAMPolicy, err := c.getOrgIAMPolicy(ctx, existingUser.ResourceOwner)
	if err != nil {
		return nil, nil, err
	}

	id, err := c.idGenerator.Next()
	if err != nil {
		return nil, nil, err
	}
	return []eventstore.EventPusher{
		user.NewDomainClaimedEvent(
			ctx,
			userAgg,
			fmt.Sprintf("%s@temporary.%s", id, c.iamDomain),
			existingUser.UserName,
			orgIAMPolicy.UserLoginMustBeDomain),
	}, changedUserGrant, nil
}

func (c *Commands) UserDomainClaimedSent(ctx context.Context, orgID, userID string) (err error) {
	existingUser, err := c.userWriteModelByID(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if !isUserStateExists(existingUser.UserState) {
		return caos_errs.ThrowNotFound(nil, "COMMAND-5m9gK", "Errors.User.NotFound")
	}

	_, err = c.eventstore.PushEvents(ctx,
		user.NewDomainClaimedSentEvent(ctx, UserAggregateFromWriteModel(&existingUser.WriteModel)))
	return err
}

func (c *Commands) checkUserExists(ctx context.Context, userID, resourceOwner string) error {
	existingUser, err := c.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}
	if !isUserStateExists(existingUser.UserState) {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-4M0fs", "Errors.User.NotFound")
	}
	return nil
}

func (c *Commands) userWriteModelByID(ctx context.Context, userID, resourceOwner string) (writeModel *UserWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewUserWriteModel(userID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
