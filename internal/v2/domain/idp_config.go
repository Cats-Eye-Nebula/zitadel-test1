package domain

import (
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"time"
)

type IDPConfig struct {
	es_models.ObjectRoot
	IDPConfigID string
	Type        IDPConfigType
	Name        string
	StylingType IDPConfigStylingType
	State       IDPConfigState
	OIDCConfig  *OIDCIDPConfig
}

type IDPConfigView struct {
	AggregateID     string
	IDPConfigID     string
	Name            string
	StylingType     IDPConfigStylingType
	State           IDPConfigState
	CreationDate    time.Time
	ChangeDate      time.Time
	Sequence        uint64
	IDPProviderType IdentityProviderType

	IsOIDC                    bool
	OIDCClientID              string
	OIDCClientSecret          *crypto.CryptoValue
	OIDCIssuer                string
	OIDCScopes                []string
	OIDCIDPDisplayNameMapping OIDCMappingField
	OIDCUsernameMapping       OIDCMappingField
}

type OIDCIDPConfig struct {
	es_models.ObjectRoot
	IDPConfigID           string
	ClientID              string
	ClientSecret          *crypto.CryptoValue
	ClientSecretString    string
	Issuer                string
	Scopes                []string
	IDPDisplayNameMapping OIDCMappingField
	UsernameMapping       OIDCMappingField
}

type IDPConfigType int32

const (
	IDPConfigTypeOIDC IDPConfigType = iota
	IDPConfigTypeSAML

	//count is for validation
	idpConfigTypeCount
)

func (f IDPConfigType) Valid() bool {
	return f >= 0 && f < idpConfigTypeCount
}

type IDPConfigState int32

const (
	IDPConfigStateUnspecified IDPConfigState = iota
	IDPConfigStateActive
	IDPConfigStateInactive
	IDPConfigStateRemoved

	idpConfigStateCount
)

func (f IDPConfigState) Valid() bool {
	return f >= 0 && f < idpConfigStateCount
}

type IDPConfigStylingType int32

const (
	IDPConfigStylingTypeUnspecified IDPConfigStylingType = iota
	IDPConfigStylingTypeGoogle

	idpConfigStylingTypeCount
)

func (f IDPConfigStylingType) Valid() bool {
	return f >= 0 && f < idpConfigStylingTypeCount
}
