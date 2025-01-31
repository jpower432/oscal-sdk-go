/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package generators

import (
	"time"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
)

const (
	// OSCALVersion is the string that represents the version of OSCAL supported.
	OSCALVersion = "1.1.2"
	// SampleRequiredString is the default string for required string data in
	// OSCAL models. This value matched the default value used in compliance-trestle.
	SampleRequiredString = "REPLACE_ME"
	defaultVersion       = "0.1.0"
)

// NewSampleMetadata returns OSCAL Metadata with default values for all required
// fields.
func NewSampleMetadata() oscalTypes.Metadata {
	return oscalTypes.Metadata{
		Title:        SampleRequiredString,
		LastModified: time.Now(),
		OscalVersion: OSCALVersion,
		Version:      defaultVersion,
	}
}
