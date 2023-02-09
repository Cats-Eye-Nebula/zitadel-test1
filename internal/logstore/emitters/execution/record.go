package execution

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/zitadel/zitadel/internal/logstore"
)

var _ logstore.LogRecord = (*Record)(nil)

type Record struct {
	LogDate    time.Time
	Took       time.Duration
	Message    string
	LogLevel   logrus.Level
	InstanceID string
	ProjectID  string
	ActionID   string
	Metadata   map[string]interface{}
}

func (e Record) Normalize() logstore.LogRecord {
	return &e
}
