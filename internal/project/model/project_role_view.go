package model

import (
	"github.com/caos/zitadel/internal/model"
	"time"
)

type ProjectRoleView struct {
	ResourceOwner string
	OrgID         string
	ProjectID     string
	Key           string
	DisplayName   string
	Group         string
	CreationDate  time.Time
	Sequence      uint64
}

type ProjectRoleSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn ProjectRoleSearchKey
	Asc           bool
	Queries       []*ProjectRoleSearchQuery
}

type ProjectRoleSearchKey int32

const (
	PROJECTROLESEARCHKEY_UNSPECIFIED ProjectRoleSearchKey = iota
	PROJECTROLESEARCHKEY_KEY
	PROJECTROLESEARCHKEY_PROJECTID
	PROJECTROLESEARCHKEY_ORGID
	PROJECTROLESEARCHKEY_RESOURCEOWNER
	PROJECTROLESEARCHKEY_DISPLAY_NAME
)

type ProjectRoleSearchQuery struct {
	Key    ProjectRoleSearchKey
	Method model.SearchMethod
	Value  string
}

type ProjectRoleSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*ProjectRoleView
}

func (r *ProjectRoleSearchRequest) AppendMyOrgQuery(orgID string) {
	r.Queries = append(r.Queries, &ProjectRoleSearchQuery{Key: PROJECTROLESEARCHKEY_ORGID, Method: model.SEARCHMETHOD_EQUALS, Value: orgID})
}

func (r *ProjectRoleSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}
