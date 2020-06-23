package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type ProjectGrant struct {
	es_models.ObjectRoot

	GrantID      string
	GrantedOrgID string
	State        ProjectGrantState
	RoleKeys     []string
	Members      []*ProjectGrantMember
}

type ProjectGrantIDs struct {
	ProjectID string
	GrantID   string
}

type ProjectGrantState int32

const (
	PROJECTGRANTSTATE_ACTIVE ProjectGrantState = iota
	PROJECTGRANTSTATE_INACTIVE
)

func NewProjectGrant(projectID, grantID string) *ProjectGrant {
	return &ProjectGrant{ObjectRoot: es_models.ObjectRoot{AggregateID: projectID}, GrantID: grantID, State: PROJECTGRANTSTATE_ACTIVE}
}

func (p *ProjectGrant) IsActive() bool {
	return p.State == PROJECTGRANTSTATE_ACTIVE
}

func (p *ProjectGrant) IsValid() bool {
	return p.GrantedOrgID != ""
}

func (p *ProjectGrant) GetMember(userID string) (int, *ProjectGrantMember) {
	for i, m := range p.Members {
		if m.UserID == userID {
			return i, m
		}
	}
	return -1, nil
}

func (p *ProjectGrant) GetRemovedRoles(roleKeys []string) []string {
	removed := make([]string, 0)
	for _, role := range p.RoleKeys {
		if !containsKey(roleKeys, role) {
			removed = append(removed, role)
		}
	}
	return removed
}

func containsKey(roles []string, key string) bool {
	for _, role := range roles {
		if role == key {
			return true
		}
	}
	return false
}
