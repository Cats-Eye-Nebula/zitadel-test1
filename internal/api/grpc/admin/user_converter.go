package admin

import (
	"github.com/caos/logging"
	user_grpc "github.com/zitadel/zitadel/internal/api/grpc/user"
	"github.com/zitadel/zitadel/internal/domain"
	admin_grpc "github.com/zitadel/zitadel/pkg/grpc/admin"
	"golang.org/x/text/language"
)

func setUpOrgHumanToDomain(human *admin_grpc.SetUpOrgRequest_Human) *domain.Human {
	return &domain.Human{
		Username: human.UserName,
		Profile:  setUpOrgHumanProfileToDomain(human.Profile),
		Email:    setUpOrgHumanEmailToDomain(human.Email),
		Phone:    setUpOrgHumanPhoneToDomain(human.Phone),
		Password: setUpOrgHumanPasswordToDomain(human.Password),
	}
}

func setUpOrgHumanProfileToDomain(profile *admin_grpc.SetUpOrgRequest_Human_Profile) *domain.Profile {
	var lang language.Tag
	lang, err := language.Parse(profile.PreferredLanguage)
	logging.Log("ADMIN-tiMWs").OnError(err).Debug("unable to parse language")

	return &domain.Profile{
		FirstName:         profile.FirstName,
		LastName:          profile.LastName,
		NickName:          profile.NickName,
		DisplayName:       profile.DisplayName,
		PreferredLanguage: lang,
		Gender:            user_grpc.GenderToDomain(profile.Gender),
	}
}

func setUpOrgHumanEmailToDomain(email *admin_grpc.SetUpOrgRequest_Human_Email) *domain.Email {
	return &domain.Email{
		EmailAddress:    email.Email,
		IsEmailVerified: email.IsEmailVerified,
	}
}

func setUpOrgHumanPhoneToDomain(phone *admin_grpc.SetUpOrgRequest_Human_Phone) *domain.Phone {
	if phone == nil {
		return nil
	}
	return &domain.Phone{
		PhoneNumber:     phone.Phone,
		IsPhoneVerified: phone.IsPhoneVerified,
	}
}

func setUpOrgHumanPasswordToDomain(password string) *domain.Password {
	if password == "" {
		return nil
	}
	return domain.NewPassword(password)
}
