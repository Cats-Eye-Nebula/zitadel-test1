package types

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	"github.com/zitadel/zitadel/internal/notification/channels/log"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	"github.com/zitadel/zitadel/internal/notification/templates"
	"github.com/zitadel/zitadel/internal/query"
)

type Notify func(
	url string,
	args map[string]interface{},
	messageType string,
	allowUnverifiedNotificationChannel bool,
) error

func SendEmail(
	ctx context.Context,
	mailhtml string,
	translator *i18n.Translator,
	user *query.NotifyUser,
	emailConfig func(ctx context.Context) (*smtp.Config, error),
	getFileSystemProvider func(ctx context.Context) (*fs.Config, error),
	getLogProvider func(ctx context.Context) (*log.Config, error),
	colors *query.LabelPolicy,
	assetsPrefix string,
	triggeringEvent eventstore.Event,
	successMetricName,
	failureMetricName string,
) Notify {
	return func(
		url string,
		args map[string]interface{},
		messageType string,
		allowUnverifiedNotificationChannel bool,
	) error {
		args = mapNotifyUserToArgs(user, args)
		data := GetTemplateData(translator, args, assetsPrefix, url, messageType, user.PreferredLanguage.String(), colors)
		template, err := templates.GetParsedTemplate(mailhtml, data)
		if err != nil {
			return err
		}
		return generateEmail(
			ctx,
			user,
			data.Subject,
			template,
			emailConfig,
			getFileSystemProvider,
			getLogProvider,
			allowUnverifiedNotificationChannel,
			triggeringEvent,
			successMetricName,
			failureMetricName,
		)
	}
}

func SendSMSTwilio(
	ctx context.Context,
	translator *i18n.Translator,
	user *query.NotifyUser,
	twilioConfig func(ctx context.Context) (*twilio.Config, error),
	getFileSystemProvider func(ctx context.Context) (*fs.Config, error),
	getLogProvider func(ctx context.Context) (*log.Config, error),
	colors *query.LabelPolicy,
	assetsPrefix string,
	triggeringEvent eventstore.Event,
	successMetricName,
	failureMetricName string,
) Notify {
	return func(
		url string,
		args map[string]interface{},
		messageType string,
		allowUnverifiedNotificationChannel bool,
	) error {
		args = mapNotifyUserToArgs(user, args)
		data := GetTemplateData(translator, args, assetsPrefix, url, messageType, user.PreferredLanguage.String(), colors)
		return generateSms(
			ctx,
			user,
			data.Text,
			twilioConfig,
			getFileSystemProvider,
			getLogProvider,
			allowUnverifiedNotificationChannel,
			triggeringEvent,
			successMetricName,
			failureMetricName,
		)
	}
}

func SendJSON(
	ctx context.Context,
	webhookConfig webhook.Config,
	getFileSystemProvider func(ctx context.Context) (*fs.Config, error),
	getLogProvider func(ctx context.Context) (*log.Config, error),
	serializable interface{},
	triggeringEvent eventstore.Event,
	successMetricName,
	failureMetricName string,
) Notify {
	return func(_ string, _ map[string]interface{}, _ string, _ bool) error {
		return handleJSON(
			ctx,
			webhookConfig,
			getFileSystemProvider,
			getLogProvider,
			serializable,
			triggeringEvent,
			successMetricName,
			failureMetricName,
		)
	}
}