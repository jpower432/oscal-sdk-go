/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	oscalValidation "github.com/defenseunicorns/go-oscal/src/pkg/validation"
)

var _ Handler = (*SchemaValidator)(nil)

/*
SchemaValidator implementation a validation.Handler and
validates the OSCAL documents with OSCAL JSON schema using `go-oscal`.
*/
type SchemaValidator struct {
	oscalVersion string
}

// NewSchemaValidator returns a new versioned SchemaValidator.
func NewSchemaValidator(version string) *SchemaValidator {
	return &SchemaValidator{oscalVersion: version}
}

func (s *SchemaValidator) Validate(data []byte) (bool, error) {
	validator, err := oscalValidation.NewValidatorDesiredVersion(data, s.oscalVersion)
	if err != nil {
		return false, err
	}

	modelType := validator.GetModelType()
	err = validator.Validate()
	if err != nil {
		return false, &ErrValidation{Type: "schema", Model: modelType, Err: err}
	}

	return true, nil
}
