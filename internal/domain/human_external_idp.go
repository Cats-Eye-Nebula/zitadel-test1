package domain

import es_models "github.com/caos/zitadel/internal/eventstore/v1/models"

type ExternalIDP struct {
	es_models.ObjectRoot

	IDPConfigID    string
	ExternalUserID string
	DisplayName    string
}

func (idp *ExternalIDP) IsValid() bool {
	return idp.IDPConfigID != "" && idp.ExternalUserID != ""
}

type ExternalIDPState int32

const (
	ExternalIDPStateUnspecified ExternalIDPState = iota
	ExternalIDPStateActive
	ExternalIDPStateRemoved

	externalIDPStateCount
)

func (s ExternalIDPState) Valid() bool {
	return s >= 0 && s < externalIDPStateCount
}
