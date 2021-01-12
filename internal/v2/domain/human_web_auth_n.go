package domain

import es_models "github.com/caos/zitadel/internal/eventstore/models"

type WebAuthNToken struct {
	es_models.ObjectRoot

	WebAuthNTokenID        string
	CredentialCreationData []byte
	State                  MFAState
	Challenge              string
	AllowedCredentialIDs   [][]byte
	UserVerification       UserVerificationRequirement
	KeyID                  []byte
	PublicKey              []byte
	AttestationType        string
	AAGUID                 []byte
	SignCount              uint32
	WebAuthNTokenName      string
}

type WebAuthNLogin struct {
	es_models.ObjectRoot

	CredentialAssertionData []byte
	Challenge               string
	AllowedCredentialIDs    [][]byte
	UserVerification        UserVerificationRequirement
	//TODO: Add Auth Request
	//*model.AuthRequest
}

type UserVerificationRequirement int32

const (
	UserVerificationRequirementUnspecified UserVerificationRequirement = iota
	UserVerificationRequirementRequired
	UserVerificationRequirementPreferred
	UserVerificationRequirementDiscouraged
)

type AuthenticatorAttachment int32

const (
	AuthenticatorAttachmentUnspecified AuthenticatorAttachment = iota
	AuthenticatorAttachmentPlattform
	AuthenticatorAttachmentCrossPlattform
)

type WebAuthNState int32

const (
	WebAuthNStateUnspecified WebAuthNState = iota
	WebAuthNStateActive
	WebAuthNStateRemoved

	webAuthNStateCount
)

func (s WebAuthNState) Valid() bool {
	return s >= 0 && s < webAuthNStateCount
}
