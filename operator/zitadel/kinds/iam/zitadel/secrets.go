package zitadel

import (
	"github.com/caos/orbos/pkg/secret"
)

func getSecretsMap(desiredKind *DesiredV0) (
	map[string]*secret.Secret,
	map[string]*secret.Existing,
) {

	var (
		secrets  = map[string]*secret.Secret{}
		existing = map[string]*secret.Existing{}
	)

	if desiredKind.Spec != nil && desiredKind.Spec.Configuration != nil {
		conf := desiredKind.Spec.Configuration
		if conf.Tracing != nil {
			if conf.Tracing.ServiceAccountJSON == nil {
				conf.Tracing.ServiceAccountJSON = &secret.Secret{}
			}
			if conf.Tracing.ExistingServiceAccountJSON == nil {
				conf.Tracing.ExistingServiceAccountJSON = &secret.Existing{}
			}
			sakey := "tracingserviceaccountjson"
			secrets[sakey] = conf.Tracing.ServiceAccountJSON
			existing[sakey] = conf.Tracing.ExistingServiceAccountJSON
		}

		if conf.Secrets != nil {
			if conf.Secrets.Keys == nil {
				conf.Secrets.Keys = &secret.Secret{}
			}
			if conf.Secrets.ExistingKeys == nil {
				conf.Secrets.ExistingKeys = &secret.Existing{}
			}
			keysKey := "keys"
			secrets[keysKey] = conf.Secrets.Keys
			existing[keysKey] = conf.Secrets.ExistingKeys
		}

		if conf.Notifications != nil {

			if conf.Notifications.GoogleChatURL == nil {
				conf.Notifications.GoogleChatURL = &secret.Secret{}
			}
			if conf.Notifications.ExistingGoogleChatURL == nil {
				conf.Notifications.ExistingGoogleChatURL = &secret.Existing{}
			}
			gchatkey := "googlechaturl"
			secrets[gchatkey] = conf.Notifications.GoogleChatURL
			existing[gchatkey] = conf.Notifications.ExistingGoogleChatURL

			if conf.Notifications.Twilio.SID == nil {
				conf.Notifications.Twilio.SID = &secret.Secret{}
			}
			if conf.Notifications.Twilio.ExistingSID == nil {
				conf.Notifications.Twilio.ExistingSID = &secret.Existing{}
			}
			twilKey := "twiliosid"
			secrets[twilKey] = conf.Notifications.Twilio.SID
			existing[twilKey] = conf.Notifications.Twilio.ExistingSID

			if conf.Notifications.Twilio.AuthToken == nil {
				conf.Notifications.Twilio.AuthToken = &secret.Secret{}
			}
			if conf.Notifications.Twilio.ExistingAuthToken == nil {
				conf.Notifications.Twilio.ExistingAuthToken = &secret.Existing{}
			}
			twilOAuthKey := "twilioauthtoken"
			secrets[twilOAuthKey] = conf.Notifications.Twilio.AuthToken
			existing[twilOAuthKey] = conf.Notifications.Twilio.ExistingAuthToken

			if conf.Notifications.Email.AppKey == nil {
				conf.Notifications.Email.AppKey = &secret.Secret{}
			}
			if conf.Notifications.Email.ExistingAppKey == nil {
				conf.Notifications.Email.ExistingAppKey = &secret.Existing{}
			}
			mailKey := "emailappkey"
			secrets[mailKey] = conf.Notifications.Email.AppKey
			existing[mailKey] = conf.Notifications.Email.ExistingAppKey
		}
	}
	return secrets, existing
}
