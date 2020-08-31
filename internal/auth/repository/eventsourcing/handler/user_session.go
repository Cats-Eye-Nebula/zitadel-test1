package handler

import (
	req_model "github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/errors"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing"
	user_events "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type UserSession struct {
	handler
	userEvents *user_events.UserEventstore
}

const (
	userSessionTable = "auth.user_sessions"
)

func (u *UserSession) ViewModel() string {
	return userSessionTable
}

func (u *UserSession) EventQuery() (*models.SearchQuery, error) {
	sequence, err := u.view.GetLatestUserSessionSequence()
	if err != nil {
		return nil, err
	}
	return eventsourcing.UserQuery(sequence.CurrentSequence), nil
}

func (u *UserSession) Reduce(event *models.Event) (err error) {
	var session *view_model.UserSessionView
	switch event.Type {
	case es_model.UserPasswordCheckSucceeded,
		es_model.UserPasswordCheckFailed,
		es_model.MfaOtpCheckSucceeded,
		es_model.MfaOtpCheckFailed,
		es_model.SignedOut,
		es_model.HumanPasswordCheckSucceeded,
		es_model.HumanPasswordCheckFailed,
		es_model.HumanMfaOtpCheckSucceeded,
		es_model.HumanMfaOtpCheckFailed,
		es_model.HumanSignedOut:
		eventData, err := view_model.UserSessionFromEvent(event)
		if err != nil {
			return err
		}
		session, err = u.view.UserSessionByIDs(eventData.UserAgentID, event.AggregateID)
		if err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
			session = &view_model.UserSessionView{
				CreationDate:  event.CreationDate,
				ResourceOwner: event.ResourceOwner,
				UserAgentID:   eventData.UserAgentID,
				UserID:        event.AggregateID,
				State:         int32(req_model.UserSessionStateActive),
			}
		}
		return u.updateSession(session, event)
	case es_model.UserPasswordChanged,
		es_model.MFAOTPRemoved,
		es_model.UserProfileChanged,
		es_model.UserLocked,
		es_model.UserDeactivated,
		es_model.HumanPasswordChanged,
		es_model.HumanMFAOTPRemoved,
		es_model.HumanProfileChanged,
		es_model.DomainClaimed,
		es_model.UserUserNameChanged:
		sessions, err := u.view.UserSessionsByUserID(event.AggregateID)
		if err != nil {
			return err
		}
		if len(sessions) == 0 {
			return u.view.ProcessedUserSessionSequence(event.Sequence)
		}
		for _, session := range sessions {
			session.AppendEvent(event)
			if err := u.fillUserInfo(session, event.AggregateID); err != nil {
				return err
			}
		}
		return u.view.PutUserSessions(sessions, event.Sequence)
	case es_model.UserRemoved:
		return u.view.DeleteUserSessions(event.AggregateID, event.Sequence)
	default:
		return u.view.ProcessedUserSessionSequence(event.Sequence)
	}
}

func (u *UserSession) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-sdfw3s", "id", event.AggregateID).WithError(err).Warn("something went wrong in user session handler")
	return spooler.HandleError(event, err, u.view.GetLatestUserSessionFailedEvent, u.view.ProcessedUserSessionFailedEvent, u.view.ProcessedUserSessionSequence, u.errorCountUntilSkip)
}

func (u *UserSession) updateSession(session *view_model.UserSessionView, event *models.Event) error {
	session.AppendEvent(event)
	if err := u.fillUserInfo(session, event.AggregateID); err != nil {
		return err
	}
	return u.view.PutUserSession(session)
}

func (u *UserSession) fillUserInfo(session *view_model.UserSessionView, id string) error {
	user, err := u.view.UserByID(id)
	if err != nil {
		return err
	}
	session.UserName = user.UserName
	session.LoginName = user.PreferredLoginName
	session.DisplayName = user.DisplayName
	return nil
}
