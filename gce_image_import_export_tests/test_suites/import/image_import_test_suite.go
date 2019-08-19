//  Copyright 2019 Google Inc. All Rights Reserved.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

// Package importtestsuites contains e2e tests for image import cli tools
package importtestsuites

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/common/utils/path"
	"github.com/GoogleCloudPlatform/compute-image-tools/gce_image_import_export_tests/compute"
	"github.com/GoogleCloudPlatform/compute-image-tools/gce_image_import_export_tests/test_suites"
	"github.com/GoogleCloudPlatform/compute-image-tools/go/e2e_test_utils/junitxml"
	"github.com/GoogleCloudPlatform/compute-image-tools/go/e2e_test_utils/test_config"
)

const (
	testSuiteName = "ImageImportTests"
)

// TestSuite is image import test suite.
func TestSuite(
	ctx context.Context, tswg *sync.WaitGroup, testSuites chan *junitxml.TestSuite,
	logger *log.Logger, testSuiteRegex, testCaseRegex *regexp.Regexp,
	testProjectConfig *testconfig.Project) {

	testTypes := []testsuiteutils.TestType{
		testsuiteutils.Wrapper,
		testsuiteutils.GcloudProdWrapperLatest,
		testsuiteutils.GcloudLatestWrapperLatest,
	}

	testsMap := map[testsuiteutils.TestType]map[*junitxml.TestCase]func(
		context.Context, *junitxml.TestCase, *log.Logger, *testconfig.Project, testsuiteutils.TestType){}

	for _, testType := range testTypes {
		imageImportDataDiskTestCase := junitxml.NewTestCase(
			testSuiteName, fmt.Sprintf("[%v][ImageImport] %v", testType, "Import data disk"))
		imageImportOSTestCase := junitxml.NewTestCase(
			testSuiteName, fmt.Sprintf("[%v][ImageImport] %v", testType, "Import OS"))
		imageImportOSFromImageTestCase := junitxml.NewTestCase(
			testSuiteName, fmt.Sprintf("[%v][ImageImport] %v", testType, "Import OS from image"))
		imageImportWithRichParamsTestCase := junitxml.NewTestCase(
			testSuiteName, fmt.Sprintf("[%v][ImageImport] %v", testType, "Import with rich params"))
		imageImportWithSubnetWithoutNetworkSpecifiedTestCase := junitxml.NewTestCase(
			testSuiteName, fmt.Sprintf("[%v][ImageImport] %v", testType, "Import with subnet but without network"))

		testsMap[testType] = map[*junitxml.TestCase]func(
			context.Context, *junitxml.TestCase, *log.Logger, *testconfig.Project, testsuiteutils.TestType){}
		testsMap[testType][imageImportDataDiskTestCase] = runImageImportDataDiskTest
		testsMap[testType][imageImportOSTestCase] = runImageImportOSTest
		testsMap[testType][imageImportOSFromImageTestCase] = runImageImportOSFromImageTest
		testsMap[testType][imageImportWithRichParamsTestCase] = runImageImportWithRichParamsTest
		testsMap[testType][imageImportWithSubnetWithoutNetworkSpecifiedTestCase] = runImageImportWithSubnetWithoutNetworkSpecified

	}
	testsuiteutils.TestSuite(ctx, tswg, testSuites, logger, testSuiteRegex, testCaseRegex,
		testProjectConfig, testSuiteName, testsMap)
}

func runImageImportDataDiskTest(ctx context.Context, testCase *junitxml.TestCase, logger *log.Logger,
	testProjectConfig *testconfig.Project, testType testsuiteutils.TestType) {

	suffix := path.RandString(5)
	imageName := "e2e-test-image-import-data-disk-" + suffix

	argsMap := map[testsuiteutils.TestType][]string{
		testsuiteutils.Wrapper: {"-client_id=e2e", fmt.Sprintf("-project=%v", testProjectConfig.TestProjectID),
			fmt.Sprintf("-image_name=%s", imageName), "-data_disk",
			fmt.Sprintf("-source_file=gs://%v-test-image/image-file-10g-vmdk", testProjectConfig.TestProjectID),
			fmt.Sprintf("-zone=%v", testProjectConfig.TestZone),
		},
		testsuiteutils.GcloudProdWrapperLatest: {"beta", "compute", "images", "import", imageName, "--quiet",
			"--docker-image-tag=latest", "--data-disk", fmt.Sprintf("--project=%v", testProjectConfig.TestProjectID),
			fmt.Sprintf("--source-file=gs://%v-test-image/image-file-10g-vmdk", testProjectConfig.TestProjectID),
			fmt.Sprintf("--zone=%v", testProjectConfig.TestZone),
		},
		testsuiteutils.GcloudLatestWrapperLatest: {"beta", "compute", "images", "import", imageName, "--quiet",
			"--docker-image-tag=latest", "--data-disk", fmt.Sprintf("--project=%v", testProjectConfig.TestProjectID),
			fmt.Sprintf("--source-file=gs://%v-test-image/image-file-10g-vmdk", testProjectConfig.TestProjectID),
			fmt.Sprintf("--zone=%v", testProjectConfig.TestZone),
		},
	}

	runImportTest(ctx, argsMap[testType], testType, testProjectConfig, imageName, logger, testCase)
}

func runImageImportOSTest(ctx context.Context, testCase *junitxml.TestCase, logger *log.Logger,
	testProjectConfig *testconfig.Project, testType testsuiteutils.TestType) {

	suffix := path.RandString(5)
	imageName := "e2e-test-image-import-os-" + suffix

	argsMap := map[testsuiteutils.TestType][]string{
		testsuiteutils.Wrapper: {"-client_id=e2e", fmt.Sprintf("-project=%v", testProjectConfig.TestProjectID),
			fmt.Sprintf("-image_name=%v", imageName), "-os=debian-9",
			fmt.Sprintf("-source_file=gs://%v-test-image/image-file-10g-vmdk", testProjectConfig.TestProjectID),
			fmt.Sprintf("-zone=%v", testProjectConfig.TestZone),
		},
		testsuiteutils.GcloudProdWrapperLatest: {"beta", "compute", "images", "import", imageName, "--quiet",
			"--docker-image-tag=latest", "--os=debian-9", fmt.Sprintf("--project=%v", testProjectConfig.TestProjectID),
			fmt.Sprintf("--source-file=gs://%v-test-image/image-file-10g-vmdk", testProjectConfig.TestProjectID),
			fmt.Sprintf("--zone=%v", testProjectConfig.TestZone),
		},
		testsuiteutils.GcloudLatestWrapperLatest: {"beta", "compute", "images", "import", imageName, "--quiet",
			"--docker-image-tag=latest", "--os=debian-9", fmt.Sprintf("--project=%v", testProjectConfig.TestProjectID),
			fmt.Sprintf("--source-file=gs://%v-test-image/image-file-10g-vmdk", testProjectConfig.TestProjectID),
			fmt.Sprintf("--zone=%v", testProjectConfig.TestZone),
		},
	}

	runImportTest(ctx, argsMap[testType], testType, testProjectConfig, imageName, logger, testCase)
}

func runImageImportOSFromImageTest(ctx context.Context, testCase *junitxml.TestCase, logger *log.Logger,
	testProjectConfig *testconfig.Project, testType testsuiteutils.TestType) {

	suffix := path.RandString(5)
	imageName := "e2e-test-image-import-os-from-image-" + suffix

	argsMap := map[testsuiteutils.TestType][]string{
		testsuiteutils.Wrapper: {"-client_id=e2e", fmt.Sprintf("-project=%v", testProjectConfig.TestProjectID),
			fmt.Sprintf("-image_name=%v", imageName), "-os=debian-9", "-source_image=e2e-test-image-10g",
			fmt.Sprintf("-zone=%v", testProjectConfig.TestZone),
		},
		testsuiteutils.GcloudProdWrapperLatest: {"beta", "compute", "images", "import", imageName, "--quiet",
			"--docker-image-tag=latest", "--os=debian-9", fmt.Sprintf("--project=%v", testProjectConfig.TestProjectID),
			"--source-image=e2e-test-image-10g",
			fmt.Sprintf("--zone=%v", testProjectConfig.TestZone),
		},
		testsuiteutils.GcloudLatestWrapperLatest: {"beta", "compute", "images", "import", imageName, "--quiet",
			"--docker-image-tag=latest", "--os=debian-9", fmt.Sprintf("--project=%v", testProjectConfig.TestProjectID),
			"--source-image=e2e-test-image-10g",
			fmt.Sprintf("--zone=%v", testProjectConfig.TestZone),
		},
	}

	runImportTest(ctx, argsMap[testType], testType, testProjectConfig, imageName, logger, testCase)
}

// Test most of params except -oauth, -compute_endpoint_override, and -scratch_bucket_gcs_path
func runImageImportWithRichParamsTest(ctx context.Context, testCase *junitxml.TestCase, logger *log.Logger,
	testProjectConfig *testconfig.Project, testType testsuiteutils.TestType) {

	family := "test-family"
	description := "test-description"
	labels := []string{"key1=value1", "key2=value2"}

	suffix := path.RandString(5)
	imageName := "e2e-test-image-import-rich-param-" + suffix

	argsMap := map[testsuiteutils.TestType][]string{
		testsuiteutils.Wrapper: {"-client_id=e2e", fmt.Sprintf("-project=%v", testProjectConfig.TestProjectID),
			fmt.Sprintf("-image_name=%s", imageName), "-data_disk",
			fmt.Sprintf("-source_file=gs://%v-test-image/image-file-10g-vmdk", testProjectConfig.TestProjectID),
			"-no_guest_environment", fmt.Sprintf("-family=%v", family), fmt.Sprintf("-description=%v", description),
			fmt.Sprintf("-network=%v-vpc-1", testProjectConfig.TestProjectID),
			fmt.Sprintf("-subnet=%v-subnet-1", testProjectConfig.TestProjectID),
			fmt.Sprintf("-zone=%v", testProjectConfig.TestZone),
			"-timeout=2h", "-disable_gcs_logging", "-disable_cloud_logging", "-disable_stdout_logging",
			"-no_external_ip", fmt.Sprintf("-labels=%v", strings.Join(labels, ",")),
		},
		testsuiteutils.GcloudProdWrapperLatest: {"beta", "compute", "images", "import", imageName, "--quiet",
			"--docker-image-tag=latest", "--data-disk", fmt.Sprintf("--project=%v", testProjectConfig.TestProjectID),
			fmt.Sprintf("--source-file=gs://%v-test-image/image-file-10g-vmdk", testProjectConfig.TestProjectID),
			"--no-guest-environment",
			fmt.Sprintf("--network=%v-vpc-1", testProjectConfig.TestProjectID),
			fmt.Sprintf("--subnet=%v-subnet-1", testProjectConfig.TestProjectID),
			fmt.Sprintf("--zone=%v", testProjectConfig.TestZone), "--timeout=2h",
		},
		testsuiteutils.GcloudLatestWrapperLatest: {"beta", "compute", "images", "import", imageName, "--quiet",
			"--docker-image-tag=latest", "--data-disk", fmt.Sprintf("--project=%v", testProjectConfig.TestProjectID),
			fmt.Sprintf("--source-file=gs://%v-test-image/image-file-10g-vmdk", testProjectConfig.TestProjectID),
			"--no-guest-environment",
			fmt.Sprintf("--network=%v-vpc-1", testProjectConfig.TestProjectID),
			fmt.Sprintf("--subnet=%v-subnet-1", testProjectConfig.TestProjectID),
			fmt.Sprintf("--zone=%v", testProjectConfig.TestZone), "--timeout=2h",
		},
	}

	runImportTestWithExtraParams(ctx, argsMap[testType], testType, testProjectConfig, imageName,
		logger, testCase, family, description, labels)
}

func runImageImportWithSubnetWithoutNetworkSpecified(ctx context.Context, testCase *junitxml.TestCase,
	logger *log.Logger, testProjectConfig *testconfig.Project, testType testsuiteutils.TestType) {

	suffix := path.RandString(5)
	imageName := "e2e-test-image-import-subnet-" + suffix

	argsMap := map[testsuiteutils.TestType][]string{
		testsuiteutils.Wrapper: {"-client_id=e2e", fmt.Sprintf("-project=%v", testProjectConfig.TestProjectID),
			fmt.Sprintf("-image_name=%s", imageName), "-data_disk",
			fmt.Sprintf("-source_file=gs://%v-test-image/image-file-10g-vmdk", testProjectConfig.TestProjectID),
			fmt.Sprintf("-subnet=%v-subnet-1", testProjectConfig.TestProjectID),
			fmt.Sprintf("-zone=%v", testProjectConfig.TestZone),
		},
		testsuiteutils.GcloudProdWrapperLatest: {"beta", "compute", "images", "import", imageName, "--quiet",
			"--docker-image-tag=latest", "--data-disk", fmt.Sprintf("--project=%v", testProjectConfig.TestProjectID),
			fmt.Sprintf("--source-file=gs://%v-test-image/image-file-10g-vmdk", testProjectConfig.TestProjectID),
			fmt.Sprintf("--subnet=%v-subnet-1", testProjectConfig.TestProjectID),
			fmt.Sprintf("--zone=%v", testProjectConfig.TestZone),
		},
		testsuiteutils.GcloudLatestWrapperLatest: {"beta", "compute", "images", "import", imageName, "--quiet",
			"--docker-image-tag=latest", "--data-disk", fmt.Sprintf("--project=%v", testProjectConfig.TestProjectID),
			fmt.Sprintf("--source-file=gs://%v-test-image/image-file-10g-vmdk", testProjectConfig.TestProjectID),
			fmt.Sprintf("--subnet=%v-subnet-1", testProjectConfig.TestProjectID),
			fmt.Sprintf("--zone=%v", testProjectConfig.TestZone),
		},
	}

	runImportTest(ctx, argsMap[testType], testType, testProjectConfig, imageName, logger, testCase)
}

func runImportTest(ctx context.Context, args []string, testType testsuiteutils.TestType,
	testProjectConfig *testconfig.Project, imageName string, logger *log.Logger, testCase *junitxml.TestCase) {

	runImportTestWithExtraParams(ctx, args, testType, testProjectConfig, imageName, logger, testCase, "", "", nil)
}

func runImportTestWithExtraParams(ctx context.Context, args []string, testType testsuiteutils.TestType,
	testProjectConfig *testconfig.Project, imageName string, logger *log.Logger, testCase *junitxml.TestCase,
	expectedFamily string, expectedDescription string, expectedLabels []string) {

	cmds := map[testsuiteutils.TestType]string{
		testsuiteutils.Wrapper:                   "./gce_vm_image_import",
		testsuiteutils.GcloudProdWrapperLatest:   "gcloud",
		testsuiteutils.GcloudLatestWrapperLatest: "gcloud",
	}

	// "family", "description" and "labels" hasn't been supported by gcloud
	if testType != testsuiteutils.Wrapper {
		expectedFamily = ""
		expectedDescription = ""
		expectedLabels = nil
	}

	if testsuiteutils.RunTestForTestType(cmds[testType], args, testType, logger, testCase) {
		verifyImportedImage(ctx, testCase, testProjectConfig, imageName, logger, expectedFamily,
			expectedDescription, expectedLabels)
	}
}

func verifyImportedImage(ctx context.Context, testCase *junitxml.TestCase,
	testProjectConfig *testconfig.Project, imageName string, logger *log.Logger,
	expectedFamily string, expectedDescription string, expectedLabels []string) {

	logger.Printf("Verifying imported image...")
	image, err := compute.CreateImageObject(ctx, testProjectConfig.TestProjectID, imageName)
	if err != nil {
		testCase.WriteFailure("Image '%v' doesn't exist after import: %v", imageName, err)
		logger.Printf("Image '%v' doesn't exist after import: %v", imageName, err)
		return
	}
	logger.Printf("Image '%v' exists! Import success.", imageName)

	if expectedFamily != "" && image.Family != expectedFamily {
		testCase.WriteFailure("Image '%v' family expect: %v, actual: %v", imageName, expectedFamily, image.Family)
		logger.Printf("Image '%v' family expect: %v, actual: %v", imageName, expectedFamily, image.Family)
	}

	if expectedDescription != "" && image.Description != expectedDescription {
		testCase.WriteFailure("Image '%v' description expect: %v, actual: %v", imageName, expectedDescription, image.Description)
		logger.Printf("Image '%v' description expect: %v, actual: %v", imageName, expectedDescription, image.Description)
	}

	if expectedLabels != nil {
		imageLabels := make([]string, 0, len(image.Labels))
		for k, v := range image.Labels {
			imageLabels = append(imageLabels, k+"="+v)
		}
		if !contains(imageLabels, expectedLabels) {
			testCase.WriteFailure("Image '%v' labels expect: %v, actual: %v", imageName, strings.Join(expectedLabels, ","), strings.Join(imageLabels, ","))
			logger.Printf("Image '%v' labels expect: %v, actual: %v", imageName, strings.Join(expectedLabels, ","), strings.Join(imageLabels, ","))
		}
	}

	if err := image.Cleanup(); err != nil {
		logger.Printf("Image '%v' failed to clean up.", imageName)
	} else {
		logger.Printf("Image '%v' cleaned up.", imageName)
	}
}

func contains(arr []string, subarr []string) bool {
	for item := range subarr {
		exists := false
		for i := range arr {
			if item == i {
				exists = true
				break
			}
		}
		if !exists {
			return false
		}
	}
	return true
}
