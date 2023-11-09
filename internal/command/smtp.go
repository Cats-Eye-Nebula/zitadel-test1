package command

import (
	"context"
	"net"
	"strings"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func (c *Commands) AddSMTPConfig(ctx context.Context, instanceID string, config *smtp.Config) (string, *domain.ObjectDetails, error) {
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}

	from := strings.TrimSpace(config.From)
	if from == "" {
		return "", nil, errors.ThrowInvalidArgument(nil, "INST-ASv2d", "Errors.Invalid.Argument")
	}

	replyTo := strings.TrimSpace(config.ReplyToAddress)
	hostAndPort := strings.TrimSpace(config.SMTP.Host)

	if _, _, err := net.SplitHostPort(hostAndPort); err != nil {
		return "", nil, errors.ThrowInvalidArgument(nil, "INST-9JdRe", "Errors.Invalid.Argument")
	}

	var smtpPassword *crypto.CryptoValue
	if config.SMTP.Password != "" {
		smtpPassword, err = crypto.Encrypt([]byte(config.SMTP.Password), c.smtpEncryption)
		if err != nil {
			return "", nil, err
		}
	}

	// TODO @n40lab
	// err = checkSenderAddress(writeModel)
	// if err != nil {
	// 	return nil, err
	// }

	smtpConfigWriteModel, err := c.getSMTPConfig(ctx, instanceID, id)
	if err != nil {
		return "", nil, err
	}

	iamAgg := InstanceAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewSMTPConfigAddedEvent(
		ctx,
		iamAgg,
		id,
		config.Tls,
		config.From,
		config.FromName,
		replyTo,
		hostAndPort,
		config.SMTP.User,
		smtpPassword,
		config.SMTP.ProviderType,
	))

	if err != nil {
		return "", nil, err
	}
	err = AppendAndReduce(smtpConfigWriteModel, pushedEvents...)
	if err != nil {
		return "", nil, err
	}
	return id, writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel), nil
}

func (c *Commands) ChangeSMTPConfig(ctx context.Context, instanceID string, id string, config *smtp.Config) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "SMS-e9jwf", "Errors.IDMissing")
	}

	from := strings.TrimSpace(config.From)
	if from == "" {
		return nil, errors.ThrowInvalidArgument(nil, "INST-ASv2d", "Errors.Invalid.Argument")
	}

	replyTo := strings.TrimSpace(config.ReplyToAddress)
	hostAndPort := strings.TrimSpace(config.SMTP.Host)
	if _, _, err := net.SplitHostPort(hostAndPort); err != nil {
		return nil, errors.ThrowInvalidArgument(nil, "INST-Kv875", "Errors.Invalid.Argument")
	}

	var smtpPassword *crypto.CryptoValue
	var err error
	if config.SMTP.Password != "" {
		smtpPassword, err = crypto.Encrypt([]byte(config.SMTP.Password), c.smtpEncryption)
		if err != nil {
			return nil, err
		}
	}

	smtpConfigWriteModel, err := c.getSMTPConfig(ctx, instanceID, id)
	if err != nil {
		return nil, err
	}

	if !smtpConfigWriteModel.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-7j8gv", "Errors.SMTPConfig.NotFound")
	}

	iamAgg := InstanceAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel)

	changedEvent, hasChanged, err := smtpConfigWriteModel.NewChangedEvent(
		ctx,
		iamAgg,
		id,
		config.Tls,
		from,
		config.FromName,
		replyTo,
		hostAndPort,
		config.SMTP.User,
		smtpPassword,
		config.SMTP.ProviderType,
	)

	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-lh3op", "Errors.NoChangesFound")
	}

	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(smtpConfigWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel), nil
}

func (c *Commands) ChangeSMTPConfigPassword(ctx context.Context, instanceID, id string, password string) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	// TODO @n40lab test
	smtpConfigWriteModel, err := c.getSMTPConfig(ctx, instanceID, id)
	if err != nil {
		return nil, err
	}
	if !smtpConfigWriteModel.State.Exists() {
		return nil, errors.ThrowNotFound(nil, "COMMAND-3n9ls", "Errors.SMTPConfig.NotFound")
	}

	var smtpPassword *crypto.CryptoValue
	if password != "" {
		smtpPassword, err = crypto.Encrypt([]byte(password), c.smtpEncryption)
		if err != nil {
			return nil, err
		}
	}

	pushedEvents, err := c.eventstore.Push(ctx, instance.NewSMTPConfigPasswordChangedEvent(
		ctx,
		&instanceAgg.Aggregate,
		id,
		smtpPassword))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(smtpConfigWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel), nil
}

func (c *Commands) ActivateSMTPConfig(ctx context.Context, instanceID, id, activatedId string) (*domain.ObjectDetails, error) {
	var pushedEvents []eventstore.Event

	if id == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "SMTP-nm56k", "Errors.IDMissing")
	}

	if len(activatedId) > 0 {
		smtpConfigWriteModel, err := c.getSMTPConfig(ctx, instanceID, activatedId)
		if err != nil {
			return nil, err
		}

		if !smtpConfigWriteModel.State.Exists() {
			return nil, caos_errs.ThrowNotFound(nil, "COMMAND-jg8ir", "Errors.SMTPConfig.NotFound")
		}

		if smtpConfigWriteModel.State == domain.SMTPConfigStateInactive {
			return nil, caos_errs.ThrowNotFound(nil, "COMMAND-eh6kd", "Errors.SMTPConfig.AlreadyDeactivated")
		}
		iamAgg := InstanceAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel)
		pushedEvents, err = c.eventstore.Push(ctx, instance.NewSMTPConfigDeactivatedEvent(
			ctx,
			iamAgg,
			activatedId))
		if err != nil {
			return nil, err
		}
	}

	smtpConfigWriteModel, err := c.getSMTPConfig(ctx, instanceID, id)
	if err != nil {
		return nil, err
	}

	if !smtpConfigWriteModel.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-kg8yr", "Errors.SMTPConfig.NotFound")
	}

	if smtpConfigWriteModel.State == domain.SMTPConfigStateActive {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-ed3lr", "Errors.SMTPConfig.AlreadyActive")
	}

	iamAgg := InstanceAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel)
	pushedEvents, err = c.eventstore.Push(ctx, instance.NewSMTPConfigActivatedEvent(
		ctx,
		iamAgg,
		id))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(smtpConfigWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel), nil
}

func (c *Commands) DeactivateSMTPConfig(ctx context.Context, instanceID, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "SMTP-98ikl", "Errors.IDMissing")
	}
	smtpConfigWriteModel, err := c.getSMTPConfig(ctx, instanceID, id)
	if err != nil {
		return nil, err
	}
	if !smtpConfigWriteModel.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-k39PJ", "Errors.SMTPConfig.NotFound")
	}
	if smtpConfigWriteModel.State == domain.SMTPConfigStateInactive {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-km8g3", "Errors.SMTPConfig.AlreadyDeactivated")
	}

	iamAgg := InstanceAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewSMTPConfigDeactivatedEvent(
		ctx,
		iamAgg,
		id))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(smtpConfigWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel), nil
}

func (c *Commands) RemoveSMTPConfig(ctx context.Context, instanceID, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "SMTP-7f5cv", "Errors.IDMissing")
	}

	smtpConfigWriteModel, err := c.getSMTPConfig(ctx, instanceID, id)
	if err != nil {
		return nil, err
	}
	if !smtpConfigWriteModel.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-kg8rt", "Errors.SMTPConfig.NotFound")
	}

	iamAgg := InstanceAggregateFromWriteModel(&smtpConfigWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewSMTPConfigRemovedEvent(
		ctx,
		iamAgg,
		id))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(smtpConfigWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&smtpConfigWriteModel.WriteModel), nil
}

func checkSenderAddress(writeModel *InstanceSMTPConfigWriteModel) error {
	if !writeModel.smtpSenderAddressMatchesInstanceDomain {
		return nil
	}
	if !writeModel.domainState.Exists() {
		return errors.ThrowInvalidArgument(nil, "INST-83nl8", "Errors.SMTPConfig.SenderAdressNotCustomDomain")
	}
	return nil
}

func (c *Commands) getSMTPConfig(ctx context.Context, instanceID, id string) (_ *InstanceSMTPConfigWriteModel, err error) {
	writeModel := NewIAMSMTPConfigWriteModel(authz.GetInstance(ctx).InstanceID(), id)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	return writeModel, nil
}
