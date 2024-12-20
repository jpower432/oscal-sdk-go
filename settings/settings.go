/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package settings

import (
	"context"
	"fmt"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/rules"
)

// Settings defines methods for tuning extensions.RuleSet options.
type Settings interface {
	// ApplyParameterSettings returns the given rule set with update parameter values based on the implementation.
	//
	// If the implementation does have parameter values or the rule set does not have a parameter, the original rule set
	// is returned.
	// The parameter value is not altered on the original rule set, it is copied and returned with the new rule set.
	ApplyParameterSettings(set extensions.RuleSet) extensions.RuleSet

	// ContainsRule returns whether the given rule id is defined in the Settings.
	ContainsRule(ruleId string) bool
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
