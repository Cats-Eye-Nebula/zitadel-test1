package crtlcrd

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/zitadel/zitadel/operator/crtlcrd/database"
	"github.com/zitadel/zitadel/operator/crtlcrd/zitadel"
	"github.com/zitadel/zitadel/pkg/databases/db"
)

func Destroy(monitor mntr.Monitor, k8sClient kubernetes.ClientInt, dbConn db.Connection, version string, features ...string) error {
	for _, feature := range features {
		switch feature {
		case Zitadel:
			if err := zitadel.Destroy(monitor, k8sClient, dbConn, version); err != nil {
				return err
			}
		case Database:
			if err := database.Destroy(monitor, k8sClient, version); err != nil {
				return err
			}
		}
	}
	return nil
}
