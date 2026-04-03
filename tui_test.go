/*
tag-monorepo: Per-module Monorepo Tagging
Copyright (C) 2026 Joseph Skubal

See the COPYING file for copyright information
*/

package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	require.Equal(t, Update_Major, Update_None.Increment())
	require.Equal(t, Update_None, Update_Patch.Increment())
	require.Equal(t, Update_Minor, Update_Patch.Decrement())
	require.Equal(t, Update_Patch, Update_None.Decrement())
}

func TestVersion(t *testing.T) {
	require.Equal(t, "v1.0.0", Version{Major: 1}.String())
	require.Equal(t, "v5.4.0", Version{Major: 5, Minor: 4}.String())
	require.Equal(t, "v5.4.8-beta", Version{Major: 5, Minor: 4, Patch: 8, Suffix: "beta"}.String())
	require.Equal(t, "v0.0.6", Version{Patch: 6}.String())
	require.Equal(t, "v0.0.4-dev", Version{Patch: 4, Suffix: "dev"}.String())
	require.Equal(t, "v2.0.3", Version{Major: 2, Patch: 3}.String())
}

func TestParseVersion(t *testing.T) {
	x, ok := ParseVersion("1.2.3")
	require.False(t, ok)

	x, ok = ParseVersion("v1.2.3")
	require.True(t, ok)
	require.Equal(t, Version{Major: 1, Minor: 2, Patch: 3}, x)

	x, ok = ParseVersion("v0.2.4-beta")
	require.True(t, ok)
	require.Equal(t, Version{Major: 0, Minor: 2, Patch: 4, Suffix: "beta"}, x)

	x, ok = ParseVersion("v8.2-beta")
	require.True(t, ok)
	require.Equal(t, Version{Major: 8, Minor: 2, Patch: 0, Suffix: "beta"}, x)

	x, ok = ParseVersion("v0-beta")
	require.True(t, ok)
	require.Equal(t, Version{Major: 0, Minor: 0, Patch: 0, Suffix: "beta"}, x)

	x, ok = ParseVersion("v6.0.0-rc2")
	require.True(t, ok)
	require.Equal(t, Version{Major: 6, Minor: 0, Patch: 0, Suffix: "rc2"}, x)
}

func TestCompare(t *testing.T) {
	require.Positive(t, Version{Major: 1}.Compare(Version{Major: 0}))
	require.Positive(t, Version{Major: 1, Minor: 2}.Compare(Version{Major: 0, Minor: 9, Patch: 12}))
	require.Negative(t, Version{Major: 0, Minor: 2}.Compare(Version{Major: 0, Minor: 9, Patch: 12}))
	require.Positive(t, Version{Major: 0, Minor: 2, Patch: 1}.Compare(Version{Major: 0, Minor: 2, Suffix: "beta"}))
	require.Positive(t, Version{Major: 0, Minor: 2, Patch: 14}.Compare(Version{Major: 0, Minor: 1, Patch: 12}))
}
