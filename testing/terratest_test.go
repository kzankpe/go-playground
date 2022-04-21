package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerraformIninitPlan(t *testing.T) {
	// Read the backend and the terraform TFvars  files
	rex := regexp.MustCompile(`(\w+)=\"(.+?)"`) // Regex to get key=value from the file
	file, errmsg := ioutil.ReadFile("../backend.tfvars")
	require.NoError(t, errmsg)

	terraformFile, errmsg := ioutil.ReadFile("../terraform.tfvars")
	require.NoError(t, errmsg)

	// Data transformation
	data := rex.FindAllStringSubmatch(string(file), -1)
	backendRaw := make(map[string]string) // Create an empty map for the value

	tfdata := rex.FindAllStringSubmatch(string(terraformFile), -1)
	tfRaw := make(map[string]string)

	// Adding the content file content respecting the regex to the map
	for _, keyval := range data {
		k := keyval[1]
		v := keyval[2]
		fmt.Println(k, v)
		backendRaw[k] = v
	}
	// Convert map to the correct type for the Backend config
	backend := make(map[string]interface{}, len(backendRaw))
	for i, y := range backendRaw {
		backend[i] = y
	}

	for _, val := range tfdata {
		k := val[1]
		v := val[2]
		fmt.Println(k, v)
		tfRaw[k] = v
	}

	terraformvars := make(map[string]interface{}, len(tfRaw))
	for i, y := range tfRaw {
		terraformvars[i] = y
	}

	PlanFile := "plan.out"
	TerraformDir := "../"

	terraformOption := &terraform.Options{
		TerraformDir:  TerraformDir,
		BackendConfig: backend,
		Vars:          terraformvars,
		PlanFilePath:  filepath.Join(TerraformDir, "plan.out"),
	}

	defer os.RemoveAll(PlanFile)

	out, err := terraform.InitAndPlanE(t, terraformOption)
	require.NoError(t, err)
	assert.Contains(t, out, fmt.Sprintf("Saved the plan to: %s", PlanFile))
	assert.FileExists(t, PlanFile, "Plan file was not saved to expected location:", PlanFile)

	showOptions := &terraform.Options{
		TerraformDir: "../",
		PlanFilePath: filepath.Join(TerraformDir, "plan.out"),
	}

	planJSON := terraform.Show(t, showOptions)
	require.Contains(t, planJSON, "null_resource.test[0]")
}
