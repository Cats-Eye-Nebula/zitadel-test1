// Code generated by protoc-gen-cli-client. DO NOT EDIT.

package metadata

import (
	cli_client "github.com/adlerhurst/cli-client"
	pflag "github.com/spf13/pflag"
	object "github.com/zitadel/zitadel/pkg/grpc/object"
	os "os"
)

type MetadataFlag struct {
	*Metadata

	changed bool
	set     *pflag.FlagSet

	detailsFlag *object.ObjectDetailsFlag
	keyFlag     *cli_client.StringParser
	valueFlag   *cli_client.BytesParser
}

func (x *MetadataFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("Metadata", pflag.ContinueOnError)

	x.keyFlag = cli_client.NewStringParser(x.set, "key", "")
	x.valueFlag = cli_client.NewBytesParser(x.set, "value", "")
	x.detailsFlag = &object.ObjectDetailsFlag{ObjectDetails: new(object.ObjectDetails)}
	x.detailsFlag.AddFlags(x.set)
	parent.AddFlagSet(x.set)
}

func (x *MetadataFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
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
		x.Metadata.Details = x.detailsFlag.ObjectDetails
	}

	if x.keyFlag.Changed() {
		x.changed = true
		x.Metadata.Key = *x.keyFlag.Value
	}
	if x.valueFlag.Changed() {
		x.changed = true
		x.Metadata.Value = *x.valueFlag.Value
	}
}

func (x *MetadataFlag) Changed() bool {
	return x.changed
}

type MetadataKeyQueryFlag struct {
	*MetadataKeyQuery

	changed bool
	set     *pflag.FlagSet

	keyFlag    *cli_client.StringParser
	methodFlag *cli_client.EnumParser[object.TextQueryMethod]
}

func (x *MetadataKeyQueryFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("MetadataKeyQuery", pflag.ContinueOnError)

	x.keyFlag = cli_client.NewStringParser(x.set, "key", "")
	x.methodFlag = cli_client.NewEnumParser[object.TextQueryMethod](x.set, "method", "")
	parent.AddFlagSet(x.set)
}

func (x *MetadataKeyQueryFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args)

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if x.keyFlag.Changed() {
		x.changed = true
		x.MetadataKeyQuery.Key = *x.keyFlag.Value
	}
	if x.methodFlag.Changed() {
		x.changed = true
		x.MetadataKeyQuery.Method = *x.methodFlag.Value
	}
}

func (x *MetadataKeyQueryFlag) Changed() bool {
	return x.changed
}

type MetadataQueryFlag struct {
	*MetadataQuery

	changed bool
	set     *pflag.FlagSet

	keyQueryFlag *MetadataKeyQueryFlag
}

func (x *MetadataQueryFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("MetadataQuery", pflag.ContinueOnError)

	x.keyQueryFlag = &MetadataKeyQueryFlag{MetadataKeyQuery: new(MetadataKeyQuery)}
	x.keyQueryFlag.AddFlags(x.set)
	parent.AddFlagSet(x.set)
}

func (x *MetadataQueryFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args, "key-query")

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if flagIdx := flagIndexes.LastByName("key-query"); flagIdx != nil {
		x.keyQueryFlag.ParseFlags(x.set, flagIdx.Args)
	}

	switch cli_client.FieldIndexes(args, "key-query").Last().Flag {
	case "key-query":
		if x.keyQueryFlag.Changed() {
			x.changed = true
			x.MetadataQuery.Query = &MetadataQuery_KeyQuery{KeyQuery: x.keyQueryFlag.MetadataKeyQuery}
		}
	}
}

func (x *MetadataQueryFlag) Changed() bool {
	return x.changed
}
