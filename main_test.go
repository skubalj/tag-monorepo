/*
tag-monorepo: Per-module Monorepo Tagging
Copyright (C) 2026 Joseph Skubal

See the COPYING file for copyright information
*/

package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getTags(t *testing.T) {
	t.SkipNow() // Manual test

	rows, err := getTags("/tmp/repo")
	require.NoError(t, err)
	fmt.Println(rows)
}
