package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/docker"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/k8s"
)

func TestDockerK8s(t *testing.T) {
	buildDockerImage(t)

	path := "../examples/docker-k8s/deployment.yml"
	options := k8s.NewKubectlOptions("", "", "default")
	defer k8s.KubectlDelete(t, options, path)

	k8s.KubectlApply(t, options, path)
	validate(t, options)
}

func buildDockerImage(t *testing.T) {
	options := &docker.BuildOptions{
		Tags: []string{"hello-world-app:v1"},
	}
	path := "../examples/docker-k8s"
	docker.Build(t, path, options)
}

func validate(t *testing.T, opts *k8s.KubectlOptions) {
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
