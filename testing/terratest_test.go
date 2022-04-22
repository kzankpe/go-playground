package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerraformRemoteBackendStorage(t *testing.T) {

	rex := regexp.MustCompile(`(\w+)=\"(.+?)"`) // Regex to get key=value from the file
	file, errmsg := ioutil.ReadFile("../backend.tfvars")
	require.NoError(t, errmsg)

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")

	// Data transformation
	data := rex.FindAllStringSubmatch(string(file), -1)
	backendRaw := make(map[string]string) // Create an empty map for the value

	// Adding the content file content respecting the regex to the map
	for _, keyval := range data {
		k := keyval[1]
		v := keyval[2]
		backendRaw[k] = v
	}
	// Convert map to the correct type for the Backend config
	backend := make(map[string]interface{}, len(backendRaw))
	for i, y := range backendRaw {
		backend[i] = y
	}

	fmt.Println(backendRaw["storage_account_name"])
	_, err := azure.StorageAccountExistsE(backendRaw["storage_account_name"], backendRaw["resource_group_name"], subscriptionID)
	require.Error(t, err)

}

func TestTerraformRemoteBackendContainer(t *testing.T) {

	rex := regexp.MustCompile(`(\w+)=\"(.+?)"`) // Regex to get key=value from the file
	file, errmsg := ioutil.ReadFile("../backend.tfvars")
	require.NoError(t, errmsg)

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")

	// Data transformation
	data := rex.FindAllStringSubmatch(string(file), -1)
	backendRaw := make(map[string]string) // Create an empty map for the value

	// Adding the content file content respecting the regex to the map
	for _, keyval := range data {
		k := keyval[1]
		v := keyval[2]
		backendRaw[k] = v
	}
	// Convert map to the correct type for the Backend config
	backend := make(map[string]interface{}, len(backendRaw))
	for i, y := range backendRaw {
		backend[i] = y
	}

	fmt.Println(backendRaw["container_name"])
	_, err := azure.StorageBlobContainerExistsE(backendRaw["container_name"], backendRaw["storage_account_name"], backendRaw["resource_group_name"], subscriptionID)
	require.Error(t, err)

}

func TestTerraformIninitPlan(t *testing.T) {

	testFolder, err := files.CopyTerraformFolderToTemp("../", t.Name())
	require.NoError(t, err)

	defer os.RemoveAll(testFolder)
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
		tfRaw[k] = v
	}

	terraformvars := make(map[string]interface{}, len(tfRaw))
	for i, y := range tfRaw {
		terraformvars[i] = y
	}

	PlanFile := filepath.Join(testFolder, "plan.out")
	TerraformDir := "../"

	terraformOption := &terraform.Options{
		TerraformDir:  TerraformDir,
		BackendConfig: backend,
		Vars:          terraformvars,
		PlanFilePath:  PlanFile,
	}

	out, err := terraform.InitAndPlanE(t, terraformOption)
	require.NoError(t, err)
	assert.Contains(t, out, fmt.Sprintf("Saved the plan to: %s", PlanFile))
	assert.FileExists(t, PlanFile, "Plan file was not saved to expected location:", PlanFile)

	showOptions := &terraform.Options{
		TerraformDir: "../",
		PlanFilePath: PlanFile,
	}

	planJSON := terraform.Show(t, showOptions)
	require.Contains(t, planJSON, "module.cce2_cluster.flexibleengine_cce_cluster_v3.cce_cluster") // Confitm the resource is present in the plan result
}
