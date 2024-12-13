/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package requirements

import (
	"context"
	"fmt"

	"github.com/oscal-compass/oscal-sdk-go/extensions"
	. "github.com/oscal-compass/oscal-sdk-go/internal/container"
	"github.com/oscal-compass/oscal-sdk-go/rules"
)

// Interface check
var _ rules.Store = (*Store)(nil)

// Store returns rule information with added contexts from a specific
// set of framework requirements.
type Store struct {
	allMappedRules        Set[string]
	allParameterOverrides map[string]string
	baseStore             rules.Store
}

// NewRequirementsStore returns a Store using the given store as base storage and
// processed requirements.
func NewRequirementsStore(store rules.Store, requirements RequirementMap) *Store {
	requirementSet := NewSet[string]()
	allParameters := make(map[string]string)
	for _, requirement := range requirements {
		for _, mappedRule := range requirement.MappedRules {
			requirementSet.Add(mappedRule)
		}
		// TODO(jpower432): Handle key collision.
		for key, value := range requirement.SelectedParameters {
			allParameters[key] = value
		}
	}
	return &Store{
		allMappedRules:        requirementSet,
		allParameterOverrides: allParameters,
		baseStore:             store,
	}
}

func (r *Store) GetByRuleID(ctx context.Context, ruleID string) (extensions.RuleSet, error) {
	ruleSet, err := r.baseStore.GetByRuleID(ctx, ruleID)
	if err != nil {
		return extensions.RuleSet{}, err
	}

	if !r.allMappedRules.Has(ruleID) {
		return extensions.RuleSet{}, fmt.Errorf("rule %s filtered out by requirements", ruleID)
	}

	r.applyParameterSettings(ruleSet)
	return ruleSet, nil
}

func (r *Store) GetByCheckID(ctx context.Context, checkID string) (extensions.RuleSet, error) {
	ruleSet, err := r.baseStore.GetByCheckID(ctx, checkID)
	if err != nil {
		return extensions.RuleSet{}, err
	}
	if !r.allMappedRules.Has(ruleSet.Rule.ID) {
		return extensions.RuleSet{}, fmt.Errorf("rule for check %s filtered out by requirements", checkID)
	}
	r.applyParameterSettings(ruleSet)
	return ruleSet, nil
}

func (r *Store) FindByComponent(ctx context.Context, componentId string) ([]extensions.RuleSet, error) {
	var resolvedRules []extensions.RuleSet
	componentRuleSets, err := r.baseStore.FindByComponent(ctx, componentId)
	if err != nil {
		return []extensions.RuleSet{}, err
	}

	for _, ruleSet := range componentRuleSets {
		if !r.allMappedRules.Has(ruleSet.Rule.ID) {
			continue
		}
		ruleSet = r.applyParameterSettings(ruleSet)
		resolvedRules = append(resolvedRules, ruleSet)
	}
	if len(resolvedRules) == 0 {
		return []extensions.RuleSet{}, fmt.Errorf("no rules found with criteria for component %s", componentId)
	}

	return resolvedRules, nil
}

// applyParameterSettings processes and framework or requirement specific parameter values.
func (r *Store) applyParameterSettings(ruleSet extensions.RuleSet) extensions.RuleSet {
	if len(r.allParameterOverrides) > 0 && ruleSet.Rule.Parameter != nil {
		selectedValue, ok := r.allParameterOverrides[ruleSet.Rule.Parameter.ID]
		if ok {
			parameterCopy := *ruleSet.Rule.Parameter
			parameterCopy.Value = selectedValue
			ruleSet.Rule.Parameter = &parameterCopy
		}
	}
	return ruleSet
}
