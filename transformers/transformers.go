/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package transformers

import (
	"context"
	"fmt"

	oscaltypes112 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"

	"github.com/oscal-compass/oscal-sdk-go/models/components"
	"github.com/oscal-compass/oscal-sdk-go/models/plans"
	"github.com/oscal-compass/oscal-sdk-go/settings"
)

// ComponentDefinitionsToAssessmentPlan transforms the data from one or more OSCAL Component Definitions to a single OSCAL Assessment Plan.
func ComponentDefinitionsToAssessmentPlan(ctx context.Context, definitions []oscaltypes112.ComponentDefinition, framework string) (*oscaltypes112.AssessmentPlan, error) {
	// Collect and aggregate all component information for each component definition
	var allComponents []components.Component
	var allImplementations []oscaltypes112.ControlImplementationSet
	for _, compDef := range definitions {
		if compDef.Components == nil {
			continue
		}
		for _, comp := range *compDef.Components {
			if comp.ControlImplementations == nil {
				continue
			}
			componentAdapter := components.NewDefinedComponentAdapter(comp)
			allComponents = append(allComponents, componentAdapter)
			allImplementations = append(allImplementations, *comp.ControlImplementations...)
		}
	}
	implementationSettings, err := settings.Framework(framework, allImplementations)
	if err != nil || implementationSettings == nil {
		return nil, fmt.Errorf("cannot transform definitions for framework %s: %w", framework, err)
	}
	return plans.GenerateAssessmentPlan(ctx, allComponents, *implementationSettings)
}
