package memelibrary

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"image/jpeg"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEveryTemplateRenders guards the embed loader end to end. A template
// whose image or metadata failed to load would surface here rather than as a
// blank meme in production.
func TestEveryTemplateRenders(t *testing.T) {
	memes := Memes()
	require.Len(t, memes, 55, "expected 55 meme templates")

	for name := range memes {
		t.Run(name, func(t *testing.T) {
			template := Template(name)
			require.NotNil(t, template)
			require.NotNil(t, template.Image, "template has no image")
			require.NotEmpty(t, template.TextSlots, "template has no text slots")

			img, err := template.Render([]string{"top text", "bottom text"})
			require.NoError(t, err)
			require.NotNil(t, img)

			var buf bytes.Buffer
			require.NoError(t, jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}))
			assert.Greater(t, buf.Len(), 1024, "encoded meme is implausibly small")
		})
	}
}

// TestFontsLoaded asserts the three bundled faces are present. A missing font
// does not fail rendering, it silently produces text in the wrong face.
func TestFontsLoaded(t *testing.T) {
	for _, name := range []string{"Anton-Regular", "ComicNeue-Bold", "ComicNeue-Regular"} {
		assert.NotNil(t, fonts[name], "font %s not loaded", name)
	}
}

// TestAssetsArePinned checksums the source assets rather than rendered output.
// A golden hash over an encoded JPEG would be brittle — Go's image/jpeg has
// changed its output across releases — but the input files must never move.
func TestAssetsArePinned(t *testing.T) {
	pinned := map[string]string{
		"images/brace-yourselves.jpg": "53d3aa4cd7788dbfd1f2b8dccee1e72151ebd3f3cfca5a73da16729726b26b70",
		"fonts/Anton-Regular.ttf":     "83be67769f0287a34b25ff70297b58ef1c0b259939cbea11a0768204237834db",
	}
	for name, want := range pinned {
		got := fmt.Sprintf("%x", sha256.Sum256(MustAsset(name)))
		assert.Equal(t, want, got, "asset %s changed", name)
	}
}
