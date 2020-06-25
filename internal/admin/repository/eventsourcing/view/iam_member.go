package view

import (
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
)

const (
	iamMemberTable = "admin_api.iam_members"
)

func (v *View) IamMemberByIDs(orgID, userID string) (*model.IamMemberView, error) {
	return view.IamMemberByIDs(v.Db, iamMemberTable, orgID, userID)
}

func (v *View) SearchIamMembers(request *iam_model.IamMemberSearchRequest) ([]*model.IamMemberView, int, error) {
	return view.SearchIamMembers(v.Db, iamMemberTable, request)
}

func (v *View) IamMembersByUserID(userID string) ([]*model.IamMemberView, error) {
	return view.IamMembersByUserID(v.Db, iamMemberTable, userID)
}

func (v *View) PutIamMember(org *model.IamMemberView, sequence uint64) error {
	err := view.PutIamMember(v.Db, iamMemberTable, org)
	if err != nil {
		return err
	}
	return v.ProcessedIamMemberSequence(sequence)
}

func (v *View) DeleteIamMember(orgID, userID string, eventSequence uint64) error {
	err := view.DeleteIamMember(v.Db, iamMemberTable, orgID, userID)
	if err != nil {
		return nil
	}
	return v.ProcessedIamMemberSequence(eventSequence)
}

func (v *View) GetLatestIamMemberSequence() (uint64, error) {
	return v.latestSequence(iamMemberTable)
}

func (v *View) ProcessedIamMemberSequence(eventSequence uint64) error {
	return v.saveCurrentSequence(iamMemberTable, eventSequence)
}

func (v *View) GetLatestIamMemberFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(iamMemberTable, sequence)
}

func (v *View) ProcessedIamMemberFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
