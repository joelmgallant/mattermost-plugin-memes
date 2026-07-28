package main

import (
	"bytes"
	"image/jpeg"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderMemeJPEGByTemplateName(t *testing.T) {
	name, data, err := renderMemeJPEG(`brace-yourselves "brace yourself." "memes are coming."`)
	require.NoError(t, err)
	assert.Equal(t, "brace-yourselves", name)

	img, err := jpeg.Decode(bytes.NewReader(data))
	require.NoError(t, err, "output must be a decodable JPEG")
	assert.Greater(t, img.Bounds().Dx(), 0)
}

func TestRenderMemeJPEGByPattern(t *testing.T) {
	name, data, err := renderMemeJPEG("brace yourself. memes are coming.")
	require.NoError(t, err)
	assert.Equal(t, "brace-yourselves", name, "pattern match should win without naming a template")
	assert.NotEmpty(t, data)
}

func TestRenderMemeJPEGUnknownTemplate(t *testing.T) {
	_, _, err := renderMemeJPEG("definitely-not-a-real-meme one two")
	require.ErrorIs(t, err, errUnknownMeme)
}

// Upstream indexes parts[0] with no length check. Quote-only input reaches that
// line with an empty slice in some shellquote versions, so guard it explicitly.
func TestRenderMemeJPEGDegenerateInput(t *testing.T) {
	for _, input := range []string{`""`, `''`, "   ", `"`} {
		t.Run(input, func(t *testing.T) {
			assert.NotPanics(t, func() {
				_, _, _ = renderMemeJPEG(input)
			})
		})
	}
}

func TestRenderMemeJPEGUnbalancedQuotes(t *testing.T) {
	_, _, err := renderMemeJPEG(`doge "unterminated`)
	require.Error(t, err)
	assert.NotErrorIs(t, err, errUnknownMeme, "a parse failure is not an unknown meme")
}
