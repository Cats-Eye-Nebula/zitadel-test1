package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/id"
)

type OIDCConfig struct {
	models.ObjectRoot

	AppID                    string
	ClientID                 string
	ClientSecret             *crypto.CryptoValue
	ClientSecretString       string
	RedirectUris             []string
	ResponseTypes            []OIDCResponseType
	GrantTypes               []OIDCGrantType
	ApplicationType          OIDCApplicationType
	AuthMethodType           OIDCAuthMethodType
	PostLogoutRedirectUris   []string
	OIDCVersion              OIDCVersion
	Compliance               *Compliance
	DevMode                  bool
	AccessTokenType          OIDCTokenType
	AccessTokenRoleAssertion bool
	IDTokenRoleAssertion     bool
	IDTokenUserinfoAssertion bool
	ClockSkew                time.Duration
}

type OIDCVersion int32

const (
	OIDCVersionV1 OIDCVersion = iota
)

type OIDCResponseType int32

const (
	OIDCResponseTypeCode OIDCResponseType = iota
	OIDCResponseTypeIDToken
	OIDCResponseTypeIDTokenToken
)

type OIDCGrantType int32

const (
	OIDCGrantTypeAuthorizationCode OIDCGrantType = iota
	OIDCGrantTypeImplicit
	OIDCGrantTypeRefreshToken
)

type OIDCApplicationType int32

const (
	OIDCApplicationTypeWeb OIDCApplicationType = iota
	OIDCApplicationTypeUserAgent
	OIDCApplicationTypeNative
)

type OIDCAuthMethodType int32

const (
	OIDCAuthMethodTypeBasic OIDCAuthMethodType = iota
	OIDCAuthMethodTypePost
	OIDCAuthMethodTypeNone
)

type Compliance struct {
	NoneCompliant bool
	Problems      []string
}

type OIDCTokenType int32

const (
	OIDCTokenTypeBearer OIDCTokenType = iota
	OIDCTokenTypeJWT
)

func (c *OIDCConfig) IsValid() bool {
	grantTypes := c.getRequiredGrantTypes()
	for _, grantType := range grantTypes {
		ok := containsOIDCGrantType(c.GrantTypes, grantType)
		if !ok {
			return false
		}
	}
	return true
}

//ClientID random_number@projectname (eg. 495894098234@zitadel)
func (c *OIDCConfig) GenerateNewClientID(idGenerator id.Generator, project *Project) error {
	rndID, err := idGenerator.Next()
	if err != nil {
		return err
	}

	c.ClientID = fmt.Sprintf("%v@%v", rndID, strings.ReplaceAll(strings.ToLower("project.Name"), " ", "_")) //TODO: project?
	return nil
}

func (c *OIDCConfig) GenerateClientSecretIfNeeded(generator crypto.Generator) (string, error) {
	if c.AuthMethodType == OIDCAuthMethodTypeNone {
		return "", nil
	}
	return c.GenerateNewClientSecret(generator)
}

func (c *OIDCConfig) GenerateNewClientSecret(generator crypto.Generator) (string, error) {
	cryptoValue, stringSecret, err := crypto.NewCode(generator)
	if err != nil {
		logging.Log("MODEL-UpnTI").OnError(err).Error("unable to create client secret")
		return "", errors.ThrowInternal(err, "MODEL-gH2Wl", "Errors.Project.CouldNotGenerateClientSecret")
	}
	c.ClientSecret = cryptoValue
	return stringSecret, nil
}

func (c *OIDCConfig) getRequiredGrantTypes() []OIDCGrantType {
	grantTypes := make([]OIDCGrantType, 0)
	implicit := false
	for _, r := range c.ResponseTypes {
		switch r {
		case OIDCResponseTypeCode:
			grantTypes = append(grantTypes, OIDCGrantTypeAuthorizationCode)
		case OIDCResponseTypeIDToken, OIDCResponseTypeIDTokenToken:
			if !implicit {
				implicit = true
				grantTypes = append(grantTypes, OIDCGrantTypeImplicit)
			}
		}
	}
	return grantTypes
}

func containsOIDCGrantType(grantTypes []OIDCGrantType, grantType OIDCGrantType) bool {
	for _, gt := range grantTypes {
		if gt == grantType {
			return true
		}
	}
	return false
}
