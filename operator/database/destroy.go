package database

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/zitadel/zitadel/operator"
)

func Destroy(
	monitor mntr.Monitor,
	gitClient *git.Client,
	adapt operator.AdaptFunc,
	k8sClient *kubernetes.Client,
) error {
	internalMonitor := monitor.WithField("operator", "database")
	internalMonitor.Info("Destroy")
	treeDesired, err := operator.Parse(gitClient, "database.yml")
	if err != nil {
		return err
	}
	treeCurrent := &tree.Tree{}

	_, destroy, _, _, _, _, err := adapt(internalMonitor, treeDesired, treeCurrent)
	if err != nil {
		return err
	}

	if err := destroy(k8sClient); err != nil {
		return err
	}
	return nil
}
