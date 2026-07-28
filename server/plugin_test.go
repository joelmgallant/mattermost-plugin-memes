package main

import (
	"net/http"
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestExecuteCommandHelp(t *testing.T) {
	p := &Plugin{}
	for _, command := range []string{"/meme", "/meme ", "/meme help"} {
		response, appErr := p.ExecuteCommand(nil, &model.CommandArgs{Command: command})
		require.Nil(t, appErr)
		require.NotNil(t, response)
		assert.Equal(t, model.CommandResponseTypeEphemeral, response.ResponseType)
		assert.Contains(t, response.Text, "Available memes:")
	}
}

// An unknown meme is a user typo, not a server fault. Upstream returned an
// AppError here, which surfaces as a generic error banner.
func TestExecuteCommandUnknownMemeIsEphemeral(t *testing.T) {
	p := &Plugin{}
	response, appErr := p.ExecuteCommand(nil, &model.CommandArgs{
		Command: "/meme not-a-real-meme hello",
	})
	require.Nil(t, appErr, "must not surface as a server error")
	assert.Equal(t, model.CommandResponseTypeEphemeral, response.ResponseType)
	assert.Contains(t, response.Text, "/meme help")
}

func TestExecuteCommandUploadsAndPosts(t *testing.T) {
	api := &plugintest.API{}
	defer api.AssertExpectations(t)

	api.On("UploadFile", mock.AnythingOfType("[]uint8"), "channel-1", "brace-yourselves.jpg").
		Return(&model.FileInfo{Id: "file-1"}, nil)

	var posted *model.Post
	api.On("CreatePost", mock.AnythingOfType("*model.Post")).
		Run(func(args mock.Arguments) { posted = args.Get(0).(*model.Post) }).
		Return(&model.Post{Id: "post-1"}, nil)

	p := &Plugin{}
	p.SetAPI(api)

	response, appErr := p.ExecuteCommand(nil, &model.CommandArgs{
		Command:   `/meme brace-yourselves "brace yourself." "memes are coming."`,
		UserId:    "user-1",
		ChannelId: "channel-1",
		RootId:    "root-1",
	})

	require.Nil(t, appErr)
	assert.Empty(t, response.Text, "the meme is the post; the response adds nothing")

	require.NotNil(t, posted)
	assert.Equal(t, "user-1", posted.UserId, "meme posts as the invoking user")
	assert.Equal(t, "channel-1", posted.ChannelId)
	assert.Equal(t, "root-1", posted.RootId, "/meme in a thread must reply in that thread")
	assert.Equal(t, model.StringArray{"file-1"}, posted.FileIds)
	assert.Empty(t, posted.Message, "no markdown image: the file IS the meme")
}

func TestExecuteCommandUploadFailureIsEphemeral(t *testing.T) {
	api := &plugintest.API{}
	defer api.AssertExpectations(t)
	api.On("UploadFile", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, model.NewAppError("UploadFile", "boom", nil, "", 500))
	api.On("LogError", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once()

	p := &Plugin{}
	p.SetAPI(api)

	response, appErr := p.ExecuteCommand(nil, &model.CommandArgs{
		Command:   "/meme doge wow such meme",
		ChannelId: "channel-1",
	})
	require.Nil(t, appErr)
	assert.Equal(t, model.CommandResponseTypeEphemeral, response.ResponseType)
	assert.Equal(t, "I rendered that meme but couldn't upload it. Try again?", response.Text)
}

// httpServer is the exact shape Mattermost dispatches plugin HTTP requests to.
// Declaring it with the real signature is what gives the assertion below teeth:
// an interface written with placeholder argument types could never match, so
// the test would pass no matter what the plugin does.
type httpServer interface {
	ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request)
}

// If this stops compiling, the real hook signature moved and the negative
// assertion below has gone vacuous.
var _ httpServer = (plugin.Hooks)(nil)

// The entire point of the rebuild's security half: no HTTP surface to leave
// unauthenticated. If someone reintroduces a route handler, this fails.
func TestPluginServesNoHTTP(t *testing.T) {
	_, hasServeHTTP := interface{}(&Plugin{}).(httpServer)
	assert.False(t, hasServeHTTP, "this build must not serve HTTP")
}
