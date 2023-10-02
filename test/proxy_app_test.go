package test

import (
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

func TestProxyAppIntegration(t *testing.T) {
	t.Parallel()

	// Since we want to be able to run multiple tests in parallel on the same modules, we need to copy them into
	// temp folders so that the state files and .terraform folders don't clash
	webServicePath := test_structure.CopyTerraformFolderToTemp(t, "../", "examples/web-service")
	proxyAppPath := test_structure.CopyTerraformFolderToTemp(t, "../", "examples/proxy-app")

	// alias for the RunTestStage func
	stage := test_structure.RunTestStage

	// Undeploy the web-service module at the end of the test
	defer stage(t, "cleanup_web_service", func() {
		cleanupWebService(t, webServicePath)
	})

	// Deploy the web-service module
	stage(t, "deploy_web_service", func() {
		deployWebService(t, webServicePath)
	})

	// Undeploy the proxy-app module at the end of the test
	defer test_structure.RunTestStage(t, "cleanup_proxy_app", func() {
		cleanupProxyApp(t, proxyAppPath)
	})

	// Deploy the proxy-app module
	test_structure.RunTestStage(t, "deploy_proxy_app", func() {
		deployProxyApp(t, webServicePath, proxyAppPath)
	})

	// Validate the proxy-app module proxies the web-service correctly
	test_structure.RunTestStage(t, "validate_proxy_app", func() {
		proxyAppOpts := test_structure.LoadTerraformOptions(t, proxyAppPath)
		validateProxyApp(t, proxyAppOpts)
	})
}

func configWebService(t *testing.T, path string) *terraform.Options {
	return &terraform.Options{
		TerraformDir: path,
	}
}

func deployWebService(t *testing.T, path string) {
	webServiceOpts := configWebService(t, path)
	test_structure.SaveTerraformOptions(t, path, webServiceOpts)
	terraform.InitAndApply(t, webServiceOpts)
}

func cleanupWebService(t *testing.T, path string) {
	opts := test_structure.LoadTerraformOptions(t, path)
	terraform.Destroy(t, opts)
}

func configProxyApp(t *testing.T, webServiceOpts *terraform.Options) *terraform.Options {
	url := terraform.Output(t, webServiceOpts, "url")

	return &terraform.Options{
		TerraformDir: "../examples/proxy-app",
		Vars: map[string]interface{}{
			"url_to_proxy": url,
		},
	}
}

func deployProxyApp(t *testing.T, webServicePath, proxyPath string) {
	webServiceOpts := test_structure.LoadTerraformOptions(t, webServicePath)
	proxyAppOpts := configProxyApp(t, webServiceOpts)
	test_structure.SaveTerraformOptions(t, proxyPath, proxyAppOpts)
	terraform.InitAndApply(t, proxyAppOpts)
}

func cleanupProxyApp(t *testing.T, path string) {
	proxyAppOpts := test_structure.LoadTerraformOptions(t, path)
	terraform.Destroy(t, proxyAppOpts)
}

func validateProxyApp(t *testing.T, opts *terraform.Options) {
	url := terraform.Output(t, opts, "url")
	http_helper.HttpGetWithRetry(
		t,
		url,
		nil,
		200,
		`{"text": "Hello, World!}"`,
		10,
		3*time.Second,
	)
}
