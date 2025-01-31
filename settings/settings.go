/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package settings

import (
	"context"
	"fmt"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/internal/set"
	"github.com/oscal-compass/oscal-sdk-go/rules"
)

// Settings defines settings for RuleSets to tune options.
type Settings struct {
	// mappedRules is a list of rule IDs that are mapped to this requirement.
	mappedRules set.Set[string]
	// selectedParameters is a map of parameter names and their selected values for this requirement.
	selectedParameters map[string]string
}

// NewSettingsFromImplementedRequirement returns an initialized Settings from an
// OSCAL Implemented Requirement.
func NewSettingsFromImplementedRequirement(implementedReq oscalTypes.ImplementedRequirementControlImplementation) Settings {
	requirement := Settings{
		selectedParameters: make(map[string]string),
		mappedRules:        set.New[string](),
	}

	if implementedReq.Props != nil {
		mappedRulesProps := extensions.FindAllProps(extensions.RuleIdProp, *implementedReq.Props)
		for _, mappedRule := range mappedRulesProps {
			requirement.mappedRules.Add(mappedRule.Value)
		}
	}

	if implementedReq.SetParameters != nil {
		setParameters(*implementedReq.SetParameters, requirement.selectedParameters)
	}
	return requirement
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

// ApplyParameterSettings returns the given rule set with update parameter values based on the implementation.
//
// If the implementation does have parameter values or the rule set does not have a parameter, the original rule set
// is returned.
// The parameter value is not altered on the original rule set, it is copied and returned with the new rule set.
func (i Settings) ApplyParameterSettings(set extensions.RuleSet) extensions.RuleSet {
	if len(i.selectedParameters) > 0 && set.Rule.Parameter != nil {
		selectedValue, ok := i.selectedParameters[set.Rule.Parameter.ID]
		if ok {
			parameterCopy := *set.Rule.Parameter
			parameterCopy.Value = selectedValue
			set.Rule.Parameter = &parameterCopy
		}
	}
	return set
}

// ContainsRule returns whether the given rule id is defined in the Settings.
func (i Settings) ContainsRule(ruleId string) bool {
	return i.mappedRules.Has(ruleId)
}

// ApplyToComponent returns a list of RuleSets for a given component with options applied from the given Settings.
//
// Only the rules that overlap between the component and the mapped rules in the implementation are returned.
// Implementation-level parameters will be applied as RuleSet selected parameter values.
func ApplyToComponent(ctx context.Context, componentId string, store rules.Store, settings Settings) ([]extensions.RuleSet, error) {
	var resolvedRules []extensions.RuleSet
	componentRuleSets, err := store.FindByComponent(ctx, componentId)
	if err != nil {
		return []extensions.RuleSet{}, err
	}

	for _, ruleSet := range componentRuleSets {
		if !settings.ContainsRule(ruleSet.Rule.ID) {
			continue
		}
		ruleSet = settings.ApplyParameterSettings(ruleSet)
		resolvedRules = append(resolvedRules, ruleSet)
	}
	if len(resolvedRules) == 0 {
		return []extensions.RuleSet{}, fmt.Errorf("no rules found with criteria for component %s", componentId)
	}
	return resolvedRules, nil
}
