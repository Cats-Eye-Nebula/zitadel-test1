package auth

import (
	authz_repo "github.com/caos/zitadel/internal/authz/repository/eventsourcing"
)

type TokenVerifier struct {
	adminID   string
	authZRepo *authz_repo.EsRepository
}

func Start(authZRepo *authz_repo.EsRepository) (v *TokenVerifier) {
	return &TokenVerifier{authZRepo: authZRepo}
}

//func (v *TokenVerifier) VerifyAccessToken(ctx context.Context, token string) (string, string, string, error) {
//	userID, clientID, agentID, err := v.authZRepo.VerifyAccessToken(ctx, token, adminName, v.adminID)
//	if clientID != "" {
//		v.adminID = clientID
//	}
//	return userID, clientID, agentID, err
//}
//
//func (v *TokenVerifier) ResolveGrant(ctx context.Context) (*authz.Grant, error) {
//	return v.authZRepo.ResolveGrants(ctx)
//}
//
//func (v *TokenVerifier) GetProjectIDByClientID(ctx context.Context, clientID string) (string, error) {
//	return v.authZRepo.ProjectIDByClientID(ctx, clientID)
//}
