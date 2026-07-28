package main

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"

	shellquote "github.com/kballard/go-shellquote"

	"github.com/joelmgallant/mattermost-plugin-memes/server/meme"
	"github.com/joelmgallant/mattermost-plugin-memes/server/memelibrary"
)

// jpegQuality matches the quality upstream used when serving memes over HTTP.
const jpegQuality = 90

// errUnknownMeme means the input named no template and matched no pattern.
// Callers turn it into a pointer at /meme help rather than a server error.
var errUnknownMeme = errors.New("i don't know that meme")

// resolveMeme maps command input onto a template and the text for its slots.
// A pattern match wins over an explicit template name, matching upstream, so
// "brace yourself. memes are coming." works without naming a template.
func resolveMeme(input string) (*meme.Template, []string, error) {
	if template, text := memelibrary.PatternMatch(input); template != nil {
		return template, text, nil
	}

	parts, err := shellquote.Split(input)
	if err != nil {
		return nil, nil, fmt.Errorf("i couldn't parse that: %w", err)
	}
	// Upstream indexes parts[0] unguarded. Degenerate input reaches here with
	// nothing in it.
	if len(parts) == 0 {
		return nil, nil, errUnknownMeme
	}
	if template := memelibrary.Template(parts[0]); template != nil {
		return template, parts[1:], nil
	}
	return nil, nil, errUnknownMeme
}

// renderMemeJPEG turns command input into an encoded JPEG. It has no
// Mattermost dependencies, which is what keeps the plugin's own logic
// testable without mocks.
func renderMemeJPEG(input string) (string, []byte, error) {
	template, text, err := resolveMeme(input)
	if err != nil {
		return "", nil, err
	}

	img, err := template.Render(text)
	if err != nil {
		return "", nil, fmt.Errorf("rendering %s: %w", template.Name, err)
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: jpegQuality}); err != nil {
		return "", nil, fmt.Errorf("encoding %s: %w", template.Name, err)
	}
	return template.Name, buf.Bytes(), nil
}
