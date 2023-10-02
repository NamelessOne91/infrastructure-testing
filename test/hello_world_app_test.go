package test

import (
	"fmt"
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestHelloWorldAppUnit(t *testing.T) {
	t.Parallel()

	// 6 characters random string
	uniqueId := random.UniqueId()
	terraformOptions := &terraform.Options{
		// path to Terraform code
		TerraformDir: "../examples/hello-world-app",
		// variables being passed to our Terraform code using -var options
		Vars: map[string]interface{}{
			"name": fmt.Sprintf("hello-world-app-%s", uniqueId),
		},
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
