/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package generators

import (
	"bytes"

	oscal112 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"

	"github.com/oscal-compass/oscal-sdk-go/validation"
)

// NewCatalogWithValidation creates a new OSCAL-based control catalog with validation.
func NewCatalogWithValidation(data []byte, oscalValidator validation.Handler) (catalog *oscal112.Catalog, err error) {
	_, err = oscalValidator.Validate(data)
	if err != nil {
		return catalog, err
	}
	reader := bytes.NewReader(data)
	return NewCatalog(reader)
}

// NewProfileWithValidation creates a new OSCAL-based profile with validation.
func NewProfileWithValidation(data []byte, oscalValidator validation.Handler) (profile *oscal112.Profile, err error) {
	_, err = oscalValidator.Validate(data)
	if err != nil {
		return profile, err
	}
	reader := bytes.NewReader(data)
	return NewProfile(reader)
}

// NewComponentDefinitionWithValidation creates a new OSCAL-based component definition with validation.
func NewComponentDefinitionWithValidation(data []byte, oscalValidator validation.Handler) (componentDefinition *oscal112.ComponentDefinition, err error) {
	_, err = oscalValidator.Validate(data)
	if err != nil {
		return componentDefinition, err
	}
	reader := bytes.NewReader(data)
	return NewComponentDefinition(reader)
}

// NewSystemSecurityPlanWithValidation creates a new OSCAL-based system security plan with validation.
func NewSystemSecurityPlanWithValidation(data []byte, oscalValidator validation.Handler) (ssp *oscal112.SystemSecurityPlan, err error) {
	_, err = oscalValidator.Validate(data)
	if err != nil {
		return ssp, err
	}
	reader := bytes.NewReader(data)
	return NewSystemSecurityPlan(reader)
}

// NewAssessmentPlanWithValidation creates a new OSCAL-based assessment plan with validation.
func NewAssessmentPlanWithValidation(data []byte, oscalValidator validation.Handler) (assessmentPlan *oscal112.AssessmentPlan, err error) {
	_, err = oscalValidator.Validate(data)
	if err != nil {
		return assessmentPlan, err
	}
	reader := bytes.NewReader(data)
	return NewAssessmentPlan(reader)
}

// NewAssessmentResultsWithValidation creates a new OSCAL-based assessment results set with validation.
func NewAssessmentResultsWithValidation(data []byte, oscalValidator validation.Handler) (assessmentResults *oscal112.AssessmentResults, err error) {
	_, err = oscalValidator.Validate(data)
	if err != nil {
		return assessmentResults, err
	}
	reader := bytes.NewReader(data)
	return NewAssessmentResults(reader)
}

// NewPOAMWithValidation creates a new OSCAL-based plan of action and milestones with validation.
func NewPOAMWithValidation(data []byte, oscalValidator validation.Handler) (pOAM *oscal112.PlanOfActionAndMilestones, err error) {
	_, err = oscalValidator.Validate(data)
	if err != nil {
		return pOAM, err
	}
	reader := bytes.NewReader(data)
	return NewPOAM(reader)
}
