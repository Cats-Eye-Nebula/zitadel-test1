package admin

import (
	"github.com/caos/logging"
	user_grpc "github.com/caos/zitadel/internal/api/grpc/user"
	"github.com/caos/zitadel/internal/v2/domain"
	admin_grpc "github.com/caos/zitadel/pkg/grpc/admin"
	"golang.org/x/text/language"
)

func setUpOrgHumanToDomain(human *admin_grpc.SetUpOrgRequest_Human) *domain.Human {
	return &domain.Human{
		Username: human.UserName,
		Profile:  setUpOrgHumanProfileToDomain(human.Profile),
		Address:  setUpOrgHumanAddressToDomain(human.Address),
		Email:    setUpOrgHumanEmailToDomain(human.Email),
		Phone:    setUpOrgHumanPhoneToDomain(human.Phone),
	}
	return nil
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

func setUpOrgHumanAddressToDomain(address *admin_grpc.SetUpOrgRequest_Human_Address) *domain.Address {
	if address == nil {
		return nil
	}
	return &domain.Address{
		Country:       address.Country,
		Locality:      address.Locality,
		PostalCode:    address.PostalCode,
		Region:        address.Region,
		StreetAddress: address.StreetAddress,
	}
}

func setUpOrgHumanEmailToDomain(email *admin_grpc.SetUpOrgRequest_Human_Email) *domain.Email {
	return &domain.Email{
		EmailAddress:    email.Email,
		IsEmailVerified: email.IsEmailVerified,
	}
}

func setUpOrgHumanPhoneToDomain(phone *admin_grpc.SetUpOrgRequest_Human_Phone) *domain.Phone {
	return &domain.Phone{
		PhoneNumber:     phone.Phone,
		IsPhoneVerified: phone.IsPhoneVerified,
	}
}
