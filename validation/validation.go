/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

/*
Package validation defines mechanisms for validation OSCAL objects in various ways.
This could include implementations like json schema validation or constraints validations.
*/
package validation

import (
	"errors"
	"fmt"
)

// ErrValidation is returned when data is not valid.
type ErrValidation struct {
	Type  string
	Model string
	Err   error
}

func (e *ErrValidation) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Err.Error())
}

// ErrStopValidation signals the validation handler to stop.
var ErrStopValidation = errors.New("stop validation")

// Handler execute validation on an OSCAL model.
type Handler interface {
	// Validate will take input data perform validation.
	// Return values will represent whether the data was valid and
	// if not, what errors occurred.
	Validate(data []byte) (bool, error)
}

// ValidateFunc function implements the Handler interface.
type ValidateFunc func(data []byte) (bool, error)

func (fn ValidateFunc) Validate(data []byte) (bool, error) {
	return fn(data)
}

// Handlers returns a handler that will run the validation handlers in sequence.
func Handlers(handlers ...Handler) ValidateFunc {
	return func(data []byte) (bool, error) {
		var errs []error
		for _, handler := range handlers {
			_, err := handler.Validate(data)
			if err != nil {
				if errors.Is(err, ErrStopValidation) {
					break
				}
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			return false, errors.Join(errs...)
		}
		return true, nil
	}
}
