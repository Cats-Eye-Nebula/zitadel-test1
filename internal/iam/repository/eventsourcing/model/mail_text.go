package model

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type MailText struct {
	es_models.ObjectRoot
	State        int32 `json:"-"`
	MailTextType string
	Language     string
	Title        string
	PreHeader    string
	Subject      string
	Greeting     string
	Text         string
	ButtonText   string
}

func GetMailText(mailTexts []*MailText, mailTextType string, language string) (int, *MailText) {
	for i, m := range mailTexts {
		if m.MailTextType == mailTextType && m.Language == language {
			return i, m
		}
	}
	return -1, nil
}

func MailTextsToModel(mailTexts []*MailText) []*iam_model.MailText {
	convertedMailTexts := make([]*iam_model.MailText, len(mailTexts))
	for i, m := range mailTexts {
		convertedMailTexts[i] = MailTextToModel(m)
	}
	return convertedMailTexts
}

func MailTextToModel(mailText *MailText) *iam_model.MailText {
	return &iam_model.MailText{
		ObjectRoot:   mailText.ObjectRoot,
		State:        iam_model.PolicyState(mailText.State),
		MailTextType: mailText.MailTextType,
		Language:     mailText.Language,
		Title:        mailText.Title,
		PreHeader:    mailText.PreHeader,
		Subject:      mailText.Subject,
		Greeting:     mailText.Greeting,
		Text:         mailText.Text,
		ButtonText:   mailText.ButtonText,
	}
}

func MailTextsFromModel(mailTexts []*iam_model.MailText) []*MailText {
	convertedMailTexts := make([]*MailText, len(mailTexts))
	for i, m := range mailTexts {
		convertedMailTexts[i] = MailTextFromModel(m)
	}
	return convertedMailTexts
}

func MailTextFromModel(mailText *iam_model.MailText) *MailText {
	return &MailText{
		ObjectRoot:   mailText.ObjectRoot,
		State:        int32(mailText.State),
		MailTextType: mailText.MailTextType,
		Language:     mailText.Language,
		Title:        mailText.Title,
		PreHeader:    mailText.PreHeader,
		Subject:      mailText.Subject,
		Greeting:     mailText.Greeting,
		Text:         mailText.Text,
		ButtonText:   mailText.ButtonText,
	}
}

func (p *MailText) Changes(changed *MailText) map[string]interface{} {
	changes := make(map[string]interface{}, 8)

	changes["mailTextType"] = changed.MailTextType

	changes["language"] = changed.Language

	if changed.Title != p.Title {
		changes["title"] = changed.Title
	}

	if changed.PreHeader != p.PreHeader {
		changes["preHeader"] = changed.PreHeader
	}

	if changed.Subject != p.Subject {
		changes["subject"] = changed.Subject
	}

	if changed.Greeting != p.Greeting {
		changes["greeting"] = changed.Greeting
	}

	if changed.Text != p.Text {
		changes["text"] = changed.Text
	}

	if changed.ButtonText != p.ButtonText {
		changes["buttonText"] = changed.ButtonText
	}

	return changes
}

func (p *MailText) SetDataLabel(event *es_models.Event) error {
	err := json.Unmarshal(event.Data, p)
	if err != nil {
		return errors.ThrowInternal(err, "MODEL-3FUV5", "unable to unmarshal data")
	}
	return nil
}
