package database

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator"
)

func AdaptFunc(
	monitor mntr.Monitor,
	repoURL string,
	repoKey string,
) (
	operator.QueryFunc,
	error,
) {

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
		dbHost, dbPort, err := GetConnectionInfo(monitor, k8sClient, repoURL, repoKey)
		if err != nil {
			return nil, err
		}

		curr := &Current{
			Host: dbHost,
			Port: dbPort,
		}

		SetDatabaseInQueried(queried, curr)

		return nil, nil
	}, nil
}
