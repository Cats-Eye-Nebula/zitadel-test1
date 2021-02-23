package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"
	"reflect"

	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/org"
)

func (r *CommandSide) AddOrgMember(ctx context.Context, member *domain.Member) (*domain.Member, error) {
	addedMember := NewOrgMemberWriteModel(member.AggregateID, member.UserID)
	orgAgg := OrgAggregateFromWriteModel(&addedMember.WriteModel)
	event, err := r.addOrgMember(ctx, orgAgg, addedMember, member)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := r.eventstore.PushEvents(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedMember, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return memberWriteModelToMember(&addedMember.MemberWriteModel), nil
}

func (r *CommandSide) addOrgMember(ctx context.Context, orgAgg *eventstore.Aggregate, addedMember *OrgMemberWriteModel, member *domain.Member) (eventstore.EventPusher, error) {
	//TODO: check if roles valid

	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-W8m4l", "Errors.Org.MemberInvalid")
	}

	err := r.eventstore.FilterToQueryReducer(ctx, addedMember)
	if err != nil {
		return nil, err
	}
	if addedMember.State == domain.MemberStateActive {
		return nil, errors.ThrowAlreadyExists(nil, "Org-PtXi1", "Errors.Org.Member.AlreadyExists")
	}

	return org.NewMemberAddedEvent(ctx, orgAgg, member.UserID, member.Roles...), nil
}

//ChangeOrgMember updates an existing member
func (r *CommandSide) ChangeOrgMember(ctx context.Context, member *domain.Member) (*domain.Member, error) {
	//TODO: check if roles valid

	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-LiaZi", "Errors.Org.MemberInvalid")
	}

	existingMember, err := r.orgMemberWriteModelByID(ctx, member.AggregateID, member.UserID)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(existingMember.Roles, member.Roles) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-LiaZi", "Errors.Org.Member.RolesNotChanged")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingMember.MemberWriteModel.WriteModel)
	pushedEvents, err := r.eventstore.PushEvents(ctx, org.NewMemberChangedEvent(ctx, orgAgg, member.UserID, member.Roles...))
	err = AppendAndReduce(existingMember, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return memberWriteModelToMember(&existingMember.MemberWriteModel), nil
}

func (r *CommandSide) RemoveOrgMember(ctx context.Context, orgID, userID string) error {
	m, err := r.orgMemberWriteModelByID(ctx, orgID, userID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if errors.IsNotFound(err) {
		return nil
	}

	orgAgg := OrgAggregateFromWriteModel(&m.MemberWriteModel.WriteModel)
	_, err = r.eventstore.PushEvents(ctx, org.NewMemberRemovedEvent(ctx, orgAgg, userID))
	return err
}

func (r *CommandSide) orgMemberWriteModelByID(ctx context.Context, orgID, userID string) (member *OrgMemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewOrgMemberWriteModel(orgID, userID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if writeModel.State == domain.MemberStateUnspecified || writeModel.State == domain.MemberStateRemoved {
		return nil, errors.ThrowNotFound(nil, "Org-D8JxR", "Errors.NotFound")
	}

	return writeModel, nil
}
