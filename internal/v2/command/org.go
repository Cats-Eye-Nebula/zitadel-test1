package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/org"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

func (r *CommandSide) GetOrg(ctx context.Context, aggregateID string) (*domain.Org, error) {
	orgWriteModel := NewOrgWriteModel(aggregateID)
	err := r.eventstore.FilterToQueryReducer(ctx, orgWriteModel)
	if err != nil {
		return nil, err
	}
	return orgWriteModelToOrg(orgWriteModel), nil
}

func (r *CommandSide) SetUpOrg(ctx context.Context, organisation *domain.Org, admin *domain.User) (*domain.Org, error) {
	orgAgg, userAgg, err := r.setUpOrg(ctx, organisation, admin)
	if err != nil {
		return nil, err
	}

	_, err = r.eventstore.PushAggregates(ctx, orgAgg[0], orgAgg[1], userAgg)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *CommandSide) setUpOrg(ctx context.Context, organisation *domain.Org, admin *domain.User) ([]*org.Aggregate, *user.Aggregate, error) {
	orgAgg, _, err := r.addOrg(ctx, organisation)
	if err != nil {
		return nil, nil, nil
	}

	userAgg, _, err := r.addHuman(ctx, orgAgg.ID(), admin.UserName, admin.Human)
	if err != nil {
		return nil, nil, nil
	}

	addedMember := NewOrgMemberWriteModel(orgAgg.ID(), userAgg.ID())
	orgAgg2 := OrgAggregateFromWriteModel(&addedMember.WriteModel)
	err = r.addOrgMember(ctx, orgAgg2, addedMember, domain.NewMember(orgAgg2.ID(), userAgg.ID(), domain.OrgOwnerRole)) //TODO: correct?
	if err != nil {
		return nil, nil, err
	}
	return []*org.Aggregate{orgAgg, orgAgg2}, userAgg, nil
}

func (r *CommandSide) addOrg(ctx context.Context, organisation *domain.Org) (_ *org.Aggregate, _ *OrgWriteModel, err error) {
	if organisation == nil || !organisation.IsValid() {
		return nil, nil, caos_errs.ThrowInvalidArgument(nil, "COMM-deLSk", "Errors.Org.Invalid")
	}

	organisation.AggregateID, err = r.idGenerator.Next()
	if err != nil {
		return nil, nil, caos_errs.ThrowInternal(err, "COMMA-OwciI", "Errors.Internal")
	}
	organisation.AddIAMDomain(r.iamDomain)
	addedOrg := NewOrgWriteModel(organisation.AggregateID)

	orgAgg := OrgAggregateFromWriteModel(&addedOrg.WriteModel)
	//TODO: uniqueness org name
	orgAgg.PushEvents(org.NewOrgAddedEvent(ctx, organisation.Name))
	for _, orgDomain := range organisation.Domains {
		if err := r.addOrgDomain(ctx, orgAgg, NewOrgDomainWriteModel(orgAgg.ID(), orgDomain.Domain), orgDomain); err != nil {
			return nil, nil, err
		}
	}
	return orgAgg, addedOrg, nil
}
