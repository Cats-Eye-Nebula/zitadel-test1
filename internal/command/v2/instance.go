package command

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/ui/console"
	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/project"
	"github.com/caos/zitadel/internal/repository/user"
)

const (
	zitadelProjectName    = "ZITADEL"
	mgmtAppName           = "Management-API"
	adminAppName          = "Admin-API"
	authAppName           = "Auth-API"
	consoleAppName        = "Console"
	consoleRedirectPath   = console.HandlerPrefix + "/auth/callback"
	consolePostLogoutPath = console.HandlerPrefix + "/signedout"
)

type InstanceSetup struct {
	Org      OrgSetup
	Zitadel  ZitadelConfig
	Features struct {
		TierName                 string
		TierDescription          string
		Retention                time.Duration
		State                    domain.FeaturesState
		StateDescription         string
		LoginPolicyFactors       bool
		LoginPolicyIDP           bool
		LoginPolicyPasswordless  bool
		LoginPolicyRegistration  bool
		LoginPolicyUsernameLogin bool
		LoginPolicyPasswordReset bool
		PasswordComplexityPolicy bool
		LabelPolicyPrivateLabel  bool
		LabelPolicyWatermark     bool
		CustomDomain             bool
		PrivacyPolicy            bool
		MetadataUser             bool
		CustomTextMessage        bool
		CustomTextLogin          bool
		LockoutPolicy            bool
		ActionsAllowed           domain.ActionsAllowed
		MaxActions               int
	}
	PasswordComplexityPolicy struct {
		MinLength    uint64
		HasLowercase bool
		HasUppercase bool
		HasNumber    bool
		HasSymbol    bool
	}
	PasswordAgePolicy struct {
		ExpireWarnDays uint64
		MaxAgeDays     uint64
	}
	DomainPolicy struct {
		UserLoginMustBeDomain bool
	}
	LoginPolicy struct {
		AllowUsernamePassword      bool
		AllowRegister              bool
		AllowExternalIDP           bool
		ForceMFA                   bool
		HidePasswordReset          bool
		PasswordlessType           domain.PasswordlessType
		PasswordCheckLifetime      time.Duration
		ExternalLoginCheckLifetime time.Duration
		MfaInitSkipLifetime        time.Duration
		SecondFactorCheckLifetime  time.Duration
		MultiFactorCheckLifetime   time.Duration
	}
	PrivacyPolicy struct {
		TOSLink     string
		PrivacyLink string
		HelpLink    string
	}
	LabelPolicy struct {
		PrimaryColor        string
		BackgroundColor     string
		WarnColor           string
		FontColor           string
		PrimaryColorDark    string
		BackgroundColorDark string
		WarnColorDark       string
		FontColorDark       string
		HideLoginNameSuffix bool
		ErrorMsgPopup       bool
		DisableWatermark    bool
	}
	LockoutPolicy struct {
		MaxAttempts              uint64
		ShouldShowLockoutFailure bool
	}
	EmailTemplate []byte
	MessageTexts  []*domain.CustomMessageText
}

type ZitadelConfig struct {
	IsDevMode bool
	BaseURL   string

	projectID       string
	mgmtID          string
	mgmtClientID    string
	adminID         string
	adminClientID   string
	authID          string
	authClientID    string
	consoleID       string
	consoleClientID string
}

func (s *InstanceSetup) generateIDs() (err error) {
	s.Zitadel.projectID, err = id.SonyFlakeGenerator.Next()
	if err != nil {
		return err
	}

	s.Zitadel.mgmtID, err = id.SonyFlakeGenerator.Next()
	if err != nil {
		return err
	}
	s.Zitadel.mgmtClientID, err = domain.NewClientID(id.SonyFlakeGenerator, zitadelProjectName)
	if err != nil {
		return err
	}

	s.Zitadel.adminID, err = id.SonyFlakeGenerator.Next()
	if err != nil {
		return err
	}
	s.Zitadel.adminClientID, err = domain.NewClientID(id.SonyFlakeGenerator, zitadelProjectName)
	if err != nil {
		return err
	}

	s.Zitadel.authID, err = id.SonyFlakeGenerator.Next()
	if err != nil {
		return err
	}
	s.Zitadel.authClientID, err = domain.NewClientID(id.SonyFlakeGenerator, zitadelProjectName)
	if err != nil {
		return err
	}

	s.Zitadel.consoleID, err = id.SonyFlakeGenerator.Next()
	if err != nil {
		return err
	}
	s.Zitadel.consoleClientID, err = domain.NewClientID(id.SonyFlakeGenerator, zitadelProjectName)
	if err != nil {
		return err
	}
	return nil
}

func (command *Command) SetUpInstance(ctx context.Context, setup *InstanceSetup) (*domain.ObjectDetails, error) {
	// TODO
	// instanceID, err := id.SonyFlakeGenerator.Next()
	// if err != nil {
	// 	return nil, err
	// }
	ctx = authz.SetCtxData(authz.WithInstanceID(ctx, "system"), authz.CtxData{OrgID: domain.IAMID, ResourceOwner: domain.IAMID})

	orgID, err := id.SonyFlakeGenerator.Next()
	if err != nil {
		return nil, err
	}

	userID, err := id.SonyFlakeGenerator.Next()
	if err != nil {
		return nil, err
	}

	if err = setup.generateIDs(); err != nil {
		return nil, err
	}

	setup.Org.Human.PasswordChangeRequired = true

	instanceAgg := instance.NewAggregate()
	orgAgg := org.NewAggregate(orgID, orgID)
	userAgg := user.NewAggregate(userID, orgID)
	projectAgg := project.NewAggregate(setup.Zitadel.projectID, orgID)

	validations := []preparation.Validation{
		SetDefaultFeatures(
			instanceAgg,
			setup.Features.TierName,
			setup.Features.TierDescription,
			setup.Features.State,
			setup.Features.StateDescription,
			setup.Features.Retention,
			setup.Features.LoginPolicyFactors,
			setup.Features.LoginPolicyIDP,
			setup.Features.LoginPolicyPasswordless,
			setup.Features.LoginPolicyRegistration,
			setup.Features.LoginPolicyUsernameLogin,
			setup.Features.LoginPolicyPasswordReset,
			setup.Features.PasswordComplexityPolicy,
			setup.Features.LabelPolicyPrivateLabel,
			setup.Features.LabelPolicyWatermark,
			setup.Features.CustomDomain,
			setup.Features.PrivacyPolicy,
			setup.Features.MetadataUser,
			setup.Features.CustomTextMessage,
			setup.Features.CustomTextLogin,
			setup.Features.LockoutPolicy,
			setup.Features.ActionsAllowed,
			setup.Features.MaxActions,
		),
		AddPasswordComplexityPolicy(
			instanceAgg,
			setup.PasswordComplexityPolicy.MinLength,
			setup.PasswordComplexityPolicy.HasLowercase,
			setup.PasswordComplexityPolicy.HasUppercase,
			setup.PasswordComplexityPolicy.HasNumber,
			setup.PasswordComplexityPolicy.HasSymbol,
		),
		AddPasswordAgePolicy(
			instanceAgg,
			setup.PasswordAgePolicy.ExpireWarnDays,
			setup.PasswordAgePolicy.MaxAgeDays,
		),
		AddDefaultDomainPolicy(
			instanceAgg,
			setup.DomainPolicy.UserLoginMustBeDomain,
		),
		AddDefaultLoginPolicy(
			instanceAgg,
			setup.LoginPolicy.AllowUsernamePassword,
			setup.LoginPolicy.AllowRegister,
			setup.LoginPolicy.AllowExternalIDP,
			setup.LoginPolicy.ForceMFA,
			setup.LoginPolicy.HidePasswordReset,
			setup.LoginPolicy.PasswordlessType,
			setup.LoginPolicy.PasswordCheckLifetime,
			setup.LoginPolicy.ExternalLoginCheckLifetime,
			setup.LoginPolicy.MfaInitSkipLifetime,
			setup.LoginPolicy.SecondFactorCheckLifetime,
			setup.LoginPolicy.MultiFactorCheckLifetime,
		),
		AddSecondFactorToDefaultLoginPolicy(instanceAgg, domain.SecondFactorTypeOTP),
		AddSecondFactorToDefaultLoginPolicy(instanceAgg, domain.SecondFactorTypeU2F),
		AddMultiFactorToDefaultLoginPolicy(instanceAgg, domain.MultiFactorTypeU2FWithPIN),

		AddPrivacyPolicy(instanceAgg, setup.PrivacyPolicy.TOSLink, setup.PrivacyPolicy.PrivacyLink, setup.PrivacyPolicy.HelpLink),
		AddDefaultLockoutPolicy(instanceAgg, setup.LockoutPolicy.MaxAttempts, setup.LockoutPolicy.ShouldShowLockoutFailure),

		AddDefaultLabelPolicy(
			instanceAgg,
			setup.LabelPolicy.PrimaryColor,
			setup.LabelPolicy.BackgroundColor,
			setup.LabelPolicy.WarnColor,
			setup.LabelPolicy.FontColor,
			setup.LabelPolicy.PrimaryColorDark,
			setup.LabelPolicy.BackgroundColorDark,
			setup.LabelPolicy.WarnColorDark,
			setup.LabelPolicy.FontColorDark,
			setup.LabelPolicy.HideLoginNameSuffix,
			setup.LabelPolicy.ErrorMsgPopup,
			setup.LabelPolicy.DisableWatermark,
		),
		ActivateDefaultLabelPolicy(instanceAgg),

		AddEmailTemplate(instanceAgg, setup.EmailTemplate),
	}

	for _, msg := range setup.MessageTexts {
		validations = append(validations, SetInstanceCustomTexts(instanceAgg, msg))
	}

	validations = append(validations,
		AddOrg(orgAgg, setup.Org.Name, command.iamDomain),
		AddHumanCommand(userAgg, &setup.Org.Human, command.userPasswordAlg),
		AddOrgMember(orgAgg, userID, domain.RoleOrgOwner),
		AddInstanceMember(instanceAgg, userID, domain.RoleIAMOwner),

		AddProject(projectAgg, zitadelProjectName, userID, false, false, false, domain.PrivateLabelingSettingUnspecified),
		SetIAMProject(instanceAgg, projectAgg.ID),

		AddAPIApp(
			*projectAgg,
			setup.Zitadel.mgmtID,
			mgmtAppName,
			setup.Zitadel.mgmtClientID,
			nil,
			domain.APIAuthMethodTypePrivateKeyJWT,
		),

		AddAPIApp(
			*projectAgg,
			setup.Zitadel.adminID,
			adminAppName,
			setup.Zitadel.adminClientID,
			nil,
			domain.APIAuthMethodTypePrivateKeyJWT,
		),

		AddAPIApp(
			*projectAgg,
			setup.Zitadel.authID,
			authAppName,
			setup.Zitadel.authClientID,
			nil,
			domain.APIAuthMethodTypePrivateKeyJWT,
		),

		AddOIDCApp(
			*projectAgg,
			domain.OIDCVersionV1,
			setup.Zitadel.consoleID,
			consoleAppName,
			setup.Zitadel.consoleClientID,
			nil,
			[]string{setup.Zitadel.BaseURL + consoleRedirectPath},
			[]domain.OIDCResponseType{domain.OIDCResponseTypeCode},
			[]domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
			domain.OIDCApplicationTypeUserAgent,
			domain.OIDCAuthMethodTypeNone,
			[]string{setup.Zitadel.BaseURL + consolePostLogoutPath},
			setup.Zitadel.IsDevMode,
			domain.OIDCTokenTypeBearer,
			false,
			false,
			false,
			0,
			nil,
		),
		SetIAMConsoleID(instanceAgg, setup.Zitadel.consoleClientID),
	)

	cmds, err := preparation.PrepareCommands(ctx, command.es.Filter, validations...)
	if err != nil {
		return nil, err
	}

	events, err := command.es.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: orgID,
	}, nil
}

//SetIAMProject defines the command to set the id of the IAM project onto the instance
func SetIAMProject(a *instance.Aggregate, projectID string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return []eventstore.Command{
				instance.NewIAMProjectSetEvent(ctx, &a.Aggregate, projectID),
			}, nil
		}, nil
	}
}

//SetIAMConsoleID defines the command to set the clientID of the Console App onto the instance
func SetIAMConsoleID(a *instance.Aggregate, clientID string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return []eventstore.Command{
				instance.NewIAMConsoleSetEvent(ctx, &a.Aggregate, clientID),
			}, nil
		}, nil
	}
}
