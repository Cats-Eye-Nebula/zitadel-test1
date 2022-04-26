package command

import (
	"context"
	"reflect"

	"github.com/zitadel/zitadel/internal/eventstore"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddOrgMember(ctx context.Context, member *domain.Member) (*domain.Member, error) {
	if member.UserID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-u8fkf", "Errors.Org.MemberInvalid")
	}
	addedMember := NewOrgMemberWriteModel(member.AggregateID, member.UserID)
	orgAgg := OrgAggregateFromWriteModel(&addedMember.WriteModel)
	err := c.checkUserExists(ctx, addedMember.UserID, "")
	if err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(err, "Org-2H8ds", "Errors.User.NotFound")
	}
	event, err := c.addOrgMember(ctx, orgAgg, addedMember, member)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedMember, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return memberWriteModelToMember(&addedMember.MemberWriteModel), nil
}

func (c *Commands) addOrgMember(ctx context.Context, orgAgg *eventstore.Aggregate, addedMember *OrgMemberWriteModel, member *domain.Member) (eventstore.Command, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-W8m4l", "Errors.Org.MemberInvalid")
	}
	if len(domain.CheckForInvalidRoles(member.Roles, domain.OrgRolePrefix, c.zitadelRoles)) > 0 && len(domain.CheckForInvalidRoles(member.Roles, domain.RoleSelfManagementGlobal, c.zitadelRoles)) > 0 {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-4N8es", "Errors.Org.MemberInvalid")
	}
	err := c.eventstore.FilterToQueryReducer(ctx, addedMember)
	if err != nil {
		return nil, err
	}
	if addedMember.State == domain.MemberStateActive {
		return nil, errors.ThrowAlreadyExists(nil, "Org-PtXi1", "Errors.Org.Member.AlreadyExists")
	}

	return org.NewMemberAddedEvent(ctx, orgAgg, member.UserID, member.Roles...), nil
}

//ChangeOrgMember updates an existing member
func (c *Commands) ChangeOrgMember(ctx context.Context, member *domain.Member) (*domain.Member, error) {
	if !member.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-LiaZi", "Errors.Org.MemberInvalid")
	}
	if len(domain.CheckForInvalidRoles(member.Roles, domain.OrgRolePrefix, c.zitadelRoles)) > 0 {
		return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-m9fG8", "Errors.Org.MemberInvalid")
	}

	existingMember, err := c.orgMemberWriteModelByID(ctx, member.AggregateID, member.UserID)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(existingMember.Roles, member.Roles) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-LiaZi", "Errors.Org.Member.RolesNotChanged")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingMember.MemberWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewMemberChangedEvent(ctx, orgAgg, member.UserID, member.Roles...))
	err = AppendAndReduce(existingMember, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return memberWriteModelToMember(&existingMember.MemberWriteModel), nil
}

func (c *Commands) RemoveOrgMember(ctx context.Context, orgID, userID string) (*domain.ObjectDetails, error) {
	m, err := c.orgMemberWriteModelByID(ctx, orgID, userID)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if errors.IsNotFound(err) {
		return nil, nil
	}

	orgAgg := OrgAggregateFromWriteModel(&m.MemberWriteModel.WriteModel)
	removeEvent := c.removeOrgMember(ctx, orgAgg, userID, false)
	pushedEvents, err := c.eventstore.Push(ctx, removeEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(m, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&m.WriteModel), nil
}

func (c *Commands) removeOrgMember(ctx context.Context, orgAgg *eventstore.Aggregate, userID string, cascade bool) eventstore.Command {
	if cascade {
		return org.NewMemberCascadeRemovedEvent(
			ctx,
			orgAgg,
			userID)
	} else {
		return org.NewMemberRemovedEvent(ctx, orgAgg, userID)
	}
}

func (c *Commands) orgMemberWriteModelByID(ctx context.Context, orgID, userID string) (member *OrgMemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewOrgMemberWriteModel(orgID, userID)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if writeModel.State == domain.MemberStateUnspecified || writeModel.State == domain.MemberStateRemoved {
		return nil, errors.ThrowNotFound(nil, "Org-D8JxR", "Errors.NotFound")
	}

	return writeModel, nil
}
