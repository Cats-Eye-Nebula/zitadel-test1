package hosts

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/ambassador/host"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ambassador/skipcrd"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/configuration"
)

const (
	AccountsHostName = "accounts"
	ApiHostName      = "api"
	ConsoleHostName  = "console"
	IssuerHostName   = "issuer"
)

func AdaptFunc(
	monitor mntr.Monitor,
	componentLabels *labels.Component,
	namespace string,
	dns *configuration.DNS,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	internalMonitor := monitor.WithField("part", "hosts")

	destroyAccounts, err := host.AdaptFuncToDestroy(namespace, AccountsHostName)
	if err != nil {
		return nil, nil, err
	}

	destroyAPI, err := host.AdaptFuncToDestroy(namespace, ApiHostName)
	if err != nil {
		return nil, nil, err
	}

	destroyConsole, err := host.AdaptFuncToDestroy(namespace, ConsoleHostName)
	if err != nil {
		return nil, nil, err
	}

	destroyIssuer, err := host.AdaptFuncToDestroy(namespace, IssuerHostName)
	if err != nil {
		return nil, nil, err
	}

	destroyers := []operator.DestroyFunc{
		operator.ResourceDestroyToZitadelDestroy(destroyAccounts),
		operator.ResourceDestroyToZitadelDestroy(destroyAPI),
		operator.ResourceDestroyToZitadelDestroy(destroyConsole),
		operator.ResourceDestroyToZitadelDestroy(destroyIssuer),
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			if skipEnsure, err := skipcrd.EnsureFunc(monitor, k8sClient, "hosts.getambassador.io"); err != nil || skipEnsure != nil {
				return skipEnsure, err
			}

			accountsDomain := dns.Subdomains.Accounts + "." + dns.Domain
			apiDomain := dns.Subdomains.API + "." + dns.Domain
			consoleDomain := dns.Subdomains.Console + "." + dns.Domain
			issuerDomain := dns.Subdomains.Issuer + "." + dns.Domain
			originCASecretName := dns.TlsSecret
			authority := dns.ACMEAuthority
			if authority == "" {
				authority = "none"
			}

			accountsSelector := map[string]string{
				"hostname": accountsDomain,
			}
			queryAccounts, err := host.AdaptFuncToEnsure(
				monitor,
				namespace,
				AccountsHostName,
				labels.MustForNameK8SMap(componentLabels, AccountsHostName),
				accountsDomain,
				authority,
				"",
				accountsSelector,
				originCASecretName,
			)
			if err != nil {
				return nil, err
			}

			apiSelector := map[string]string{
				"hostname": apiDomain,
			}
			queryAPI, err := host.AdaptFuncToEnsure(
				monitor,
				namespace,
				ApiHostName,
				labels.MustForNameK8SMap(componentLabels, ApiHostName),
				apiDomain,
				authority,
				"",
				apiSelector,
				originCASecretName,
			)
			if err != nil {
				return nil, err
			}

			consoleSelector := map[string]string{
				"hostname": consoleDomain,
			}
			queryConsole, err := host.AdaptFuncToEnsure(
				monitor,
				namespace,
				ConsoleHostName,
				labels.MustForNameK8SMap(componentLabels, ConsoleHostName),
				consoleDomain,
				authority,
				"",
				consoleSelector,
				originCASecretName,
			)
			if err != nil {
				return nil, err
			}

			issuerSelector := map[string]string{
				"hostname": issuerDomain,
			}
			queryIssuer, err := host.AdaptFuncToEnsure(
				monitor,
				namespace,
				IssuerHostName,
				labels.MustForNameK8SMap(componentLabels, IssuerHostName),
				issuerDomain,
				authority,
				"",
				issuerSelector,
				originCASecretName,
			)
			if err != nil {
				return nil, err
			}

			queriers := []operator.QueryFunc{
				operator.ResourceQueryToZitadelQuery(queryAccounts),
				operator.ResourceQueryToZitadelQuery(queryAPI),
				operator.ResourceQueryToZitadelQuery(queryConsole),
				operator.ResourceQueryToZitadelQuery(queryIssuer),
			}

			return operator.QueriersToEnsureFunc(internalMonitor, false, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(internalMonitor, destroyers),
		nil
}
