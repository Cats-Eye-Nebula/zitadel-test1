package configuration

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/pkg/errors"
)

func GetReadyFunc(
	monitor mntr.Monitor,
	namespace string,
	secretName string,
	secretVarsName string,
	secretPasswordName string,
	cmName string,
	consoleCMName string,
) func(k8sClient kubernetes.ClientInt) error {
	return func(k8sClient kubernetes.ClientInt) error {
		monitor.Debug("Waiting for configuration to be created")
		if err := k8sClient.WaitForSecret(namespace, secretName, 60); err != nil {
			return errors.Wrap(err, "error while waiting for secret")
		}

		if err := k8sClient.WaitForSecret(namespace, secretVarsName, 60); err != nil {
			return errors.Wrap(err, "error while waiting for vars secret ")
		}

		if err := k8sClient.WaitForSecret(namespace, secretPasswordName, 60); err != nil {
			return errors.Wrap(err, "error while waiting for password secret")
		}

		if err := k8sClient.WaitForConfigMap(namespace, cmName, 60); err != nil {
			return errors.Wrap(err, "error while waiting for configmap")
		}

		if err := k8sClient.WaitForConfigMap(namespace, consoleCMName, 60); err != nil {
			return errors.Wrap(err, "error while waiting for console configmap")
		}
		monitor.Debug("configuration is created")
		return nil
	}
}
