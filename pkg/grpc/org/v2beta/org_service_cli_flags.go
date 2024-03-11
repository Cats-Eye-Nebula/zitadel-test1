// Code generated by protoc-gen-cli-client. DO NOT EDIT.

package org

import (
	cli_client "github.com/adlerhurst/cli-client"
	pflag "github.com/spf13/pflag"
	v2beta1 "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	v2beta "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
	os "os"
)

type AddOrganizationRequestFlag struct {
	*AddOrganizationRequest

	changed bool
	set     *pflag.FlagSet

	nameFlag   *cli_client.StringParser
	adminsFlag []*AddOrganizationRequest_AdminFlag
}

func (x *AddOrganizationRequestFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("AddOrganizationRequest", pflag.ContinueOnError)

	x.nameFlag = cli_client.NewStringParser(x.set, "name", "")
	x.adminsFlag = []*AddOrganizationRequest_AdminFlag{}
	parent.AddFlagSet(x.set)
}

func (x *AddOrganizationRequestFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args, "admins")

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	for _, flagIdx := range flagIndexes.ByName("admins") {
		x.adminsFlag = append(x.adminsFlag, &AddOrganizationRequest_AdminFlag{AddOrganizationRequest_Admin: new(AddOrganizationRequest_Admin)})
		x.adminsFlag[len(x.adminsFlag)-1].AddFlags(x.set)
		x.adminsFlag[len(x.adminsFlag)-1].ParseFlags(x.set, flagIdx.Args)
	}
	if x.nameFlag.Changed() {
		x.changed = true
		x.AddOrganizationRequest.Name = *x.nameFlag.Value
	}
	if len(x.adminsFlag) > 0 {
		x.changed = true
		x.Admins = make([]*AddOrganizationRequest_Admin, len(x.adminsFlag))
		for i, value := range x.adminsFlag {
			x.AddOrganizationRequest.Admins[i] = value.AddOrganizationRequest_Admin
		}
	}

}

func (x *AddOrganizationRequestFlag) Changed() bool {
	return x.changed
}

type AddOrganizationRequest_AdminFlag struct {
	*AddOrganizationRequest_Admin

	changed bool
	set     *pflag.FlagSet

	userIdFlag *cli_client.StringParser
	humanFlag  *v2beta.AddHumanUserRequestFlag
	rolesFlag  *cli_client.StringSliceParser
}

func (x *AddOrganizationRequest_AdminFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("AddOrganizationRequest_Admin", pflag.ContinueOnError)

	x.userIdFlag = cli_client.NewStringParser(x.set, "user-id", "")
	x.rolesFlag = cli_client.NewStringSliceParser(x.set, "roles", "")
	x.humanFlag = &v2beta.AddHumanUserRequestFlag{AddHumanUserRequest: new(v2beta.AddHumanUserRequest)}
	x.humanFlag.AddFlags(x.set)
	parent.AddFlagSet(x.set)
}

func (x *AddOrganizationRequest_AdminFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args, "human")

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if flagIdx := flagIndexes.LastByName("human"); flagIdx != nil {
		x.humanFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if x.rolesFlag.Changed() {
		x.changed = true
		x.AddOrganizationRequest_Admin.Roles = *x.rolesFlag.Value
	}

	switch cli_client.FieldIndexes(args, "user-id", "human").Last().Flag {
	case "user-id":
		if x.userIdFlag.Changed() {
			x.changed = true
			x.AddOrganizationRequest_Admin.UserType = &AddOrganizationRequest_Admin_UserId{UserId: *x.userIdFlag.Value}
		}
	case "human":
		if x.humanFlag.Changed() {
			x.changed = true
			x.AddOrganizationRequest_Admin.UserType = &AddOrganizationRequest_Admin_Human{Human: x.humanFlag.AddHumanUserRequest}
		}
	}
}

func (x *AddOrganizationRequest_AdminFlag) Changed() bool {
	return x.changed
}

type AddOrganizationResponseFlag struct {
	*AddOrganizationResponse

	changed bool
	set     *pflag.FlagSet

	detailsFlag        *v2beta1.DetailsFlag
	organizationIdFlag *cli_client.StringParser
	createdAdminsFlag  []*AddOrganizationResponse_CreatedAdminFlag
}

func (x *AddOrganizationResponseFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("AddOrganizationResponse", pflag.ContinueOnError)

	x.organizationIdFlag = cli_client.NewStringParser(x.set, "organization-id", "")
	x.createdAdminsFlag = []*AddOrganizationResponse_CreatedAdminFlag{}
	x.detailsFlag = &v2beta1.DetailsFlag{Details: new(v2beta1.Details)}
	x.detailsFlag.AddFlags(x.set)
	parent.AddFlagSet(x.set)
}

func (x *AddOrganizationResponseFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args, "details", "created-admins")

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if flagIdx := flagIndexes.LastByName("details"); flagIdx != nil {
		x.detailsFlag.ParseFlags(x.set, flagIdx.Args)
	}

	for _, flagIdx := range flagIndexes.ByName("created-admins") {
		x.createdAdminsFlag = append(x.createdAdminsFlag, &AddOrganizationResponse_CreatedAdminFlag{AddOrganizationResponse_CreatedAdmin: new(AddOrganizationResponse_CreatedAdmin)})
		x.createdAdminsFlag[len(x.createdAdminsFlag)-1].AddFlags(x.set)
		x.createdAdminsFlag[len(x.createdAdminsFlag)-1].ParseFlags(x.set, flagIdx.Args)
	}

	if x.detailsFlag.Changed() {
		x.changed = true
		x.AddOrganizationResponse.Details = x.detailsFlag.Details
	}

	if x.organizationIdFlag.Changed() {
		x.changed = true
		x.AddOrganizationResponse.OrganizationId = *x.organizationIdFlag.Value
	}
	if len(x.createdAdminsFlag) > 0 {
		x.changed = true
		x.CreatedAdmins = make([]*AddOrganizationResponse_CreatedAdmin, len(x.createdAdminsFlag))
		for i, value := range x.createdAdminsFlag {
			x.AddOrganizationResponse.CreatedAdmins[i] = value.AddOrganizationResponse_CreatedAdmin
		}
	}

}

func (x *AddOrganizationResponseFlag) Changed() bool {
	return x.changed
}

type AddOrganizationResponse_CreatedAdminFlag struct {
	*AddOrganizationResponse_CreatedAdmin

	changed bool
	set     *pflag.FlagSet

	userIdFlag    *cli_client.StringParser
	emailCodeFlag *cli_client.StringParser
	phoneCodeFlag *cli_client.StringParser
}

func (x *AddOrganizationResponse_CreatedAdminFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("AddOrganizationResponse_CreatedAdmin", pflag.ContinueOnError)

	x.userIdFlag = cli_client.NewStringParser(x.set, "user-id", "")
	x.emailCodeFlag = cli_client.NewStringParser(x.set, "email-code", "")
	x.phoneCodeFlag = cli_client.NewStringParser(x.set, "phone-code", "")
	parent.AddFlagSet(x.set)
}

func (x *AddOrganizationResponse_CreatedAdminFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args)

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if x.userIdFlag.Changed() {
		x.changed = true
		x.AddOrganizationResponse_CreatedAdmin.UserId = *x.userIdFlag.Value
	}
	if x.emailCodeFlag.Changed() {
		x.changed = true
		x.AddOrganizationResponse_CreatedAdmin.EmailCode = x.emailCodeFlag.Value
	}
	if x.phoneCodeFlag.Changed() {
		x.changed = true
		x.AddOrganizationResponse_CreatedAdmin.PhoneCode = x.phoneCodeFlag.Value
	}
}

func (x *AddOrganizationResponse_CreatedAdminFlag) Changed() bool {
	return x.changed
}
