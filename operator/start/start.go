package start

import (
	"context"
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/databases"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	orbconfig "github.com/caos/orbos/pkg/orb"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/kinds/orb"
	kubernetes2 "github.com/caos/zitadel/pkg/kubernetes"
	"runtime/debug"
	"time"
)

func Operator(monitor mntr.Monitor, orbConfigPath string, k8sClient *kubernetes.Client, migrationsPath string) error {
	takeoffChan := make(chan struct{})
	go func() {
		takeoffChan <- struct{}{}
	}()

	for range takeoffChan {
		orbConfig, err := orbconfig.ParseOrbConfig(orbConfigPath)
		if err != nil {
			monitor.Error(err)
			return err
		}

		gitClient := git.New(context.Background(), monitor, "orbos", "orbos@caos.ch")
		if err := gitClient.Configure(orbConfig.URL, []byte(orbConfig.Repokey)); err != nil {
			monitor.Error(err)
			return err
		}

		takeoff := operator.Takeoff(monitor, gitClient, orb.AdaptFunc(orbConfig, "ensure", migrationsPath, []string{"iam"}), k8sClient)

		go func() {
			started := time.Now()
			takeoff()

			monitor.WithFields(map[string]interface{}{
				"took": time.Since(started),
			}).Info("Iteration done")
			debug.FreeOSMemory()

			takeoffChan <- struct{}{}
		}()
	}

	return nil
}

func Restore(monitor mntr.Monitor, gitClient *git.Client, k8sClient *kubernetes.Client, backup, migrationsPath string) error {
	databasesList := []string{
		"notification",
		"adminapi",
		"auth",
		"authz",
		"eventstore",
		"management",
	}
	emptyOrbConfig := &orbconfig.Orb{}

	if err := kubernetes2.ScaleZitadelOperator(monitor, k8sClient, 0); err != nil {
		return err
	}

	if err := operator.Takeoff(monitor, gitClient, orb.AdaptFunc(emptyOrbConfig, "scaledown", migrationsPath, []string{"scaledown"}), k8sClient)(); err != nil {
		return err
	}

	if err := databases.Clear(monitor, k8sClient, gitClient, databasesList); err != nil {
		return err
	}

	if err := operator.Takeoff(monitor, gitClient, orb.AdaptFunc(emptyOrbConfig, "migration", migrationsPath, []string{"migration"}), k8sClient)(); err != nil {
		return err
	}

	if err := databases.Restore(
		monitor,
		k8sClient,
		gitClient,
		backup,
		databasesList,
	); err != nil {
		return err
	}

	if err := operator.Takeoff(monitor, gitClient, orb.AdaptFunc(emptyOrbConfig, "scaleup", migrationsPath, []string{"scaleup"}), k8sClient)(); err != nil {
		return err
	}

	if err := kubernetes2.ScaleZitadelOperator(monitor, k8sClient, 1); err != nil {
		return err
	}

	return nil
}
