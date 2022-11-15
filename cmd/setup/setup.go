package setup

import (
	"context"
	_ "embed"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/build"
	"github.com/zitadel/zitadel/cmd/key"
	"github.com/zitadel/zitadel/cmd/tls"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/migration"
	"github.com/zitadel/zitadel/internal/query/projection"
)

var (
	//go:embed steps.yaml
	defaultSteps []byte
	stepFiles    []string
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "setup ZITADEL instance",
		Long: `sets up data to start ZITADEL.
Requirements:
- cockroachdb`,
		Run: func(cmd *cobra.Command, args []string) {
			err := tls.ModeFromFlag(cmd)
			logging.OnError(err).Fatal("invalid tlsMode")

			config := MustNewConfig(viper.GetViper())
			steps := MustNewSteps(viper.New())

			masterKey, err := key.MasterKey(cmd)
			logging.OnError(err).Panic("No master key provided")

			Setup(config, steps, masterKey)
		},
	}

	Flags(cmd)

	return cmd
}

func Flags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringArrayVar(&stepFiles, "steps", nil, "paths to step files to overwrite default steps")
	key.AddMasterKeyFlag(cmd)
	tls.AddTLSModeFlag(cmd)
}

func Setup(config *Config, steps *Steps, masterKey string) {
	ctx := context.Background()
	logging.Info("setup started")

	dbClient, err := database.Connect(config.Database, false)
	logging.OnError(err).Fatal("unable to connect to database")

	eventstoreClient, err := eventstore.Start(dbClient)
	logging.OnError(err).Fatal("unable to start eventstore")
	migration.RegisterMappers(eventstoreClient)

	steps.s1ProjectionTable = &ProjectionTable{dbClient: dbClient}
	steps.s2AssetsTable = &AssetTable{dbClient: dbClient}

	steps.FirstInstance.instanceSetup = config.DefaultInstance
	steps.FirstInstance.userEncryptionKey = config.EncryptionKeys.User
	steps.FirstInstance.smtpEncryptionKey = config.EncryptionKeys.SMTP
	steps.FirstInstance.masterKey = masterKey
	steps.FirstInstance.db = dbClient
	steps.FirstInstance.es = eventstoreClient
	steps.FirstInstance.defaults = config.SystemDefaults
	steps.FirstInstance.zitadelRoles = config.InternalAuthZ.RolePermissionMappings
	steps.FirstInstance.externalDomain = config.ExternalDomain
	steps.FirstInstance.externalSecure = config.ExternalSecure
	steps.FirstInstance.externalPort = config.ExternalPort

	steps.s4EventstoreIndexes = &EventstoreIndexes{dbClient: dbClient, dbType: config.Database.Type()}

	err = projection.Create(ctx, dbClient, eventstoreClient, config.Projections, nil, nil)
	logging.OnError(err).Fatal("unable to start projections")

	repeatableSteps := []migration.RepeatableMigration{
		&externalConfigChange{
			es:             eventstoreClient,
			ExternalDomain: config.ExternalDomain,
			ExternalPort:   config.ExternalPort,
			ExternalSecure: config.ExternalSecure,
		},
		&projectionTables{
			es:      eventstoreClient,
			Version: build.Version(),
		},
	}

	err = migration.Migrate(ctx, eventstoreClient, steps.s1ProjectionTable)
	logging.OnError(err).Fatal("unable to migrate step 1")
	err = migration.Migrate(ctx, eventstoreClient, steps.s2AssetsTable)
	logging.OnError(err).Fatal("unable to migrate step 2")
	err = migration.Migrate(ctx, eventstoreClient, steps.FirstInstance)
	logging.OnError(err).Fatal("unable to migrate step 3")
	err = migration.Migrate(ctx, eventstoreClient, steps.s4EventstoreIndexes)
	logging.OnError(err).Fatal("unable to migrate step 4")

	for _, repeatableStep := range repeatableSteps {
		err = migration.Migrate(ctx, eventstoreClient, repeatableStep)
		logging.OnError(err).Fatalf("unable to migrate repeatable step: %s", repeatableStep.String())
	}
}
