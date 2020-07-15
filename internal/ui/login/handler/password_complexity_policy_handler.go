package handler

import (
	"github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/errors"
	policy_model "github.com/caos/zitadel/internal/policy/model"
	"net/http"
	"regexp"
	"strconv"
)

const (
	LowerCaseRegex = `[a-z]`
	UpperCaseRegex = `[A-Z]`
	NumberRegex    = `[0-9]`
	SymbolRegex    = `[^A-Za-z0-9]`
)

var (
	hasStringLowerCase = regexp.MustCompile(LowerCaseRegex).MatchString
	hasStringUpperCase = regexp.MustCompile(UpperCaseRegex).MatchString
	hasNumber          = regexp.MustCompile(NumberRegex).MatchString
	hasSymbol          = regexp.MustCompile(SymbolRegex).MatchString
)

func (l *Login) getPasswordComplexityPolicy(r *http.Request, authReq *model.AuthRequest) (*policy_model.PasswordComplexityPolicy, string, error) {
	policy, err := l.authRepo.GetMyPasswordComplexityPolicy(setContext(r.Context(), authReq.UserOrgID))
	if err != nil {
		return nil, err.Error(), err
	}
	description := "<ul id=\"passwordcomplexity\">"
	minLength := l.renderer.LocalizeFromRequest(r, "Password.MinLength", nil)
	description += "<li id=\"minlength\">" + minLength + " " + strconv.Itoa(int(policy.MinLength)) + " <i class=\"material-icons\">cancel</i></li>"
	if policy.HasUppercase {
		uppercase := l.renderer.LocalizeFromRequest(r, "Password.HasUppercase", nil)
		description += "<li id=\"uppercase\" class=\"invalid\">" + uppercase + " <i class=\"material-icons\">cancel</i></li>"
	}
	if policy.HasLowercase {
		lowercase := l.renderer.LocalizeFromRequest(r, "Password.HasLowercase", nil)
		description += "<li id=\"lowercase\">" + lowercase + " <i class=\"material-icons\">cancel</i></li>"
	}
	if policy.HasNumber {
		hasnumber := l.renderer.LocalizeFromRequest(r, "Password.HasNumber", nil)
		description += "<li id=\"number\">" + hasnumber + " <i class=\"material-icons\">cancel</i></li>"
	}
	if policy.HasSymbol {
		hassymbol := l.renderer.LocalizeFromRequest(r, "Password.HasSymbol", nil)
		description += "<li id=\"symbol\">" + hassymbol + " <i class=\"material-icons\">cancel</i></li>"
	}

	description += "</ul>"
	return policy, description, nil
}

func (l *Login) checkPasswordComplexityPolicy(password string, r *http.Request, authReq *model.AuthRequest) error {
	policy, err := l.authRepo.GetMyPasswordComplexityPolicy(setContext(r.Context(), authReq.UserOrgID))
	if err != nil {
		return nil
	}
	if policy.MinLength != 0 && uint64(len(password)) < policy.MinLength {
		return errors.ThrowInvalidArgument(nil, "POLICY-LSo0p", "Errors.User.PasswordComplexityPolicy.MinLength")
	}

	if policy.HasLowercase && !hasStringLowerCase(password) {
		return errors.ThrowInvalidArgument(nil, "POLICY-4Sjsf", "Errors.User.PasswordComplexityPolicy.HasLower")
	}

	if policy.HasUppercase && !hasStringUpperCase(password) {
		return errors.ThrowInvalidArgument(nil, "POLICY-6Sjc9", "Errors.User.PasswordComplexityPolicy.HasUpper")
	}

	if policy.HasNumber && !hasNumber(password) {
		return errors.ThrowInvalidArgument(nil, "POLICY-2Fksi", "Errors.User.PasswordComplexityPolicy.HasNumber")
	}

	if policy.HasSymbol && !hasSymbol(password) {
		return errors.ThrowInvalidArgument(nil, "POLICY-0Js6e", "Errors.User.PasswordComplexityPolicy.HasSymbol")
	}
	return nil
}
