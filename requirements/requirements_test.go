/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package requirements

import (
	"os"
	"path/filepath"
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/stretchr/testify/require"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/generators"
)

func TestGetFrameworkShortName(t *testing.T) {
	tests := []struct {
		name                string
		inputImplementation oscalTypes.ControlImplementationSet
		wantName            string
		wantFound           bool
	}{
		{
			name: "Valid/ShortNameFromProps",
			inputImplementation: oscalTypes.ControlImplementationSet{
				Props: &[]oscalTypes.Property{
					{
						Name:  extensions.FrameworkProp,
						Value: "propFramework",
						Ns:    extensions.TrestleNameSpace,
					},
				},
				Source: "profiles/framework/profile.json",
			},
			wantName:  "propFramework",
			wantFound: true,
		},
		{
			name: "Valid/ShortNameFromSource",
			inputImplementation: oscalTypes.ControlImplementationSet{
				Source: "profiles/sourceFramework/profile.json",
			},
			wantName:  "sourceFramework",
			wantFound: true,
		},
		{
			name:                "Valid/NoShortName",
			inputImplementation: oscalTypes.ControlImplementationSet{},
			wantName:            "",
			wantFound:           false,
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			name, found := GetFrameworkShortName(c.inputImplementation)
			require.Equal(t, c.wantName, name)
			require.Equal(t, c.wantFound, found)
		})
	}

}

func TestListByFramework(t *testing.T) {
	testDataPath := filepath.Join("../testdata", "component-definition-test-reqs.json")
	file, err := os.Open(testDataPath)
	require.NoError(t, err)
	definition, err := generators.NewComponentDefinition(file)
	require.NoError(t, err)

	require.NotNil(t, definition.Components)

	var allImplementations []oscalTypes.ControlImplementationSet
	for _, component := range *definition.Components {
		if component.ControlImplementations == nil {
			continue
		}
		allImplementations = append(allImplementations, *component.ControlImplementations...)
	}

	implementationsMap := ListByFramework(allImplementations)
	expectedMap := map[string]RequirementMap{
		"cis": {
			"CIS-2.1": {
				MappedRules: []string{
					"etcd_cert_file",
					"etcd_key_file",
				},
				SelectedParameters: map[string]string{},
			},
		},
		"example": {
			"ex-1": {
				MappedRules: []string{
					"etcd_key_file",
				},
				SelectedParameters: map[string]string{
					"temperature_tolerance": "10%",
				},
			},
		},
	}

	require.Equal(t, expectedMap, implementationsMap)
}
