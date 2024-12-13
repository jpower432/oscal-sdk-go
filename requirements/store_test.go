package requirements

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/rules"
)

func TestStore_GetByRuleID(t *testing.T) {
	baseStore := newFakeStore()
	exampleRequirements := RequirementMap{
		"ex-1": RequirementSettings{
			MappedRules: []string{"testRule1"},
		},
	}
	testStore := NewRequirementsStore(baseStore, exampleRequirements)
	ruleSet, err := testStore.GetByRuleID(context.Background(), "testRule1")
	require.NoError(t, err)
	require.Equal(t, testSet1, ruleSet)

	_, err = testStore.GetByRuleID(context.Background(), "notARule")
	require.EqualError(t, err, "rule notARule not found")

	_, err = testStore.GetByRuleID(context.Background(), "testRule2")
	require.EqualError(t, err, "rule testRule2 filtered out by requirements")
}

func TestStore_GetByCheckID(t *testing.T) {
	baseStore := newFakeStore()
	exampleRequirements := RequirementMap{
		"ex-1": RequirementSettings{
			MappedRules: []string{"testRule1", "testRule2"},
		},
	}
	testStore := NewRequirementsStore(baseStore, exampleRequirements)
	ruleSet, err := testStore.GetByCheckID(context.Background(), "testCheck2")
	require.NoError(t, err)
	require.Equal(t, testSet2, ruleSet)

	_, err = testStore.GetByCheckID(context.Background(), "notACheck")
	require.EqualError(t, err, "rule not found for notACheck")

	_, err = testStore.GetByCheckID(context.Background(), "testCheck3")
	require.EqualError(t, err, "rule for check testCheck3 filtered out by requirements")
}

func TestStore_FindByComponent(t *testing.T) {
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
			testStore := NewRequirementsStore(baseStore, c.requirements)

			foundRules, err := testStore.FindByComponent(testCtx, c.componentID)
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
