/*
 Copyright 2025 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package components

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/oscal-compass/oscal-sdk-go/generators"
)

func TestSystemComponentAdapter(t *testing.T) {
	testDataPath := filepath.Join("../../testdata", "test-ssp.json")

	file, err := os.Open(testDataPath)
	require.NoError(t, err)
	ssp, err := generators.NewSystemSecurityPlan(file)
	require.NoError(t, err)
	require.NotNil(t, ssp)

	require.Len(t, ssp.SystemImplementation.Components, 2)
	adapter := NewSystemComponentAdapter(ssp.SystemImplementation.Components[0])
	require.Equal(t, "Example Service", adapter.Title())
	require.Equal(t, Service, adapter.Type())
	require.Equal(t, "4e19131e-b361-4f0e-8262-02bf4456202e", adapter.UUID())
	require.Len(t, adapter.Props(), 7)
	systemComp, ok := adapter.AsSystemComponent()
	require.True(t, ok)
	require.Equal(t, adapter.UUID(), systemComp.UUID)
	definedComp, ok := adapter.AsDefinedComponent()
	require.True(t, ok)
	require.Equal(t, adapter.UUID(), definedComp.UUID)
}

// TODO(jpower432): Fill in tests
func TestControlImplementationAdapter(t *testing.T) {

}
