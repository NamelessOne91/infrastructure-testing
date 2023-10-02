package test

import (
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestHelloWorldAppUnit(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/hello-world-app",
	}
	// cleanup
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)
	validateHelloWorldApp(t, terraformOptions)
}

func validateHelloWorldApp(t *testing.T, opts *terraform.Options) {
	url := terraform.Output(t, opts, "url")
	http_helper.HttpGetWithRetry(t,
		url,             // URL to test
		nil,             // TLS config
		200,             // Expected status
		"Hello, World!", //Expected body
		10,              // Max retries
		3*time.Second,   // Time between retries
	)
}
