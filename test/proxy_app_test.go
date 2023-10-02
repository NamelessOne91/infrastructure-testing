package test

import (
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestProxyApp(t *testing.T) {
	webServiceOtps := configWebService(t)
	defer terraform.Destroy(t, webServiceOtps)
	terraform.InitAndApply(t, webServiceOtps)

	proxyAppOpts := configProxyApp(t, webServiceOtps)
	defer terraform.Destroy(t, proxyAppOpts)
	terraform.InitAndApply(t, proxyAppOpts)

	validateProxyApp(t, proxyAppOpts)
}

func configWebService(t *testing.T) *terraform.Options {
	return &terraform.Options{
		TerraformDir: "../examples/web-service",
	}
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
