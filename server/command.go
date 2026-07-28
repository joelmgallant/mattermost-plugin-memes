package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"

	"github.com/joelmgallant/mattermost-plugin-memes/server/memelibrary"
)

// memeCommand is the slash command trigger. Users type /meme.
const memeCommand = "meme"

// availableMemes returns every template name and alias, sorted, so the help
// listing and the autocomplete agree on ordering.
func availableMemes() []string {
	var names []string
	for name, metadata := range memelibrary.Memes() {
		names = append(names, name)
		names = append(names, metadata.Aliases...)
	}
	sort.Strings(names)
	return names
}

// memeCommandData builds the /meme registration with one autocomplete entry
// per template, so typing "/meme " offers the full catalogue.
func memeCommandData() *model.Command {
	root := model.NewAutocompleteData(memeCommand, "[meme-name]", "Create a meme")
	for _, name := range availableMemes() {
		item := model.NewAutocompleteData(name, "", fmt.Sprintf("sends the %s meme", name))
		item.AddTextArgument("text to place on the meme", "[text]", "")
		root.AddCommand(item)
	}

	return &model.Command{
		Trigger:          memeCommand,
		AutoComplete:     true,
		AutoCompleteDesc: "Renders custom memes so you can express yourself with culture.",
		AutocompleteData: root,
	}
}

// helpText is the ephemeral response to /meme and /meme help. Upstream embedded
// a sample image here pointing at its own HTTP route; this build serves no
// routes, so the sample is described rather than rendered.
func helpText() string {
	return strings.Join([]string{
		"You can meme in one of two ways.",
		"",
		"If your meme has well-defined phrasing, just type it:",
		"",
		"`/meme brace yourself. memes are coming.`",
		"",
		"Or name a template and fill its slots:",
		"",
		"`/meme brace-yourselves \"brace yourself.\" \"memes are coming.\"`",
		"",
		"Available memes: " + strings.Join(availableMemes(), ", "),
	}, "\n")
}
