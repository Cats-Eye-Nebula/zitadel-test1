package grpc

import (
	"github.com/caos/logging"
	admin_model "github.com/caos/zitadel/internal/admin/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/golang/protobuf/ptypes"
)

func setUpRequestToModel(setUp *OrgSetUpRequest) *admin_model.SetupOrg {
	return &admin_model.SetupOrg{
		Org: orgCreateRequestToModel(setUp.Org),
	}
}

func orgCreateRequestToModel(org *CreateOrgRequest) *org_model.Org {
	return &org_model.Org{
		Domain: org.Domain,
		Name:   org.Name,
	}
}

func setUpOrgResponseFromModel(setUp *admin_model.SetupOrg) *OrgSetUpResponse {
	return &OrgSetUpResponse{
		Org: orgFromModel(setUp.Org),
	}
}

func orgsFromModel(orgs []*org_model.Org) []*Org {
	orgList := make([]*Org, len(orgs))
	for i, org := range orgs {
		orgList[i] = orgFromModel(org)
	}
	return orgList
}

func orgFromModel(org *org_model.Org) *Org {
	creationDate, err := ptypes.TimestampProto(org.CreationDate)
	logging.Log("GRPC-GTHsZ").OnError(err).Debug("unable to get timestamp from time")

	changeDate, err := ptypes.TimestampProto(org.ChangeDate)
	logging.Log("GRPC-dVnoj").OnError(err).Debug("unable to get timestamp from time")

	return &Org{
		Domain:       org.Domain,
		ChangeDate:   changeDate,
		CreationDate: creationDate,
		Id:           org.AggregateID,
		Name:         org.Name,
		State:        orgStateFromModel(org.State),
	}
}

func orgStateFromModel(state org_model.OrgState) OrgState {
	switch state {
	case org_model.ORGSTATE_ACTIVE:
		return OrgState_ORGSTATE_ACTIVE
	case org_model.ORGSTATE_INACTIVE:
		return OrgState_ORGSTATE_INACTIVE
	default:
		return OrgState_ORGSTATE_UNSPECIFIED
	}
}
