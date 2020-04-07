package auth

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/caos/zitadel/internal/errors"
)

const (
	authenticated = "authenticated"
)

func CheckUserAuthorization(ctx context.Context, req interface{}, token, orgID string, verifier TokenVerifier, authConfig *Config, requiredAuthOption Option) (context.Context, error) {
	ctx, err := VerifyTokenAndWriteCtxData(ctx, token, orgID, verifier)
	if err != nil {
		return nil, err
	}

	var perms []string
	//TODO: Remove as soon as authentification is implemented
	if CheckInternal(ctx) {
		return ctx, nil
	}

	if requiredAuthOption.Permission == authenticated {
		return ctx, nil
	}

	ctx, perms, err = getUserMethodPermissions(ctx, verifier, requiredAuthOption.Permission, authConfig)
	if err != nil {
		return nil, err
	}

	err = checkUserPermissions(req, perms, requiredAuthOption)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func checkUserPermissions(req interface{}, userPerms []string, authOpt Option) error {
	if len(userPerms) == 0 {
		return errors.ThrowPermissionDenied(nil, "AUTH-5mWD2", "No matching permissions found")
	}

	if authOpt.CheckParam == "" {
		return nil
	}

	if HasGlobalPermission(userPerms) {
		return nil
	}

	if hasContextPermission(req, authOpt.CheckParam, userPerms) {
		return nil
	}

	return errors.ThrowPermissionDenied(nil, "AUTH-3jknH", "No matching permissions found")
}

func SplitPermission(perm string) (string, string) {
	splittedPerm := strings.Split(perm, ":")
	if len(splittedPerm) == 1 {
		return splittedPerm[0], ""
	}
	return splittedPerm[0], splittedPerm[1]
}

func hasContextPermission(req interface{}, fieldName string, permissions []string) bool {
	for _, perm := range permissions {
		_, ctxID := SplitPermission(perm)
		if checkPermissionContext(req, fieldName, ctxID) {
			return true
		}
	}
	return false
}

func checkPermissionContext(req interface{}, fieldName, roleContextID string) bool {
	field := getFieldFromReq(req, fieldName)
	return field != "" && field == roleContextID
}

func getFieldFromReq(req interface{}, field string) string {
	v := reflect.Indirect(reflect.ValueOf(req)).FieldByName(field)
	if reflect.ValueOf(v).IsZero() {
		return ""
	}
	return fmt.Sprintf("%v", v.Interface())
}

func HasGlobalPermission(perms []string) bool {
	for _, perm := range perms {
		_, ctxID := SplitPermission(perm)
		if ctxID == "" {
			return true
		}
	}
	return false
}

func GetPermissionCtxIDs(perms []string) []string {
	ctxIDs := make([]string, 0)
	for _, perm := range perms {
		_, ctxID := SplitPermission(perm)
		if ctxID != "" {
			ctxIDs = append(ctxIDs, ctxID)
		}
	}
	return ctxIDs
}
