package web_auth_n

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/user/human/mfa"
)

const (
	u2fEventPrefix                    = eventstore.EventType("user.human.mfa.u2f.token.")
	HumanU2FTokenAddedType            = u2fEventPrefix + "added"
	HumanU2FTokenVerifiedType         = u2fEventPrefix + "verified"
	HumanU2FTokenSignCountChangedType = u2fEventPrefix + "signcount.changed"
	HumanU2FTokenRemovedType          = u2fEventPrefix + "removed"
	HumanU2FTokenBeginLoginType       = u2fEventPrefix + "begin.login"
	HumanU2FTokenCheckSucceededType   = u2fEventPrefix + "check.succeeded"
	HumanU2FTokenCheckFailedType      = u2fEventPrefix + "check.failed"

	passwordlessEventPrefix                    = eventstore.EventType("user.human.mfa.passwordless.token.")
	HumanPasswordlessTokenAddedType            = passwordlessEventPrefix + "added"
	HumanPasswordlessTokenVerifiedType         = passwordlessEventPrefix + "verified"
	HumanPasswordlessTokenSignCountChangedType = passwordlessEventPrefix + "signcount.changed"
	HumanPasswordlessTokenRemovedType          = passwordlessEventPrefix + "removed"
	HumanPasswordlessTokenBeginLoginType       = passwordlessEventPrefix + "begin.login"
	HumanPasswordlessTokenCheckSucceededType   = passwordlessEventPrefix + "check.succeeded"
	HumanPasswordlessTokenCheckFailedType      = passwordlessEventPrefix + "check.failed"
)

type HumanWebAuthNAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	WebAuthNTokenID string       `json:"webAuthNTokenId"`
	Challenge       string       `json:"challenge"`
	State           mfa.MFAState `json:"-"`
}

func (e *HumanWebAuthNAddedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanWebAuthNAddedEvent) Data() interface{} {
	return e
}

func NewHumanWebAuthNU2FAddedEvent(
	ctx context.Context,
	webAuthNTokenID,
	challenge string,
) *HumanWebAuthNAddedEvent {
	return &HumanWebAuthNAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanU2FTokenAddedType,
		),
		WebAuthNTokenID: webAuthNTokenID,
		Challenge:       challenge,
	}
}

func NewHumanWebAuthNPasswordlessAddedEvent(
	ctx context.Context,
	webAuthNTokenID,
	challenge string,
) *HumanWebAuthNAddedEvent {
	return &HumanWebAuthNAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanPasswordlessTokenAddedType,
		),
		WebAuthNTokenID: webAuthNTokenID,
		Challenge:       challenge,
	}
}

func HumanWebAuthNAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webAuthNAdded := &HumanWebAuthNAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
		State:     mfa.MFAStateNotReady,
	}
	err := json.Unmarshal(event.Data, webAuthNAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-tB8sf", "unable to unmarshal human webAuthN added")
	}
	return webAuthNAdded, nil
}

type HumanWebAuthNVerifiedEvent struct {
	eventstore.BaseEvent `json:"-"`

	WebAuthNTokenID   string       `json:"webAuthNTokenId"`
	KeyID             []byte       `json:"keyId"`
	PublicKey         []byte       `json:"publicKey"`
	AttestationType   string       `json:"attestationType"`
	AAGUID            []byte       `json:"aaguid"`
	SignCount         uint32       `json:"signCount"`
	WebAuthNTokenName string       `json:"webAuthNTokenName"`
	State             mfa.MFAState `json:"-"`
}

func (e *HumanWebAuthNVerifiedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanWebAuthNVerifiedEvent) Data() interface{} {
	return e
}

func NewHumanWebAuthNU2FVerifiedEvent(
	ctx context.Context,
	webAuthNTokenID,
	webAuthNTokenName,
	attestationType string,
	keyID,
	publicKey,
	aaguid []byte,
	signCount uint32,
) *HumanWebAuthNVerifiedEvent {
	return &HumanWebAuthNVerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanU2FTokenVerifiedType,
		),
		WebAuthNTokenID:   webAuthNTokenID,
		KeyID:             keyID,
		PublicKey:         publicKey,
		AttestationType:   attestationType,
		AAGUID:            aaguid,
		SignCount:         signCount,
		WebAuthNTokenName: webAuthNTokenName,
	}
}

func NewHumanWebAuthNPasswordlessVerifiedEvent(
	ctx context.Context,
	webAuthNTokenID,
	webAuthNTokenName,
	attestationType string,
	keyID,
	publicKey,
	aaguid []byte,
	signCount uint32,
) *HumanWebAuthNVerifiedEvent {
	return &HumanWebAuthNVerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanPasswordlessTokenVerifiedType,
		),
		WebAuthNTokenID:   webAuthNTokenID,
		KeyID:             keyID,
		PublicKey:         publicKey,
		AttestationType:   attestationType,
		AAGUID:            aaguid,
		SignCount:         signCount,
		WebAuthNTokenName: webAuthNTokenName,
	}
}

func HumanWebAuthNVerifiedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webauthNVerified := &HumanWebAuthNVerifiedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
		State:     mfa.MFAStateReady,
	}
	err := json.Unmarshal(event.Data, webauthNVerified)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-B0zDs", "unable to unmarshal human webAuthN verified")
	}
	return webauthNVerified, nil
}

type HumanWebAuthNSignCountChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	WebAuthNTokenID string       `json:"webAuthNTokenId"`
	SignCount       uint32       `json:"signCount"`
	State           mfa.MFAState `json:"-"`
}

func (e *HumanWebAuthNSignCountChangedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanWebAuthNSignCountChangedEvent) Data() interface{} {
	return e
}

func NewHumanWebAuthNU2FSignCountChangedEvent(
	ctx context.Context,
	webAuthNTokenID string,
	signCount uint32,
) *HumanWebAuthNSignCountChangedEvent {
	return &HumanWebAuthNSignCountChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanU2FTokenSignCountChangedType,
		),
		WebAuthNTokenID: webAuthNTokenID,
		SignCount:       signCount,
	}
}

func NewHumanWebAuthNPasswordlessSignCountChangedEvent(
	ctx context.Context,
	webAuthNTokenID string,
	signCount uint32,
) *HumanWebAuthNSignCountChangedEvent {
	return &HumanWebAuthNSignCountChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanPasswordlessTokenSignCountChangedType,
		),
		WebAuthNTokenID: webAuthNTokenID,
		SignCount:       signCount,
	}
}

func HumanWebAuthNSignCountChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webauthNVerified := &HumanWebAuthNSignCountChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, webauthNVerified)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-5Gm0s", "unable to unmarshal human webAuthN sign count")
	}
	return webauthNVerified, nil
}

type HumanWebAuthNRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	WebAuthNTokenID string       `json:"webAuthNTokenId"`
	State           mfa.MFAState `json:"-"`
}

func (e *HumanWebAuthNRemovedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanWebAuthNRemovedEvent) Data() interface{} {
	return e
}

func NewHumanWebAuthNU2FRemovedEvent(
	ctx context.Context,
	webAuthNTokenID string,
) *HumanWebAuthNRemovedEvent {
	return &HumanWebAuthNRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanU2FTokenRemovedType,
		),
		WebAuthNTokenID: webAuthNTokenID,
	}
}

func NewHumanWebAuthNPasswordlessRemovedEvent(
	ctx context.Context,
	webAuthNTokenID string,
) *HumanWebAuthNRemovedEvent {
	return &HumanWebAuthNRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanPasswordlessTokenRemovedType,
		),
		WebAuthNTokenID: webAuthNTokenID,
	}
}

func HumanWebAuthNRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webauthNVerified := &HumanWebAuthNRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, webauthNVerified)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-gM9sd", "unable to unmarshal human webAuthN token removed")
	}
	return webauthNVerified, nil
}

type HumanWebAuthNBeginLoginEvent struct {
	eventstore.BaseEvent `json:"-"`

	WebAuthNTokenID string `json:"webAuthNTokenId"`
	Challenge       string `json:"challenge"`
	//TODO: Handle Auth Req??
	//*AuthRequest
}

func (e *HumanWebAuthNBeginLoginEvent) CheckPrevious() bool {
	return true
}

func (e *HumanWebAuthNBeginLoginEvent) Data() interface{} {
	return e
}

func NewHumanWebAuthNU2FBeginLoginEvent(
	ctx context.Context,
	webAuthNTokenID,
	challenge string,
) *HumanWebAuthNBeginLoginEvent {
	return &HumanWebAuthNBeginLoginEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanU2FTokenRemovedType,
		),
		WebAuthNTokenID: webAuthNTokenID,
		Challenge:       challenge,
	}
}

func NewHumanWebAuthNPasswordlessBeginLoginEvent(
	ctx context.Context,
	webAuthNTokenID,
	challenge string,
) *HumanWebAuthNBeginLoginEvent {
	return &HumanWebAuthNBeginLoginEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanPasswordlessTokenRemovedType,
		),
		WebAuthNTokenID: webAuthNTokenID,
		Challenge:       challenge,
	}
}

func HumanWebAuthNBeginLoginEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webAuthNAdded := &HumanWebAuthNBeginLoginEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, webAuthNAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-rMb8x", "unable to unmarshal human webAuthN begin login")
	}
	return webAuthNAdded, nil
}

type HumanWebAuthNCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`

	//TODO: Handle Auth Req??
	//*AuthRequest
}

func (e *HumanWebAuthNCheckSucceededEvent) CheckPrevious() bool {
	return true
}

func (e *HumanWebAuthNCheckSucceededEvent) Data() interface{} {
	return e
}

func NewHumanWebAuthNU2FCheckSucceededEvent(ctx context.Context) *HumanWebAuthNCheckSucceededEvent {
	return &HumanWebAuthNCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanU2FTokenCheckSucceededType,
		),
	}
}

func NewHumanWebAuthNPasswordlessCheckSucceededEvent(ctx context.Context) *HumanWebAuthNCheckSucceededEvent {
	return &HumanWebAuthNCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanPasswordlessTokenCheckSucceededType,
		),
	}
}

func HumanWebAuthNCheckSucceededEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webAuthNAdded := &HumanWebAuthNCheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, webAuthNAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-2M0fg", "unable to unmarshal human webAuthN check succeeded")
	}
	return webAuthNAdded, nil
}

type HumanWebAuthNCheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`

	//TODO: Handle Auth Req??
	//*AuthRequest
}

func (e *HumanWebAuthNCheckFailedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanWebAuthNCheckFailedEvent) Data() interface{} {
	return e
}

func NewHumanWebAuthNU2FCheckFailedEvent(ctx context.Context) *HumanWebAuthNCheckFailedEvent {
	return &HumanWebAuthNCheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanU2FTokenCheckFailedType,
		),
	}
}

func NewHumanWebAuthNPasswordlessCheckFailedEvent(ctx context.Context) *HumanWebAuthNCheckFailedEvent {
	return &HumanWebAuthNCheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanPasswordlessTokenCheckFailedType,
		),
	}
}

func HumanWebAuthNCheckFailedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webAuthNAdded := &HumanWebAuthNCheckFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, webAuthNAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-O0dse", "unable to unmarshal human webAuthN check failed")
	}
	return webAuthNAdded, nil
}
