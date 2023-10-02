package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/docker"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

func TestDockerK8sUnit(t *testing.T) {
	t.Parallel()

	buildDockerImage(t)
	// path to K8s resources
	path := "../examples/docker-k8s/deployment.yml"
	// random namespace
	namespace := strings.ToLower(random.UniqueId())
	options := k8s.NewKubectlOptions("", "", namespace)

	defer k8s.DeleteNamespace(t, options, namespace)
	k8s.CreateNamespace(t, options, namespace)

	defer k8s.KubectlDelete(t, options, path)
	k8s.KubectlApply(t, options, path)

	validateDockerK8s(t, options)
}

func buildDockerImage(t *testing.T) {
	options := &docker.BuildOptions{
		Tags: []string{"hello-world-app:v1"},
	}
	path := "../examples/docker-k8s"
	docker.Build(t, path, options)
}

func validateDockerK8s(t *testing.T, opts *k8s.KubectlOptions) {
	k8s.WaitUntilServiceAvailable(t, opts, "hello-world-app-service", 10, 3*time.Second)

	http_helper.HttpGetWithRetry(
		t,
		serviceUrl(t, opts),
		nil,
		200,
		"Hello, World!",
		10,
		3*time.Second,
	)
}

func serviceUrl(t *testing.T, opts *k8s.KubectlOptions) string {
	service := k8s.GetService(t, opts, "hello-world-app-service")
	endpoint := k8s.GetServiceEndpoint(t, opts, service, 8080)
	return fmt.Sprintf("http://%s", endpoint)
}
