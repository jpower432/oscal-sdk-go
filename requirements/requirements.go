/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package requirements

import (
	"path/filepath"
	"strings"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
)

// RequirementSettings represents a specific implementation of an existing RuleSet.
type RequirementSettings struct {
	// MappedRules is a list of rule IDs that are mapped to this requirement.
	MappedRules []string
	// SelectedParameters is a map of parameter names and their selected values for this requirement.
	SelectedParameters map[string]string
}

// RequirementMap is a mapping of requirement IDs to their corresponding Requirement.
type RequirementMap map[string]RequirementSettings

// GetFrameworkShortName returns the human-readable short name for the control source in a
// control implementation set and whether this value is populated.
//
// This function checks the associated properties and falls back to the implementation
// Source reference.
func GetFrameworkShortName(implementation oscalTypes.ControlImplementationSet) (string, bool) {
	// Looks for the property, fallback to parsing it out of the control source href.
	if implementation.Props != nil {
		property, found := extensions.GetTrestleProp(extensions.FrameworkProp, *implementation.Props)
		if found {
			return property.Value, true
		}
	}

	// Fallback to the control source string based on trestle
	// workspace conventions of $MODEL/$MODEL_ID/$MODEL.json.
	var expectedParts, modelIdPos, filePos = 3, 1, 2
	cleanedSource := filepath.Clean(implementation.Source)
	parts := strings.Split(cleanedSource, "/")
	if len(parts) == expectedParts && strings.HasSuffix(parts[filePos], ".json") {
		return parts[modelIdPos], true
	}

	return "", false
}

// ListByFramework returns a mapping of framework short names to a corresponding RequirementMap from
// a list of OSCAL ControlImplementationSets.
func ListByFramework(controlImplementations []oscalTypes.ControlImplementationSet) map[string]RequirementMap {
	requirementsByFramework := make(map[string]RequirementMap)

	for _, controlImplementation := range controlImplementations {
		frameworkShortName, found := GetFrameworkShortName(controlImplementation)
		if !found {
			continue
		}

		parameterOverrides := make(map[string]string)
		if controlImplementation.SetParameters != nil {
			setParameters(*controlImplementation.SetParameters, parameterOverrides)
		}

		requirementsMap, ok := requirementsByFramework[frameworkShortName]
		if !ok {
			requirementsMap = make(RequirementMap)
		}

		for _, implementedReq := range controlImplementation.ImplementedRequirements {
			if implementedReq.Props == nil {
				continue
			}
			mappedRulesProps := extensions.FindAllProps(extensions.RuleIdProp, *implementedReq.Props)
			var mappedRules []string
			for _, mappedRule := range mappedRulesProps {
				mappedRules = append(mappedRules, mappedRule.Value)
			}

			// Create a copy of the implementation level map, we
			// only want to overwrite values for this control
			impReqParams := make(map[string]string)
			for key, value := range parameterOverrides {
				impReqParams[key] = value
			}
			if implementedReq.SetParameters != nil {
				setParameters(*implementedReq.SetParameters, impReqParams)
			}

			requirement := RequirementSettings{
				MappedRules:        mappedRules,
				SelectedParameters: impReqParams,
			}
			requirementsMap[implementedReq.ControlId] = requirement
		}
		requirementsByFramework[frameworkShortName] = requirementsMap
	}

	return requirementsByFramework
}

// setParameters updates the paramMap with the input list of SetParameters.
func setParameters(parameters []oscalTypes.SetParameter, paramMap map[string]string) {
	for _, prm := range parameters {
		// Parameter values set for trestle Rule selection
		// should only map to a single value.
		if len(prm.Values) != 1 {
			continue
		}
		paramMap[prm.ParamId] = prm.Values[0]
	}
}
