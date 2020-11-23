package grpc

import (
	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/configuration"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	apixv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

func SetReturnResourceVersion(
	k8sClient *kubernetesmock.MockClientInt,
	group,
	version,
	kind,
	namespace,
	name string,
	resourceVersion string,
) {
	ret := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"resourceVersion": resourceVersion,
			},
		},
	}
	k8sClient.EXPECT().GetNamespacedCRDResource(group, version, kind, namespace, name).Return(ret, nil)
}

func TestGrpc_Adapt(t *testing.T) {
	monitor := mntr.Monitor{}
	namespace := "test"
	labels := map[string]string{"test": "test"}
	url := "url"
	dns := &configuration.DNS{
		Domain:    "",
		TlsSecret: "",
		Subdomains: &configuration.Subdomains{
			Accounts: "",
			API:      "",
			Console:  "",
			Issuer:   "",
		},
	}
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	k8sClient.EXPECT().CheckCRD("mappings.getambassador.io").Times(1).Return(&apixv1beta1.CustomResourceDefinition{}, nil)

	group := "getambassador.io"
	version := "v2"
	kind := "Mapping"

	cors := map[string]interface{}{
		"origins":         "*",
		"methods":         "POST, GET, OPTIONS, DELETE, PUT",
		"headers":         "*",
		"credentials":     true,
		"exposed_headers": "*",
		"max_age":         "86400",
	}
	adminM := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": group + "/" + version,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"labels":    labels,
				"name":      AdminMName,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"connect_timeout_ms": 30000,
				"host":               ".",
				"prefix":             "/caos.zitadel.admin.api.v1.AdminService/",
				"rewrite":            "",
				"service":            url,
				"timeout_ms":         30000,
				"cors":               cors,
				"grpc":               true,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, AdminMName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, AdminMName, adminM).Times(1)

	authM := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": group + "/" + version,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"labels":    labels,
				"name":      AuthMName,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"connect_timeout_ms": 30000,
				"host":               ".",
				"prefix":             "/caos.zitadel.auth.api.v1.AuthService/",
				"rewrite":            "",
				"service":            url,
				"timeout_ms":         30000,
				"cors":               cors,
				"grpc":               true,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, AuthMName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, AuthMName, authM).Times(1)

	mgmtM := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": group + "/" + version,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"labels":    labels,
				"name":      MgmtMName,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"connect_timeout_ms": 30000,
				"host":               ".",
				"prefix":             "/caos.zitadel.management.api.v1.ManagementService/",
				"rewrite":            "",
				"service":            url,
				"timeout_ms":         30000,
				"cors":               cors,
				"grpc":               true,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, MgmtMName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, MgmtMName, mgmtM).Times(1)

	query, _, err := AdaptFunc(monitor, namespace, labels, url, dns)
	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))
}

func TestGrpc_Adapt2(t *testing.T) {
	monitor := mntr.Monitor{}
	namespace := "test"
	labels := map[string]string{"test": "test"}
	url := "url"
	dns := &configuration.DNS{
		Domain:    "domain",
		TlsSecret: "tls",
		Subdomains: &configuration.Subdomains{
			Accounts: "accounts",
			API:      "api",
			Console:  "console",
			Issuer:   "issuer",
		},
	}
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	k8sClient.EXPECT().CheckCRD("mappings.getambassador.io").Times(1).Return(&apixv1beta1.CustomResourceDefinition{}, nil)

	group := "getambassador.io"
	version := "v2"
	kind := "Mapping"

	cors := map[string]interface{}{
		"origins":         "*",
		"methods":         "POST, GET, OPTIONS, DELETE, PUT",
		"headers":         "*",
		"credentials":     true,
		"exposed_headers": "*",
		"max_age":         "86400",
	}
	adminM := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": group + "/" + version,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"labels":    labels,
				"name":      AdminMName,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"connect_timeout_ms": 30000,
				"host":               "api.domain",
				"prefix":             "/caos.zitadel.admin.api.v1.AdminService/",
				"rewrite":            "",
				"service":            url,
				"timeout_ms":         30000,
				"cors":               cors,
				"grpc":               true,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, AdminMName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, AdminMName, adminM).Times(1)

	authM := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": group + "/" + version,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"labels":    labels,
				"name":      AuthMName,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"connect_timeout_ms": 30000,
				"host":               "api.domain",
				"prefix":             "/caos.zitadel.auth.api.v1.AuthService/",
				"rewrite":            "",
				"service":            url,
				"timeout_ms":         30000,
				"cors":               cors,
				"grpc":               true,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, AuthMName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, AuthMName, authM).Times(1)

	mgmtM := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": group + "/" + version,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"labels":    labels,
				"name":      MgmtMName,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"connect_timeout_ms": 30000,
				"host":               "api.domain",
				"prefix":             "/caos.zitadel.management.api.v1.ManagementService/",
				"rewrite":            "",
				"service":            url,
				"timeout_ms":         30000,
				"cors":               cors,
				"grpc":               true,
			},
		},
	}
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, MgmtMName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, MgmtMName, mgmtM).Times(1)

	query, _, err := AdaptFunc(monitor, namespace, labels, url, dns)
	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))
}
