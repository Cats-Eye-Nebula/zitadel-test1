package eventstore

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/v1/sdk"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_view "github.com/caos/zitadel/internal/iam/repository/view"
	key_model "github.com/caos/zitadel/internal/key/model"
	key_view_model "github.com/caos/zitadel/internal/key/repository/view/model"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	proj_model "github.com/caos/zitadel/internal/project/model"
	proj_view "github.com/caos/zitadel/internal/project/repository/view"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_view "github.com/caos/zitadel/internal/user/repository/view"
	usr_es_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type ProjectRepo struct {
	v1.Eventstore
	SearchLimit     uint64
	View            *view.View
	Roles           []string
	IAMID           string
	PrefixAvatarURL string
}

func (repo *ProjectRepo) ProjectMemberByID(ctx context.Context, projectID, userID string) (*proj_model.ProjectMemberView, error) {
	member, err := repo.View.ProjectMemberByIDs(projectID, userID)
	if err != nil {
		return nil, err
	}
	return model.ProjectMemberToModel(member, repo.PrefixAvatarURL), nil
}

func (repo *ProjectRepo) SearchProjectMembers(ctx context.Context, request *proj_model.ProjectMemberSearchRequest) (*proj_model.ProjectMemberSearchResponse, error) {
	err := request.EnsureLimit(repo.SearchLimit)
	if err != nil {
		return nil, err
	}
	sequence, sequenceErr := repo.View.GetLatestProjectMemberSequence()
	logging.Log("EVENT-3dgt6").OnError(sequenceErr).Warn("could not read latest project member sequence")
	members, count, err := repo.View.SearchProjectMembers(request)
	if err != nil {
		return nil, err
	}
	result := &proj_model.ProjectMemberSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      model.ProjectMembersToModel(members, repo.PrefixAvatarURL),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *ProjectRepo) ProjectChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool, retention time.Duration) (*proj_model.ProjectChanges, error) {
	changes, err := repo.getProjectChanges(ctx, id, lastSequence, limit, sortAscending, retention)
	if err != nil {
		return nil, err
	}
	for _, change := range changes.Changes {
		change.ModifierName = change.ModifierId
		change.ModifierLoginName = change.ModifierId
		user, _ := repo.userByID(ctx, change.ModifierId)
		if user != nil {
			change.ModifierLoginName = user.PreferredLoginName
			if user.HumanView != nil {
				change.ModifierName = user.HumanView.DisplayName
				change.ModifierAvatarURL = user.HumanView.AvatarURL
			}
			if user.MachineView != nil {
				change.ModifierName = user.MachineView.Name
			}
		}
	}
	return changes, nil
}

func (repo *ProjectRepo) ApplicationByID(ctx context.Context, projectID, appID string) (*proj_model.ApplicationView, error) {
	app, viewErr := repo.View.ApplicationByID(projectID, appID)
	if viewErr != nil && !caos_errs.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if caos_errs.IsNotFound(viewErr) {
		app = new(model.ApplicationView)
		app.ID = appID
	}

	events, esErr := repo.getProjectEvents(ctx, projectID, app.Sequence)
	if caos_errs.IsNotFound(viewErr) && len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-Fshu8", "Errors.Application.NotFound")
	}

	if esErr != nil {
		logging.Log("EVENT-SLCo9").WithError(viewErr).Debug("error retrieving new events")
		return model.ApplicationViewToModel(app), nil
	}

	viewApp := *app
	for _, event := range events {
		err := app.AppendEventIfMyApp(event)
		if err != nil {
			return model.ApplicationViewToModel(&viewApp), nil
		}
		if app.State == int32(proj_model.AppStateRemoved) {
			return nil, caos_errs.ThrowNotFound(nil, "EVENT-Msl96", "Errors.Application.NotFound")
		}
	}
	return model.ApplicationViewToModel(app), nil
}

func (repo *ProjectRepo) SearchApplications(ctx context.Context, request *proj_model.ApplicationSearchRequest) (*proj_model.ApplicationSearchResponse, error) {
	err := request.EnsureLimit(repo.SearchLimit)
	if err != nil {
		return nil, err
	}
	sequence, sequenceErr := repo.View.GetLatestApplicationSequence()
	logging.Log("EVENT-SKe8s").OnError(sequenceErr).Warn("could not read latest application sequence")
	apps, count, err := repo.View.SearchApplications(request)
	if err != nil {
		return nil, err
	}
	result := &proj_model.ApplicationSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      model.ApplicationViewsToModel(apps),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *ProjectRepo) ApplicationChanges(ctx context.Context, projectID string, appID string, lastSequence uint64, limit uint64, sortAscending bool, retention time.Duration) (*proj_model.ApplicationChanges, error) {
	changes, err := repo.getApplicationChanges(ctx, projectID, appID, lastSequence, limit, sortAscending, retention)
	if err != nil {
		return nil, err
	}
	for _, change := range changes.Changes {
		change.ModifierName = change.ModifierId
		change.ModifierLoginName = change.ModifierId
		user, _ := repo.userByID(ctx, change.ModifierId)
		if user != nil {
			change.ModifierLoginName = user.PreferredLoginName
			if user.HumanView != nil {
				change.ModifierName = user.HumanView.DisplayName
				change.ModifierAvatarURL = user.HumanView.AvatarURL
			}
			if user.MachineView != nil {
				change.ModifierName = user.MachineView.Name
			}
		}
	}
	return changes, nil
}

func (repo *ProjectRepo) SearchClientKeys(ctx context.Context, request *key_model.AuthNKeySearchRequest) (*key_model.AuthNKeySearchResponse, error) {
	err := request.EnsureLimit(repo.SearchLimit)
	if err != nil {
		return nil, err
	}
	sequence, sequenceErr := repo.View.GetLatestAuthNKeySequence()
	logging.Log("EVENT-ADwgw").OnError(sequenceErr).Warn("could not read latest authn key sequence")
	keys, count, err := repo.View.SearchAuthNKeys(request)
	if err != nil {
		return nil, err
	}
	result := &key_model.AuthNKeySearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      key_view_model.AuthNKeysToModel(keys),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *ProjectRepo) GetClientKey(ctx context.Context, projectID, applicationID, keyID string) (*key_model.AuthNKeyView, error) {
	key, viewErr := repo.View.AuthNKeyByIDs(applicationID, keyID)
	if viewErr != nil {
		return nil, viewErr
	}

	events, esErr := repo.getProjectEvents(ctx, projectID, key.Sequence)
	if caos_errs.IsNotFound(viewErr) && len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-SFf2g", "Errors.User.KeyNotFound")
	}

	if esErr != nil {
		logging.Log("EVENT-ADbf2").WithError(viewErr).Debug("error retrieving new events")
		return key_view_model.AuthNKeyToModel(key), nil
	}

	viewKey := *key
	for _, event := range events {
		err := key.AppendEventIfMyClientKey(event)
		if err != nil {
			return key_view_model.AuthNKeyToModel(&viewKey), nil
		}
		if key.State != int32(proj_model.AppStateActive) {
			return nil, caos_errs.ThrowNotFound(nil, "EVENT-Adfg3", "Errors.User.KeyNotFound")
		}
	}
	return key_view_model.AuthNKeyToModel(key), nil
}

func (repo *ProjectRepo) ProjectGrantByID(ctx context.Context, grantID string) (*proj_model.ProjectGrantView, error) {
	grant, err := repo.View.ProjectGrantByID(grantID)
	if err != nil {
		return nil, err
	}
	return model.ProjectGrantToModel(grant), nil
}

func (repo *ProjectRepo) ProjectGrantMemberByID(ctx context.Context, projectID, userID string) (*proj_model.ProjectGrantMemberView, error) {
	member, err := repo.View.ProjectGrantMemberByIDs(projectID, userID)
	if err != nil {
		return nil, err
	}
	return model.ProjectGrantMemberToModel(member, repo.PrefixAvatarURL), nil
}

func (repo *ProjectRepo) SearchProjectGrantRoles(ctx context.Context, projectID, grantID string, request *proj_model.ProjectRoleSearchRequest) (*proj_model.ProjectRoleSearchResponse, error) {
	projectGrant, err := repo.ProjectGrantByID(ctx, grantID)
	if err != nil {
		return nil, err
	}
	err = request.EnsureLimit(repo.SearchLimit)
	if err != nil {
		return nil, err
	}
	request.AppendProjectQuery(projectID)
	request.AppendRoleKeysQuery(projectGrant.GrantedRoleKeys)
	sequence, sequenceErr := repo.View.GetLatestProjectRoleSequence()
	logging.Log("EVENT-3M9fs").OnError(sequenceErr).Warn("could not read latest project role sequence")
	roles, count, err := repo.View.SearchProjectRoles(request)
	if err != nil {
		return nil, err
	}

	result := &proj_model.ProjectRoleSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      model.ProjectRolesToModel(roles),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *ProjectRepo) SearchProjectGrantMembers(ctx context.Context, request *proj_model.ProjectGrantMemberSearchRequest) (*proj_model.ProjectGrantMemberSearchResponse, error) {
	err := request.EnsureLimit(repo.SearchLimit)
	if err != nil {
		return nil, err
	}
	sequence, sequenceErr := repo.View.GetLatestProjectGrantMemberSequence()
	logging.Log("EVENT-Du8sk").OnError(sequenceErr).Warn("could not read latest project grant sequence")
	members, count, err := repo.View.SearchProjectGrantMembers(request)
	if err != nil {
		return nil, err
	}
	result := &proj_model.ProjectGrantMemberSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      model.ProjectGrantMembersToModel(members, repo.PrefixAvatarURL),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *ProjectRepo) GetProjectMemberRoles(ctx context.Context) ([]string, error) {
	iam, err := repo.GetIAMByID(ctx)
	if err != nil {
		return nil, err
	}
	roles := make([]string, 0)
	global := authz.GetCtxData(ctx).OrgID == iam.GlobalOrgID
	for _, roleMap := range repo.Roles {
		if strings.HasPrefix(roleMap, "PROJECT") && !strings.HasPrefix(roleMap, "PROJECT_GRANT") {
			if global && !strings.HasSuffix(roleMap, "GLOBAL") {
				continue
			}
			roles = append(roles, roleMap)
		}
	}
	return roles, nil
}

func (repo *ProjectRepo) GetProjectGrantMemberRoles() []string {
	roles := make([]string, 0)
	for _, roleMap := range repo.Roles {
		if strings.HasPrefix(roleMap, "PROJECT_GRANT") {
			roles = append(roles, roleMap)
		}
	}
	return roles
}

func (repo *ProjectRepo) userByID(ctx context.Context, id string) (*usr_model.UserView, error) {
	user, viewErr := repo.View.UserByID(id)
	if viewErr != nil && !caos_errs.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if caos_errs.IsNotFound(viewErr) {
		user = new(usr_es_model.UserView)
	}
	events, esErr := repo.getUserEvents(ctx, id, user.Sequence)
	if caos_errs.IsNotFound(viewErr) && len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-4n8Fs", "Errors.User.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-PSoc3").WithError(esErr).Debug("error retrieving new events")
		return usr_es_model.UserToModel(user, repo.PrefixAvatarURL), nil
	}
	userCopy := *user
	for _, event := range events {
		if err := userCopy.AppendEvent(event); err != nil {
			return usr_es_model.UserToModel(user, repo.PrefixAvatarURL), nil
		}
	}
	if userCopy.State == int32(usr_model.UserStateDeleted) {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-2m0Fs", "Errors.User.NotFound")
	}
	return usr_es_model.UserToModel(&userCopy, repo.PrefixAvatarURL), nil
}

func (r *ProjectRepo) getUserEvents(ctx context.Context, userID string, sequence uint64) ([]*models.Event, error) {
	query, err := usr_view.UserByIDQuery(userID, sequence)
	if err != nil {
		return nil, err
	}
	return r.Eventstore.FilterEvents(ctx, query)
}

func (repo *ProjectRepo) getProjectChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool, retention time.Duration) (*proj_model.ProjectChanges, error) {
	query := proj_view.ChangesQuery(id, lastSequence, limit, sortAscending, retention)

	events, err := repo.Eventstore.FilterEvents(context.Background(), query)
	if err != nil {
		logging.Log("EVENT-ZRffs").WithError(err).Warn("eventstore unavailable")
		return nil, caos_errs.ThrowInternal(err, "EVENT-328b1", "Errors.Internal")
	}
	if len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-FpQqK", "Errors.Changes.NotFound")
	}

	changes := make([]*proj_model.ProjectChange, len(events))

	for i, event := range events {
		creationDate, err := ptypes.TimestampProto(event.CreationDate)
		logging.Log("EVENT-qxIR7").OnError(err).Debug("unable to parse timestamp")
		change := &proj_model.ProjectChange{
			ChangeDate: creationDate,
			EventType:  event.Type.String(),
			ModifierId: event.EditorUser,
			Sequence:   event.Sequence,
		}

		if event.Data != nil {
			var data interface{}
			if strings.Contains(change.EventType, "application") {
				data = new(proj_model.Application)
			} else {
				data = new(proj_model.Project)
			}
			err = json.Unmarshal(event.Data, data)
			logging.Log("EVENT-NCkpN").OnError(err).Debug("unable to unmarshal data")
			change.Data = data
		}

		changes[i] = change
		if lastSequence < event.Sequence {
			lastSequence = event.Sequence
		}
	}

	return &proj_model.ProjectChanges{
		Changes:      changes,
		LastSequence: lastSequence,
	}, nil
}

func (repo *ProjectRepo) getProjectEvents(ctx context.Context, id string, sequence uint64) ([]*models.Event, error) {
	query, err := proj_view.ProjectByIDQuery(id, sequence)
	if err != nil {
		return nil, err
	}
	return repo.Eventstore.FilterEvents(ctx, query)
}

func (repo *ProjectRepo) getApplicationChanges(ctx context.Context, projectID string, appID string, lastSequence uint64, limit uint64, sortAscending bool, retention time.Duration) (*proj_model.ApplicationChanges, error) {
	query := proj_view.ChangesQuery(projectID, lastSequence, limit, sortAscending, retention)

	events, err := repo.Eventstore.FilterEvents(ctx, query)
	if err != nil {
		logging.Log("EVENT-ZRffs").WithError(err).Warn("eventstore unavailable")
		return nil, caos_errs.ThrowInternal(err, "EVENT-sw6Ku", "Errors.Internal")
	}
	if len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-9IHLP", "Errors.Changes.NotFound")
	}

	result := make([]*proj_model.ApplicationChange, 0)
	for _, event := range events {
		if !strings.Contains(event.Type.String(), "application") || event.Data == nil {
			continue
		}

		app := new(proj_model.Application)
		err := json.Unmarshal(event.Data, app)
		logging.Log("EVENT-GIiKD").OnError(err).Debug("unable to unmarshal data")
		if app.AppID != appID {
			continue
		}

		creationDate, err := ptypes.TimestampProto(event.CreationDate)
		logging.Log("EVENT-MJzeN").OnError(err).Debug("unable to parse timestamp")

		result = append(result, &proj_model.ApplicationChange{
			ChangeDate: creationDate,
			EventType:  event.Type.String(),
			ModifierId: event.EditorUser,
			Sequence:   event.Sequence,
			Data:       app,
		})
		if lastSequence < event.Sequence {
			lastSequence = event.Sequence
		}
	}

	return &proj_model.ApplicationChanges{
		Changes:      result,
		LastSequence: lastSequence,
	}, nil
}

func (u *ProjectRepo) GetIAMByID(ctx context.Context) (*iam_model.IAM, error) {
	query, err := iam_view.IAMByIDQuery(domain.IAMID, 0)
	if err != nil {
		return nil, err
	}
	iam := &iam_es_model.IAM{
		ObjectRoot: models.ObjectRoot{
			AggregateID: domain.IAMID,
		},
	}
	err = es_sdk.Filter(ctx, u.Eventstore.FilterEvents, iam.AppendEvents, query)
	if err != nil && caos_errs.IsNotFound(err) && iam.Sequence == 0 {
		return nil, err
	}
	return iam_es_model.IAMToModel(iam), nil
}
