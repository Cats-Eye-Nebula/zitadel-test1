package orb

import (
	"github.com/caos/orbos/mntr"
	kubernetes2 "github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/pkg/kubernetes"
	"github.com/pkg/errors"
)

func Reconcile(monitor mntr.Monitor, desiredTree *tree.Tree) operator.EnsureFunc {
	return func(k8sClient kubernetes2.ClientInt) (err error) {
		defer func() {
			err = errors.Wrapf(err, "building %s failed", desiredTree.Common.Kind)
		}()

		desiredKind, err := parseDesiredV0(desiredTree)
		if err != nil {
			return errors.Wrap(err, "parsing desired state failed")
		}
		desiredTree.Parsed = desiredKind

		recMonitor := monitor.WithField("version", desiredKind.Spec.Version)

		if desiredKind.Spec.Version == "" {
			err := errors.New("No version set in zitadel.yml")
			recMonitor.Error(err)
			return err
		}

		if desiredKind.Spec.SelfReconciling {
			if err := kubernetes.EnsureZitadelOperatorArtifacts(monitor, k8sClient, desiredKind.Spec.Version, desiredKind.Spec.NodeSelector, desiredKind.Spec.Tolerations); err != nil {
				recMonitor.Error(errors.Wrap(err, "Failed to deploy zitadel-operator into k8s-cluster"))
				return err
			}

			recMonitor.Info("Applied zitadel-operator")
		}
		return nil

	}
}
