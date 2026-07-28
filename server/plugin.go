package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

// Plugin renders memes and posts them as file attachments.
//
// Upstream posted a markdown image pointing at a route the plugin served on
// demand. The mobile clients fall back to rendering a bare link when a post
// carries no image metadata, and a plugin returning a CommandResponse cannot
// influence that metadata, so memes never appeared on phones. Uploading the
// render as a real attachment takes the same path as any photo and needs no
// metadata at all.
//
// Deliberately no ServeHTTP: upstream's handler never checked
// Mattermost-User-Id, leaving an unauthenticated image renderer exposed.
type Plugin struct {
	plugin.MattermostPlugin
}

// OnActivate registers /meme. No routes, no background work, no config.
func (p *Plugin) OnActivate() error {
	if err := p.API.RegisterCommand(memeCommandData()); err != nil {
		return fmt.Errorf("failed to register /%s: %w", memeCommand, err)
	}
	return nil
}

// ExecuteCommand renders the meme, uploads it, and posts it as the invoking
// user. Every failure is reported to that user alone as an ephemeral message.
func (p *Plugin) ExecuteCommand(_ *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	input := strings.TrimSpace(strings.TrimPrefix(args.Command, "/"+memeCommand))
	if input == "" || input == "help" {
		return ephemeral(helpText()), nil
	}

	name, data, err := renderMemeJPEG(input)
	if err != nil {
		return ephemeral(err.Error() + " — try `/meme help` for the full list."), nil
	}

	fileInfo, appErr := p.API.UploadFile(data, args.ChannelId, name+".jpg")
	if appErr != nil {
		p.API.LogError("meme upload failed", "template", name, "err", appErr.Error())
		return ephemeral("I rendered that meme but couldn't upload it. Try again?"), nil
	}

	post := &model.Post{
		UserId:    args.UserId,
		ChannelId: args.ChannelId,
		RootId:    args.RootId,
		FileIds:   []string{fileInfo.Id},
	}
	if _, appErr := p.API.CreatePost(post); appErr != nil {
		p.API.LogError("meme post failed", "template", name, "err", appErr.Error())
		return ephemeral("I uploaded that meme but couldn't post it. Try again?"), nil
	}

	return &model.CommandResponse{}, nil
}

// ephemeral builds a response only the invoking user sees.
func ephemeral(text string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         text,
	}
}
