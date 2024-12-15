/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package requirements

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/stretchr/testify/require"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/generators"
	"github.com/oscal-compass/oscal-sdk-go/rules"
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

func TestComponents(t *testing.T) {
	tests := []struct {
		name               string
		requirements       RequirementMap
		componentID        string
		expError           string
		wantRules          []extensions.RuleSet
		postValidationFunc func(store rules.Store) bool
	}{
		{
			name:        "Valid/WithMappedRules",
			componentID: "testComponent1",
			requirements: RequirementMap{
				"ex-1": RequirementSettings{
					MappedRules: []string{"testRule1", "testRule2"},
				},
			},
			wantRules: []extensions.RuleSet{testSet2},
		},
		{
			name:        "Valid/WithParameterOverrides",
			componentID: "testComponent2",
			requirements: RequirementMap{
				"ex-1": RequirementSettings{
					MappedRules: []string{"testRule1", "testRule2"},
					SelectedParameters: map[string]string{
						"testParam1": "updatedValue",
					},
				},
			},
			wantRules: []extensions.RuleSet{
				{
					Rule: extensions.Rule{
						ID:          "testRule1",
						Description: "Test Rule",
						Parameter: &extensions.Parameter{
							ID:          "testParam1",
							Description: "Test Parameter",
							Value:       "updatedValue",
						},
					},
					Checks: []extensions.Check{
						{
							ID:          "testCheck1",
							Description: "Test Check",
						},
					},
				},
				testSet2,
			},
			postValidationFunc: func(store rules.Store) bool {
				ruleSet, _ := store.GetByRuleID(context.TODO(), "testRule1")
				return ruleSet.Rule.Parameter != nil && ruleSet.Rule.Parameter.Value == ""
			},
		},
		{
			name:        "Invalid/InvalidSettings",
			componentID: "testComponent1",
			requirements: RequirementMap{
				"ex-1": RequirementSettings{
					MappedRules: []string{"doesnotexist"},
				},
			},
			expError: "no rules found with criteria for component testComponent1",
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			testCtx := context.Background()
			baseStore := newFakeStore()

			foundRules, err := Components(testCtx, c.componentID, baseStore, c.requirements)
			sort.SliceStable(foundRules, func(i, j int) bool {
				return foundRules[i].Rule.ID < foundRules[j].Rule.ID
			})

			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.wantRules, foundRules)
			}

			if c.postValidationFunc != nil {
				require.True(t, c.postValidationFunc(baseStore))
			}
		})
	}
}

var (
	testSet1 = extensions.RuleSet{
		Rule: extensions.Rule{
			ID:          "testRule1",
			Description: "Test Rule",
			Parameter: &extensions.Parameter{
				ID:          "testParam1",
				Description: "Test Parameter",
			},
		},
		Checks: []extensions.Check{
			{
				ID:          "testCheck1",
				Description: "Test Check",
			},
		},
	}
	testSet2 = extensions.RuleSet{
		Rule: extensions.Rule{
			ID:          "testRule2",
			Description: "Test Rule",
		},
		Checks: []extensions.Check{
			{
				ID:          "testCheck2",
				Description: "Test Check",
			},
		},
	}
	testSet3 = extensions.RuleSet{
		Rule: extensions.Rule{
			ID:          "testRule3",
			Description: "Test Rule",
			Parameter: &extensions.Parameter{
				ID:          "testParam3",
				Description: "Test Parameter",
				Value:       "default",
			},
		},
		Checks: []extensions.Check{
			{
				ID:          "testCheck3",
				Description: "Test Check",
			},
		},
	}
)

// fakeStore is a fake implementation of a rules.Store with static data
type fakeStore struct {
	staticRuleData map[string]extensions.RuleSet
}

func newFakeStore() *fakeStore {
	return &fakeStore{
		staticRuleData: map[string]extensions.RuleSet{
			"testRule1": testSet1,
			"testRule2": testSet2,
			"testRule3": testSet3,
		},
	}
}

func (f fakeStore) GetByRuleID(ctx context.Context, ruleID string) (extensions.RuleSet, error) {
	ruleSet, ok := f.staticRuleData[ruleID]
	if !ok {
		return extensions.RuleSet{}, fmt.Errorf("rule %s not found", ruleID)
	}
	return ruleSet, nil
}

func (f fakeStore) GetByCheckID(ctx context.Context, checkID string) (extensions.RuleSet, error) {
	switch checkID {
	case "testCheck1":
		return f.staticRuleData["testRule1"], nil
	case "testCheck2":
		return f.staticRuleData["testRule2"], nil
	case "testCheck3":
		return f.staticRuleData["testRule3"], nil
	default:
		return extensions.RuleSet{}, fmt.Errorf("rule not found for %s", checkID)
	}
}

func (f fakeStore) FindByComponent(ctx context.Context, componentId string) ([]extensions.RuleSet, error) {
	switch componentId {
	case "testComponent1":
		return []extensions.RuleSet{testSet2, testSet3}, nil
	case "testComponent2":
		return []extensions.RuleSet{testSet1, testSet2}, nil
	default:
		return []extensions.RuleSet{}, fmt.Errorf("invalid component id: %s", componentId)
	}
}
