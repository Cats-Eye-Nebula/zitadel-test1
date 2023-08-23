package command

import (
	"context"
	"io"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/repository/session"
)

func (c *Commands) CreateOTPSMSChallengeReturnCode(dst *string) SessionCommand {
	return c.createOTPSMSChallenge(true, dst)
}

func (c *Commands) CreateOTPSMSChallenge() SessionCommand {
	return c.createOTPSMSChallenge(false, nil)
}

func (c *Commands) createOTPSMSChallenge(returnCode bool, dst *string) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) error {
		writeModel := NewHumanOTPSMSWriteModel(cmd.sessionWriteModel.UserID, "")
		if err := cmd.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
			return err
		}
		if !writeModel.OTPAdded() {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-BJ2g3", "Errors.User.MFA.OTP.NotReady")
		}
		code, err := cmd.createCode(ctx, cmd.eventstore.Filter, domain.SecretGeneratorTypeOTPSMS, cmd.otpAlg, //TODO: get from config
			&crypto.GeneratorConfig{
				Length:              7,
				Expiry:              5 * time.Minute,
				IncludeLowerLetters: false,
				IncludeUpperLetters: false,
				IncludeDigits:       true,
				IncludeSymbols:      false,
			})
		if err != nil {
			return err
		}
		if returnCode {
			*dst = code.Plain
		}
		cmd.OTPSMSChallenged(ctx, code.Crypted, code.Expiry, returnCode)
		return nil
	}
}

func (c *Commands) OTPSMSSent(ctx context.Context, sessionID, resourceOwner string) error {
	sessionWriteModel := NewSessionWriteModel(sessionID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, sessionWriteModel)
	if err != nil {
		return err
	}
	if sessionWriteModel.OTPSMSCodeChallenge == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-G3t31", "Errors.User.Code.NotFound")
	}
	return c.pushAppendAndReduce(ctx, sessionWriteModel,
		session.NewOTPSMSSentEvent(ctx, &session.NewAggregate(sessionID, sessionWriteModel.ResourceOwner).Aggregate),
	)
}

func (c *Commands) CreateOTPEmailChallengeURLTemplate(urlTmpl string) (SessionCommand, error) {
	if err := domain.RenderOTPEmailURLTemplate(io.Discard, urlTmpl, "userID", "code"); err != nil {
		return nil, err
	}
	return c.createOTPEmailChallenge(false, urlTmpl, nil), nil
}

func (c *Commands) CreateOTPEmailChallengeReturnCode(dst *string) SessionCommand {
	return c.createOTPEmailChallenge(true, "", dst)
}

func (c *Commands) CreateOTPEmailChallenge() SessionCommand {
	return c.createOTPEmailChallenge(false, "", nil)
}

func (c *Commands) createOTPEmailChallenge(returnCode bool, urlTmpl string, dst *string) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) error {
		writeModel := NewHumanOTPEmailWriteModel(cmd.sessionWriteModel.UserID, "")
		if err := cmd.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
			return err
		}
		if !writeModel.OTPAdded() {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-JKLJ3", "Errors.User.MFA.OTP.NotReady")
		}
		code, err := cmd.createCode(ctx, cmd.eventstore.Filter, domain.SecretGeneratorTypeOTPEmail, cmd.otpAlg, //TODO: get from config
			&crypto.GeneratorConfig{
				Length:              7,
				Expiry:              5 * time.Minute,
				IncludeLowerLetters: false,
				IncludeUpperLetters: false,
				IncludeDigits:       true,
				IncludeSymbols:      false,
			})
		if err != nil {
			return err
		}
		if returnCode {
			*dst = code.Plain
		}
		cmd.OTPEmailChallenged(ctx, code.Crypted, code.Expiry, returnCode, urlTmpl)
		return nil
	}
}

func (c *Commands) OTPEmailSent(ctx context.Context, sessionID, resourceOwner string) error {
	sessionWriteModel := NewSessionWriteModel(sessionID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, sessionWriteModel)
	if err != nil {
		return err
	}
	if sessionWriteModel.OTPEmailCodeChallenge == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-SLr02", "Errors.User.Code.NotFound")
	}
	return c.pushAppendAndReduce(ctx, sessionWriteModel,
		session.NewOTPEmailSentEvent(ctx, &session.NewAggregate(sessionID, sessionWriteModel.ResourceOwner).Aggregate),
	)
}

func CheckOTPSMS(code string) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) (err error) {
		if cmd.sessionWriteModel.UserID == "" {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-VDrh3", "Errors.User.UserIDMissing")
		}
		challenge := cmd.sessionWriteModel.OTPSMSCodeChallenge
		if challenge == nil {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-SF3tv", "Errors.User.Code.NotFound")
		}
		err = crypto.VerifyCodeWithAlgorithm(challenge.CreationDate, challenge.Expiry, challenge.Code, code, cmd.otpAlg)
		if err != nil {
			return err
		}
		cmd.OTPSMSChecked(ctx, cmd.now())
		return nil
	}
}

func CheckOTPEmail(code string) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) (err error) {
		if cmd.sessionWriteModel.UserID == "" {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-ejo2w", "Errors.User.UserIDMissing")
		}
		challenge := cmd.sessionWriteModel.OTPEmailCodeChallenge
		if challenge == nil {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-zF3g3", "Errors.User.Code.NotFound")
		}
		err = crypto.VerifyCodeWithAlgorithm(challenge.CreationDate, challenge.Expiry, challenge.Code, code, cmd.otpAlg)
		if err != nil {
			return err
		}
		cmd.OTPEmailChecked(ctx, cmd.now())
		return nil
	}
}
