package view

import (
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	iam_model "github.com/zitadel/zitadel/internal/iam/model"
	"github.com/zitadel/zitadel/internal/iam/repository/view/model"
	"github.com/zitadel/zitadel/internal/view/repository"

	"github.com/jinzhu/gorm"
)

func GetStylingByAggregateIDAndState(db *gorm.DB, table, aggregateID, instanceID string, state int32) (*model.LabelPolicyView, error) {
	policy := new(model.LabelPolicyView)
	aggregateIDQuery := &model.LabelPolicySearchQuery{Key: iam_model.LabelPolicySearchKeyAggregateID, Value: aggregateID, Method: domain.SearchMethodEquals}
	stateQuery := &model.LabelPolicySearchQuery{Key: iam_model.LabelPolicySearchKeyState, Value: state, Method: domain.SearchMethodEquals}
	instanceIDQuery := &model.LabelPolicySearchQuery{Key: iam_model.LabelPolicySearchKeyInstanceID, Value: instanceID, Method: domain.SearchMethodEquals}
	query := repository.PrepareGetByQuery(table, aggregateIDQuery, stateQuery, instanceIDQuery)
	err := query(db, policy)
	if caos_errs.IsNotFound(err) {
		return nil, caos_errs.ThrowNotFound(nil, "VIEW-68G11", "Errors.IAM.LabelPolicy.NotExisting")
	}
	return policy, err
}

func PutStyling(db *gorm.DB, table string, policy *model.LabelPolicyView) error {
	save := repository.PrepareSave(table)
	return save(db, policy)
}

func DeleteInstanceStyling(db *gorm.DB, table, instanceID string) error {
	delete := repository.PrepareDeleteByKey(table, model.LabelPolicySearchKey(iam_model.LabelPolicySearchKeyInstanceID), instanceID)
	return delete(db)
}
