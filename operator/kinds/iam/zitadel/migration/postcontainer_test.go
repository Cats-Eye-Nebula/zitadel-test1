package migration

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestMigration_DropUserCommand(t *testing.T) {
	user := "test"
	pw := "test"
	file := "test"
	equals := strings.Join([]string{
		"echo -n 'DROP USER IF EXISTS ' > " + file,
		"echo -n ${" + user + "} >> " + file,
		"echo -n ';' >> " + file,
	}, ";")

	cmd := createUserCommand(user, pw, file)
	assert.Equal(t, cmd, equals)
}
