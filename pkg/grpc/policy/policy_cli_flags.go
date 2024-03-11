// Code generated by protoc-gen-cli-client. DO NOT EDIT.

package policy

import (
	cli_client "github.com/adlerhurst/cli-client"
	pflag "github.com/spf13/pflag"
	idp "github.com/zitadel/zitadel/pkg/grpc/idp"
	object "github.com/zitadel/zitadel/pkg/grpc/object"
	os "os"
)

type DomainPolicyFlag struct {
	*DomainPolicy

	changed bool
	set     *pflag.FlagSet

	detailsFlag                                *object.ObjectDetailsFlag
	userLoginMustBeDomainFlag                  *cli_client.BoolParser
	isDefaultFlag                              *cli_client.BoolParser
	validateOrgDomainsFlag                     *cli_client.BoolParser
	smtpSenderAddressMatchesInstanceDomainFlag *cli_client.BoolParser
}

func (x *DomainPolicyFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("DomainPolicy", pflag.ContinueOnError)

	x.userLoginMustBeDomainFlag = cli_client.NewBoolParser(x.set, "user-login-must-be-domain", "")
	x.isDefaultFlag = cli_client.NewBoolParser(x.set, "is-default", "")
	x.validateOrgDomainsFlag = cli_client.NewBoolParser(x.set, "validate-org-domains", "")
	x.smtpSenderAddressMatchesInstanceDomainFlag = cli_client.NewBoolParser(x.set, "smtp-sender-address-matches-instance-domain", "")
	x.detailsFlag = &object.ObjectDetailsFlag{ObjectDetails: new(object.ObjectDetails)}
	x.detailsFlag.AddFlags(x.set)
	parent.AddFlagSet(x.set)
}

func (x *DomainPolicyFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args, "details")

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if flagIdx := flagIndexes.LastByName("details"); flagIdx != nil {
		x.detailsFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if x.detailsFlag.Changed() {
		x.changed = true
		x.DomainPolicy.Details = x.detailsFlag.ObjectDetails
	}

	if x.userLoginMustBeDomainFlag.Changed() {
		x.changed = true
		x.DomainPolicy.UserLoginMustBeDomain = *x.userLoginMustBeDomainFlag.Value
	}
	if x.isDefaultFlag.Changed() {
		x.changed = true
		x.DomainPolicy.IsDefault = *x.isDefaultFlag.Value
	}
	if x.validateOrgDomainsFlag.Changed() {
		x.changed = true
		x.DomainPolicy.ValidateOrgDomains = *x.validateOrgDomainsFlag.Value
	}
	if x.smtpSenderAddressMatchesInstanceDomainFlag.Changed() {
		x.changed = true
		x.DomainPolicy.SmtpSenderAddressMatchesInstanceDomain = *x.smtpSenderAddressMatchesInstanceDomainFlag.Value
	}
}

func (x *DomainPolicyFlag) Changed() bool {
	return x.changed
}

type LabelPolicyFlag struct {
	*LabelPolicy

	changed bool
	set     *pflag.FlagSet

	detailsFlag             *object.ObjectDetailsFlag
	primaryColorFlag        *cli_client.StringParser
	isDefaultFlag           *cli_client.BoolParser
	hideLoginNameSuffixFlag *cli_client.BoolParser
	warnColorFlag           *cli_client.StringParser
	backgroundColorFlag     *cli_client.StringParser
	fontColorFlag           *cli_client.StringParser
	primaryColorDarkFlag    *cli_client.StringParser
	backgroundColorDarkFlag *cli_client.StringParser
	warnColorDarkFlag       *cli_client.StringParser
	fontColorDarkFlag       *cli_client.StringParser
	disableWatermarkFlag    *cli_client.BoolParser
	logoUrlFlag             *cli_client.StringParser
	iconUrlFlag             *cli_client.StringParser
	logoUrlDarkFlag         *cli_client.StringParser
	iconUrlDarkFlag         *cli_client.StringParser
	fontUrlFlag             *cli_client.StringParser
	themeModeFlag           *cli_client.EnumParser[ThemeMode]
}

func (x *LabelPolicyFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("LabelPolicy", pflag.ContinueOnError)

	x.primaryColorFlag = cli_client.NewStringParser(x.set, "primary-color", "")
	x.isDefaultFlag = cli_client.NewBoolParser(x.set, "is-default", "")
	x.hideLoginNameSuffixFlag = cli_client.NewBoolParser(x.set, "hide-login-name-suffix", "")
	x.warnColorFlag = cli_client.NewStringParser(x.set, "warn-color", "")
	x.backgroundColorFlag = cli_client.NewStringParser(x.set, "background-color", "")
	x.fontColorFlag = cli_client.NewStringParser(x.set, "font-color", "")
	x.primaryColorDarkFlag = cli_client.NewStringParser(x.set, "primary-color-dark", "")
	x.backgroundColorDarkFlag = cli_client.NewStringParser(x.set, "background-color-dark", "")
	x.warnColorDarkFlag = cli_client.NewStringParser(x.set, "warn-color-dark", "")
	x.fontColorDarkFlag = cli_client.NewStringParser(x.set, "font-color-dark", "")
	x.disableWatermarkFlag = cli_client.NewBoolParser(x.set, "disable-watermark", "")
	x.logoUrlFlag = cli_client.NewStringParser(x.set, "logo-url", "")
	x.iconUrlFlag = cli_client.NewStringParser(x.set, "icon-url", "")
	x.logoUrlDarkFlag = cli_client.NewStringParser(x.set, "logo-url-dark", "")
	x.iconUrlDarkFlag = cli_client.NewStringParser(x.set, "icon-url-dark", "")
	x.fontUrlFlag = cli_client.NewStringParser(x.set, "font-url", "")
	x.themeModeFlag = cli_client.NewEnumParser[ThemeMode](x.set, "theme-mode", "")
	x.detailsFlag = &object.ObjectDetailsFlag{ObjectDetails: new(object.ObjectDetails)}
	x.detailsFlag.AddFlags(x.set)
	parent.AddFlagSet(x.set)
}

func (x *LabelPolicyFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args, "details")

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if flagIdx := flagIndexes.LastByName("details"); flagIdx != nil {
		x.detailsFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if x.detailsFlag.Changed() {
		x.changed = true
		x.LabelPolicy.Details = x.detailsFlag.ObjectDetails
	}

	if x.primaryColorFlag.Changed() {
		x.changed = true
		x.LabelPolicy.PrimaryColor = *x.primaryColorFlag.Value
	}
	if x.isDefaultFlag.Changed() {
		x.changed = true
		x.LabelPolicy.IsDefault = *x.isDefaultFlag.Value
	}
	if x.hideLoginNameSuffixFlag.Changed() {
		x.changed = true
		x.LabelPolicy.HideLoginNameSuffix = *x.hideLoginNameSuffixFlag.Value
	}
	if x.warnColorFlag.Changed() {
		x.changed = true
		x.LabelPolicy.WarnColor = *x.warnColorFlag.Value
	}
	if x.backgroundColorFlag.Changed() {
		x.changed = true
		x.LabelPolicy.BackgroundColor = *x.backgroundColorFlag.Value
	}
	if x.fontColorFlag.Changed() {
		x.changed = true
		x.LabelPolicy.FontColor = *x.fontColorFlag.Value
	}
	if x.primaryColorDarkFlag.Changed() {
		x.changed = true
		x.LabelPolicy.PrimaryColorDark = *x.primaryColorDarkFlag.Value
	}
	if x.backgroundColorDarkFlag.Changed() {
		x.changed = true
		x.LabelPolicy.BackgroundColorDark = *x.backgroundColorDarkFlag.Value
	}
	if x.warnColorDarkFlag.Changed() {
		x.changed = true
		x.LabelPolicy.WarnColorDark = *x.warnColorDarkFlag.Value
	}
	if x.fontColorDarkFlag.Changed() {
		x.changed = true
		x.LabelPolicy.FontColorDark = *x.fontColorDarkFlag.Value
	}
	if x.disableWatermarkFlag.Changed() {
		x.changed = true
		x.LabelPolicy.DisableWatermark = *x.disableWatermarkFlag.Value
	}
	if x.logoUrlFlag.Changed() {
		x.changed = true
		x.LabelPolicy.LogoUrl = *x.logoUrlFlag.Value
	}
	if x.iconUrlFlag.Changed() {
		x.changed = true
		x.LabelPolicy.IconUrl = *x.iconUrlFlag.Value
	}
	if x.logoUrlDarkFlag.Changed() {
		x.changed = true
		x.LabelPolicy.LogoUrlDark = *x.logoUrlDarkFlag.Value
	}
	if x.iconUrlDarkFlag.Changed() {
		x.changed = true
		x.LabelPolicy.IconUrlDark = *x.iconUrlDarkFlag.Value
	}
	if x.fontUrlFlag.Changed() {
		x.changed = true
		x.LabelPolicy.FontUrl = *x.fontUrlFlag.Value
	}
	if x.themeModeFlag.Changed() {
		x.changed = true
		x.LabelPolicy.ThemeMode = *x.themeModeFlag.Value
	}
}

func (x *LabelPolicyFlag) Changed() bool {
	return x.changed
}

type LockoutPolicyFlag struct {
	*LockoutPolicy

	changed bool
	set     *pflag.FlagSet

	detailsFlag             *object.ObjectDetailsFlag
	maxPasswordAttemptsFlag *cli_client.Uint64Parser
	isDefaultFlag           *cli_client.BoolParser
}

func (x *LockoutPolicyFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("LockoutPolicy", pflag.ContinueOnError)

	x.maxPasswordAttemptsFlag = cli_client.NewUint64Parser(x.set, "max-password-attempts", "")
	x.isDefaultFlag = cli_client.NewBoolParser(x.set, "is-default", "")
	x.detailsFlag = &object.ObjectDetailsFlag{ObjectDetails: new(object.ObjectDetails)}
	x.detailsFlag.AddFlags(x.set)
	parent.AddFlagSet(x.set)
}

func (x *LockoutPolicyFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args, "details")

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if flagIdx := flagIndexes.LastByName("details"); flagIdx != nil {
		x.detailsFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if x.detailsFlag.Changed() {
		x.changed = true
		x.LockoutPolicy.Details = x.detailsFlag.ObjectDetails
	}

	if x.maxPasswordAttemptsFlag.Changed() {
		x.changed = true
		x.LockoutPolicy.MaxPasswordAttempts = *x.maxPasswordAttemptsFlag.Value
	}
	if x.isDefaultFlag.Changed() {
		x.changed = true
		x.LockoutPolicy.IsDefault = *x.isDefaultFlag.Value
	}
}

func (x *LockoutPolicyFlag) Changed() bool {
	return x.changed
}

type LoginPolicyFlag struct {
	*LoginPolicy

	changed bool
	set     *pflag.FlagSet

	detailsFlag                    *object.ObjectDetailsFlag
	allowUsernamePasswordFlag      *cli_client.BoolParser
	allowRegisterFlag              *cli_client.BoolParser
	allowExternalIdpFlag           *cli_client.BoolParser
	forceMfaFlag                   *cli_client.BoolParser
	passwordlessTypeFlag           *cli_client.EnumParser[PasswordlessType]
	isDefaultFlag                  *cli_client.BoolParser
	hidePasswordResetFlag          *cli_client.BoolParser
	ignoreUnknownUsernamesFlag     *cli_client.BoolParser
	defaultRedirectUriFlag         *cli_client.StringParser
	passwordCheckLifetimeFlag      *cli_client.DurationParser
	externalLoginCheckLifetimeFlag *cli_client.DurationParser
	mfaInitSkipLifetimeFlag        *cli_client.DurationParser
	secondFactorCheckLifetimeFlag  *cli_client.DurationParser
	multiFactorCheckLifetimeFlag   *cli_client.DurationParser
	secondFactorsFlag              *cli_client.EnumSliceParser[SecondFactorType]
	multiFactorsFlag               *cli_client.EnumSliceParser[MultiFactorType]
	idpsFlag                       []*idp.IDPLoginPolicyLinkFlag
	allowDomainDiscoveryFlag       *cli_client.BoolParser
	disableLoginWithEmailFlag      *cli_client.BoolParser
	disableLoginWithPhoneFlag      *cli_client.BoolParser
	forceMfaLocalOnlyFlag          *cli_client.BoolParser
}

func (x *LoginPolicyFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("LoginPolicy", pflag.ContinueOnError)

	x.allowUsernamePasswordFlag = cli_client.NewBoolParser(x.set, "allow-username-password", "")
	x.allowRegisterFlag = cli_client.NewBoolParser(x.set, "allow-register", "")
	x.allowExternalIdpFlag = cli_client.NewBoolParser(x.set, "allow-external-idp", "")
	x.forceMfaFlag = cli_client.NewBoolParser(x.set, "force-mfa", "")
	x.passwordlessTypeFlag = cli_client.NewEnumParser[PasswordlessType](x.set, "passwordless-type", "")
	x.isDefaultFlag = cli_client.NewBoolParser(x.set, "is-default", "")
	x.hidePasswordResetFlag = cli_client.NewBoolParser(x.set, "hide-password-reset", "")
	x.ignoreUnknownUsernamesFlag = cli_client.NewBoolParser(x.set, "ignore-unknown-usernames", "")
	x.defaultRedirectUriFlag = cli_client.NewStringParser(x.set, "default-redirect-uri", "")
	x.passwordCheckLifetimeFlag = cli_client.NewDurationParser(x.set, "password-check-lifetime", "")
	x.externalLoginCheckLifetimeFlag = cli_client.NewDurationParser(x.set, "external-login-check-lifetime", "")
	x.mfaInitSkipLifetimeFlag = cli_client.NewDurationParser(x.set, "mfa-init-skip-lifetime", "")
	x.secondFactorCheckLifetimeFlag = cli_client.NewDurationParser(x.set, "second-factor-check-lifetime", "")
	x.multiFactorCheckLifetimeFlag = cli_client.NewDurationParser(x.set, "multi-factor-check-lifetime", "")
	x.secondFactorsFlag = cli_client.NewEnumSliceParser[SecondFactorType](x.set, "second-factors", "")
	x.multiFactorsFlag = cli_client.NewEnumSliceParser[MultiFactorType](x.set, "multi-factors", "")
	x.idpsFlag = []*idp.IDPLoginPolicyLinkFlag{}
	x.allowDomainDiscoveryFlag = cli_client.NewBoolParser(x.set, "allow-domain-discovery", "")
	x.disableLoginWithEmailFlag = cli_client.NewBoolParser(x.set, "disable-login-with-email", "")
	x.disableLoginWithPhoneFlag = cli_client.NewBoolParser(x.set, "disable-login-with-phone", "")
	x.forceMfaLocalOnlyFlag = cli_client.NewBoolParser(x.set, "force-mfa-local-only", "")
	x.detailsFlag = &object.ObjectDetailsFlag{ObjectDetails: new(object.ObjectDetails)}
	x.detailsFlag.AddFlags(x.set)
	parent.AddFlagSet(x.set)
}

func (x *LoginPolicyFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args, "details", "idps")

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if flagIdx := flagIndexes.LastByName("details"); flagIdx != nil {
		x.detailsFlag.ParseFlags(x.set, flagIdx.Args)
	}

	for _, flagIdx := range flagIndexes.ByName("idps") {
		x.idpsFlag = append(x.idpsFlag, &idp.IDPLoginPolicyLinkFlag{IDPLoginPolicyLink: new(idp.IDPLoginPolicyLink)})
		x.idpsFlag[len(x.idpsFlag)-1].AddFlags(x.set)
		x.idpsFlag[len(x.idpsFlag)-1].ParseFlags(x.set, flagIdx.Args)
	}

	if x.detailsFlag.Changed() {
		x.changed = true
		x.LoginPolicy.Details = x.detailsFlag.ObjectDetails
	}

	if x.allowUsernamePasswordFlag.Changed() {
		x.changed = true
		x.LoginPolicy.AllowUsernamePassword = *x.allowUsernamePasswordFlag.Value
	}
	if x.allowRegisterFlag.Changed() {
		x.changed = true
		x.LoginPolicy.AllowRegister = *x.allowRegisterFlag.Value
	}
	if x.allowExternalIdpFlag.Changed() {
		x.changed = true
		x.LoginPolicy.AllowExternalIdp = *x.allowExternalIdpFlag.Value
	}
	if x.forceMfaFlag.Changed() {
		x.changed = true
		x.LoginPolicy.ForceMfa = *x.forceMfaFlag.Value
	}
	if x.passwordlessTypeFlag.Changed() {
		x.changed = true
		x.LoginPolicy.PasswordlessType = *x.passwordlessTypeFlag.Value
	}
	if x.isDefaultFlag.Changed() {
		x.changed = true
		x.LoginPolicy.IsDefault = *x.isDefaultFlag.Value
	}
	if x.hidePasswordResetFlag.Changed() {
		x.changed = true
		x.LoginPolicy.HidePasswordReset = *x.hidePasswordResetFlag.Value
	}
	if x.ignoreUnknownUsernamesFlag.Changed() {
		x.changed = true
		x.LoginPolicy.IgnoreUnknownUsernames = *x.ignoreUnknownUsernamesFlag.Value
	}
	if x.defaultRedirectUriFlag.Changed() {
		x.changed = true
		x.LoginPolicy.DefaultRedirectUri = *x.defaultRedirectUriFlag.Value
	}
	if x.passwordCheckLifetimeFlag.Changed() {
		x.changed = true
		x.LoginPolicy.PasswordCheckLifetime = x.passwordCheckLifetimeFlag.Value
	}
	if x.externalLoginCheckLifetimeFlag.Changed() {
		x.changed = true
		x.LoginPolicy.ExternalLoginCheckLifetime = x.externalLoginCheckLifetimeFlag.Value
	}
	if x.mfaInitSkipLifetimeFlag.Changed() {
		x.changed = true
		x.LoginPolicy.MfaInitSkipLifetime = x.mfaInitSkipLifetimeFlag.Value
	}
	if x.secondFactorCheckLifetimeFlag.Changed() {
		x.changed = true
		x.LoginPolicy.SecondFactorCheckLifetime = x.secondFactorCheckLifetimeFlag.Value
	}
	if x.multiFactorCheckLifetimeFlag.Changed() {
		x.changed = true
		x.LoginPolicy.MultiFactorCheckLifetime = x.multiFactorCheckLifetimeFlag.Value
	}
	if x.secondFactorsFlag.Changed() {
		x.changed = true
		x.LoginPolicy.SecondFactors = *x.secondFactorsFlag.Value
	}
	if x.multiFactorsFlag.Changed() {
		x.changed = true
		x.LoginPolicy.MultiFactors = *x.multiFactorsFlag.Value
	}
	if len(x.idpsFlag) > 0 {
		x.changed = true
		x.Idps = make([]*idp.IDPLoginPolicyLink, len(x.idpsFlag))
		for i, value := range x.idpsFlag {
			x.LoginPolicy.Idps[i] = value.IDPLoginPolicyLink
		}
	}

	if x.allowDomainDiscoveryFlag.Changed() {
		x.changed = true
		x.LoginPolicy.AllowDomainDiscovery = *x.allowDomainDiscoveryFlag.Value
	}
	if x.disableLoginWithEmailFlag.Changed() {
		x.changed = true
		x.LoginPolicy.DisableLoginWithEmail = *x.disableLoginWithEmailFlag.Value
	}
	if x.disableLoginWithPhoneFlag.Changed() {
		x.changed = true
		x.LoginPolicy.DisableLoginWithPhone = *x.disableLoginWithPhoneFlag.Value
	}
	if x.forceMfaLocalOnlyFlag.Changed() {
		x.changed = true
		x.LoginPolicy.ForceMfaLocalOnly = *x.forceMfaLocalOnlyFlag.Value
	}
}

func (x *LoginPolicyFlag) Changed() bool {
	return x.changed
}

type NotificationPolicyFlag struct {
	*NotificationPolicy

	changed bool
	set     *pflag.FlagSet

	detailsFlag        *object.ObjectDetailsFlag
	isDefaultFlag      *cli_client.BoolParser
	passwordChangeFlag *cli_client.BoolParser
}

func (x *NotificationPolicyFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("NotificationPolicy", pflag.ContinueOnError)

	x.isDefaultFlag = cli_client.NewBoolParser(x.set, "is-default", "")
	x.passwordChangeFlag = cli_client.NewBoolParser(x.set, "password-change", "")
	x.detailsFlag = &object.ObjectDetailsFlag{ObjectDetails: new(object.ObjectDetails)}
	x.detailsFlag.AddFlags(x.set)
	parent.AddFlagSet(x.set)
}

func (x *NotificationPolicyFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args, "details")

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if flagIdx := flagIndexes.LastByName("details"); flagIdx != nil {
		x.detailsFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if x.detailsFlag.Changed() {
		x.changed = true
		x.NotificationPolicy.Details = x.detailsFlag.ObjectDetails
	}

	if x.isDefaultFlag.Changed() {
		x.changed = true
		x.NotificationPolicy.IsDefault = *x.isDefaultFlag.Value
	}
	if x.passwordChangeFlag.Changed() {
		x.changed = true
		x.NotificationPolicy.PasswordChange = *x.passwordChangeFlag.Value
	}
}

func (x *NotificationPolicyFlag) Changed() bool {
	return x.changed
}

type OrgIAMPolicyFlag struct {
	*OrgIAMPolicy

	changed bool
	set     *pflag.FlagSet

	detailsFlag               *object.ObjectDetailsFlag
	userLoginMustBeDomainFlag *cli_client.BoolParser
	isDefaultFlag             *cli_client.BoolParser
}

func (x *OrgIAMPolicyFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("OrgIAMPolicy", pflag.ContinueOnError)

	x.userLoginMustBeDomainFlag = cli_client.NewBoolParser(x.set, "user-login-must-be-domain", "")
	x.isDefaultFlag = cli_client.NewBoolParser(x.set, "is-default", "")
	x.detailsFlag = &object.ObjectDetailsFlag{ObjectDetails: new(object.ObjectDetails)}
	x.detailsFlag.AddFlags(x.set)
	parent.AddFlagSet(x.set)
}

func (x *OrgIAMPolicyFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args, "details")

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if flagIdx := flagIndexes.LastByName("details"); flagIdx != nil {
		x.detailsFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if x.detailsFlag.Changed() {
		x.changed = true
		x.OrgIAMPolicy.Details = x.detailsFlag.ObjectDetails
	}

	if x.userLoginMustBeDomainFlag.Changed() {
		x.changed = true
		x.OrgIAMPolicy.UserLoginMustBeDomain = *x.userLoginMustBeDomainFlag.Value
	}
	if x.isDefaultFlag.Changed() {
		x.changed = true
		x.OrgIAMPolicy.IsDefault = *x.isDefaultFlag.Value
	}
}

func (x *OrgIAMPolicyFlag) Changed() bool {
	return x.changed
}

type PasswordAgePolicyFlag struct {
	*PasswordAgePolicy

	changed bool
	set     *pflag.FlagSet

	detailsFlag        *object.ObjectDetailsFlag
	maxAgeDaysFlag     *cli_client.Uint64Parser
	expireWarnDaysFlag *cli_client.Uint64Parser
	isDefaultFlag      *cli_client.BoolParser
}

func (x *PasswordAgePolicyFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("PasswordAgePolicy", pflag.ContinueOnError)

	x.maxAgeDaysFlag = cli_client.NewUint64Parser(x.set, "max-age-days", "")
	x.expireWarnDaysFlag = cli_client.NewUint64Parser(x.set, "expire-warn-days", "")
	x.isDefaultFlag = cli_client.NewBoolParser(x.set, "is-default", "")
	x.detailsFlag = &object.ObjectDetailsFlag{ObjectDetails: new(object.ObjectDetails)}
	x.detailsFlag.AddFlags(x.set)
	parent.AddFlagSet(x.set)
}

func (x *PasswordAgePolicyFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args, "details")

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if flagIdx := flagIndexes.LastByName("details"); flagIdx != nil {
		x.detailsFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if x.detailsFlag.Changed() {
		x.changed = true
		x.PasswordAgePolicy.Details = x.detailsFlag.ObjectDetails
	}

	if x.maxAgeDaysFlag.Changed() {
		x.changed = true
		x.PasswordAgePolicy.MaxAgeDays = *x.maxAgeDaysFlag.Value
	}
	if x.expireWarnDaysFlag.Changed() {
		x.changed = true
		x.PasswordAgePolicy.ExpireWarnDays = *x.expireWarnDaysFlag.Value
	}
	if x.isDefaultFlag.Changed() {
		x.changed = true
		x.PasswordAgePolicy.IsDefault = *x.isDefaultFlag.Value
	}
}

func (x *PasswordAgePolicyFlag) Changed() bool {
	return x.changed
}

type PasswordComplexityPolicyFlag struct {
	*PasswordComplexityPolicy

	changed bool
	set     *pflag.FlagSet

	detailsFlag      *object.ObjectDetailsFlag
	minLengthFlag    *cli_client.Uint64Parser
	hasUppercaseFlag *cli_client.BoolParser
	hasLowercaseFlag *cli_client.BoolParser
	hasNumberFlag    *cli_client.BoolParser
	hasSymbolFlag    *cli_client.BoolParser
	isDefaultFlag    *cli_client.BoolParser
}

func (x *PasswordComplexityPolicyFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("PasswordComplexityPolicy", pflag.ContinueOnError)

	x.minLengthFlag = cli_client.NewUint64Parser(x.set, "min-length", "")
	x.hasUppercaseFlag = cli_client.NewBoolParser(x.set, "has-uppercase", "")
	x.hasLowercaseFlag = cli_client.NewBoolParser(x.set, "has-lowercase", "")
	x.hasNumberFlag = cli_client.NewBoolParser(x.set, "has-number", "")
	x.hasSymbolFlag = cli_client.NewBoolParser(x.set, "has-symbol", "")
	x.isDefaultFlag = cli_client.NewBoolParser(x.set, "is-default", "")
	x.detailsFlag = &object.ObjectDetailsFlag{ObjectDetails: new(object.ObjectDetails)}
	x.detailsFlag.AddFlags(x.set)
	parent.AddFlagSet(x.set)
}

func (x *PasswordComplexityPolicyFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args, "details")

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if flagIdx := flagIndexes.LastByName("details"); flagIdx != nil {
		x.detailsFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if x.detailsFlag.Changed() {
		x.changed = true
		x.PasswordComplexityPolicy.Details = x.detailsFlag.ObjectDetails
	}

	if x.minLengthFlag.Changed() {
		x.changed = true
		x.PasswordComplexityPolicy.MinLength = *x.minLengthFlag.Value
	}
	if x.hasUppercaseFlag.Changed() {
		x.changed = true
		x.PasswordComplexityPolicy.HasUppercase = *x.hasUppercaseFlag.Value
	}
	if x.hasLowercaseFlag.Changed() {
		x.changed = true
		x.PasswordComplexityPolicy.HasLowercase = *x.hasLowercaseFlag.Value
	}
	if x.hasNumberFlag.Changed() {
		x.changed = true
		x.PasswordComplexityPolicy.HasNumber = *x.hasNumberFlag.Value
	}
	if x.hasSymbolFlag.Changed() {
		x.changed = true
		x.PasswordComplexityPolicy.HasSymbol = *x.hasSymbolFlag.Value
	}
	if x.isDefaultFlag.Changed() {
		x.changed = true
		x.PasswordComplexityPolicy.IsDefault = *x.isDefaultFlag.Value
	}
}

func (x *PasswordComplexityPolicyFlag) Changed() bool {
	return x.changed
}

type PrivacyPolicyFlag struct {
	*PrivacyPolicy

	changed bool
	set     *pflag.FlagSet

	detailsFlag      *object.ObjectDetailsFlag
	tosLinkFlag      *cli_client.StringParser
	privacyLinkFlag  *cli_client.StringParser
	isDefaultFlag    *cli_client.BoolParser
	helpLinkFlag     *cli_client.StringParser
	supportEmailFlag *cli_client.StringParser
}

func (x *PrivacyPolicyFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("PrivacyPolicy", pflag.ContinueOnError)

	x.tosLinkFlag = cli_client.NewStringParser(x.set, "tos-link", "")
	x.privacyLinkFlag = cli_client.NewStringParser(x.set, "privacy-link", "")
	x.isDefaultFlag = cli_client.NewBoolParser(x.set, "is-default", "")
	x.helpLinkFlag = cli_client.NewStringParser(x.set, "help-link", "")
	x.supportEmailFlag = cli_client.NewStringParser(x.set, "support-email", "")
	x.detailsFlag = &object.ObjectDetailsFlag{ObjectDetails: new(object.ObjectDetails)}
	x.detailsFlag.AddFlags(x.set)
	parent.AddFlagSet(x.set)
}

func (x *PrivacyPolicyFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args, "details")

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if flagIdx := flagIndexes.LastByName("details"); flagIdx != nil {
		x.detailsFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if x.detailsFlag.Changed() {
		x.changed = true
		x.PrivacyPolicy.Details = x.detailsFlag.ObjectDetails
	}

	if x.tosLinkFlag.Changed() {
		x.changed = true
		x.PrivacyPolicy.TosLink = *x.tosLinkFlag.Value
	}
	if x.privacyLinkFlag.Changed() {
		x.changed = true
		x.PrivacyPolicy.PrivacyLink = *x.privacyLinkFlag.Value
	}
	if x.isDefaultFlag.Changed() {
		x.changed = true
		x.PrivacyPolicy.IsDefault = *x.isDefaultFlag.Value
	}
	if x.helpLinkFlag.Changed() {
		x.changed = true
		x.PrivacyPolicy.HelpLink = *x.helpLinkFlag.Value
	}
	if x.supportEmailFlag.Changed() {
		x.changed = true
		x.PrivacyPolicy.SupportEmail = *x.supportEmailFlag.Value
	}
}

func (x *PrivacyPolicyFlag) Changed() bool {
	return x.changed
}
