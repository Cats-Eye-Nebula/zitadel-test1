package types

import (
	"html"
	"strings"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/notification/templates"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type DomainClaimedData struct {
	templates.TemplateData
	URL string
}

func SendDomainClaimed(mailhtml string, text *iam_model.MailTextView, user *view_model.NotifyUser, username string, systemDefaults systemdefaults.SystemDefaults, colors *iam_model.LabelPolicyView) error {
	url, err := templates.ParseTemplateText(systemDefaults.Notifications.Endpoints.DomainClaimed, &UrlData{UserID: user.ID})
	if err != nil {
		return err
	}
	var args = map[string]interface{}{
		"FirstName":    user.FirstName,
		"LastName":     user.LastName,
		"Username":     user.LastEmail,
		"TempUsername": username,
		"Domain":       strings.Split(user.LastEmail, "@")[1],
	}

	text.Greeting, err = templates.ParseTemplateText(text.Greeting, args)
	text.Text, err = templates.ParseTemplateText(text.Text, args)
	text.Text = html.UnescapeString(text.Text)

	emailCodeData := &DomainClaimedData{
		TemplateData: templates.TemplateData{
			Title:          text.Title,
			PreHeader:      text.PreHeader,
			Subject:        text.Subject,
			Greeting:       text.Greeting,
			Text:           html.UnescapeString(text.Text),
			Href:           url,
			ButtonText:     text.ButtonText,
			PrimaryColor:   colors.PrimaryColor,
			SecondaryColor: colors.SecondaryColor,
		},
		URL: url,
	}
	template, err := templates.GetParsedTemplate(mailhtml, emailCodeData)
	if err != nil {
		return err
	}
	return generateEmail(user, text.Subject, template, systemDefaults.Notifications, true)
}
