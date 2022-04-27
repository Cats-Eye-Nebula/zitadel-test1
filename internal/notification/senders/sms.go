package senders

import (
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
)

func SMSChannels(config systemdefaults.Notifications) (channels.NotificationChannel, error) {

	debug, err := debugChannels(config)
	if err != nil {
		return nil, err
	}

	if !config.DebugMode {
		return chainChannels(debug, twilio.InitTwilioChannel(config.Providers.Twilio)), nil
	}

	return debug, nil
}
