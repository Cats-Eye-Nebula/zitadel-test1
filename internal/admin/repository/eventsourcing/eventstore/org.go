package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/org/repository/view"
	"github.com/caos/zitadel/internal/telemetry/tracing"

	"github.com/caos/logging"
	admin_view "github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/eventstore"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	"github.com/caos/zitadel/internal/org/repository/view/model"
)

type OrgRepo struct {
	Eventstore    eventstore.Eventstore
	OrgEventstore *org_es.OrgEventstore

	View *admin_view.View

	SearchLimit    uint64
	SystemDefaults systemdefaults.SystemDefaults
}

func (repo *OrgRepo) OrgByID(ctx context.Context, id string) (*org_model.OrgView, error) {
	org, viewErr := repo.View.OrgByID(id)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		org = new(model.OrgView)
	}
	events, esErr := repo.getOrgEvents(ctx, id, org.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-Lsoj7", "Errors.Org.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-PSoc3").WithError(esErr).Debug("error retrieving new events")
		return model.OrgToModel(org), nil
	}
	orgCopy := *org
	for _, event := range events {
		if err := orgCopy.AppendEvent(event); err != nil {
			return model.OrgToModel(&orgCopy), nil
		}
	}
	return model.OrgToModel(&orgCopy), nil
}

func (repo *OrgRepo) SearchOrgs(ctx context.Context, query *org_model.OrgSearchRequest) (*org_model.OrgSearchResult, error) {
	query.EnsureLimit(repo.SearchLimit)
	sequence, err := repo.View.GetLatestOrgSequence()
	logging.Log("EVENT-LXo9w").OnError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Warn("could not read latest iam sequence")
	orgs, count, err := repo.View.SearchOrgs(query)
	if err != nil {
		return nil, err
	}
	result := &org_model.OrgSearchResult{
		Offset:      query.Offset,
		Limit:       query.Limit,
		TotalResult: count,
		Result:      model.OrgsToModel(orgs),
	}
	if err == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *OrgRepo) IsOrgUnique(ctx context.Context, name, domain string) (isUnique bool, err error) {
	return repo.OrgEventstore.IsOrgUnique(ctx, name, domain)
}

func (repo *OrgRepo) GetOrgIAMPolicyByID(ctx context.Context, id string) (*iam_model.OrgIAMPolicyView, error) {
	policy, err := repo.View.OrgIAMPolicyByAggregateID(id)
	if errors.IsNotFound(err) {
		return repo.GetDefaultOrgIAMPolicy(ctx)
	}
	if err != nil {
		return nil, err
	}
	return iam_es_model.OrgIAMViewToModel(policy), err
}

func (repo *OrgRepo) GetDefaultOrgIAMPolicy(ctx context.Context) (*iam_model.OrgIAMPolicyView, error) {
	policy, err := repo.View.OrgIAMPolicyByAggregateID(repo.SystemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	policy.Default = true
	return iam_es_model.OrgIAMViewToModel(policy), err
}

func (repo *OrgRepo) getOrgEvents(ctx context.Context, orgID string, sequence uint64) ([]*models.Event, error) {
	query, err := view.OrgByIDQuery(orgID, sequence)
	if err != nil {
		return nil, err
	}
	return repo.Eventstore.FilterEvents(ctx, query)
}
