package plugin

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"

	"github.com/kaakaa/mattermost-plugin-reacji/server/reacji"
)

// ReactionHasBeenAdded shares a post if registered reaction is attached to a post
func (p *Plugin) ReactionHasBeenAdded(_ *plugin.Context, reaction *model.Reaction) {
	postID := reaction.PostId
	emojiName := reaction.EmojiName

	post, appErr := p.API.GetPost(postID)
	if appErr != nil {
		p.API.LogWarn("failed to get post", "post_id", postID)
		return
	}
	channelID := post.ChannelId
	channel, appErr := p.API.GetChannel(channelID)
	if appErr != nil {
		p.API.LogWarn("failed to get channel", "channel_id", channelID)
		return
	}
	// Don't share the post if reacji-bot doesn't have read permission for channel
	if !p.HasPermissionToChannel(channel, p.botUserID) {
		p.API.LogDebug("stop sharing because reaji-bot doesn't have read permission")
		return
	}

	teamID := channel.TeamId

	var reacjis []*reacji.Reacji
	for _, reacji := range p.reacjiList.Reacjis {
		if reacji.TeamID != teamID {
			continue
		}
		if reacji.FromChannelID == FromAllChannelKeyword || reacji.FromChannelID == channelID {
			if reacji.EmojiName == emojiName {
				reacjis = append(reacjis, reacji)
			}
		}
	}

	go p.sharePost(reacjis, post, channelID)
}

// MessageWillBePosted expand contents of permalink of local post
func (p *Plugin) MessageWillBePosted(_ *plugin.Context, post *model.Post) (*model.Post, string) {
	if _, ok := post.GetProps()[SharedPostPropKey]; !ok {
		return post, ""
	}
	siteURL := p.API.GetConfig().ServiceSettings.SiteURL
	channel, appErr := p.API.GetChannel(post.ChannelId)
	if appErr != nil {
		return post, ""
	}

	if channel.Type == model.ChannelTypeDirect || channel.Type == model.ChannelTypeGroup {
		return post, ""
	}

	team, appErr := p.API.GetTeam(channel.TeamId)
	if appErr != nil {
		return post, ""
	}

	selfLink := fmt.Sprintf("%s/%s", *siteURL, team.Name)
	selfLinkPattern, err := regexp.Compile(fmt.Sprintf("%s%s", selfLink, `/[\w/]+`))
	if err != nil {
		return post, ""
	}

	matches := selfLinkPattern.FindAllString(post.Message, -1)
	if len(matches) != 0 {
		// Only first post matched the pattern is expanded, because can't deal with files that have more than five total attachments.
		match := matches[0]

		separated := strings.Split(match, "/")
		postID := separated[len(separated)-1]
		oldPost, appErr := p.API.GetPost(postID)
		if appErr != nil {
			return post, ""
		}

		newFileIds, appErr := p.API.CopyFileInfos(post.UserId, oldPost.FileIds)
		if appErr != nil {
			p.API.LogWarn("Failed to copy file ids", "error", appErr.Error())
			return post, ""
		}
		// NOTES: if attaching over 5 files, error will occur
		post.FileIds = append(post.FileIds, newFileIds...)

		oldchannel, appErr := p.API.GetChannel(oldPost.ChannelId)
		if appErr != nil {
			return post, ""
		}

		postUser, appErr := p.API.GetUser(oldPost.UserId)
		if appErr != nil {
			return post, ""
		}
		oldPostCreateAt := time.Unix(oldPost.CreateAt/1000, 0)

		AuthorName := postUser.GetDisplayNameWithPrefix(model.ShowNicknameFullName, "@")
		if postUser.IsBot {
			botUser := model.BotFromUser(postUser)
			AuthorName = botUser.DisplayName
		}

		attachment := []*model.SlackAttachment{
			{
				Timestamp:  oldPost.CreateAt,
				AuthorName: AuthorName,
				Text:       oldPost.Message,
				Footer: fmt.Sprintf("Posted in ~%s %s",
					oldchannel.Name,
					oldPostCreateAt.Format("on Mon 2 Jan 2006 at 15:04:05 MST"),
				),
			},
			nil,
		}
		model.ParseSlackAttachment(post, attachment)
	}

	return post, ""
}
