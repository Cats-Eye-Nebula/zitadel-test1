// Code generated by protoc-gen-cli-client. DO NOT EDIT.

package session

import (
	cli_client "github.com/adlerhurst/cli-client"
	pflag "github.com/spf13/pflag"
	object "github.com/zitadel/zitadel/pkg/grpc/object"
	os "os"
)

type CreationDateQueryFlag struct {
	*CreationDateQuery

	changed bool
	set     *pflag.FlagSet

	creationDateFlag *cli_client.TimestampParser
	methodFlag       *cli_client.EnumParser[object.TimestampQueryMethod]
}

func (x *CreationDateQueryFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("CreationDateQuery", pflag.ContinueOnError)

	x.creationDateFlag = cli_client.NewTimestampParser(x.set, "creation-date", "")
	x.methodFlag = cli_client.NewEnumParser[object.TimestampQueryMethod](x.set, "method", "")
	parent.AddFlagSet(x.set)
}

func (x *CreationDateQueryFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args)

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if x.creationDateFlag.Changed() {
		x.changed = true
		x.CreationDateQuery.CreationDate = x.creationDateFlag.Value
	}
	if x.methodFlag.Changed() {
		x.changed = true
		x.CreationDateQuery.Method = *x.methodFlag.Value
	}
}

func (x *CreationDateQueryFlag) Changed() bool {
	return x.changed
}

type FactorsFlag struct {
	*Factors

	changed bool
	set     *pflag.FlagSet

	userFlag     *UserFactorFlag
	passwordFlag *PasswordFactorFlag
	webAuthNFlag *WebAuthNFactorFlag
	intentFlag   *IntentFactorFlag
	totpFlag     *TOTPFactorFlag
	otpSmsFlag   *OTPFactorFlag
	otpEmailFlag *OTPFactorFlag
}

func (x *FactorsFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("Factors", pflag.ContinueOnError)

	x.userFlag = &UserFactorFlag{UserFactor: new(UserFactor)}
	x.userFlag.AddFlags(x.set)
	x.passwordFlag = &PasswordFactorFlag{PasswordFactor: new(PasswordFactor)}
	x.passwordFlag.AddFlags(x.set)
	x.webAuthNFlag = &WebAuthNFactorFlag{WebAuthNFactor: new(WebAuthNFactor)}
	x.webAuthNFlag.AddFlags(x.set)
	x.intentFlag = &IntentFactorFlag{IntentFactor: new(IntentFactor)}
	x.intentFlag.AddFlags(x.set)
	x.totpFlag = &TOTPFactorFlag{TOTPFactor: new(TOTPFactor)}
	x.totpFlag.AddFlags(x.set)
	x.otpSmsFlag = &OTPFactorFlag{OTPFactor: new(OTPFactor)}
	x.otpSmsFlag.AddFlags(x.set)
	x.otpEmailFlag = &OTPFactorFlag{OTPFactor: new(OTPFactor)}
	x.otpEmailFlag.AddFlags(x.set)
	parent.AddFlagSet(x.set)
}

func (x *FactorsFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args, "user", "password", "web-auth-n", "intent", "totp", "otp-sms", "otp-email")

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if flagIdx := flagIndexes.LastByName("user"); flagIdx != nil {
		x.userFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if flagIdx := flagIndexes.LastByName("password"); flagIdx != nil {
		x.passwordFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if flagIdx := flagIndexes.LastByName("web-auth-n"); flagIdx != nil {
		x.webAuthNFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if flagIdx := flagIndexes.LastByName("intent"); flagIdx != nil {
		x.intentFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if flagIdx := flagIndexes.LastByName("totp"); flagIdx != nil {
		x.totpFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if flagIdx := flagIndexes.LastByName("otp-sms"); flagIdx != nil {
		x.otpSmsFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if flagIdx := flagIndexes.LastByName("otp-email"); flagIdx != nil {
		x.otpEmailFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if x.userFlag.Changed() {
		x.changed = true
		x.Factors.User = x.userFlag.UserFactor
	}

	if x.passwordFlag.Changed() {
		x.changed = true
		x.Factors.Password = x.passwordFlag.PasswordFactor
	}

	if x.webAuthNFlag.Changed() {
		x.changed = true
		x.Factors.WebAuthN = x.webAuthNFlag.WebAuthNFactor
	}

	if x.intentFlag.Changed() {
		x.changed = true
		x.Factors.Intent = x.intentFlag.IntentFactor
	}

	if x.totpFlag.Changed() {
		x.changed = true
		x.Factors.Totp = x.totpFlag.TOTPFactor
	}

	if x.otpSmsFlag.Changed() {
		x.changed = true
		x.Factors.OtpSms = x.otpSmsFlag.OTPFactor
	}

	if x.otpEmailFlag.Changed() {
		x.changed = true
		x.Factors.OtpEmail = x.otpEmailFlag.OTPFactor
	}

}

func (x *FactorsFlag) Changed() bool {
	return x.changed
}

type IDsQueryFlag struct {
	*IDsQuery

	changed bool
	set     *pflag.FlagSet

	idsFlag *cli_client.StringSliceParser
}

func (x *IDsQueryFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("IDsQuery", pflag.ContinueOnError)

	x.idsFlag = cli_client.NewStringSliceParser(x.set, "ids", "")
	parent.AddFlagSet(x.set)
}

func (x *IDsQueryFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args)

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if x.idsFlag.Changed() {
		x.changed = true
		x.IDsQuery.Ids = *x.idsFlag.Value
	}
}

func (x *IDsQueryFlag) Changed() bool {
	return x.changed
}

type IntentFactorFlag struct {
	*IntentFactor

	changed bool
	set     *pflag.FlagSet

	verifiedAtFlag *cli_client.TimestampParser
}

func (x *IntentFactorFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("IntentFactor", pflag.ContinueOnError)

	x.verifiedAtFlag = cli_client.NewTimestampParser(x.set, "verified-at", "")
	parent.AddFlagSet(x.set)
}

func (x *IntentFactorFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args)

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if x.verifiedAtFlag.Changed() {
		x.changed = true
		x.IntentFactor.VerifiedAt = x.verifiedAtFlag.Value
	}
}

func (x *IntentFactorFlag) Changed() bool {
	return x.changed
}

type OTPFactorFlag struct {
	*OTPFactor

	changed bool
	set     *pflag.FlagSet

	verifiedAtFlag *cli_client.TimestampParser
}

func (x *OTPFactorFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("OTPFactor", pflag.ContinueOnError)

	x.verifiedAtFlag = cli_client.NewTimestampParser(x.set, "verified-at", "")
	parent.AddFlagSet(x.set)
}

func (x *OTPFactorFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args)

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if x.verifiedAtFlag.Changed() {
		x.changed = true
		x.OTPFactor.VerifiedAt = x.verifiedAtFlag.Value
	}
}

func (x *OTPFactorFlag) Changed() bool {
	return x.changed
}

type PasswordFactorFlag struct {
	*PasswordFactor

	changed bool
	set     *pflag.FlagSet

	verifiedAtFlag *cli_client.TimestampParser
}

func (x *PasswordFactorFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("PasswordFactor", pflag.ContinueOnError)

	x.verifiedAtFlag = cli_client.NewTimestampParser(x.set, "verified-at", "")
	parent.AddFlagSet(x.set)
}

func (x *PasswordFactorFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args)

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if x.verifiedAtFlag.Changed() {
		x.changed = true
		x.PasswordFactor.VerifiedAt = x.verifiedAtFlag.Value
	}
}

func (x *PasswordFactorFlag) Changed() bool {
	return x.changed
}

type SearchQueryFlag struct {
	*SearchQuery

	changed bool
	set     *pflag.FlagSet

	idsQueryFlag          *IDsQueryFlag
	userIdQueryFlag       *UserIDQueryFlag
	creationDateQueryFlag *CreationDateQueryFlag
}

func (x *SearchQueryFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("SearchQuery", pflag.ContinueOnError)

	x.idsQueryFlag = &IDsQueryFlag{IDsQuery: new(IDsQuery)}
	x.idsQueryFlag.AddFlags(x.set)
	x.userIdQueryFlag = &UserIDQueryFlag{UserIDQuery: new(UserIDQuery)}
	x.userIdQueryFlag.AddFlags(x.set)
	x.creationDateQueryFlag = &CreationDateQueryFlag{CreationDateQuery: new(CreationDateQuery)}
	x.creationDateQueryFlag.AddFlags(x.set)
	parent.AddFlagSet(x.set)
}

func (x *SearchQueryFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args, "ids-query", "user-id-query", "creation-date-query")

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if flagIdx := flagIndexes.LastByName("ids-query"); flagIdx != nil {
		x.idsQueryFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if flagIdx := flagIndexes.LastByName("user-id-query"); flagIdx != nil {
		x.userIdQueryFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if flagIdx := flagIndexes.LastByName("creation-date-query"); flagIdx != nil {
		x.creationDateQueryFlag.ParseFlags(x.set, flagIdx.Args)
	}

	switch cli_client.FieldIndexes(args, "ids-query", "user-id-query", "creation-date-query").Last().Flag {
	case "ids-query":
		if x.idsQueryFlag.Changed() {
			x.changed = true
			x.SearchQuery.Query = &SearchQuery_IdsQuery{IdsQuery: x.idsQueryFlag.IDsQuery}
		}
	case "user-id-query":
		if x.userIdQueryFlag.Changed() {
			x.changed = true
			x.SearchQuery.Query = &SearchQuery_UserIdQuery{UserIdQuery: x.userIdQueryFlag.UserIDQuery}
		}
	case "creation-date-query":
		if x.creationDateQueryFlag.Changed() {
			x.changed = true
			x.SearchQuery.Query = &SearchQuery_CreationDateQuery{CreationDateQuery: x.creationDateQueryFlag.CreationDateQuery}
		}
	}
}

func (x *SearchQueryFlag) Changed() bool {
	return x.changed
}

type SessionFlag struct {
	*Session

	changed bool
	set     *pflag.FlagSet

	idFlag             *cli_client.StringParser
	creationDateFlag   *cli_client.TimestampParser
	changeDateFlag     *cli_client.TimestampParser
	sequenceFlag       *cli_client.Uint64Parser
	factorsFlag        *FactorsFlag
	userAgentFlag      *UserAgentFlag
	expirationDateFlag *cli_client.TimestampParser
}

func (x *SessionFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("Session", pflag.ContinueOnError)

	x.idFlag = cli_client.NewStringParser(x.set, "id", "")
	x.creationDateFlag = cli_client.NewTimestampParser(x.set, "creation-date", "")
	x.changeDateFlag = cli_client.NewTimestampParser(x.set, "change-date", "")
	x.sequenceFlag = cli_client.NewUint64Parser(x.set, "sequence", "")
	x.expirationDateFlag = cli_client.NewTimestampParser(x.set, "expiration-date", "")
	x.factorsFlag = &FactorsFlag{Factors: new(Factors)}
	x.factorsFlag.AddFlags(x.set)
	x.userAgentFlag = &UserAgentFlag{UserAgent: new(UserAgent)}
	x.userAgentFlag.AddFlags(x.set)
	parent.AddFlagSet(x.set)
}

func (x *SessionFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args, "factors", "user-agent")

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if flagIdx := flagIndexes.LastByName("factors"); flagIdx != nil {
		x.factorsFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if flagIdx := flagIndexes.LastByName("user-agent"); flagIdx != nil {
		x.userAgentFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if x.idFlag.Changed() {
		x.changed = true
		x.Session.Id = *x.idFlag.Value
	}
	if x.creationDateFlag.Changed() {
		x.changed = true
		x.Session.CreationDate = x.creationDateFlag.Value
	}
	if x.changeDateFlag.Changed() {
		x.changed = true
		x.Session.ChangeDate = x.changeDateFlag.Value
	}
	if x.sequenceFlag.Changed() {
		x.changed = true
		x.Session.Sequence = *x.sequenceFlag.Value
	}

	if x.factorsFlag.Changed() {
		x.changed = true
		x.Session.Factors = x.factorsFlag.Factors
	}

	if x.userAgentFlag.Changed() {
		x.changed = true
		x.Session.UserAgent = x.userAgentFlag.UserAgent
	}

	if x.expirationDateFlag.Changed() {
		x.changed = true
		x.Session.ExpirationDate = x.expirationDateFlag.Value
	}
}

func (x *SessionFlag) Changed() bool {
	return x.changed
}

type TOTPFactorFlag struct {
	*TOTPFactor

	changed bool
	set     *pflag.FlagSet

	verifiedAtFlag *cli_client.TimestampParser
}

func (x *TOTPFactorFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("TOTPFactor", pflag.ContinueOnError)

	x.verifiedAtFlag = cli_client.NewTimestampParser(x.set, "verified-at", "")
	parent.AddFlagSet(x.set)
}

func (x *TOTPFactorFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args)

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if x.verifiedAtFlag.Changed() {
		x.changed = true
		x.TOTPFactor.VerifiedAt = x.verifiedAtFlag.Value
	}
}

func (x *TOTPFactorFlag) Changed() bool {
	return x.changed
}

type UserAgentFlag struct {
	*UserAgent

	changed bool
	set     *pflag.FlagSet

	fingerprintIdFlag *cli_client.StringParser
	ipFlag            *cli_client.StringParser
	descriptionFlag   *cli_client.StringParser
}

func (x *UserAgentFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("UserAgent", pflag.ContinueOnError)

	x.fingerprintIdFlag = cli_client.NewStringParser(x.set, "fingerprint-id", "")
	x.ipFlag = cli_client.NewStringParser(x.set, "ip", "")
	x.descriptionFlag = cli_client.NewStringParser(x.set, "description", "")
	parent.AddFlagSet(x.set)
}

func (x *UserAgentFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args)

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if x.fingerprintIdFlag.Changed() {
		x.changed = true
		x.UserAgent.FingerprintId = x.fingerprintIdFlag.Value
	}
	if x.ipFlag.Changed() {
		x.changed = true
		x.UserAgent.Ip = x.ipFlag.Value
	}
	if x.descriptionFlag.Changed() {
		x.changed = true
		x.UserAgent.Description = x.descriptionFlag.Value
	}
}

func (x *UserAgentFlag) Changed() bool {
	return x.changed
}

type UserAgent_HeaderValuesFlag struct {
	*UserAgent_HeaderValues

	changed bool
	set     *pflag.FlagSet

	valuesFlag *cli_client.StringSliceParser
}

func (x *UserAgent_HeaderValuesFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("UserAgent_HeaderValues", pflag.ContinueOnError)

	x.valuesFlag = cli_client.NewStringSliceParser(x.set, "values", "")
	parent.AddFlagSet(x.set)
}

func (x *UserAgent_HeaderValuesFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args)

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if x.valuesFlag.Changed() {
		x.changed = true
		x.UserAgent_HeaderValues.Values = *x.valuesFlag.Value
	}
}

func (x *UserAgent_HeaderValuesFlag) Changed() bool {
	return x.changed
}

type UserFactorFlag struct {
	*UserFactor

	changed bool
	set     *pflag.FlagSet

	verifiedAtFlag     *cli_client.TimestampParser
	idFlag             *cli_client.StringParser
	loginNameFlag      *cli_client.StringParser
	displayNameFlag    *cli_client.StringParser
	organisationIdFlag *cli_client.StringParser
	organizationIdFlag *cli_client.StringParser
}

func (x *UserFactorFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("UserFactor", pflag.ContinueOnError)

	x.verifiedAtFlag = cli_client.NewTimestampParser(x.set, "verified-at", "")
	x.idFlag = cli_client.NewStringParser(x.set, "id", "")
	x.loginNameFlag = cli_client.NewStringParser(x.set, "login-name", "")
	x.displayNameFlag = cli_client.NewStringParser(x.set, "display-name", "")
	x.organisationIdFlag = cli_client.NewStringParser(x.set, "organisation-id", "")
	x.organizationIdFlag = cli_client.NewStringParser(x.set, "organization-id", "")
	parent.AddFlagSet(x.set)
}

func (x *UserFactorFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args)

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if x.verifiedAtFlag.Changed() {
		x.changed = true
		x.UserFactor.VerifiedAt = x.verifiedAtFlag.Value
	}
	if x.idFlag.Changed() {
		x.changed = true
		x.UserFactor.Id = *x.idFlag.Value
	}
	if x.loginNameFlag.Changed() {
		x.changed = true
		x.UserFactor.LoginName = *x.loginNameFlag.Value
	}
	if x.displayNameFlag.Changed() {
		x.changed = true
		x.UserFactor.DisplayName = *x.displayNameFlag.Value
	}
	if x.organisationIdFlag.Changed() {
		x.changed = true
		x.UserFactor.OrganisationId = *x.organisationIdFlag.Value
	}
	if x.organizationIdFlag.Changed() {
		x.changed = true
		x.UserFactor.OrganizationId = *x.organizationIdFlag.Value
	}
}

func (x *UserFactorFlag) Changed() bool {
	return x.changed
}

type UserIDQueryFlag struct {
	*UserIDQuery

	changed bool
	set     *pflag.FlagSet

	idFlag *cli_client.StringParser
}

func (x *UserIDQueryFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("UserIDQuery", pflag.ContinueOnError)

	x.idFlag = cli_client.NewStringParser(x.set, "id", "")
	parent.AddFlagSet(x.set)
}

func (x *UserIDQueryFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args)

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if x.idFlag.Changed() {
		x.changed = true
		x.UserIDQuery.Id = *x.idFlag.Value
	}
}

func (x *UserIDQueryFlag) Changed() bool {
	return x.changed
}

type WebAuthNFactorFlag struct {
	*WebAuthNFactor

	changed bool
	set     *pflag.FlagSet

	verifiedAtFlag   *cli_client.TimestampParser
	userVerifiedFlag *cli_client.BoolParser
}

func (x *WebAuthNFactorFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("WebAuthNFactor", pflag.ContinueOnError)

	x.verifiedAtFlag = cli_client.NewTimestampParser(x.set, "verified-at", "")
	x.userVerifiedFlag = cli_client.NewBoolParser(x.set, "user-verified", "")
	parent.AddFlagSet(x.set)
}

func (x *WebAuthNFactorFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args)

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if x.verifiedAtFlag.Changed() {
		x.changed = true
		x.WebAuthNFactor.VerifiedAt = x.verifiedAtFlag.Value
	}
	if x.userVerifiedFlag.Changed() {
		x.changed = true
		x.WebAuthNFactor.UserVerified = *x.userVerifiedFlag.Value
	}
}

func (x *WebAuthNFactorFlag) Changed() bool {
	return x.changed
}
