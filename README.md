# Mattermost Memes Plugin (attachment build)

A rebuild of [mattermost-community/mattermost-plugin-memes](https://github.com/mattermost-community/mattermost-plugin-memes)
that posts memes as **file attachments** instead of markdown image URLs, so they
render on the mobile apps (verified on iOS; Android uses the same attachment path).

```
/meme brace-yourselves "brace yourself." "memes are coming."
/meme memes. memes everywhere.
```

## Why this fork exists

Upstream returns a post whose body is a markdown image pointing at a URL the
plugin renders on demand. The Mattermost mobile clients fall back to rendering a
bare link when a post carries no image metadata, and a plugin returning a
`CommandResponse` cannot influence that metadata — so memes have never appeared
on iOS or Android. That is upstream issue #41, open since 2020.

This build renders the meme, uploads it with the file API, and attaches it to the
post. Attachments carry their own dimensions and need no post metadata, so a meme
renders wherever an ordinary uploaded photo does.

Two things follow from the change:

- **No HTTP routes.** Upstream's handler never checked `Mattermost-User-Id`,
  which left an unauthenticated on-demand image renderer exposed to the internet.
  This build registers no routes at all.
- **Memes are real files.** They appear in the channel's file list and in search,
  and cost roughly 20–100 KB each (median ~50 KB).

## Install

1. Download the bundle from [Releases](https://github.com/joelmgallant/mattermost-plugin-memes/releases).
2. Upload it in **System Console → Plugins → Plugin Management**, which requires
   **Enable Plugin Uploads** to be set to true.
3. Enable the plugin. There is nothing to configure.

Requires Mattermost v9.0.0 or later.

## Build

```bash
make test   # gofmt gate, then go test ./server/...
make dist   # linux-amd64 bundle in dist/
```

Preview a meme without a server:

```bash
go run ./server -out demo.jpg 'memes. memes everywhere'
```

## Differences from upstream

| | Upstream | This build |
|---|---|---|
| Delivery | Markdown image → plugin URL | Uploaded file attachment |
| Mobile | Does not render | Renders (verified on iOS) |
| HTTP routes | One, unauthenticated | None |
| Assets | 7.9 MB go-bindata blob | `//go:embed`, 2.1 MB raw |
| SDK | `mattermost-server/v5` | `mattermost/server/public` |
| Bundle | 34 MB, three architectures | ~9.0 MB, linux-amd64 |
| Unknown meme | `AppError` | Ephemeral message |

## Licence and attribution

Apache-2.0, inherited from upstream. See `LICENSE` and `NOTICE`.

The bundled meme templates are third-party photographs redistributed from the
upstream project, which has an open issue about their licensing
([#20](https://github.com/mattermost-community/mattermost-plugin-memes/issues/20)).
They are included here on the same basis. If you are deploying somewhere that
matters, satisfy yourself about that first.
