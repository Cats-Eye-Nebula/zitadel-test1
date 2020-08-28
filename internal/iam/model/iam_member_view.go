package model

import (
	"time"

	"github.com/caos/zitadel/internal/model"
)

type IAMMemberView struct {
	UserID       string
	IAMID        string
	UserName     string
	Email        string
	FirstName    string
	LastName     string
	DisplayName  string
	Roles        []string
	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
}

type IAMMemberSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn IAMMemberSearchKey
	Asc           bool
	Queries       []*IAMMemberSearchQuery
}

type IAMMemberSearchKey int32

const (
	IAMMemberSearchKeyUnspecified IAMMemberSearchKey = iota
	IAMMemberSearchKeyUserName
	IAMMemberSearchKeyEmail
	IAMMemberSearchKeyFirstName
	IAMMemberSearchKeyLastName
	IAMMemberSearchKeyIamID
	IAMMemberSearchKeyUserID
)

type IAMMemberSearchQuery struct {
	Key    IAMMemberSearchKey
	Method model.SearchMethod
	Value  interface{}
}

type IAMMemberSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*IAMMemberView
	Sequence    uint64
	Timestamp   time.Time
}

func (r *IAMMemberSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}
