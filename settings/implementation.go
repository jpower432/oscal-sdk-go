/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package settings

import (
	"fmt"
	"path/filepath"
	"strings"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/internal/set"
)

// GetFrameworkShortName returns the human-readable short name for the control source in a
// control implementation set and whether this value is populated.
//
// This function checks the associated properties and falls back to the implementation
// Source reference.
func GetFrameworkShortName(implementation oscalTypes.ControlImplementationSet) (string, bool) {
	const (
		expectedPathParts = 3
		modelIDIndex      = 1
		filenameIndex     = 2
	)
	// Looks for the property, fallback to parsing it out of the control source href.
	if implementation.Props != nil {
		property, found := extensions.GetTrestleProp(extensions.FrameworkProp, *implementation.Props)
		if found {
			return property.Value, true
		}
	}

	// Fallback to the control source string based on trestle
	// workspace conventions of $MODEL/$MODEL_ID/$MODEL.json.
	cleanedSource := filepath.Clean(implementation.Source)
	parts := strings.Split(cleanedSource, "/")
	if len(parts) == expectedPathParts && strings.HasSuffix(parts[filenameIndex], ".json") {
		return parts[modelIDIndex], true
	}

	return "", false
}

// ImplementationSettings defines settings for RuleSets defined at the control
// implementation level.
type ImplementationSettings struct {
	// requirementSettings define settings for RuleSets at the
	// implemented requirement/individual control level.
	requirementSettings map[string]RequirementSettings
	// parameterOverrides define parameters for RuleSets defined at the
	// control implementation.
	parameterOverrides map[string]string
	// allMappedRules
	allMappedRules set.Set[string]
}

// NewImplementationSettings returned ImplementationSettings for an OSCAL ControlImplementationSet.
func NewImplementationSettings(controlImplementation oscalTypes.ControlImplementationSet) *ImplementationSettings {
	implementation := &ImplementationSettings{
		allMappedRules:      set.New[string](),
		requirementSettings: make(map[string]RequirementSettings),
		parameterOverrides:  make(map[string]string),
	}
	if controlImplementation.SetParameters != nil {
		setParameters(*controlImplementation.SetParameters, implementation.parameterOverrides)
	}

	for _, implementedReq := range controlImplementation.ImplementedRequirements {
		requirement := NewRequirementSettings(implementedReq)
		for mappedRule := range requirement.mappedRules {
			implementation.allMappedRules.Add(mappedRule)
		}
		implementation.requirementSettings[implementedReq.ControlId] = requirement
	}

	return implementation
}

// ApplyParameterSettings returns the given rule set with update parameter values based on the implementation.
//
// If the implementation does have parameter values or the rule set does not have a parameter, the original rule set
// is returned.
// The parameter value is not altered on the original rule set, it is copied and returned with the new rule set.
func (i *ImplementationSettings) ApplyParameterSettings(set extensions.RuleSet) extensions.RuleSet {
	if len(i.parameterOverrides) > 0 && set.Rule.Parameter != nil {
		selectedValue, ok := i.parameterOverrides[set.Rule.Parameter.ID]
		if ok {
			parameterCopy := *set.Rule.Parameter
			parameterCopy.Value = selectedValue
			set.Rule.Parameter = &parameterCopy
		}
	}
	return set
}

// ContainsRule returns whether the given rule set is mapped to the implementation.
func (i *ImplementationSettings) ContainsRule(set extensions.RuleSet) bool {
	return i.allMappedRules.Has(set.Rule.ID)
}

// Framework returns ImplementationSettings from a list of OSCAL Control Implementations for a given framework.
func Framework(framework string, controlImplementations []oscalTypes.ControlImplementationSet) (*ImplementationSettings, error) {
	var requirements *ImplementationSettings

	for _, controlImplementation := range controlImplementations {
		frameworkShortName, found := GetFrameworkShortName(controlImplementation)
		if found && frameworkShortName == framework {
			requirements = NewImplementationSettings(controlImplementation)
			break
		}
	}

	if requirements == nil {
		return requirements, fmt.Errorf("framework %s is not in control implementations", framework)
	}
	return requirements, nil
}
