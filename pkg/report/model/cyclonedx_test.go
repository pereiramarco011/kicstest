package model

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/Checkmarx/kics/internal/constants"
	"github.com/Checkmarx/kics/pkg/model"
	"github.com/Checkmarx/kics/test"
	"github.com/stretchr/testify/assert"
)

var metadata Metadata = Metadata{
	Timestamp: time.Now().Format(time.RFC3339),
	Tools: &[]Tool{
		{
			Vendor:  "Checkmarx",
			Name:    "KICS",
			Version: constants.Version,
		},
	},
}

var initCycloneDxReport CycloneDxReport = CycloneDxReport{
	XMLNS:        "http://cyclonedx.org/schema/bom/1.3",
	XMLNSV:       "http://cyclonedx.org/schema/ext/vulnerability/1.0",
	SerialNumber: "urn:uuid:", // set to "urn:uuid:" because it will be different for every report
	Version:      1,
	Metadata:     &metadata,
}

// TestInitCycloneDxReport tests the InitCycloneDxReport function
func TestInitCycloneDxReport(t *testing.T) {
	tests := []struct {
		name string
		want *CycloneDxReport
	}{
		{
			name: "Init CycloneDX report",
			want: &initCycloneDxReport,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InitCycloneDxReport()
			got.SerialNumber = "urn:uuid:" // set to "urn:uuid:" because it will be different for every report
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitCycloneDxReport() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestBuildCycloneDxReport tests the BuildCycloneDxReport function
func TestBuildCycloneDxReport(t *testing.T) {
	var cycloneDx CycloneDxReport = initCycloneDxReport
	var cycloneDxCWE CycloneDxReport = initCycloneDxReport
	var vulnsC1, vulnsC2, vulnsC3 []Vulnerability
	var positiveSha, negativeSha string

	var sha256TestMap = map[string]map[string]string{
		"positive": {
			"Unix":    "487d5879d7ec205b4dcd037d5ce0075ed4fedb9dd5b8e45390ffdfa3442f15f7",
			"Windows": "bd4ac2f61e7c623477b5d200b3662fd7caac5e89e042960fd1adb008e0962635",
		},
		"negative": {
			"Unix":    "cd10cef2b154363f32ca4018426982509efbc9e1a8ea6bca587e68ffaef09c37",
			"Windows": "68b4caecf5d5130426a8b8f0222cdd7f31232b5c99a5bf0daf19099e26e2ec29",
		},
	}

	if runtime.GOOS == "windows" {
		positiveSha = sha256TestMap["positive"]["Windows"]
		negativeSha = sha256TestMap["negative"]["Windows"]
	} else {
		positiveSha = sha256TestMap["positive"]["Unix"]
		negativeSha = sha256TestMap["negative"]["Unix"]
	}

	v1 := Vulnerability{
		Ref: fmt.Sprintf("pkg:generic/../../../assets/queries/terraform/aws/guardduty_detector_disabled/test/positive.tf@0.0.0-%se38a8e0a-b88b-4902-b3fe-b0fcb17d5c10", positiveSha[0:12]),
		ID:  "e38a8e0a-b88b-4902-b3fe-b0fcb17d5c10",
		CWE: "",
		Source: Source{
			Name: "KICS",
			URL:  "https://kics.io/",
		},
		Ratings: []Rating{
			{
				Severity: "None",
				Method:   "Other",
			},
		},
		Description: "[Terraform].[Resource Not Using Tags]: AWS services resource tags are an essential part of managing components",
		Recommendations: []Recommendation{
			{
				Recommendation: "Problem found in line 1. Expected value: aws_guardduty_detector[{{positive1}}].tags is defined and not null. Actual value: aws_guardduty_detector[{{positive1}}].tags is undefined or null.",
			},
		},
	}

	v2 := Vulnerability{
		Ref: fmt.Sprintf("pkg:generic/../../../assets/queries/terraform/aws/guardduty_detector_disabled/test/positive.tf@0.0.0-%s704dadd3-54fc-48ac-b6a0-02f170011473", positiveSha[0:12]),
		ID:  "704dadd3-54fc-48ac-b6a0-02f170011473",
		CWE: "",
		Source: Source{
			Name: "KICS",
			URL:  "https://kics.io/",
		},
		Ratings: []Rating{
			{
				Severity: "Medium",
				Method:   "Other",
			},
		},
		Description: "[Terraform].[GuardDuty Detector Disabled]: Make sure that Amazon GuardDuty is Enabled",
		Recommendations: []Recommendation{
			{
				Recommendation: "Problem found in line 2. Expected value: GuardDuty Detector should be Enabled. Actual value: GuardDuty Detector is not Enabled.",
			},
		},
	}

	v3 := Vulnerability{
		Ref: fmt.Sprintf("pkg:generic/../../../assets/queries/terraform/aws/guardduty_detector_disabled/test/negative.tf@0.0.0-%se38a8e0a-b88b-4902-b3fe-b0fcb17d5c10", negativeSha[0:12]),
		ID:  "e38a8e0a-b88b-4902-b3fe-b0fcb17d5c10",
		CWE: "",
		Source: Source{
			Name: "KICS",
			URL:  "https://kics.io/",
		},
		Ratings: []Rating{
			{
				Severity: "None",
				Method:   "Other",
			},
		},
		Description: "[Terraform].[Resource Not Using Tags]: AWS services resource tags are an essential part of managing components",
		Recommendations: []Recommendation{
			{
				Recommendation: "Problem found in line 1. Expected value: aws_guardduty_detector[{{negative1}}].tags is defined and not null. Actual value: aws_guardduty_detector[{{negative1}}].tags is undefined or null.",
			},
		},
	}

	v4 := Vulnerability{
		Ref: fmt.Sprintf("pkg:generic/../../../assets/queries/terraform/aws/guardduty_detector_disabled/test/negative.tf@0.0.0-%s704dadd3-54fc-48ac-b6a0-02f170011473", negativeSha[0:12]),
		ID:  "704dadd3-54fc-48ac-b6a0-02f170011473",
		CWE: "22",
		Source: Source{
			Name: "KICS",
			URL:  "https://kics.io/",
		},
		Ratings: []Rating{
			{
				Severity: "Medium",
				Method:   "Other",
			},
		},
		Description: "[Terraform].[GuardDuty Detector Disabled]: Make sure that Amazon GuardDuty is Enabled",
		Recommendations: []Recommendation{
			{
				Recommendation: "Problem found in line 2. Expected value: GuardDuty Detector should be Enabled. Actual value: GuardDuty Detector is not Enabled.",
			},
		},
	}

	vulnsC1 = append(vulnsC1, v1)
	vulnsC1 = append(vulnsC1, v2)

	c1 := Component{
		Type:    "file",
		BomRef:  fmt.Sprintf("pkg:generic/../../../assets/queries/terraform/aws/guardduty_detector_disabled/test/positive.tf@0.0.0-%s", positiveSha[0:12]),
		Name:    "../../../assets/queries/terraform/aws/guardduty_detector_disabled/test/positive.tf",
		Version: fmt.Sprintf("0.0.0-%s", positiveSha[0:12]),
		Purl:    fmt.Sprintf("pkg:generic/../../../assets/queries/terraform/aws/guardduty_detector_disabled/test/positive.tf@0.0.0-%s", positiveSha[0:12]),
		Hashes: []Hash{
			{
				Alg:     "SHA-256",
				Content: positiveSha,
			},
		},
		Vulnerabilities: vulnsC1,
	}

	vulnsC2 = append(vulnsC2, v3)

	c2 := Component{
		Type:    "file",
		BomRef:  fmt.Sprintf("pkg:generic/../../../assets/queries/terraform/aws/guardduty_detector_disabled/test/negative.tf@0.0.0-%s", negativeSha[0:12]),
		Name:    "../../../assets/queries/terraform/aws/guardduty_detector_disabled/test/negative.tf",
		Version: fmt.Sprintf("0.0.0-%s", negativeSha[0:12]),
		Purl:    fmt.Sprintf("pkg:generic/../../../assets/queries/terraform/aws/guardduty_detector_disabled/test/negative.tf@0.0.0-%s", negativeSha[0:12]),
		Hashes: []Hash{
			{
				Alg:     "SHA-256",
				Content: negativeSha,
			},
		},
		Vulnerabilities: vulnsC2,
	}

	vulnsC3 = append(vulnsC3, v4)

	c3 := Component{
		Type:    "file",
		BomRef:  fmt.Sprintf("pkg:generic/../../../assets/queries/terraform/aws/guardduty_detector_disabled/test/negative.tf@0.0.0-%s", negativeSha[0:12]),
		Name:    "../../../assets/queries/terraform/aws/guardduty_detector_disabled/test/negative.tf",
		Version: fmt.Sprintf("0.0.0-%s", negativeSha[0:12]),
		Purl:    fmt.Sprintf("pkg:generic/../../../assets/queries/terraform/aws/guardduty_detector_disabled/test/negative.tf@0.0.0-%s", negativeSha[0:12]),
		Hashes: []Hash{
			{
				Alg:     "SHA-256",
				Content: negativeSha,
			},
		},
		Vulnerabilities: vulnsC3,
	}

	cycloneDx.Components.Components = append(cycloneDx.Components.Components, c2)
	cycloneDx.Components.Components = append(cycloneDx.Components.Components, c1)
	cycloneDxCWE.Components.Components = append(cycloneDxCWE.Components.Components, c3)

	filePaths := make(map[string]string)

	file1 := filepath.Join("..", "..", "..", "assets", "queries", "terraform", "aws", "guardduty_detector_disabled", "test", "positive.tf")
	file2 := filepath.Join("..", "..", "..", "assets", "queries", "terraform", "aws", "guardduty_detector_disabled", "test", "negative.tf")
	filePaths[file1] = file1
	filePaths[file2] = file2

	type args struct {
		summary   *model.Summary
		filePaths map[string]string
	}
	tests := []struct {
		name string
		args args
		want *CycloneDxReport
	}{
		{
			name: "Build CycloneDX report",
			args: args{
				summary:   &test.ExampleSummaryMock,
				filePaths: filePaths,
			},
			want: &cycloneDx,
		},
		{
			name: "Build CycloneDX report with cwe field complete",
			args: args{
				summary:   &test.ExampleSummaryMockCWE,
				filePaths: map[string]string{file2: file2},
			},
			want: &cycloneDxCWE,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queries := tt.args.summary.Queries
			for idx := range queries {
				for i := range queries[idx].Files {
					queries[idx].Files[i].FileName = filepath.Join("..", "..", "..", queries[idx].Files[i].FileName)
				}
			}
			got := BuildCycloneDxReport(tt.args.summary, tt.args.filePaths)
			got.SerialNumber = "urn:uuid:" // set to "urn:uuid:" because it will be different for every report
			assert.Equal(t, len(got.Components.Components), len(tt.want.Components.Components), "Comparing number of components")
			for idx := range got.Components.Components {
				assert.Equal(t, got.Components.Components[idx].BomRef, tt.want.Components.Components[idx].BomRef, "Comparing BomRef of components")
				assert.Equal(t, got.Components.Components[idx].Version, tt.want.Components.Components[idx].Version, "Comparing Version of components")
				assert.Equal(t, got.Components.Components[idx].Purl, tt.want.Components.Components[idx].Purl, "Comparing Purl of components")
				assert.Equal(t, got.Components.Components[idx].Hashes, tt.want.Components.Components[idx].Hashes, "Comparing Hashes of components")
				for idx2 := range got.Components.Components[idx].Vulnerabilities {
					assert.Equal(t, got.Components.Components[idx].Vulnerabilities[idx2].Ref, tt.want.Components.Components[idx].Vulnerabilities[idx2].Ref, "Comparing Vulnerabilities Ref of components")
					assert.Equal(t, got.Components.Components[idx].Vulnerabilities[idx2].Description, tt.want.Components.Components[idx].Vulnerabilities[idx2].Description, "Comparing Vulnerabilities Description of components")
					assert.Equal(t, got.Components.Components[idx].Vulnerabilities[idx2].Ratings, tt.want.Components.Components[idx].Vulnerabilities[idx2].Ratings, "Comparing Vulnerabilities Ratings of components")
					assert.Equal(t, got.Components.Components[idx].Vulnerabilities[idx2].Recommendations, tt.want.Components.Components[idx].Vulnerabilities[idx2].Recommendations, "Comparing Vulnerabilities Recommendations of components")
				}
			}
		})
	}
}
