package domain

import (
	"regexp"
	"time"

	"github.com/zitadel/zitadel/internal/errors"

	"github.com/zitadel/zitadel/internal/crypto"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

var (
	EmailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Email struct {
	es_models.ObjectRoot

	EmailAddress    string
	IsEmailVerified bool
}

type EmailCode struct {
	es_models.ObjectRoot

	Code   *crypto.CryptoValue
	Expiry time.Duration
}

func (e *Email) IsValid() error {
	if e == nil || e.EmailAddress == "" {
		return errors.ThrowInvalidArgument(nil, "EMAIL-spblu", "Errors.User.Email.Empty")
	}
	if !EmailRegex.MatchString(e.EmailAddress) {
		return errors.ThrowInvalidArgument(nil, "EMAIL-spblu", "Errors.User.Email.Invalid")
	}
	return nil
}

func NewEmailCode(emailGenerator crypto.Generator) (*EmailCode, error) {
	emailCodeCrypto, _, err := crypto.NewCode(emailGenerator)
	if err != nil {
		return nil, err
	}
	return &EmailCode{
		Code:   emailCodeCrypto,
		Expiry: emailGenerator.Expiry(),
	}, nil
}
