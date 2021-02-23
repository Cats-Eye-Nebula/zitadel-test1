package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"
	"reflect"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (r *CommandSide) AddIAMMember(ctx context.Context, member *domain.Member) (*domain.Member, error) {
	addedMember := NewIAMMemberWriteModel(member.UserID)
	iamAgg := IAMAggregateFromWriteModel(&addedMember.MemberWriteModel.WriteModel)
	event, err := r.addIAMMember(ctx, iamAgg, addedMember, member)
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

func (r *CommandSide) addIAMMember(ctx context.Context, iamAgg *eventstore.Aggregate, addedMember *IAMMemberWriteModel, member *domain.Member) (eventstore.EventPusher, error) {
	//TODO: check if roles valid

	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-GR34U", "Errors.IAM.MemberInvalid")
	}

	err := r.eventstore.FilterToQueryReducer(ctx, addedMember)
	if err != nil {
		return nil, err
	}
	if addedMember.State == domain.MemberStateActive {
		return nil, errors.ThrowAlreadyExists(nil, "IAM-sdgQ4", "Errors.IAM.Member.AlreadyExists")
	}

	return iam_repo.NewMemberAddedEvent(ctx, iamAgg, member.UserID, member.Roles...), nil
}

//ChangeIAMMember updates an existing member
func (r *CommandSide) ChangeIAMMember(ctx context.Context, member *domain.Member) (*domain.Member, error) {
	//TODO: check if roles valid

	if !member.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-LiaZi", "Errors.IAM.MemberInvalid")
	}

	existingMember, err := r.iamMemberWriteModelByID(ctx, member.UserID)
	if err != nil {
		return nil, err
	}

	if reflect.DeepEqual(existingMember.Roles, member.Roles) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-LiaZi", "Errors.IAM.Member.RolesNotChanged")
	}
	iamAgg := IAMAggregateFromWriteModel(&existingMember.MemberWriteModel.WriteModel)
	pushedEvents, err := r.eventstore.PushEvents(ctx, iam_repo.NewMemberChangedEvent(ctx, iamAgg, member.UserID, member.Roles...))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingMember, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return memberWriteModelToMember(&existingMember.MemberWriteModel), nil
}

func (r *CommandSide) RemoveIAMMember(ctx context.Context, userID string) error {
	m, err := r.iamMemberWriteModelByID(ctx, userID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if errors.IsNotFound(err) {
		return nil
	}

	iamAgg := IAMAggregateFromWriteModel(&m.MemberWriteModel.WriteModel)
	_, err = r.eventstore.PushEvents(ctx, iam_repo.NewMemberRemovedEvent(ctx, iamAgg, userID))
	return err
}

func (r *CommandSide) iamMemberWriteModelByID(ctx context.Context, userID string) (member *IAMMemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMMemberWriteModel(userID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if writeModel.State == domain.MemberStateUnspecified || writeModel.State == domain.MemberStateRemoved {
		return nil, errors.ThrowNotFound(nil, "IAM-D8JxR", "Errors.NotFound")
	}

	return writeModel, nil
}
