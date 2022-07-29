package initialise

import (
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"
)

var (
	searchUser = "SELECT username FROM [show roles] WHERE username = $1"
	//go:embed sql/01_user.sql
	createUserStmt string
)

func newUser() *cobra.Command {
	return &cobra.Command{
		Use:   "user",
		Short: "initialize only the database user",
		Long: `Sets up the ZITADEL database user.

Prereqesits:
- cockroachdb

The user provided by flags needs priviledge to 
- create the database if it does not exist
- see other users and create a new one if the user does not exist
- grant all rights of the ZITADEL database to the user created if not yet set
`,
		Run: func(cmd *cobra.Command, args []string) {
			config := MustNewConfig(viper.New())

			err := initialise(config.Database, VerifyUser(config.Database.Username(), config.Database.Password()))
			logging.OnError(err).Fatal("unable to init user")
		},
	}
}

func VerifyUser(username, password string) func(*sql.DB) error {
	return func(db *sql.DB) error {
		logging.WithFields("username", username).Info("verify user")
		return verify(db,
			exists(searchUser, username),
			exec(fmt.Sprintf(createUserStmt, username), &sql.NullString{String: password, Valid: password != ""}),
		)
	}
}
