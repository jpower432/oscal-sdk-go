/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package settings

import (
	"os"
	"path/filepath"
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/stretchr/testify/require"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/generators"
)

func TestMerge(t *testing.T) {
	tests := []struct {
		name                string
		inputImplementation oscalTypes.ControlImplementationSet
		assertFunc          func(t *testing.T, settings *ImplementationSettings)
	}{
		{
			name: "Valid/ImplementationOnly",
			inputImplementation: oscalTypes.ControlImplementationSet{
				SetParameters: &[]oscalTypes.SetParameter{
					{
						ParamId: "my-test-param",
						Values:  []string{"test-value"},
					},
				},
				ImplementedRequirements: []oscalTypes.ImplementedRequirementControlImplementation{
					{
						ControlId: "ex-1",
						Props: &[]oscalTypes.Property{
							{
								Name:  extensions.RuleIdProp,
								Value: "my-test-rule",
								Ns:    extensions.TrestleNameSpace,
							},
						},
					},
				},
			},
			assertFunc: func(t *testing.T, impSettings *ImplementationSettings) {
				settings := impSettings.AllSettings()
				require.NotNil(t, settings)
				require.True(t, settings.ContainsRule("my-test-rule"))
				applicationControls, err := impSettings.ApplicableControls("my-test-rule")
				require.NoError(t, err)
				require.Len(t, applicationControls, 1)
				require.Equal(t, applicationControls[0], "ex-1")
				require.Contains(t, settings.selectedParameters, "my-test-param")
			},
		},
		{
			name: "Valid/ControlLevel",
			inputImplementation: oscalTypes.ControlImplementationSet{
				ImplementedRequirements: []oscalTypes.ImplementedRequirementControlImplementation{
					{
						ControlId: "ex-1",
						SetParameters: &[]oscalTypes.SetParameter{
							{
								ParamId: "my-test-param",
								Values:  []string{"test-value"},
							},
						},
						Props: &[]oscalTypes.Property{
							{
								Name:  extensions.RuleIdProp,
								Value: "my-test-rule",
								Ns:    extensions.TrestleNameSpace,
							},
						},
					},
				},
			},
			assertFunc: func(t *testing.T, impSettings *ImplementationSettings) {
				settings := impSettings.AllSettings()
				require.NotNil(t, settings)
				require.True(t, settings.ContainsRule("my-test-rule"))
				applicationControls, err := impSettings.ApplicableControls("my-test-rule")
				require.NoError(t, err)
				require.Len(t, applicationControls, 1)
				require.Equal(t, applicationControls[0], "ex-1")
				impRequirementSettings, err := impSettings.ByControlID("ex-1")
				require.NoError(t, err)
				require.Contains(t, impRequirementSettings.selectedParameters, "my-test-param")
				require.True(t, impRequirementSettings.ContainsRule("my-test-rule"))
			},
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			testSettings := prepSettings(t)
			testSettings.merge(c.inputImplementation)
			c.assertFunc(t, testSettings)
		})
	}
}

func prepSettings(t *testing.T) *ImplementationSettings {
	testDataPath := filepath.Join("../testdata", "component-definition-test-reqs.json")

	file, err := os.Open(testDataPath)
	require.NoError(t, err)
	definition, err := generators.NewComponentDefinition(file)
	require.NoError(t, err)
	require.NotNil(t, definition)
	var allImplementations []oscalTypes.ControlImplementationSet
	for _, component := range *definition.Components {
		if component.ControlImplementations == nil {
			continue
		}
		allImplementations = append(allImplementations, *component.ControlImplementations...)
	}
	impSettings, err := Framework("cis", allImplementations)
	require.NoError(t, err)
	return impSettings
}
