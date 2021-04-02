package plugin

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/kaakaa/mattermost-plugin-reacji/server/reacji"
	"github.com/kaakaa/mattermost-plugin-reacji/server/store"
	"github.com/kaakaa/mattermost-plugin-reacji/server/store/kvstore"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const (
	CommandNameReacji = "reacji"
	botUserName       = "reacji-bot"
	botDisplayName    = "Reacji Bot"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin
	botUserID  string
	reacjiList *reacji.List
	Store      store.Store

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
	ServerConfig  *model.Config
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func (p *Plugin) OnActivate() error {
	p.API.LogDebug("Activate plugin")

	if p.ServerConfig.ServiceSettings.SiteURL == nil {
		return errors.New("siuteURL is not set. Please set a siteURL and restart the plugin")
	}

	bot := &model.Bot{
		Username:    botUserName,
		DisplayName: botDisplayName,
	}
	options := []plugin.EnsureBotOption{
		plugin.ProfileImagePath("assets/logo.dio.png"),
	}
	botUserID, appErr := p.Helpers.EnsureBot(bot, options...)
	if appErr != nil {
		return errors.New(appErr.Error())
	}
	p.botUserID = botUserID

	p.Store = kvstore.NewStore(p.API, p.Helpers)
	reacjiList, err := p.Store.Reacji().Get()
	if err != nil {
		return err
	}
	p.reacjiList = reacjiList
	p.API.LogDebug("store is initialized", "registered", fmt.Sprintf("%v", p.reacjiList.Reacjis))

	if err := p.registerCommand(); err != nil {
		return err
	}
	p.configuration = p.getConfiguration()
	p.API.LogDebug("slash command is initialized")
	return nil
}

func (p *Plugin) OnDeactivate() error {
	if err := p.API.UnregisterCommand("", CommandNameReacji); err != nil {
		return err
	}
	return nil
}

func (p *Plugin) sharePost(reacjis []*reacji.Reacji, post *model.Post, userID string) {
	for _, r := range reacjis {
		shared, err := p.Store.Shared().Get(post.Id, r.ToChannelID, r.DeleteKey)
		if err != nil {
			p.API.LogWarn("failed to get shared post", "post_id", post.Id, "to_channel_id", r.ToChannelID, "details", err.Error())
			continue
		}

		if shared != nil && !p.getConfiguration().AllowDuplicatedSharing {
			// Skip sharing because this post have already been shared
			p.API.LogDebug("this reacji has already been fired", "post_id", post.Id, "to_channel_id", r.ToChannelID, "delete_key", shared.Reacji.DeleteKey)
			continue
		}

		fromChannel, appErr := p.API.GetChannel(r.FromChannelID)
		if appErr != nil {
			p.API.LogWarn("failed to get channel", "channel_id", r.FromChannelID, "error", appErr.Error())
			continue
		}
		team, appErr := p.API.GetTeam(fromChannel.TeamId)
		if appErr != nil {
			p.API.LogWarn("failed to get team", "team_id", fromChannel.TeamId, "error", appErr.Error())
			continue
		}

		p.API.LogDebug("share post", "channel_id", r.ToChannelID, "post_id", post.Id, "user_id", p.botUserID)
		newPost := &model.Post{
			Type:      model.POST_DEFAULT,
			UserId:    p.botUserID,
			ChannelId: r.ToChannelID,
			Message:   fmt.Sprintf("> Shared from ~%s. ([original post](%s))", fromChannel.Name, p.makePostLink(team.Name, post.Id)),
		}
		if newPost, appErr = p.API.CreatePost(newPost); appErr != nil {
			p.API.LogWarn("failed to create post", "error", appErr.Error())
		}

		new := &reacji.SharedPost{
			PostID:       post.Id,
			ToChannelID:  r.ToChannelID,
			SharedPostID: newPost.Id,
			Reacji:       *r,
		}
		if err := p.Store.Shared().Set(new, p.getConfiguration().DaysToKeepSharedRecord); err != nil {
			p.API.LogWarn("failed to set share post", "post_id", post.Id, "to_channel_id", r.ToChannelID, err.Error())
		}
	}
}

func (p *Plugin) makePostLink(teamName, postID string) string {
	return fmt.Sprintf("%s/%s/pl/%s", *p.ServerConfig.ServiceSettings.SiteURL, teamName, postID)
}

func (p *Plugin) ConvertUserIDToDisplayName(userID string) (string, *model.AppError) {
	user, appErr := p.API.GetUser(userID)
	if appErr != nil {
		return "", appErr
	}
	return "@" + user.GetDisplayName(model.SHOW_USERNAME), nil
}

func (p *Plugin) HasAdminPermission(reacji *reacji.Reacji, issuerID string) (bool, *model.AppError) {
	if reacji != nil && issuerID == reacji.Creator {
		return true, nil
	}

	user, appErr := p.API.GetUser(issuerID)
	if appErr != nil {
		return false, appErr
	}
	if user.IsInRole(model.SYSTEM_ADMIN_ROLE_ID) {
		return true, nil
	}
	return false, nil
}

func (p *Plugin) HasPermissionToPrivateChannel(from, to *model.Channel, issuerID string) bool {
	if from.Type != model.CHANNEL_OPEN {
		if !p.API.HasPermissionToChannel(issuerID, from.Id, model.PERMISSION_READ_CHANNEL) {
			return false
		}
	}
	if to.Type != model.CHANNEL_OPEN {
		if !p.API.HasPermissionToChannel(issuerID, to.Id, model.PERMISSION_READ_CHANNEL) {
			return false
		}
	}
	return true
}

// See https://developers.mattermost.com/extend/plugins/server/reference/
