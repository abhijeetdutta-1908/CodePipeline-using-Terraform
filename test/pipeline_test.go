package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codepipeline"
	"github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAwsCodePipeline(t *testing.T) {
	t.Parallel()

	// The Terraform directory to test
	terraformDir := "../terraform"

	// Configure Terraform options
	terraformOptions := &terraform.Options{
		TerraformDir: terraformDir,
		// We don't need to specify vars here if using terraform.tfvars
	}

	// At the end of the test, run `terraform destroy`
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`
	terraform.InitAndApply(t, terraformOptions)

	// Get outputs from Terraform
	pipelineName := terraform.Output(t, terraformOptions, "codepipeline_name")
	awsRegion := terraform.Output(t, terraformOptions, "aws_region") // Assuming you add this output
	ec2PublicIp := terraform.Output(t, terraformOptions, "ec2_public_ip")

	// --- Validation Checks ---

	// 1. Check if the CodePipeline exists and has the correct name
	sess, err := session.NewSession(&aws.Config{Region: aws.String(awsRegion)})
	assert.NoError(t, err)
	cpClient := codepipeline.New(sess)

	getPipelineInput := &codepipeline.GetPipelineInput{
		Name: aws.String(pipelineName),
	}

	_, err = cpClient.GetPipeline(getPipelineInput)
	assert.NoError(t, err, "Failed to find CodePipeline: %s", pipelineName)

	// 2. Trigger the pipeline (optional, but good for end-to-end testing)
	// For this test, we'll just check if the initial deployment worked.
	// A more advanced test would commit a change and wait for the deployment.

	// 3. Check if the EC2 instance is serving the correct content
	// It can take a few minutes for the first deployment to complete.
	url := fmt.Sprintf("http://%s", ec2PublicIp)
	expectedText := "Your AWS CodePipeline deployment is working."

	// Retry checking the URL until we get the expected response or timeout.
	http_helper.HttpGetWithRetry(
		t,
		url,
		nil,
		200,
		expectedText,
		30,           // Number of retries
		10*time.Second, // Delay between retries
	)
}