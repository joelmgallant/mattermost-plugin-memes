package main

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAvailableMemesIsSortedAndIncludesAliases(t *testing.T) {
	names := availableMemes()
	require.Greater(t, len(names), 55, "aliases should push the count past the 55 templates")
	assert.True(t, sort.StringsAreSorted(names), "names must be sorted for a stable help listing")
	assert.Contains(t, names, "brace-yourselves")
}

func TestMemeCommandData(t *testing.T) {
	command := memeCommandData()
	assert.Equal(t, "meme", command.Trigger)
	assert.True(t, command.AutoComplete)
	require.NotNil(t, command.AutocompleteData)
	// Pinned to the literal count rather than len(availableMemes()): asserting
	// against the same call the implementation uses is circular, and the
	// literal also locks in the alias-dedupe fix (99 unique names, not 101).
	assert.Len(t, command.AutocompleteData.SubCommands, 99)
}

// The help text must not embed an image: this build serves no HTTP routes, so
// any URL it printed would 404.
func TestHelpTextCarriesNoImage(t *testing.T) {
	text := helpText()
	assert.Contains(t, text, "Available memes:")
	assert.Contains(t, text, "brace-yourselves")
	assert.NotContains(t, text, "![", "help must not embed a markdown image")
	assert.NotContains(t, strings.ToLower(text), "/plugins/", "help must not reference a plugin URL")
}
